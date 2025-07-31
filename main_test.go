package main

import (
	"bytes"
	"strings"
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

func TestProcessInput(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedJSON    string
		expectedNonJSON string
	}{
		{
			name:            "single line JSON",
			input:           `{"message": "hello"}`,
			expectedJSON:    `{"message": "hello"}` + "\n",
			expectedNonJSON: "",
		},
		{
			name:            "multiline JSON object",
			input:           "{\n  \"name\": \"test\",\n  \"value\": 123\n}",
			expectedJSON:    "{\n  \"name\": \"test\",\n  \"value\": 123\n}\n",
			expectedNonJSON: "",
		},
		{
			name:            "multiline JSON array",
			input:           "[\n  \"item1\",\n  \"item2\",\n  \"item3\"\n]",
			expectedJSON:    "[\n  \"item1\",\n  \"item2\",\n  \"item3\"\n]\n",
			expectedNonJSON: "",
		},
		{
			name:            "mixed content",
			input:           "Starting app\n{\"level\":\"info\",\"msg\":\"initialized\"}\nProcessing...\n[\n  1,\n  2,\n  3\n]\nDone",
			expectedJSON:    "{\"level\":\"info\",\"msg\":\"initialized\"}\n[\n  1,\n  2,\n  3\n]\n",
			expectedNonJSON: "Starting app\nProcessing...\nDone\n",
		},
		{
			name:            "incomplete JSON treated as non-JSON",
			input:           "{\n  \"incomplete\": true",
			expectedJSON:    "",
			expectedNonJSON: "{\n  \"incomplete\": true\n",
		},
		{
			name:            "nested JSON",
			input:           "{\n  \"outer\": {\n    \"inner\": [\n      {\"id\": 1},\n      {\"id\": 2}\n    ]\n  }\n}",
			expectedJSON:    "{\n  \"outer\": {\n    \"inner\": [\n      {\"id\": 1},\n      {\"id\": 2}\n    ]\n  }\n}\n",
			expectedNonJSON: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var jsonBuf, nonJSONBuf bytes.Buffer

			processInput(reader, &jsonBuf, &nonJSONBuf)

			if jsonBuf.String() != tt.expectedJSON {
				t.Errorf("JSON output mismatch\ngot:\n%q\nwant:\n%q", jsonBuf.String(), tt.expectedJSON)
			}

			if nonJSONBuf.String() != tt.expectedNonJSON {
				t.Errorf("Non-JSON output mismatch\ngot:\n%q\nwant:\n%q", nonJSONBuf.String(), tt.expectedNonJSON)
			}
		})
	}
}

func TestIsCompleteJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "complete object",
			input:    `{"key": "value"}`,
			expected: true,
		},
		{
			name:     "complete array",
			input:    `[1, 2, 3]`,
			expected: true,
		},
		{
			name:     "incomplete object",
			input:    `{"key": "value"`,
			expected: false,
		},
		{
			name:     "incomplete array",
			input:    `[1, 2, 3`,
			expected: false,
		},
		{
			name:     "multiline complete JSON",
			input:    "{\n  \"key\": \"value\"\n}",
			expected: true,
		},
		{
			name:     "multiline incomplete JSON",
			input:    "{\n  \"key\": \"value\"",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCompleteJSON(tt.input)
			if result != tt.expected {
				t.Errorf("isCompleteJSON(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
