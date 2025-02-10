package k8s

import (
	"os"
	"path/filepath"

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
