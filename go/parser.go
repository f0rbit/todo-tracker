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

func parse_dir(dirPath string, config *Config) ([]ParsedTask, error) {
	// Convert ignore patterns to regexes
	var ignoreRegexes []*regexp.Regexp
	for _, pattern := range config.Ignore {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		ignoreRegexes = append(ignoreRegexes, regex)
	}

	var tasks []ParsedTask
	var tasksMutex sync.Mutex // To safely append to tasks from multiple goroutines
	var wg sync.WaitGroup

	// filepath.Walk to traverse directories and files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking the path %q: %v\n", path, err)
			return err
		}

		// Skip directories and ignored files
		if info.IsDir() {
			for _, regex := range ignoreRegexes {
				if regex.MatchString(path) {
					return filepath.SkipDir // Skip the entire directory
				}
			}
			return nil
		}

		for _, regex := range ignoreRegexes {
			if regex.MatchString(path) {
				return nil // Skip this file
			}
		}

		// Skip binary files
		if isBinaryFile(path) {
			return nil // Skip this file
		}

		// Process files in parallel
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Placeholder for file processing logic
			foundTasks, err := processFile(path, config.Tags, dirPath+"/")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", path, err)
				return
			}

			tasksMutex.Lock()
			tasks = append(tasks, foundTasks...)
			tasksMutex.Unlock()
		}()

		return nil
	})

	wg.Wait()

	return tasks, err
}

func isBinaryFile(path string) bool {
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

func processFile(filePath string, config []Tag, dirPath string) ([]ParsedTask, error) {
	file, err := os.Open(filePath)
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

	// Remove the common prefix from the file dirPath
	var file_name = strings.TrimPrefix(filePath, dirPath)

	// Iterating through each line and check against each tag's match criteria
	for i, line := range lines {
		for _, tagConfig := range config {
			for _, match := range tagConfig.Match {
				if strings.Contains(line, match) {
					uuid, err := generateUUID()
					if err != nil {
						return nil, err
					}

					start := i - 3
					if start < 0 {
						start = 0
					}
					end := i + 4
					if end > len(lines) {
						end = len(lines)
					}
					context := lines[start:end]

					tasks = append(tasks, ParsedTask{
						ID:      uuid,
						File:    file_name,
						Line:    i + 1, // Line numbers are usually 1-indexed
						Tag:     tagConfig.Name,
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
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
