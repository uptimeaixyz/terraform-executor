package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// NamespaceExists checks if a namespace exists
func (c *K8sClient) NamespaceExists(ctx context.Context, name string) (bool, error) {
	_, err := c.clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if statusError, ok := err.(*errors.StatusError); ok && statusError.Status().Code == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateNamespace creates a new namespace
func (c *K8sClient) CreateNamespace(ctx context.Context, name string) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := c.clientset.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	return err
}

// DeleteNamespace deletes a namespace
func (c *K8sClient) DeleteNamespace(ctx context.Context, name string) error {
	return c.clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
}

// GetConfigMap retrieves a ConfigMap from the specified namespace
func (c *K8sClient) GetConfigMap(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error) {
	return c.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateConfigMap creates a new ConfigMap in the specified namespace
func (c *K8sClient) CreateConfigMap(ctx context.Context, namespace string, cm *corev1.ConfigMap) error {
	_, err := c.clientset.CoreV1().ConfigMaps(namespace).Create(ctx, cm, metav1.CreateOptions{})
	return err
}

// UpdateConfigMap updates an existing ConfigMap in the specified namespace
func (c *K8sClient) UpdateConfigMap(ctx context.Context, namespace string, cm *corev1.ConfigMap) error {
	_, err := c.clientset.CoreV1().ConfigMaps(namespace).Update(ctx, cm, metav1.UpdateOptions{})
	return err
}

// DeleteConfigMap deletes a ConfigMap from the specified namespace
func (c *K8sClient) DeleteConfigMap(ctx context.Context, namespace, name string) error {
	return c.clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// GetSecret retrieves a Secret from the specified namespace
func (c *K8sClient) GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error) {
	return c.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

// CreateSecret creates a new Secret in the specified namespace
func (c *K8sClient) CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret) error {
	_, err := c.clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	return err
}

// UpdateSecret updates an existing Secret in the specified namespace
func (c *K8sClient) UpdateSecret(ctx context.Context, namespace string, secret *corev1.Secret) error {
	_, err := c.clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	return err
}

// DeleteSecret deletes a Secret from the specified namespace
func (c *K8sClient) DeleteSecret(ctx context.Context, namespace, name string) error {
	return c.clientset.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// CreateAWSCredsSecret creates a secret containing AWS credentials in INI format
func (c *K8sClient) CreateAWSCredsSecret(ctx context.Context, namespace, accessKey, secretKey, sessionToken string, expirationTime time.Time) error {
	credentialsContent := fmt.Sprintf("[tfstate]\naws_access_key_id = %s\naws_secret_access_key = %s\naws_session_token = %s",
		accessKey, secretKey, sessionToken)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-profile",
			Namespace: namespace,
			Labels: map[string]string{
				"expirationDate": expirationTime.Format("20060102-150405"),
			},
		},
		StringData: map[string]string{
			"credentials": credentialsContent,
		},
	}

	_, err := c.clientset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		// If secret already exists, update it
		_, err = c.clientset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	}
	return err
}

// CreateJob creates a new Job in the specified namespace
func (c *K8sClient) CreateJob(ctx context.Context, namespace string, job *batchv1.Job) (*batchv1.Job, error) {
	return c.clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
}

// GetJob retrieves a Job from the specified namespace
func (c *K8sClient) GetJob(ctx context.Context, namespace, name string) (*batchv1.Job, error) {
	return c.clientset.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
}

// DeleteJob deletes a Job from the specified namespace
func (c *K8sClient) DeleteJob(ctx context.Context, namespace, name string) error {
	return c.clientset.BatchV1().Jobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// WatchJob watches for job events in the specified namespace
func (c *K8sClient) WatchJob(ctx context.Context, namespace, name string) (watch.Interface, error) {
	fmt.Printf("🔍 Starting WatchJob for namespace=%s, name=%s\n", namespace, name)

	// Check if job exists, but don't fail if it doesn't
	job, err := c.GetJob(ctx, namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("⏳ Job %s not found yet, but will watch for it\n", name)
		} else {
			fmt.Printf("⚠️ Error checking job existence (non-fatal): %v\n", err)
		}
	} else {
		fmt.Printf("✅ Job found before watching. Current status: %+v\n", job.Status)
	}

	// Create the watch interface - this will catch the job creation event when it happens
	watcher, err := c.clientset.BatchV1().Jobs(namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
	})
	if err != nil {
		fmt.Printf("❌ Error creating watcher: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ Watch interface created successfully for job %s\n", name)

	return watcher, nil
}

// GetPodLogs retrieves logs from a pod in the specified namespace
func (c *K8sClient) GetPodLogs(ctx context.Context, namespace, podName, containerName string) (string, error) {
	req := c.clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: containerName,
	})
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ListPods lists all pods in namespace matching the label selector
func (c *K8sClient) ListPods(ctx context.Context, namespace string, labelSelector string) (*corev1.PodList, error) {
	fmt.Printf("🔍 Listing pods in namespace=%s with labels=%s\n", namespace, labelSelector)
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		fmt.Printf("❌ Error listing pods: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ Found %d pods\n", len(pods.Items))
	return pods, nil
}

// GetJobPod returns the pod associated with a job
func (c *K8sClient) GetJobPod(ctx context.Context, namespace, jobName string) (*corev1.Pod, error) {
	fmt.Printf("🔍 Looking for pod for job=%s\n", jobName)
	pods, err := c.ListPods(ctx, namespace, fmt.Sprintf("job-name=%s", jobName))
	if err != nil {
		return nil, err
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no pods found for job %s", jobName)
	}
	fmt.Printf("✅ Found pod: %s\n", pods.Items[0].Name)
	return &pods.Items[0], nil
}

// Create pvc creates a new PersistentVolumeClaim in the specified namespace
func (c *K8sClient) CreatePVC(ctx context.Context, namespace string, pvc *corev1.PersistentVolumeClaim) error {
	_, err := c.clientset.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, metav1.CreateOptions{})
	return err
}
