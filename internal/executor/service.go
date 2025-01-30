package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	pb "terraform-executor/api/proto"
	"terraform-executor/pkg/utils"
)

// ExecutorService implements the ExecutorServer interface.
type ExecutorService struct {
	pb.UnimplementedExecutorServer
}

// AppendCode appends the provided code to the main.tf file in the workspace directory.
func (s *ExecutorService) AppendCode(ctx context.Context, req *pb.AppendCodeRequest) (*pb.AppendCodeResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.AppendCodeResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Append the code to main.tf
	if req.Code == "" {
		return &pb.AppendCodeResponse{Success: false, Error: "code is empty"}, nil
	}

	mainTfPath := filepath.Join(workspaceDir, "main.tf")
	if err := utils.AppendToFile(mainTfPath, req.Code); err != nil {
		return &pb.AppendCodeResponse{Success: false, Error: fmt.Sprintf("failed to write to main.tf: %v", err)}, nil
	}

	return &pb.AppendCodeResponse{Success: true}, nil
}

// Plan generates a Terraform plan and returns the result.
func (s *ExecutorService) Plan(ctx context.Context, req *pb.PlanRequest) (*pb.PlanResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.PlanResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Run `terraform init`
	_, err := runCommand(workspaceDir, "terraform", "init")
	if err != nil {
		clearMainTf(workspaceDir)
		return &pb.PlanResponse{Success: false, Error: fmt.Sprintf("failed to initialize: %v", err)}, nil
	}

	// Run `terraform plan`
	output, err := runCommand(workspaceDir, "terraform", "plan", "-no-color")
	if err != nil {
		return &pb.PlanResponse{Success: false, PlanOutput: output, Error: err.Error()}, nil
	}
	return &pb.PlanResponse{Success: true, PlanOutput: output}, nil
}

// Apply applies a Terraform plan or configuration and returns the result.
func (s *ExecutorService) Apply(ctx context.Context, req *pb.ApplyRequest) (*pb.ApplyResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.ApplyResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Run `terraform init`
	_, err := runCommand(workspaceDir, "terraform", "init")
	if err != nil {
		return &pb.ApplyResponse{Success: false, Error: fmt.Sprintf("failed to initialize: %v", err)}, nil
	}

	// Run `terraform apply`
	output, err := runCommand(workspaceDir, "terraform", "apply", "-auto-approve", "-no-color")
	if err != nil {
		return &pb.ApplyResponse{Success: false, ApplyOutput: output, Error: err.Error()}, nil
	}

	return &pb.ApplyResponse{Success: true, ApplyOutput: output}, nil
}

// Destroy destroys the Terraform-managed infrastructure and returns the result.
func (s *ExecutorService) Destroy(ctx context.Context, req *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.DestroyResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Run `terraform init`
	_, err := runCommand(workspaceDir, "terraform", "init")
	if err != nil {
		return &pb.DestroyResponse{Success: false, Error: fmt.Sprintf("failed to initialize: %v", err)}, nil
	}

	// Run `terraform destroy`
	output, err := runCommand(workspaceDir, "terraform", "destroy", "-auto-approve", "-no-color")
	if err != nil {
		return &pb.DestroyResponse{Success: false, DestroyOutput: output, Error: err.Error()}, nil
	}

	return &pb.DestroyResponse{Success: true, DestroyOutput: output}, nil
}

// GetStateList returns output of "terraform state list" command
func (s *ExecutorService) GetStateList(ctx context.Context, req *pb.GetStateListRequest) (*pb.GetStateListResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.GetStateListResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Run `terraform init`
	_, err := runCommand(workspaceDir, "terraform", "init")
	if err != nil {
		return &pb.GetStateListResponse{Success: false, Error: fmt.Sprintf("failed to initialize: %v", err)}, nil
	}

	// Run `terraform state list`
	output, err := runCommand(workspaceDir, "terraform", "state", "list")
	if err != nil {
		return &pb.GetStateListResponse{Success: false, StateListOutput: output, Error: err.Error()}, nil
	}

	return &pb.GetStateListResponse{Success: true, StateListOutput: output}, nil
}

