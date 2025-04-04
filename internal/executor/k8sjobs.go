package executor

import (
	"context"
	"fmt"
	"log"
	"strings"
	pb "terraform-executor/api/proto"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// Helper function to create Terraform job
func (s *ExecutorService) createTerraformJobTemplate(ctx context.Context, name, namespace, project string, runType string, args []string) (*batchv1.Job, error) {
	// get env vars from secret
	secretName := fmt.Sprintf("%s.%s", project, "env")
	envVars := []corev1.EnvVar{}
	if secret, err := s.K8sClient.GetSecret(ctx, namespace, secretName); err == nil {
		for key, value := range secret.Data {
			envVars = append(envVars, corev1.EnvVar{
				Name:  key,
				Value: string(value),
			})
		}
		// Ignore error if secret does not exist
	}
	envVars = append(
		envVars,
		corev1.EnvVar{
			Name:  "TF_IN_AUTOMATION",
			Value: "true",
		},
		corev1.EnvVar{
			Name:  "TF_PLUGIN_CACHE_DIR",
			Value: "/root/.terraform.d/plugin-cache",
		},
	)

	// dynamic volume mounts
	volumeMounts := []corev1.VolumeMount{}
	pluginsVolumeMount := corev1.VolumeMount{
		Name:      "plugin-cache",
		MountPath: "/root/.terraform.d/plugin-cache",
	}
	awsCredsMount := corev1.VolumeMount{
		Name:      "aws-creds",
		MountPath: "/root/.aws",
	}
	volumeMounts = append(volumeMounts, pluginsVolumeMount, awsCredsMount)
	for _, vol := range []string{"main.tf", "versions.tf", "variables.tf"} {
		// add volume mounts for main.tf, versions.tf, and variables.tf if config maps exist
		if _, err := s.K8sClient.GetConfigMap(ctx, namespace, fmt.Sprintf("%s.%s", project, vol)); err == nil {
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      strings.TrimSuffix(vol, ".tf"),
				MountPath: fmt.Sprintf("/root/%s", vol),
				SubPath:   vol,
			})
		}
	}

	// dynamic volumes
	volumes := []corev1.Volume{
		{
			Name: "plugin-cache",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: fmt.Sprintf("%s-plugin-cache", namespace),
					ReadOnly:  false,
				},
			},
		},
		{
			Name: "aws-creds",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "aws-profile",
					Items: []corev1.KeyToPath{
						{
							Key:  "credentials",
							Path: "credentials",
						},
					},
				},
			},
		},
	}

	for _, vol := range []string{"main.tf", "versions.tf", "variables.tf"} {
		if _, err := s.K8sClient.GetConfigMap(ctx, namespace, fmt.Sprintf("%s.%s", project, vol)); err == nil {
			volumes = append(volumes, corev1.Volume{
				Name: strings.TrimSuffix(vol, ".tf"),
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: fmt.Sprintf("%s.%s", project, vol),
						},
					},
				},
			})
		}
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app":     "terraform-executor",
				"user":    namespace,
				"project": project,
				"type":    runType,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            ptr.To[int32](0),   // Don't retry on failure
			TTLSecondsAfterFinished: ptr.To[int32](120), // Automatically delete job after 120 seconds
			Completions:             ptr.To[int32](1),   // Only one successful pod needed
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     "terraform-executor",
						"user":    namespace,
						"project": project,
						"type":    runType,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:       "runner",
							WorkingDir: "/root",
							Image:      "hashicorp/terraform:latest",
							Command: []string{
								"/bin/sh",
								"-c",
								fmt.Sprintf(
									"terraform init -no-color -input=false && terraform %s",
									strings.Join(args, " "),
								),
							},
							Env:          envVars,
							VolumeMounts: volumeMounts,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("1Gi"),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("200m"),
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
						},
					},
					Volumes: volumes,
				},
			},
		},
	}, nil
}

// getPodLogs gets logs from both init and runner containers if there is init container
func (s *ExecutorService) getPodLogs(ctx context.Context, userId string, pod *corev1.Pod) (string, error) {
	// check if pod has init container
	var allLogs []string

	// Get init container logs
	initLogs, err := s.K8sClient.GetPodLogs(ctx, userId, pod.Name, "init")
	if err == nil && strings.TrimSpace(initLogs) != "" {
		allLogs = append(allLogs, fmt.Sprintf("Init container logs:\n%s", initLogs))
	}

	// Get runner container logs
	runnerLogs, err := s.K8sClient.GetPodLogs(ctx, userId, pod.Name, "runner")
	if err == nil && strings.TrimSpace(runnerLogs) != "" {
		allLogs = append(allLogs, fmt.Sprintf("Runner container logs:\n%s", runnerLogs))
	}

	if len(allLogs) == 0 {
		return "", nil
	}

	return strings.Join(allLogs, "\n\n"), nil
}

