package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnvFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer CloseFile(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Split into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // or return error if you prefer strictness
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
		fmt.Println("Set environment variable:", key, "=", value)
	}

	return scanner.Err()
}
