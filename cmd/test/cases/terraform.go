package cases

import (
	"context"
	"fmt"
	"log"
	pb "terraform-executor/api/proto"
	"terraform-executor/cmd/test/utils"
	"terraform-executor/internal/executor"
)

func GetTerraformTests(ctx context.Context, svc *executor.ExecutorService, userId, contextName, workspaceName string) []utils.TestCase {
	return []utils.TestCase{
		{
			Name:     "Plan infrastructure",
			Category: "Terraform",
			Fn: func() error {
				resp, err := svc.Plan(ctx, &pb.PlanRequest{
					UserId:    userId,
					Context:   contextName,
					Workspace: workspaceName,
				})

				// First check if the RPC call itself failed
				if err != nil {
					return fmt.Errorf("RPC error: %v", err)
				}

				// Always log the plan output if it exists
				if resp.PlanOutput != "" {
					log.Printf("\n🔍 Terraform Plan Output:\n%s\n", resp.PlanOutput)
				}

				// Then check the response status
				if !resp.Success {
					return fmt.Errorf("❌ plan failed: %s", resp.Error)
				}

				// Verify we have some output
				if resp.PlanOutput == "" {
					return fmt.Errorf("plan succeeded but output is empty")
				}

				return nil
			},
		},
		{
			Name:     "Apply infrastructure",
			Category: "Terraform",
			Fn: func() error {
				resp, err := svc.Apply(ctx, &pb.ApplyRequest{
					UserId:    userId,
					Context:   contextName,
					Workspace: workspaceName,
				})

				// First check if the RPC call itself failed
				if err != nil {
					return fmt.Errorf("RPC error: %v", err)
				}

				// Always log the apply output if it exists
				if resp.ApplyOutput != "" {
					log.Printf("\n🚀 Terraform Apply Output:\n%s\n", resp.ApplyOutput)
				}

				// Then check the response status
				if !resp.Success {
					return fmt.Errorf("❌ apply failed: %s", resp.Error)
				}

				// Verify we have some output
				if resp.ApplyOutput == "" {
					return fmt.Errorf("apply succeeded but output is empty")
				}

				return nil
			},
		},
		{
			Name:     "Get State List",
			Category: "Terraform",
			Fn: func() error {
				resp, err := svc.GetStateList(ctx, &pb.GetStateListRequest{
					UserId:    userId,
					Context:   contextName,
					Workspace: workspaceName,
				})

				// First check if the RPC call itself failed
				if err != nil {
					return fmt.Errorf("RPC error: %v", err)
				}

				// Always log the state list output if it exists
				if resp.StateListOutput != "" {
					log.Printf("\n📋 Terraform State List Output:\n%s\n", resp.StateListOutput)
				}

				// Then check the response status
				if !resp.Success {
					return fmt.Errorf("❌ state list failed: %s", resp.Error)
				}

				// Verify we have some output
				if resp.StateListOutput == "" {
					return fmt.Errorf("state list succeeded but output is empty")
				}

				return nil
			},
		},
		{
			Name:     "Destroy infrastructure",
			Category: "Terraform",
			Fn: func() error {
				resp, err := svc.Destroy(ctx, &pb.DestroyRequest{
					UserId:    userId,
					Context:   contextName,
					Workspace: workspaceName,
				})

				// First check if the RPC call itself failed
				if err != nil {
					return fmt.Errorf("RPC error: %v", err)
				}

				// Always log the destroy output if it exists
				if resp.DestroyOutput != "" {
					log.Printf("\n💥 Terraform Destroy Output:\n%s\n", resp.DestroyOutput)
				}

				// Then check the response status
				if !resp.Success {
					return fmt.Errorf("❌ destroy failed: %s", resp.Error)
				}

				// Verify we have some output
				if resp.DestroyOutput == "" {
					return fmt.Errorf("destroy succeeded but output is empty")
				}

				return nil
			},
		},
	}
}
