package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	// Verify that we have at least one argument for the operation type
	if len(flag.Args()) < 1 {
		fmt.Println("Error: No operation type specified. Expected 'parse' or 'diff'")
		os.Exit(1)
	}

	operation := flag.Arg(0)

	switch operation {
	case "parse":
		parse_args()
	case "diff":
		diff_args()
	default:
		fmt.Printf("Error: Unknown operation '%s'. Expected 'parse' or 'diff'\n", operation)
		os.Exit(1)
	}
}
func parse_args() {
	if len(flag.Args()) != 3 { // Because the first arg is "parse"
		fmt.Println("Error: Incorrect number of arguments for 'parse'. Expected <directory> <config.json location>")
		os.Exit(1)
	}
	directory := flag.Arg(1)
	configFile := flag.Arg(2)
	config, err := read_config(configFile)

	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	tasks, err := parse_dir(directory, config)
	if err != nil {
		fmt.Printf("Error parsing directory: %v\n", err)
		os.Exit(1)
	}

	// tasks is a slice of ParsedTask structs, we want to print out this in JSON format
	tasksJSON, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling tasks to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(tasksJSON))
}

func diff_args() {
	if len(flag.Args()) != 3 { // Because the first arg is "diff"
		fmt.Println("Error: Incorrect number of arguments for 'diff'. Expected <previous_json> <new_json>")
		os.Exit(1)
	}
	previousJSON := flag.Arg(1)
	newJSON := flag.Arg(2)

	diff, err := generate_diff(previousJSON, newJSON)
	if err != nil {
		fmt.Printf("Error generating diff: %v\n", err)
		os.Exit(1)
	}

	// diff is a slice of DiffResult structs, we want to print out this in JSON format
	diffJSON, err := json.MarshalIndent(diff, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling diff to JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(diffJSON))
}

func read_config(filePath string) (*Config, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("error parsing config JSON: %v", err)
	}
	return &config, nil
}
