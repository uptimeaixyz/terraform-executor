package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	pb "terraform-executor/api/proto"
	"terraform-executor/pkg/utils"
)

// ExecutorService implements the ExecutorServer interface.
type ExecutorService struct {
	pb.UnimplementedExecutorServer
}

// Plan generates a Terraform plan and returns the result.
func (s *ExecutorService) Plan(ctx context.Context, req *pb.PlanRequest) (*pb.PlanResponse, error) {
	workspaceDir := filepath.Join("./", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.PlanResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Append the code to main.tf
	mainTfPath := filepath.Join(workspaceDir, "main.tf")
	if err := utils.AppendToFile(mainTfPath, req.Code); err != nil {
		return &pb.PlanResponse{Success: false, Error: fmt.Sprintf("failed to write to main.tf: %v", err)}, nil
	}

	// Run `terraform init`
	_, err := runCommand(workspaceDir, "terraform", "init")

	// Run `terraform plan`
	output, err := runCommand(workspaceDir, "terraform", "plan", "-no-color")
	if err != nil {
		return &pb.PlanResponse{Success: false, PlanOutput: output, Error: err.Error()}, nil
	}

	// Clear the main.tf file
	clearMainTf(workspaceDir)
	
	return &pb.PlanResponse{Success: true, PlanOutput: output}, nil
}

// Apply applies a Terraform plan or configuration and returns the result.
func (s *ExecutorService) Apply(ctx context.Context, req *pb.ApplyRequest) (*pb.ApplyResponse, error) {
	workspaceDir := filepath.Join("./", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.ApplyResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Append code to main.tf if provided
	if req.Code != "" {
		mainTfPath := filepath.Join(workspaceDir, "main.tf")
		if err := utils.AppendToFile(mainTfPath, req.Code); err != nil {
			return &pb.ApplyResponse{Success: false, Error: fmt.Sprintf("failed to write to main.tf: %v", err)}, nil
		}
	}

	// Run `terraform init`
	_, err := runCommand(workspaceDir, "terraform", "init")

	// Run `terraform apply`
	output, err := runCommand(workspaceDir, "terraform", "apply", "-auto-approve", "-no-color")
	if err != nil {
		return &pb.ApplyResponse{Success: false, ApplyOutput: output, Error: err.Error()}, nil
	}

	return &pb.ApplyResponse{Success: true, ApplyOutput: output}, nil
}

// Destroy destroys the Terraform-managed infrastructure and returns the result.
func (s *ExecutorService) Destroy(ctx context.Context, req *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	workspaceDir := filepath.Join("./", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.DestroyResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Run `terraform init`
	_, err := runCommand(workspaceDir, "terraform", "init")

	// Run `terraform destroy`
	output, err := runCommand(workspaceDir, "terraform", "destroy", "-auto-approve", "-no-color")
	if err != nil {
		return &pb.DestroyResponse{Success: false, DestroyOutput: output, Error: err.Error()}, nil
	}

	return &pb.DestroyResponse{Success: true, DestroyOutput: output}, nil
}

// GetStateList returns output of "terraform state list" command
func (s *ExecutorService) GetStateList(ctx context.Context, req *pb.GetStateListRequest) (*pb.GetStateListResponse, error) {
	workspaceDir := filepath.Join("./", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.GetStateListResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Run `terraform init`
	_, err := runCommand(workspaceDir, "terraform", "init")

	// Run `terraform state list`
	output, err := runCommand(workspaceDir, "terraform", "state", "list")
	if err != nil {
		return &pb.GetStateListResponse{Success: false, StateListOutput: output, Error: err.Error()}, nil
	}

	return &pb.GetStateListResponse{Success: true, StateListOutput: output}, nil
}

// Clear removes all created files in the workspace directory.
func (s *ExecutorService) Clear(ctx context.Context, req *pb.ClearRequest) (*pb.ClearResponse, error) {
	workspaceDir := filepath.Join("./", req.Workspace)

	// Remove main.tf
	if err := clearMainTf(workspaceDir); err != nil {
		return &pb.ClearResponse{Success: false, Error: fmt.Sprintf("failed to clear main.tf: %v", err)}, nil
	}

	return &pb.ClearResponse{Success: true}, nil
}

// clear main.tf file
func clearMainTf(workspaceDir string) error {
	mainTfPath := filepath.Join(workspaceDir, "main.tf")
	if err := os.Remove(mainTfPath); err != nil {
		return fmt.Errorf("failed to remove main.tf: %v", err)
	}
	return nil
}

// runCommand runs a command in the specified directory and returns its output.
func runCommand(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	return string(output), err
}
