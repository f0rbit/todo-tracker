package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func parse_dir(dir_path string, config *Config) ([]ParsedTask, error) {
	// Convert ignore patterns to regexes
	var ignore_regexes []*regexp.Regexp
	for _, pattern := range config.Ignore {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		ignore_regexes = append(ignore_regexes, regex)
	}

	var tasks []ParsedTask
	var tasks_mutex sync.Mutex // To safely append to tasks from multiple goroutines
	var wg sync.WaitGroup

	// filepath.Walk to traverse directories and files
	err := filepath.Walk(dir_path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking the path %q: %v\n", path, err)
			return err
		}

		// Skip directories and ignored files
		if info.IsDir() {
			for _, regex := range ignore_regexes {
				if regex.MatchString(path) {
					return filepath.SkipDir // Skip the entire directory
				}
			}
			return nil
		}

		for _, regex := range ignore_regexes {
			if regex.MatchString(path) {
				return nil // Skip this file
			}
		}

		// Skip binary files
		if is_binary_file(path) {
			return nil // Skip this file
		}

		// Process files in parallel
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Placeholder for file processing logic
			found_tasks, err := process_file(path, config.Tags, dir_path+"/")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", path, err)
				return
			}

			tasks_mutex.Lock()
			tasks = append(tasks, found_tasks...)
			tasks_mutex.Unlock()
		}()

		return nil
	})

	wg.Wait()

	return tasks, err
}

func is_binary_file(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return false
	}

	for _, b := range buf {
		if b == 0 {
			return true
		}
	}
	return false
}

func process_file(file_path string, config []Tag, dir_path string) ([]ParsedTask, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tasks []ParsedTask
	scanner := bufio.NewScanner(file)
	var lines []string // Store all lines to use for context

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Remove the common prefix from the file dir_path
	var file_name = strings.TrimPrefix(file_path, dir_path)

	// Iterating through each line and check against each tag's match criteria
	for i, line := range lines {
		for _, tag_config := range config {
			for _, match := range tag_config.Match {
				if strings.Contains(line, match) {
					uuid, err := generate_uuid()
					if err != nil {
						return nil, err
					}
					/** @todo make the context start/end configurable in .json */
					start := i - 4
					if start < 0 {
						start = 0
					}
					end := i + 6
					if end > len(lines) {
						end = len(lines)
					}
					context := lines[start:end]

					tasks = append(tasks, ParsedTask{
						ID:      uuid,
						File:    file_name,
						Line:    i + 1, // Line numbers are usually 1-indexed
						Tag:     tag_config.Name,
						Text:    line,
						Context: context,
					})
				}
			}
		}
	}

	return tasks, nil
}

/** @todo replace this with google's UUID library */
func generate_uuid() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
