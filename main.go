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
	)

	flag.BoolVar(&nonJSONOnly, "n", false, "Output only non-JSON lines to stdout")
	flag.StringVar(&jsonFile, "json", "", "Output JSON lines to specified file")
	flag.StringVar(&nonJSONFile, "non-json", "", "Output non-JSON lines to specified file")
	flag.Parse()

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
