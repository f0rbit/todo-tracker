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
		is_binary, err := is_binary_file(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking if file %s is binary: %v\n", path, err)
			return nil
		}
		if is_binary {
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

func is_binary_file(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read a small chunk from the file (e.g., 1024 bytes).
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil && err.Error() != "EOF" {
		// If EOF is reached with fewer bytes, that's not really an error,
		// but if something else happened, handle it.
		return false, err
	}

	// Only check the bytes we actually read.
	buf = buf[:n]

	// Simple heuristic: if we see a NULL byte, we consider it binary.
	for _, b := range buf {
		if b == 0x00 {
			return true, nil
		}
	}

	// No NULL byte found, likely a text file.
	return false, nil
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
				if match == "" {
					continue
				}
				matched_index := strings.Index(line, match)
				if matched_index == -1 {
					continue
				}
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

				// extract the text after the matched index
				after_index := matched_index + len(match)
				if after_index >= len(line) {
					after_index = len(line)
				}

				// if the line is "empty" after the match, we should start from the beginning of line to include as much context
				if len(strings.TrimSpace(line[after_index:])) < 3 {
					after_index = 0
				}

				text := strings.TrimSpace(line[after_index:])

				// if text ends with "*/" we can remove the ending
				text = strings.TrimSuffix(text, "*/")

				// then remove spaces again
				text = strings.TrimSpace(text)

				tasks = append(tasks, ParsedTask{
					ID:      uuid,
					File:    file_name,
					Line:    i + 1, // Line numbers are usually 1-indexed
					Tag:     tag_config.Name,
					Text:    text,
					Context: context,
				})
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
