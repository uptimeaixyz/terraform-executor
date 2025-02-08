package utils

import (
	"fmt"
	"log"
	"runtime/debug"
)

// TestCase represents a single test case with its status
type TestCase struct {
	Name     string
	Fn       func() error
	Category string
}

// LogTestCase executes and logs a test case with proper formatting
func LogTestCase(tc TestCase) error {
	log.Printf("=== TEST CASE [%s]: %s ===", tc.Category, tc.Name)
	if err := tc.Fn(); err != nil {
		log.Printf("\n❌ FAILED: %s\nError Details:\n%+v\nStack Trace:\n%s\n",
			tc.Name,
			err,
			string(debug.Stack()),
		)
		return fmt.Errorf("%s:\n%+v", tc.Name, err)
	}
	log.Printf("✅ PASSED: %s", tc.Name)
	return nil
}
