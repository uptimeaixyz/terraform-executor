package health

import (
	"context"
	pb "terraform-executor/api/proto"
	"terraform-executor/internal/awsclient"
	"terraform-executor/internal/k8s"
)

type HealthService struct {
	pb.UnimplementedHealthServer
	k8sClient *k8s.K8sClient
	awsClient *awsclient.AWSClient
}

func NewHealthService(k8sClient *k8s.K8sClient, awsClient *awsclient.AWSClient) *HealthService {
	return &HealthService{
		k8sClient: k8sClient,
		awsClient: awsClient,
	}
}

func (s *HealthService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	response := &pb.HealthCheckResponse{
		Status:     pb.HealthCheckResponse_SERVING,
		Components: make(map[string]pb.HealthCheckResponse_ServingStatus),
		Errors:     make(map[string]string),
	}

	// Check Kubernetes connectivity
	if err := s.k8sClient.HealthCheck(ctx); err != nil {
		response.Components["kubernetes"] = pb.HealthCheckResponse_NOT_SERVING
		response.Errors["kubernetes"] = err.Error()
		response.Status = pb.HealthCheckResponse_NOT_SERVING
	} else {
		response.Components["kubernetes"] = pb.HealthCheckResponse_SERVING
	}

	// Check AWS connectivity
	if _, err := s.awsClient.GetAccountID(ctx); err != nil {
		response.Components["aws"] = pb.HealthCheckResponse_NOT_SERVING
		response.Errors["aws"] = err.Error()
		response.Status = pb.HealthCheckResponse_NOT_SERVING
	} else {
		response.Components["aws"] = pb.HealthCheckResponse_SERVING
	}

	return response, nil
}
