package executor

import (
	"context"
	"fmt"
	"os"
	pb "terraform-executor/api/proto"
	"terraform-executor/internal/awsclient"
	"terraform-executor/internal/k8s"
)

// ExecutorService implements the ExecutorServer interface.
type ExecutorService struct {
	pb.UnimplementedExecutorServer
	K8sClient *k8s.K8sClient
	AWSClient *awsclient.AWSClient
	LogStream *pb.Executor_StreamLogsServer
	Bucket    string
	Debug     bool
	ctx       context.Context
}

func NewExecutorService(ctx context.Context) (*ExecutorService, error) {
	// Initialize Kubernetes client
	kubeconfig := os.Getenv("KUBECONFIG")
	k8sClient, err := k8s.NewK8sClient(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}
	fmt.Println("Kubernetes client created")

	// Initialize AWS client with application context
	awsClient, err := awsclient.NewAWSClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %v", err)
	}
	fmt.Println("AWS client created")
	bucket := os.Getenv("BUCKET_NAME")
	// fallback to default bucket
	if bucket == "" {
		bucket = "uptimeai-test-bucket"
	}
	return &ExecutorService{
		K8sClient: k8sClient,
		AWSClient: awsClient,
		ctx:       ctx,
		Bucket:    bucket,
	}, nil
}
