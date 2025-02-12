package executor

import (
	"context"
	"fmt"
	pb "terraform-executor/api/proto"
	"time"

	"terraform-executor/internal/awsclient"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ensureNamespace ensures that a namespace exists for the given user
func (s *ExecutorService) ensureNamespace(ctx context.Context, userId string) error {
	exists, err := s.K8sClient.NamespaceExists(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to check namespace: %v", err)
	}
	if !exists {
		if err := s.K8sClient.CreateNamespace(ctx, userId); err != nil {
			return fmt.Errorf("failed to create namespace: %v", err)
		}
	}
	return nil
}

// ensureUserRole ensures that the IAM role for the user exists, creates it if it doesn't
func (s *ExecutorService) ensureUserRole(ctx context.Context, userId string) error {
	exists, err := s.AWSClient.RoleExists(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}

	if !exists {
		// Get caller identity to get the ARN of the current user
		identity, err := s.AWSClient.STSClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return fmt.Errorf("failed to get caller identity: %w", err)
		}

		trustPolicy := fmt.Sprintf(`{
            "Version": "2012-10-17",
            "Statement": [{
                "Effect": "Allow",
                "Principal": {
                    "AWS": "%s"
                },
                "Action": "sts:AssumeRole"
            }]
        }`, *identity.Arn) // Using the actual caller's ARN

		_, err = s.AWSClient.CreateRole(ctx, userId, trustPolicy)
		if err != nil {
			return fmt.Errorf("failed to create role: %w", err)
		}
	}

	return nil
}

// ensureAWSCredentials ensures that AWS credentials secret exists and is not expired
func (s *ExecutorService) ensureAWSCredentials(ctx context.Context, userId string) error {
	// Check if secret exists and get its expiration time
	secret, err := s.K8sClient.GetSecret(ctx, userId, "aws-profile")
	if err == nil {
		// Secret exists, check expiration
		expirationStr := secret.Labels["expirationDate"]
		expiration, err := time.Parse("20060102-150405", expirationStr)
		if err != nil {
			return fmt.Errorf("invalid expiration date format: %w", err)
		}

		// If credentials expire in less than 10 minutes, refresh them
		if time.Until(expiration) > 10*time.Minute {
			return nil // Credentials are still valid
		}
		// Continue to refresh credentials
	}

	// Get role ARN for the user
	id, _ := s.AWSClient.GetAccountID(ctx)
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role%s%s", id, awsclient.RolePath, userId)
	accessKey, secretKey, sessionToken, expiration, err := s.AWSClient.GetTemporaryCredentials(
		ctx,
		roleArn,
		fmt.Sprintf("terraform-%s", userId),
	)
	if err != nil {
		return fmt.Errorf("failed to get temporary credentials: %w", err)
	}

	// Include session token in credentials file
	err = s.K8sClient.CreateAWSCredsSecret(ctx, userId, accessKey, secretKey, sessionToken, expiration)
	if err != nil {
		return fmt.Errorf("failed to create AWS credentials secret: %w", err)
	}

	return nil
}

// ensurePVC ensures that a PVC exists for plugin cache
func (s *ExecutorService) ensurePVC(ctx context.Context, namespace string) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-plugin-cache", namespace),
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	err := s.K8sClient.CreatePVC(ctx, namespace, pvc)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create plugins PVC: %v", err)
	}
	return nil
}

// ensureResources ensures all required resources exist for the given user
func (s *ExecutorService) ensureResources(ctx context.Context, userId string) error {
	// Ensure namespace exists
	if err := s.ensureNamespace(ctx, userId); err != nil {
		return fmt.Errorf("namespace error: %w", err)
	}

	// Ensure AWS role exists
	if err := s.ensureUserRole(ctx, userId); err != nil {
		return fmt.Errorf("AWS role error: %w", err)
	}

	// Ensure AWS credentials are fresh
	if err := s.ensureAWSCredentials(ctx, userId); err != nil {
		return fmt.Errorf("AWS credentials error: %w", err)
	}

	// Ensure PVC exists
	if err := s.ensurePVC(ctx, userId); err != nil {
		return fmt.Errorf("PVC error: %w", err)
	}

	return nil
}

