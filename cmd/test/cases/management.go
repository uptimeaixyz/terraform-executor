package cases

import (
	"context"
	"fmt"
	"strings"
	pb "terraform-executor/api/proto"
	"terraform-executor/cmd/test/utils"
	"terraform-executor/internal/executor"
)

func GetManagementTests(ctx context.Context, svc *executor.ExecutorService, userId, projectName string) []utils.TestCase {
	return []utils.TestCase{
		// Provider management
		{
			Name:     "Add provider",
			Category: "Management",
			Fn: func() error {
				resp, err := svc.AddProviders(ctx, &pb.AddProvidersRequest{
					UserId:  userId,
					Project: projectName,
					Providers: []*pb.AddProvidersRequest_Provider{
						{
							Name:    "null",
							Source:  "hashicorp/null",
							Version: "~> 3.0",
						},
					},
				})
				if err != nil || !resp.Success {
					return fmt.Errorf("failed to add provider: %v", err)
				}

				// Verify ConfigMap
				configMapName := fmt.Sprintf("%s.versions.tf", projectName)
				cm, err := svc.K8sClient.GetConfigMap(ctx, userId, configMapName) // Updated to use ctx
				if err != nil {
					return fmt.Errorf("failed to get ConfigMap: %v", err)
				}
				if content, ok := cm.Data["versions.tf"]; !ok || !strings.Contains(content, "hashicorp/null") {
					return fmt.Errorf("ConfigMap does not contain expected provider configuration")
				}
				return nil
			},
		},
		// Environment variables
		{
			Name:     "Add environment variables",
			Category: "Management",
			Fn: func() error {
				resp, err := svc.AddSecretEnv(ctx, &pb.AddSecretEnvRequest{
					UserId:  userId,
					Project: projectName,
					Secrets: []*pb.AddSecretEnvRequest_Secret{
						{
							Name:  "AWS_REGION",
							Value: "us-west-2",
						},
						{
							Name:  "AWS_ACCESS_KEY_ID",
							Value: "test-access-key",
						},
					},
				})
				if err != nil || !resp.Success {
					return fmt.Errorf("failed to add secret env: %v", err)
				}
				return nil
			},
		},
		// Terraform variables
		{
			Name:     "Add Terraform variables",
			Category: "Management",
			Fn: func() error {
				resp, err := svc.AddSecretVar(ctx, &pb.AddSecretVarRequest{
					UserId:  userId,
					Project: projectName,
					Secrets: []*pb.AddSecretVarRequest_Secret{
						{
							Name:  "do_token",
							Value: "my-digital-ocean-token",
						},
						{
							Name:  "db_password",
							Value: "super-secret-password",
						},
					},
				})
				if err != nil || !resp.Success {
					return fmt.Errorf("failed to add secret var: %v", err)
				}
				return nil
			},
		},
		// Code management
		{
			Name:     "Append and verify Terraform code",
			Category: "Management",
			Fn: func() error {
				// Add code
				tfCode := `
resource "null_resource" "example" {
  provisioner "local-exec" {
	command = "echo 'Hello, World!'"
  }
}
`
				resp, err := svc.AppendCode(ctx, &pb.AppendCodeRequest{
					UserId:  userId,
					Project: projectName,
					Code:    tfCode,
				})
				if err != nil || !resp.Success {
					return fmt.Errorf("failed to append code: %v", err)
				}

				// Verify ConfigMap
				configMapName := fmt.Sprintf("%s.main.tf", projectName)
				cm, err := svc.K8sClient.GetConfigMap(ctx, userId, configMapName)
				if err != nil {
					return fmt.Errorf("failed to get ConfigMap: %v", err)
				}
				if content, ok := cm.Data["main.tf"]; !ok || !strings.Contains(content, "null_resource") {
					return fmt.Errorf("ConfigMap does not contain expected content")
				}
				return nil
			},
		},
		// Verification
		{
			Name:     "Verify main.tf content",
			Category: "Management",
			Fn: func() error {
				resp, err := svc.GetMainTf(ctx, &pb.GetMainTfRequest{
					UserId:  userId,
					Project: projectName,
				})
				if err != nil || !resp.Success {
					return fmt.Errorf("failed to get main.tf content: %v", err)
				}
				if !strings.Contains(resp.Content, "null_resource") {
					return fmt.Errorf("main.tf does not contain expected content")
				}
				return nil
			},
		},
	}
}
