package executor

import (
	"context"
	"fmt"
	"strings"
	pb "terraform-executor/api/proto"
	"terraform-executor/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppendCode appends the provided code to the main.tf file in the workspace directory and ConfigMap.
func (s *ExecutorService) AppendCode(ctx context.Context, req *pb.AppendCodeRequest) (*pb.AppendCodeResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.AppendCodeResponse{Success: false, Error: err.Error()}, nil
	}

	configMapName := fmt.Sprintf("%s.%s", req.Project, "main.tf")

	// Get existing ConfigMap or create new one
	cm, err := s.K8sClient.GetConfigMap(ctx, req.UserId, configMapName)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new ConfigMap
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: configMapName,
				},
				Data: map[string]string{
					"main.tf": req.Code,
				},
			}
			if err := s.K8sClient.CreateConfigMap(ctx, req.UserId, cm); err != nil {
				return &pb.AppendCodeResponse{Success: false, Error: fmt.Sprintf("failed to create ConfigMap: %v", err)}, nil
			}
		} else {
			return &pb.AppendCodeResponse{Success: false, Error: fmt.Sprintf("failed to get ConfigMap: %v", err)}, nil
		}
	} else {
		// Update existing ConfigMap
		if cm.Data == nil {
			cm.Data = make(map[string]string)
		}
		cm.Data["main.tf"] = cm.Data["main.tf"] + "\n" + req.Code
		if err := s.K8sClient.UpdateConfigMap(ctx, req.UserId, cm); err != nil {
			return &pb.AppendCodeResponse{Success: false, Error: fmt.Sprintf("failed to update ConfigMap: %v", err)}, nil
		}
	}
	return &pb.AppendCodeResponse{Success: true}, nil
}

// Clear removes all created files in the workspace directory and ConfigMap.
func (s *ExecutorService) ClearCode(ctx context.Context, req *pb.ClearCodeRequest) (*pb.ClearCodeResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.ClearCodeResponse{Success: false, Error: err.Error()}, nil
	}

	configMapName := fmt.Sprintf("%s.%s", req.Project, "main.tf")

	// Delete ConfigMap
	if err := s.K8sClient.DeleteConfigMap(ctx, req.UserId, configMapName); err != nil {
		if !errors.IsNotFound(err) {
			return &pb.ClearCodeResponse{Success: false, Error: fmt.Sprintf("failed to delete ConfigMap: %v", err)}, nil
		}
		// Ignore not found error
	}
	return &pb.ClearCodeResponse{Success: true}, nil
}

// Add providers to the Terraform configuration
func (s *ExecutorService) AddProviders(ctx context.Context, req *pb.AddProvidersRequest) (*pb.AddProvidersResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.AddProvidersResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to ensure namespace: %v", err),
		}, err
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
		Bucket:    s.Bucket,
		UserID:    req.UserId,
		Project:   req.Project,
		Providers: providers,
	}

	// Generate the Terraform configuration
	config, err := utils.GenerateTerraformConfig(data)
	if err != nil {
		return &pb.AddProvidersResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to generate Terraform config: %v", err),
		}, err
	}

	// Create or update ConfigMap
	configMapName := fmt.Sprintf("%s.%s", req.Project, "versions.tf")
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: configMapName,
		},
		Data: map[string]string{
			"versions.tf": config,
		},
	}

	existingCm, err := s.K8sClient.GetConfigMap(ctx, req.UserId, configMapName)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new ConfigMap
			if err := s.K8sClient.CreateConfigMap(ctx, req.UserId, cm); err != nil {
				return &pb.AddProvidersResponse{
					Success: false,
					Error:   fmt.Sprintf("failed to create ConfigMap: %v", err),
				}, err
			}
		} else {
			return &pb.AddProvidersResponse{
				Success: false,
				Error:   fmt.Sprintf("failed to get ConfigMap: %v", err),
			}, err
		}
	} else {
		// Update existing ConfigMap
		existingCm.Data = cm.Data
		if err := s.K8sClient.UpdateConfigMap(ctx, req.UserId, existingCm); err != nil {
			return &pb.AddProvidersResponse{
				Success: false,
				Error:   fmt.Sprintf("failed to update ConfigMap: %v", err),
			}, nil
		}
	}

	return &pb.AddProvidersResponse{Success: true}, nil
}