// Plan generates a Terraform plan and returns the result.
func (s *ExecutorService) Plan(ctx context.Context, req *pb.PlanRequest) (*pb.PlanResponse, error) {
	if err := s.ensureResources(ctx, req.UserId); err != nil {
		return &pb.PlanResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	jobName := fmt.Sprintf("terraform-plan-%s-%s", req.UserId, time.Now().Format("20060102150405"))
	job, err := s.createTerraformJobTemplate(ctx, jobName, req.UserId, req.Project, "plan", []string{"plan", "-input=false", "-no-color"})
	if err != nil {
		if s.Debug {
			fmt.Printf("Job template that failed:\n%+v\n", job)
		}
		return &pb.PlanResponse{
			Success: false,
			Error:   fmt.Sprintf("job creation error: %v", err),
		}, nil
	}

	_, err = s.K8sClient.CreateJob(ctx, req.UserId, job)
	if err != nil {
		if s.Debug {
			fmt.Printf("Job template that failed:\n%+v\n", job)
		}
		return &pb.PlanResponse{
			Success: false,
			Error:   fmt.Sprintf("kubernetes job error: %v", err),
		}, nil
	}

	// Wait for job completion and get logs
	output, err := s.waitForJobAndGetLogs(ctx, req.UserId, jobName)
	if err != nil {
		return &pb.PlanResponse{
			Success:    false,
			Error:      err.Error(),
			PlanOutput: output,
		}, nil
	}

	return &pb.PlanResponse{
		Success:    true,
		PlanOutput: output,
	}, nil
}

// Apply applies a Terraform plan or configuration and returns the result.
func (s *ExecutorService) Apply(ctx context.Context, req *pb.ApplyRequest) (*pb.ApplyResponse, error) {
	if err := s.ensureResources(ctx, req.UserId); err != nil {
		return &pb.ApplyResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	jobName := fmt.Sprintf("terraform-apply-%s", time.Now().Format("20060102150405"))
	job, err := s.createTerraformJobTemplate(ctx, jobName, req.UserId, req.Project, "apply", []string{"apply", "-auto-approve", "-input=false", "-no-color"})
	if err != nil {
		return &pb.ApplyResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create job template: %v", err),
		}, nil
	}
	_, err = s.K8sClient.CreateJob(ctx, req.UserId, job)
	if err != nil {
		return &pb.ApplyResponse{
			Success: false,
			Error:   fmt.Sprintf("kubernetes job creation failed: %v", err),
		}, nil
	}
	output, err := s.waitForJobAndGetLogs(ctx, req.UserId, jobName)
	if err != nil {
		if s.Debug {
			fmt.Printf("⚠️ Job execution completed with error: %v\nOutput: %s\n", err, output)
		}
		return &pb.ApplyResponse{
			Success:     false,
			Error:       fmt.Sprintf("job execution failed: %v", err),
			ApplyOutput: output,
		}, nil
	}
	return &pb.ApplyResponse{
		Success:     true,
		ApplyOutput: output,
	}, nil
}

// Destroy destroys the Terraform-managed infrastructure and returns the result.
func (s *ExecutorService) Destroy(ctx context.Context, req *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	if err := s.ensureResources(ctx, req.UserId); err != nil {
		return &pb.DestroyResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	jobName := fmt.Sprintf("terraform-destroy-%s", time.Now().Format("20060102150405"))
	job, err := s.createTerraformJobTemplate(ctx, jobName, req.UserId, req.Project, "destroy", []string{"destroy", "-auto-approve", "-input=false", "-no-color"})
	if err != nil {
		return &pb.DestroyResponse{Success: false, Error: err.Error()}, nil
	}

	_, err = s.K8sClient.CreateJob(ctx, req.UserId, job)
	if err != nil {
		return &pb.DestroyResponse{
			Success: false,
			Error:   fmt.Sprintf("kubernetes job creation failed: %v", err),
		}, nil
	}

	// Wait for job completion and get logs
	output, err := s.waitForJobAndGetLogs(ctx, req.UserId, jobName)
	if err != nil {
		return &pb.DestroyResponse{
			Success:       false,
			Error:         fmt.Sprintf("job execution failed: %v", err),
			DestroyOutput: output,
		}, nil
	}

	return &pb.DestroyResponse{Success: true, DestroyOutput: output}, nil
}

// GetStateList returns output of "terraform state list" command
func (s *ExecutorService) GetStateList(ctx context.Context, req *pb.GetStateListRequest) (*pb.GetStateListResponse, error) {
	if err := s.ensureResources(ctx, req.UserId); err != nil {
		return &pb.GetStateListResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	jobName := fmt.Sprintf("terraform-state-list-%s-%s", req.UserId, time.Now().Format("20060102150405"))
	job, err := s.createTerraformJobTemplate(ctx, jobName, req.UserId, req.Project, "state-list", []string{"state", "list", "-no-color"})
	if err != nil {
		return &pb.GetStateListResponse{
			Success: false,
			Error:   fmt.Sprintf("job creation error: %v", err),
		}, nil
	}

	_, err = s.K8sClient.CreateJob(ctx, req.UserId, job)
	if err != nil {
		return &pb.GetStateListResponse{
			Success: false,
			Error:   fmt.Sprintf("kubernetes job error: %v", err),
		}, nil
	}

	output, err := s.waitForJobAndGetLogs(ctx, req.UserId, jobName)
	if err != nil {
		return &pb.GetStateListResponse{
			Success:         false,
			Error:           fmt.Sprintf("job execution failed: %v", err),
			StateListOutput: output,
		}, nil
	}
	return &pb.GetStateListResponse{
		Success:         true,
		StateListOutput: output,
	}, nil
}

// GetTFShow returns output of "terraform state show" command
func (s *ExecutorService) GetTFShow(ctx context.Context, req *pb.GetTFShowRequest) (*pb.GetTFShowResponse, error) {
	if err := s.ensureResources(ctx, req.UserId); err != nil {
		return &pb.GetTFShowResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	jobName := fmt.Sprintf("terraform-show-%s-%s", req.UserId, time.Now().Format("20060102150405"))
	job, err := s.createTerraformJobTemplate(ctx, jobName, req.UserId, req.Project, "show", []string{"show", "-json", "-no-color"})
	if err != nil {
		return &pb.GetTFShowResponse{
			Success: false,
			Error:   fmt.Sprintf("job creation error: %v", err),
		}, nil
	}

	_, err = s.K8sClient.CreateJob(ctx, req.UserId, job)
	if err != nil {
		return &pb.GetTFShowResponse{
			Success: false,
			Error:   fmt.Sprintf("kubernetes job error: %v", err),
		}, nil
	}

	output, err := s.waitForJobAndGetLogs(ctx, req.UserId, jobName)
	if err != nil {
		return &pb.GetTFShowResponse{
			Success: false,
			Error:   fmt.Sprintf("job execution failed: %v", err),
			Content: output,
		}, nil
	}
	return &pb.GetTFShowResponse{
		Success: true,
		Content: output,
	}, nil
}