// Clear removes all created files in the workspace directory.
func (s *ExecutorService) ClearCode(ctx context.Context, req *pb.ClearCodeRequest) (*pb.ClearCodeResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Remove main.tf
	if err := clearMainTf(workspaceDir); err != nil {
		return &pb.ClearCodeResponse{Success: false, Error: fmt.Sprintf("failed to clear main.tf: %v", err)}, nil
	}

	return &pb.ClearCodeResponse{Success: true}, nil
}

// Create new context
func (s *ExecutorService) CreateContext(ctx context.Context, req *pb.CreateContextRequest) (*pb.CreateContextResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)

	// Ensure the context directory exists
	if err := os.MkdirAll(contextDir, os.ModePerm); err != nil {
		return &pb.CreateContextResponse{Success: false, Error: fmt.Sprintf("failed to create context: %v", err)}, nil
	}

	return &pb.CreateContextResponse{Success: true}, nil
}

// Delete context
func (s *ExecutorService) DeleteContext(ctx context.Context, req *pb.DeleteContextRequest) (*pb.DeleteContextResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)

	// Remove the context directory
	if err := os.RemoveAll(contextDir); err != nil {
		return &pb.DeleteContextResponse{Success: false, Error: fmt.Sprintf("failed to delete context: %v", err)}, nil
	}

	return &pb.DeleteContextResponse{Success: true}, nil
}

// Create workspace
func (s *ExecutorService) CreateWorkspace(ctx context.Context, req *pb.CreateWorkspaceRequest) (*pb.CreateWorkspaceResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.CreateWorkspaceResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	return &pb.CreateWorkspaceResponse{Success: true}, nil
}

// Delete workspace
func (s *ExecutorService) DeleteWorkspace(ctx context.Context, req *pb.DeleteWorkspaceRequest) (*pb.DeleteWorkspaceResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Remove the workspace directory
	if err := os.RemoveAll(workspaceDir); err != nil {
		return &pb.DeleteWorkspaceResponse{Success: false, Error: fmt.Sprintf("failed to delete workspace: %v", err)}, nil
	}

	return &pb.DeleteWorkspaceResponse{Success: true}, nil
}

// Add providers to the Terraform configuration
func (s *ExecutorService) AddProviders(ctx context.Context, req *pb.AddProvidersRequest) (*pb.AddProvidersResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.AddProvidersResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Initialize providers slice
	providers := make([]utils.ProviderConfig, 0, len(req.Providers))

	// Loop through the providers from the request
	for _, p := range req.Providers {
		provider := utils.ProviderConfig{
			Name:    p.Name,
			Source:  p.Source,
			Version: p.Version,
		}
		providers = append(providers, provider)
	}

	// Fill the struct with the provider data
	data := utils.TerraformTemplateData{
		BackendPath: "terraform.tfstate",
		Providers:   providers,
	}

	// Generate the Terraform configuration
	config, err := utils.GenerateTerraformConfig(data)
	if err != nil {
		return &pb.AddProvidersResponse{Success: false, Error: fmt.Sprintf("failed to generate Terraform config: %v", err)}, nil
	}

	// Write the Terraform configuration to versions.tf
	versionsTfPath := filepath.Join(workspaceDir, "versions.tf")
	if err := utils.AppendToFile(versionsTfPath, config); err != nil {
		return &pb.AddProvidersResponse{Success: false, Error: fmt.Sprintf("failed to write to versions.tf: %v", err)}, nil
	}

	return &pb.AddProvidersResponse{Success: true}, nil
}

// Clear providers from the Terraform configuration
func (s *ExecutorService) ClearProviders(ctx context.Context, req *pb.ClearProvidersRequest) (*pb.ClearProvidersResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Remove versions.tf
	versionsTfPath := filepath.Join(workspaceDir, "versions.tf")
	if err := os.Remove(versionsTfPath); err != nil {
		return &pb.ClearProvidersResponse{Success: false, Error: fmt.Sprintf("failed to remove versions.tf: %v", err)}, nil
	}

	return &pb.ClearProvidersResponse{Success: true}, nil
}

