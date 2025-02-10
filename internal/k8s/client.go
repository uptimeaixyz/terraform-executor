package k8s

import (
	"context"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}

func NewK8sClient(kubeconfigPath string) (*K8sClient, error) {
	var config *rest.Config
	var err error

	if kubeconfigPath != "" {
		// Use provided kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		// Try in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			// Fall back to default kubeconfig location
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			defaultKubeconfig := filepath.Join(homeDir, ".kube", "config")
			config, _ = clientcmd.BuildConfigFromFlags("", defaultKubeconfig)
		}
	}

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sClient{
		clientset: clientset,
	}, nil
}

// HealthCheck verifies connectivity to the Kubernetes cluster
func (c *K8sClient) HealthCheck(ctx context.Context) error {
	// Try to list namespaces as a basic connectivity test
	_, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	return err
}
