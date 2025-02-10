package cases

import (
	"context"
	"fmt"
	pb "terraform-executor/api/proto"
	"terraform-executor/cmd/test/utils"
	"terraform-executor/internal/executor"
)

func GetSetupTests(ctx context.Context, svc *executor.ExecutorService, userId, projectName string) []utils.TestCase {
	return []utils.TestCase{
		{
			Name:     "Verify namespace creation",
			Category: "Setup",
			Fn: func() error {
				exists, err := svc.K8sClient.NamespaceExists(ctx, userId) // Updated to use ctx
				if err != nil {
					return fmt.Errorf("failed to check namespace: %v", err)
				}
				if !exists {
					return fmt.Errorf("namespace was not created")
				}
				return nil
			},
		},
		{
			Name:     "Create Project",
			Category: "Setup",
			Fn: func() error {
				resp, err := svc.CreateProject(ctx, &pb.CreateProjectRequest{
					UserId:  userId,
					Project: projectName,
				})
				if err != nil || !resp.Success {
					return fmt.Errorf("project creation failed: %v", err)
				}
				return nil
			},
		},
		// ...add other setup tests...
	}
}