func (s *ExecutorService) waitForJobAndGetLogs(ctx context.Context, userId, project, requestId, jobName string) (string, error) {
	var result struct {
		logs string
		err  error
	}

	// First check if job is already completed
	job, err := s.K8sClient.GetJob(ctx, userId, jobName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			fmt.Printf("Job %s not found, waiting for creation\n", jobName)
		} else {
			return "", fmt.Errorf("failed to get initial job status: %v", err)
		}
	} else {
		if job.Status.Failed > 0 {
			pod, err := s.K8sClient.GetJobPod(ctx, userId, jobName)
			if err != nil {
				return "", fmt.Errorf("failed to get job pod: %v", err)
			}
			return s.getPodLogs(ctx, userId, pod)
		}
		if job.Status.Succeeded > 0 {
			pod, err := s.K8sClient.GetJobPod(ctx, userId, jobName)
			if err != nil {
				return "", fmt.Errorf("failed to get job pod: %v", err)
			}
			return s.getPodLogs(ctx, userId, pod)
		}
	}

	err = s.streamPodLogsAndSendRPC(ctx, userId, project, requestId, jobName, *s.LogStream)
	if err != nil {
		return "", fmt.Errorf("streamPodLogsAndSendRPC failed to stream pod logs: %v", err)
	}
	watcher, err := s.K8sClient.WatchJob(ctx, userId, jobName)
	if err != nil {
		return "", fmt.Errorf("failed to watch job: %v", err)
	}
	defer watcher.Stop()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for event := range watcher.ResultChan() {
			job, ok := event.Object.(*batchv1.Job)
			if !ok {
				continue
			}

			if job.Status.Failed > 0 || job.Status.Succeeded > 0 {
				fmt.Printf("Job %s completed with status: Failed=%d, Succeeded=%d\n",
					jobName, job.Status.Failed, job.Status.Succeeded)

				pod, err := s.K8sClient.GetJobPod(ctx, userId, jobName)
				if err != nil {
					result.err = fmt.Errorf("failed to get job pod: %v", err)
					return
				}
				logs, err := s.getPodLogs(ctx, userId, pod)
				if err != nil {
					result.err = fmt.Errorf("failed to get pod logs: %v", err)
					return
				}
				result.logs = logs
				return
			}
		}
		result.err = fmt.Errorf("watch ended unexpectedly")
	}()

	select {
	case <-done:
		return result.logs, result.err
	case <-time.After(15 * time.Minute):
		return "", fmt.Errorf("job timed out after 15 minutes")
	}
}

func (s *ExecutorService) streamPodLogsAndSendRPC(ctx context.Context, userId, project, requestId string, jobName string, stream pb.Executor_StreamLogsServer) error {
	// Set up polling mechanism for logs (every 1 second)
	logTicker := time.NewTicker(1 * time.Second)
	defer logTicker.Stop()
	var lastLog string // Keeps track of previously sent logs to avoid duplication

	completionChan := make(chan error, 1) // Channel to signal completion

	go func() {
		failedAttempts := 0 // Counter for failed pod retrieval attempts

		for {
			select {
			case <-ctx.Done():
				// Stop streaming if context is cancelled
				completionChan <- ctx.Err()
				return
			case <-logTicker.C:
				podFresh, err := s.K8sClient.GetJobPod(ctx, userId, jobName)
				if err != nil {
					log.Printf("Error getting pod: %v", err)
					failedAttempts++
					if failedAttempts >= 5 {
						log.Print("failed to get pod 5 times. exiting")
						completionChan <- nil
						return
					}
					continue // Corrected: Continue to the next tick, don't return
				}

				failedAttempts = 0 // Reset counter on successful retrieval

				if podFresh.Status.Phase == corev1.PodSucceeded || podFresh.Status.Phase == corev1.PodFailed {
					log.Printf("pod completed, stopping ticker")
					completionChan <- nil
					return
				}

				if podFresh.Status.Phase != corev1.PodRunning {
					// Pod isn't in a state where logs can be retrieved
					fmt.Println("Pod is not in a running state. Will retry.")
					continue
				}
				// Pod is running, now try to fetch logs
				runnerLogs, err := s.getPodLogs(ctx, userId, podFresh)
				if err != nil {
					// Log the error and continue retrying
					fmt.Printf("Error retrieving logs: %v\n", err)
					continue // Retry after the next tick
				}
				// Check if there are new logs (by comparing with last sent logs)
				if runnerLogs != lastLog {
					// Extract the new portion of logs by removing the prefix that matches previous logs
					logDiff := strings.TrimPrefix(runnerLogs, lastLog)
					lastLog = runnerLogs // Update the last seen logs

					if err := stream.Send(&pb.LogStreamResponse{LogLine: logDiff, UserId: userId, Project: project, RequestId: requestId}); err != nil {
						log.Printf("Error sending logs via stream: %v\n", err)
						continue
					}
				}
			}
		}
	}()
	return <-completionChan // Wait for completion signal
}