// Clear providers from the Terraform configuration
func (s *ExecutorService) ClearProviders(ctx context.Context, req *pb.ClearProvidersRequest) (*pb.ClearProvidersResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.ClearProvidersResponse{Success: false, Error: err.Error()}, nil
	}

	// Remove ConfigMap
	configMapName := fmt.Sprintf("%s.%s", req.Project, "versions.tf")
	if err := s.K8sClient.DeleteConfigMap(ctx, req.UserId, configMapName); err != nil {
		if !errors.IsNotFound(err) {
			return &pb.ClearProvidersResponse{Success: false, Error: fmt.Sprintf("failed to delete ConfigMap: %v", err)}, nil
		}
		// Ignore not found error
	}
	return &pb.ClearProvidersResponse{Success: true}, nil
}

// Add secret env variables to the Terraform configuration
func (s *ExecutorService) AddSecretEnv(ctx context.Context, req *pb.AddSecretEnvRequest) (*pb.AddSecretEnvResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.AddSecretEnvResponse{Success: false, Error: err.Error()}, nil
	}

	secretName := fmt.Sprintf("%s.%s", req.Project, "env")
	// Get existing Secret or create new one
	secret, err := s.K8sClient.GetSecret(ctx, req.UserId, secretName)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new Secret from the req.Secrets
			secret = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: secretName,
				},
				StringData: make(map[string]string),
			}
			for _, s := range req.Secrets {
				secret.StringData[s.Name] = s.Value
			}
			if err := s.K8sClient.CreateSecret(ctx, req.UserId, secret); err != nil {
				return &pb.AddSecretEnvResponse{Success: false, Error: fmt.Sprintf("failed to create Secret: %v", err)}, nil
			}
		} else {
			return &pb.AddSecretEnvResponse{Success: false, Error: fmt.Sprintf("failed to get Secret: %v", err)}, nil
		}
	} else {
		// Update existing Secret by adding new key-value pairs
		if secret.StringData == nil {
			secret.StringData = make(map[string]string)
		}
		for _, s := range req.Secrets {
			secret.StringData[s.Name] = s.Value
		}
		if err := s.K8sClient.UpdateSecret(ctx, req.UserId, secret); err != nil {
			return &pb.AddSecretEnvResponse{Success: false, Error: fmt.Sprintf("failed to update Secret: %v", err)}, nil
		}
	}
	return &pb.AddSecretEnvResponse{Success: true}, nil
}

// Clear secret env variables from the Terraform configuration
func (s *ExecutorService) ClearSecretEnv(ctx context.Context, req *pb.ClearSecretEnvRequest) (*pb.ClearSecretEnvResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.ClearSecretEnvResponse{Success: false, Error: err.Error()}, nil
	}

	// Remove Secret
	secretName := fmt.Sprintf("%s.%s", req.Project, "env")
	if err := s.K8sClient.DeleteSecret(ctx, req.UserId, secretName); err != nil {
		if !errors.IsNotFound(err) {
			return &pb.ClearSecretEnvResponse{Success: false, Error: fmt.Sprintf("failed to delete Secret: %v", err)}, nil
		}
		// Ignore not found error
	}
	return &pb.ClearSecretEnvResponse{Success: true}, nil
}

// Add secret terraform variables to the Terraform configuration
func (s *ExecutorService) AddSecretVar(ctx context.Context, req *pb.AddSecretVarRequest) (*pb.AddSecretVarResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.AddSecretVarResponse{Success: false, Error: err.Error()}, nil
	}

	configMapName := fmt.Sprintf("%s.%s", req.Project, "variables.tf")
	// Get existing ConfigMap or create new one
	cm, err := s.K8sClient.GetConfigMap(ctx, req.UserId, configMapName)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create new ConfigMap
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: configMapName,
				},
				Data: map[string]string{
					"variables.tf": "",
				},
			}
			if err := s.K8sClient.CreateConfigMap(ctx, req.UserId, cm); err != nil {
				return &pb.AddSecretVarResponse{Success: false, Error: fmt.Sprintf("failed to create ConfigMap: %v", err)}, nil
			}
		} else {
			return &pb.AddSecretVarResponse{Success: false, Error: fmt.Sprintf("failed to get ConfigMap: %v", err)}, nil
		}
	}

	// Update existing ConfigMap
	varsContent := cm.Data["variables.tf"]
	for _, secret := range req.Secrets {
		varDef := fmt.Sprintf("variable \"%s\" {\n  type = string\n  default = \"%s\"\n}\n", secret.Name, secret.Value)
		varsContent += varDef
	}
	cm.Data["variables.tf"] = varsContent
	if err := s.K8sClient.UpdateConfigMap(ctx, req.UserId, cm); err != nil {
		return &pb.AddSecretVarResponse{Success: false, Error: fmt.Sprintf("failed to update ConfigMap: %v", err)}, nil
	}
	return &pb.AddSecretVarResponse{Success: true}, nil
}

