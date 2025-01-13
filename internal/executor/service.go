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

	// Run `terraform plan`
	output, err := runCommand(workspaceDir, "terraform", "plan", "-no-color", "-out", "plan.out")
	if err != nil {
		return &pb.PlanResponse{Success: false, PlanOutput: output, Error: err.Error()}, nil
	}

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

	// Prepare arguments for terraform apply
	args := []string{"apply"}
	if req.PlanFile != "" {
		args = append(args, req.PlanFile)
	} else {
		args = append(args, "-auto-approve")
	}

	// Run `terraform apply`
	output, err := runCommand(workspaceDir, "terraform", args...)
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

	// Run `terraform destroy`
	output, err := runCommand(workspaceDir, "terraform", "destroy", "-auto-approve")
	if err != nil {
		return &pb.DestroyResponse{Success: false, DestroyOutput: output, Error: err.Error()}, nil
	}

	return &pb.DestroyResponse{Success: true, DestroyOutput: output}, nil
}

// runCommand runs a command in the specified directory and returns its output.
func runCommand(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	return string(output), err
}
