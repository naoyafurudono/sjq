package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
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
		fmt.Fprintf(os.Stderr, `sjq - Stream JSON Separator (v1.0.0)

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
`)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Println("sjq version 1.0.0")
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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if isJSON(line) {
			fmt.Fprintln(jsonWriter, line)
		} else {
			fmt.Fprintln(nonJSONWriter, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}
}

func isJSON(line string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(line), &js) == nil
}
