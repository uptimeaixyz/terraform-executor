package cases

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	pb "terraform-executor/api/proto"
	"terraform-executor/cmd/test/utils"
	"terraform-executor/internal/executor"
)

func GetCleanupTests(ctx context.Context, svc *executor.ExecutorService, userId, contextName, workspaceName string) []utils.TestCase {
	return []utils.TestCase{
		{
			Name:     "Clear Workspace",
			Category: "Cleanup",
			Fn: func() error {
				// Clean up workspace
				resp, err := svc.ClearWorkspace(ctx, &pb.ClearWorkspaceRequest{
					UserId:    userId,
					Context:   contextName,
					Workspace: workspaceName,
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
				configMapName := fmt.Sprintf("%s.%s.main.tf", contextName, workspaceName)
				if _, err := svc.K8sClient.GetConfigMap(ctx, userId, configMapName); err == nil {
					errors = append(errors, "main.tf ConfigMap was not deleted")
				}

				// Check if versions.tf was deleted
				configMapName = fmt.Sprintf("%s.%s.versions.tf", contextName, workspaceName)
				if _, err := svc.K8sClient.GetConfigMap(ctx, userId, configMapName); err == nil {
					errors = append(errors, "versions.tf ConfigMap was not deleted")
				}

				// Check if variables.tf was deleted
				configMapName = fmt.Sprintf("%s.%s.variables.tf", contextName, workspaceName)
				if _, err := svc.K8sClient.GetConfigMap(ctx, userId, configMapName); err == nil {
					errors = append(errors, "variables.tf ConfigMap was not deleted")
				}

				// Check if secret was deleted
				secretName := fmt.Sprintf("%s.%s.env", contextName, workspaceName)
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

func CleanupTestResources(ctx context.Context, svc *executor.ExecutorService, userId, contextName, workspaceName string) error {
	log.Println("Cleaning up test resources...")
	// Clean up local files
	dataDir := "./data"
	pathsToClean := []string{
		filepath.Join(dataDir, "users", userId),                             // User directory
		filepath.Join(dataDir, contextName),                                 // Context directory
		filepath.Join(dataDir, "users", userId, contextName, workspaceName), // Workspace directory
	}

	for _, path := range pathsToClean {
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: failed to delete directory %s: %v", path, err)
		} else if !os.IsNotExist(err) {
			log.Printf("Deleted directory: %s", path)
		}
	}

	return nil
}
