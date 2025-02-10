package cases

import (
	"context"
	"fmt"
	pb "terraform-executor/api/proto"
	"terraform-executor/cmd/test/utils"
	"terraform-executor/internal/executor"
)

func GetCleanupTests(ctx context.Context, svc *executor.ExecutorService, userId, projectName string) []utils.TestCase {
	return []utils.TestCase{
		{
			Name:     "Delete project",
			Category: "Cleanup",
			Fn: func() error {
				// Clean up workspace
				resp, err := svc.DeleteProject(ctx, &pb.DeleteProjectRequest{
					UserId:  userId,
					Project: projectName,
				})
				if err != nil {
					return fmt.Errorf("failed to clear code, error: %+v", err)
				}
				if !resp.Success {
					return fmt.Errorf("failed to clear code, response unsuccessful: %+v", resp)
				}

				// Check all resources and collect errors
				errors := []string{}

				// Check if main.tf was deleted
				configMapName := fmt.Sprintf("%s.main.tf", projectName)
				if _, err := svc.K8sClient.GetConfigMap(ctx, userId, configMapName); err == nil {
					errors = append(errors, "main.tf ConfigMap was not deleted")
				}

				// Check if versions.tf was deleted
				configMapName = fmt.Sprintf("%s.versions.tf", projectName)
				if _, err := svc.K8sClient.GetConfigMap(ctx, userId, configMapName); err == nil {
					errors = append(errors, "versions.tf ConfigMap was not deleted")
				}

				// Check if variables.tf was deleted
				configMapName = fmt.Sprintf("%s.variables.tf", projectName)
				if _, err := svc.K8sClient.GetConfigMap(ctx, userId, configMapName); err == nil {
					errors = append(errors, "variables.tf ConfigMap was not deleted")
				}

				// Check if secret was deleted
				secretName := fmt.Sprintf("%s.env", projectName)
				if _, err := svc.K8sClient.GetSecret(ctx, userId, secretName); err == nil {
					errors = append(errors, "Secret was not deleted")
				}

				if len(errors) > 0 {
					return fmt.Errorf("deletion failures: %v", errors)
				}

				return nil
			},
		},
	}
}
