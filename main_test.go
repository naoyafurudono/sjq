package main

import (
	"testing"
)

func TestIsJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid JSON object",
			input:    `{"timestamp": "2023-04-01T12:00:00Z", "message": "Hello, world!"}`,
			expected: true,
		},
		{
			name:     "valid JSON array",
			input:    `["item1", "item2", "item3"]`,
			expected: true,
		},
		{
			name:     "valid JSON string",
			input:    `"hello world"`,
			expected: true,
		},
		{
			name:     "valid JSON number",
			input:    `42`,
			expected: true,
		},
		{
			name:     "valid JSON boolean",
			input:    `true`,
			expected: true,
		},
		{
			name:     "valid JSON null",
			input:    `null`,
			expected: true,
		},
		{
			name:     "invalid JSON - plain text",
			input:    `application started`,
			expected: false,
		},
		{
			name:     "invalid JSON - malformed object",
			input:    `{"key": "value"`,
			expected: false,
		},
		{
			name:     "empty string",
			input:    ``,
			expected: false,
		},
		{
			name:     "whitespace only",
			input:    `   `,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSON(tt.input)
			if result != tt.expected {
				t.Errorf("isJSON(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