// Clear secret variables from the Terraform configuration
func (s *ExecutorService) ClearSecretVars(ctx context.Context, req *pb.ClearSecretVarsRequest) (*pb.ClearSecretVarsResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.ClearSecretVarsResponse{Success: false, Error: err.Error()}, nil
	}

	// Remove ConfigMap
	configMapName := fmt.Sprintf("%s.%s", req.Project, "variables.tf")
	if err := s.K8sClient.DeleteConfigMap(ctx, req.UserId, configMapName); err != nil {
		if !errors.IsNotFound(err) {
			return &pb.ClearSecretVarsResponse{Success: false, Error: fmt.Sprintf("failed to delete ConfigMap: %v", err)}, nil
		}
		// Ignore not found error
	}
	return &pb.ClearSecretVarsResponse{Success: true}, nil
}

// Create a new project with the provided name.
func (s *ExecutorService) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.CreateProjectResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.CreateProjectResponse{Success: true}, nil
}

// DeleteProject removes all resources associated with the project.
func (s *ExecutorService) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.DeleteProjectResponse{Success: false, Error: err.Error()}, nil
	}

	var errors []string

	// Clear code
	if _, err := s.ClearCode(ctx, &pb.ClearCodeRequest{
		UserId:  req.UserId,
		Project: req.Project,
	}); err != nil {
		errors = append(errors, fmt.Sprintf("failed to clear code: %v", err))
	}

	// Clear providers
	if _, err := s.ClearProviders(ctx, &pb.ClearProvidersRequest{
		UserId:  req.UserId,
		Project: req.Project,
	}); err != nil {
		errors = append(errors, fmt.Sprintf("failed to clear providers: %v", err))
	}

	// Clear secret variables
	if _, err := s.ClearSecretVars(ctx, &pb.ClearSecretVarsRequest{
		UserId:  req.UserId,
		Project: req.Project,
	}); err != nil {
		errors = append(errors, fmt.Sprintf("failed to clear secret variables: %v", err))
	}

	// Clear secret env variables
	if _, err := s.ClearSecretEnv(ctx, &pb.ClearSecretEnvRequest{
		UserId:  req.UserId,
		Project: req.Project,
	}); err != nil {
		errors = append(errors, fmt.Sprintf("failed to clear secret env variables: %v", err))
	}

	if len(errors) > 0 {
		return &pb.DeleteProjectResponse{Success: false, Error: strings.Join(errors, "; ")}, nil
	}

	return &pb.DeleteProjectResponse{Success: true}, nil
}

// GetMainTf returns the content of main.tf from ConfigMap
func (s *ExecutorService) GetMainTf(ctx context.Context, req *pb.GetMainTfRequest) (*pb.GetMainTfResponse, error) {
	if err := s.ensureNamespace(ctx, req.UserId); err != nil {
		return &pb.GetMainTfResponse{Success: false, Error: err.Error()}, nil
	}

	configMapName := fmt.Sprintf("%s.%s", req.Project, "main.tf")

	// Get ConfigMap
	cm, err := s.K8sClient.GetConfigMap(ctx, req.UserId, configMapName)
	if err != nil {
		if errors.IsNotFound(err) {
			return &pb.GetMainTfResponse{
				Success: false,
				Error:   fmt.Sprintf("main.tf does not exist for project %s", req.Project),
			}, nil
		}
		return &pb.GetMainTfResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get ConfigMap: %v", err),
		}, nil
	}

	// Get content from ConfigMap
	content, ok := cm.Data["main.tf"]
	if !ok {
		return &pb.GetMainTfResponse{
			Success: false,
			Error:   "main.tf not found in ConfigMap",
		}, nil
	}

	return &pb.GetMainTfResponse{
		Success: true,
		Content: content,
	}, nil
}
