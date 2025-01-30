package main

import (
	"context"
	"fmt"
	"log"
	pb "terraform-executor/api/proto"
	"terraform-executor/internal/executor"
)

func main() {
	// Create a new executor service
	svc := &executor.ExecutorService{}
	ctx := context.Background()

	// Test creating a context
	contextName := "test-context"
	workspaceName := "test-workspace"

	// Test Clear Workspace
	if err := clearWorkspace(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	if err := createContext(ctx, svc, contextName); err != nil {
		log.Fatal(err)
	}

	// Test creating a workspace
	if err := createWorkspace(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	// Test adding AWS provider
	if err := addAwsProvider(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	// Test clear providers
	if err := clearProviders(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	// Add provider again
	if err := addAwsProvider(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	// Test adding environment variables
	if err := addSecretEnv(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	// Test adding secret variables
	if err := addSecretVar(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	// Test clear secret variables
	if _, err := svc.ClearSecretVars(ctx, &pb.ClearSecretVarsRequest{
		Context:   contextName,
		Workspace: workspaceName,
	}); err != nil {
		log.Fatal(err)
	}

	// Add secret variable again
	if err := addSecretVar(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	// Test adding Terraform code
	if err := appendTerraformCode(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}

	// Test plan
	if err := planInfrastructure(ctx, svc, contextName, workspaceName); err != nil {
		log.Fatal(err)
	}
	fmt.Println("All tests completed successfully!")
}

// test functions

// createContext creates a new context
func createContext(ctx context.Context, svc *executor.ExecutorService, contextName string) error {
	resp, err := svc.CreateContext(ctx, &pb.CreateContextRequest{
		Context: contextName,
	})
	if err != nil {
		return fmt.Errorf("failed to create context: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("context creation failed: %s", resp.Error)
	}
	fmt.Println("Context created successfully")
	return nil
}

// createWorkspace creates a new workspace
func createWorkspace(ctx context.Context, svc *executor.ExecutorService, contextName, workspaceName string) error {
	resp, err := svc.CreateWorkspace(ctx, &pb.CreateWorkspaceRequest{
		Context:   contextName,
		Workspace: workspaceName,
	})
	if err != nil {
		return fmt.Errorf("failed to create workspace: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("workspace creation failed: %s", resp.Error)
	}
	fmt.Println("Workspace created successfully")
	return nil
}

// addAwsProvider adds the AWS provider to the workspace
func addAwsProvider(ctx context.Context, svc *executor.ExecutorService, contextName, workspaceName string) error {
	resp, err := svc.AddProviders(ctx, &pb.AddProvidersRequest{
		Context:   contextName,
		Workspace: workspaceName,
		Providers: []*pb.AddProvidersRequest_Provider{
			{
				Name:    "aws",
				Source:  "hashicorp/aws",
				Version: "~> 4.0",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add provider: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("provider addition failed: %s", resp.Error)
	}
	fmt.Println("Provider added successfully")
	return nil
}

// addSecretEnv adds a secret environment variable to the workspace
func addSecretEnv(ctx context.Context, svc *executor.ExecutorService, contextName, workspaceName string) error {
	resp, err := svc.AddSecretEnv(ctx, &pb.AddSecretEnvRequest{
		Context:     contextName,
		Workspace:   workspaceName,
		SecretName:  "AWS_REGION",
		SecretValue: "us-west-2",
	})
	if err != nil {
		return fmt.Errorf("failed to add secret env: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("secret env addition failed: %s", resp.Error)
	}
	fmt.Println("Secret env added successfully")
	return nil
}

// addSecretVar adds a secret variable to the workspace
func addSecretVar(ctx context.Context, svc *executor.ExecutorService, contextName, workspaceName string) error {
	resp, err := svc.AddSecretVar(ctx, &pb.AddSecretVarRequest{
		Context:     contextName,
		Workspace:   workspaceName,
		SecretName:  "do_token",
		SecretValue: "my-digital-ocean-token",
	})

	if err != nil {
		return fmt.Errorf("failed to add secret var: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("secret var addition failed: %s", resp.Error)
	}
	fmt.Println("Secret var added successfully")
	return nil
}

// appendTerraformCode appends Terraform code to the workspace
func appendTerraformCode(ctx context.Context, svc *executor.ExecutorService, contextName, workspaceName string) error {
	// Example Terraform code for an S3 bucket
	tfCode := `
resource "aws_s3_bucket" "example" {
  bucket = "my-test-bucket"
}
`
	resp, err := svc.AppendCode(ctx, &pb.AppendCodeRequest{
		Context:   contextName,
		Workspace: workspaceName,
		Code:      tfCode,
	})
	if err != nil {
		return fmt.Errorf("failed to append code: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("code append failed: %s", resp.Error)
	}
	fmt.Println("Terraform code added successfully")
	return nil
}

// planInfrastructure plans the infrastructure
func planInfrastructure(ctx context.Context, svc *executor.ExecutorService, contextName, workspaceName string) error {
	resp, err := svc.Plan(ctx, &pb.PlanRequest{
		Context:   contextName,
		Workspace: workspaceName,
	})
	if err != nil {
		return fmt.Errorf("failed to plan: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("plan failed: %s", resp.Error)
	}
	fmt.Printf("Plan output:\n%s\n", resp.PlanOutput)
	return nil
}

// clearProviders clears the providers
func clearProviders(ctx context.Context, svc *executor.ExecutorService, contextName, workspaceName string) error {
	resp, err := svc.ClearProviders(ctx, &pb.ClearProvidersRequest{
		Context:   contextName,
		Workspace: workspaceName,
	})
	if err != nil {
		return fmt.Errorf("failed to clear providers: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("provider clear failed: %s", resp.Error)
	}
	fmt.Println("Providers cleared successfully")
	return nil
}

// clearWorkspace clears the workspace
func clearWorkspace(ctx context.Context, svc *executor.ExecutorService, contextName, workspaceName string) error {
	resp, err := svc.ClearWorkspace(ctx, &pb.ClearWorkspaceRequest{
		Context:   contextName,
		Workspace: workspaceName,
	})

	if err != nil {
		return fmt.Errorf("failed to clear workspace: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("workspace clear failed: %s", resp.Error)
	}
	fmt.Println("Workspace cleared successfully")
	return nil
}
