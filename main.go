package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var (
		nonJSONOnly bool
		jsonFile    string
		nonJSONFile string
		showHelp    bool
		showVersion bool
	)

	flag.BoolVar(&showHelp, "h", false, "Show help message")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&nonJSONOnly, "n", false, "Output only non-JSON lines to stdout")
	flag.StringVar(&jsonFile, "json", "", "Output JSON lines to specified file")
	flag.StringVar(&nonJSONFile, "non-json", "", "Output non-JSON lines to specified file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `sjq - Stream JSON Separator (v%s)

A lightweight tool to separate JSON and non-JSON lines from mixed log streams.
Zero dependencies, pure Go implementation.

Usage:
  sjq [OPTIONS]

Examples:
  # Extract only JSON lines
  cat app.log | sjq
  
  # Extract only non-JSON lines
  cat app.log | sjq -n
  
  # Separate into different files
  cat app.log | sjq --json structured.log --non-json plain.log

Options:
`, Version)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("sjq version %s\n", Version)
		os.Exit(0)
	}

	var jsonWriter io.Writer = os.Stdout
	var nonJSONWriter io.Writer = io.Discard

	if nonJSONOnly {
		jsonWriter = io.Discard
		nonJSONWriter = os.Stdout
	}

	if jsonFile != "" {
		f, err := os.Create(jsonFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating JSON file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		jsonWriter = f
	}

	if nonJSONFile != "" {
		f, err := os.Create(nonJSONFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating non-JSON file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		nonJSONWriter = f
	}

	processInput(os.Stdin, jsonWriter, nonJSONWriter)
}

func isJSON(line string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(line), &js) == nil
}

func processInput(reader io.Reader, jsonWriter, nonJSONWriter io.Writer) {
	scanner := bufio.NewScanner(reader)
	var buffer []string
	var inJSON bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Check if this line could be the start of a JSON object or array
		if !inJSON && (strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")) {
			buffer = append(buffer, line)
			inJSON = true
			// Check if it's a complete JSON on one line
			if isCompleteJSON(strings.Join(buffer, "\n")) {
				fmt.Fprintln(jsonWriter, strings.Join(buffer, "\n"))
				buffer = nil
				inJSON = false
			}
		} else if inJSON {
			buffer = append(buffer, line)
			// Check if we have a complete JSON
			if isCompleteJSON(strings.Join(buffer, "\n")) {
				fmt.Fprintln(jsonWriter, strings.Join(buffer, "\n"))
				buffer = nil
				inJSON = false
			}
		} else {
			// Single line, check if it's JSON
			if isJSON(line) {
				fmt.Fprintln(jsonWriter, line)
			} else {
				fmt.Fprintln(nonJSONWriter, line)
			}
		}
	}

	// Handle any remaining buffer content
	if len(buffer) > 0 {
		// Not a complete JSON, treat as non-JSON lines
		for _, line := range buffer {
			fmt.Fprintln(nonJSONWriter, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func isCompleteJSON(s string) bool {
	var js json.RawMessage
	err := json.Unmarshal([]byte(s), &js)
	return err == nil
}
