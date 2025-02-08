package executor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	pb "terraform-executor/api/proto"
)

// Create new context
func (s *ExecutorService) CreateContext(ctx context.Context, req *pb.CreateContextRequest) (*pb.CreateContextResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.CreateContextResponse{Success: false, Error: err.Error()}, nil
	}
	contextDir := filepath.Join("./data/", req.Context)

	// Ensure the context directory exists
	if err := os.MkdirAll(contextDir, os.ModePerm); err != nil {
		return &pb.CreateContextResponse{Success: false, Error: fmt.Sprintf("failed to create context: %v", err)}, nil
	}

	return &pb.CreateContextResponse{Success: true}, nil
}

// Delete context
func (s *ExecutorService) DeleteContext(ctx context.Context, req *pb.DeleteContextRequest) (*pb.DeleteContextResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.DeleteContextResponse{Success: false, Error: err.Error()}, nil
	}
	contextDir := filepath.Join("./data/", req.Context)

	// Remove the context directory
	if err := os.RemoveAll(contextDir); err != nil {
		return &pb.DeleteContextResponse{Success: false, Error: fmt.Sprintf("failed to delete context: %v", err)}, nil
	}

	return &pb.DeleteContextResponse{Success: true}, nil
}

// Create workspace
func (s *ExecutorService) CreateWorkspace(ctx context.Context, req *pb.CreateWorkspaceRequest) (*pb.CreateWorkspaceResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.CreateWorkspaceResponse{Success: false, Error: err.Error()}, nil
	}
	workspaceDir := filepath.Join("./data/users", req.UserId, req.Context, req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.CreateWorkspaceResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	return &pb.CreateWorkspaceResponse{Success: true}, nil
}

// Delete workspace
func (s *ExecutorService) DeleteWorkspace(ctx context.Context, req *pb.DeleteWorkspaceRequest) (*pb.DeleteWorkspaceResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.DeleteWorkspaceResponse{Success: false, Error: err.Error()}, nil
	}
	workspaceDir := filepath.Join("./data/users", req.UserId, req.Context, req.Workspace)

	// Remove the workspace directory
	if err := os.RemoveAll(workspaceDir); err != nil {
		return &pb.DeleteWorkspaceResponse{Success: false, Error: fmt.Sprintf("failed to delete workspace: %v", err)}, nil
	}

	return &pb.DeleteWorkspaceResponse{Success: true}, nil
}
