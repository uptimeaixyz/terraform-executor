package utils

import (
	"os"
)

// appendToFile appends text to a file, creating it if it doesn't exist.
func AppendToFile(filename, text string) error {
	// Open file in append mode, create it if necessary
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the text to the file
	if _, err := file.WriteString(text + "\n"); err != nil {
		return err
	}

	return nil
}
