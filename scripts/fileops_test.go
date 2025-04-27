package main

import (
	"testing"
)

func TestStandard(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError bool
	}{
		{
			name:          "Valid date prefix",
			input:         "2023-10-15-example.md",
			expected:      "2023-10-15-example.md",
			expectedError: false,
		},
		{
			name:          "Valid date without delimiters",
			input:         "202504231736-pub-example.md",
			expected:      "2025-04-23-example.md",
			expectedError: false,
		},
		{
			name:          "Unpublished file",
			input:         "202504231736-example.md",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "Invalid filename pattern",
			input:         "example.txt",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "Invalid date format",
			input:         "2023-1015-example.txt",
			expected:      "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Standard(tt.input)
			if (err != nil) != tt.expectedError {
				t.Errorf("Standard() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if result != tt.expected {
				t.Errorf("Standard() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
