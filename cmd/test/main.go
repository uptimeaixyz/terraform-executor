package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"terraform-executor/cmd/test/cases"
	"terraform-executor/cmd/test/utils"
	"terraform-executor/internal/executor"
)

func main() {
	log.Println("Starting test execution...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Pass ctx to NewExecutorService for AWS client initialization
	svc, err := executor.NewExecutorService(ctx)
	if err != nil {
		log.Fatalf("Failed to create executor service: %v", err)
	}
	log.Println("Executor service created successfully")

	// enable debug mode
	svc.Debug = true
	log.Println("Debug mode enabled")

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run tests in a separate goroutine
	go func() {
		if err := runTests(ctx, svc); err != nil {
			log.Printf("Tests failed: %v", err)
			cancel()
		} else {
			log.Println("All tests completed successfully!")
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Tests finished or context cancelled")
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		cancel()
	}

	log.Println("Test execution completed")
}

func runTests(ctx context.Context, svc *executor.ExecutorService) error {
	userId := "test-user"
	contextName := "test-context"
	workspaceName := "test-workspace"

	// Initial cleanup before running tests
	log.Println("\n🧹 Cleaning up previous test resources...")
	if err := cases.CleanupTestResources(ctx, svc, userId, contextName, workspaceName); err != nil {
		log.Printf("Warning: initial cleanup failed: %v", err)
	}
	log.Println("Cleanup completed")

	log.Println("\n🚀 Starting test suite...")

	// Define test categories in order
	testCategories := []struct {
		name  string
		tests []utils.TestCase
	}{
		{"Cleanup", cases.GetCleanupTests(ctx, svc, userId, contextName, workspaceName)},
		{"Setup", cases.GetSetupTests(ctx, svc, userId, contextName, workspaceName)},
		{"Management", cases.GetManagementTests(ctx, svc, userId, contextName, workspaceName)},
		{"Terraform", cases.GetTerraformTests(ctx, svc, userId, contextName, workspaceName)},
	}

	// Run all test categories
	for _, category := range testCategories {
		log.Printf("\n=== Running %s Tests ===\n", category.name)
		for _, tc := range category.tests {
			if err := utils.LogTestCase(tc); err != nil {
				return err
			}
		}
	}

	log.Println("\n✨ Test suite completed. Resources are preserved for manual inspection.")
	return nil
}