// Add secret env variables to the Terraform configuration
func (s *ExecutorService) AddSecretEnv(ctx context.Context, req *pb.AddSecretEnvRequest) (*pb.AddSecretEnvResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.AddSecretEnvResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Append secret to .env file
	envPath := filepath.Join(workspaceDir, ".env")
	if err := utils.AppendToFile(envPath, fmt.Sprintf("%s=%s", req.SecretName, req.SecretValue)); err != nil {
		return &pb.AddSecretEnvResponse{Success: false, Error: fmt.Sprintf("failed to write to .env: %v", err)}, nil
	}

	return &pb.AddSecretEnvResponse{Success: true}, nil
}

// Add secret terraform variables to the Terraform configuration
func (s *ExecutorService) AddSecretVar(ctx context.Context, req *pb.AddSecretVarRequest) (*pb.AddSecretVarResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Ensure the workspace directory exists
	if err := os.MkdirAll(workspaceDir, os.ModePerm); err != nil {
		return &pb.AddSecretVarResponse{Success: false, Error: fmt.Sprintf("failed to create workspace: %v", err)}, nil
	}

	// Append secret to variables.tf file
	varsPath := filepath.Join(workspaceDir, "variables.tf")
	if err := utils.AppendToFile(varsPath, fmt.Sprintf("variable \"%s\" {\n  type = string\n  default = \"%s\"\n}\n", req.SecretName, req.SecretValue)); err != nil {
		return &pb.AddSecretVarResponse{Success: false, Error: fmt.Sprintf("failed to write to variables.tf: %v", err)}, nil
	}

	return &pb.AddSecretVarResponse{Success: true}, nil
}

// Clear secret variables from the Terraform configuration
func (s *ExecutorService) ClearSecretVars(ctx context.Context, req *pb.ClearSecretVarsRequest) (*pb.ClearSecretVarsResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Remove .env file
	envPath := filepath.Join(workspaceDir, "variables.tf")
	if err := os.Remove(envPath); err != nil {
		return &pb.ClearSecretVarsResponse{Success: false, Error: fmt.Sprintf("failed to remove .env: %v", err)}, nil
	}

	return &pb.ClearSecretVarsResponse{Success: true}, nil
}

// Clear Workspace
func (s *ExecutorService) ClearWorkspace(ctx context.Context, req *pb.ClearWorkspaceRequest) (*pb.ClearWorkspaceResponse, error) {
	contextDir := filepath.Join("./data/", req.Context)
	workspaceDir := filepath.Join(contextDir, "/", req.Workspace)

	// Remove the workspace directory
	if err := os.RemoveAll(workspaceDir); err != nil {
		return &pb.ClearWorkspaceResponse{Success: false, Error: fmt.Sprintf("failed to delete workspace: %v", err)}, nil
	}

	return &pb.ClearWorkspaceResponse{Success: true}, nil
}

// clear main.tf file
func clearMainTf(workspaceDir string) error {
	mainTfPath := filepath.Join(workspaceDir, "main.tf")
	if err := os.Remove(mainTfPath); err != nil {
		return fmt.Errorf("failed to remove main.tf: %v", err)
	}
	return nil
}

// loadEnvFile loads environment variables from .env file in the workspace directory
func loadEnvFile(workspaceDir string) (map[string]string, error) {
	envPath := filepath.Join(workspaceDir, ".env")
	envVars := make(map[string]string)

	// Check if .env file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return envVars, nil
	}

	content, err := os.ReadFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		envVars[parts[0]] = parts[1]
	}

	return envVars, nil
}

// runCommand runs a command in the specified directory and returns its output.
func runCommand(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	// Load environment variables from .env file
	envVars, err := loadEnvFile(dir)
	if err != nil {
		return "", fmt.Errorf("failed to load environment variables: %v", err)
	}

	// Get current environment
	cmd.Env = os.Environ()

	// Add variables from .env file
	for k, v := range envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	output, err := cmd.CombinedOutput()
	return string(output), err
}
