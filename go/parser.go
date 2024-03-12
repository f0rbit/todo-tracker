package main

import (
    "bufio"
    "os"
    "path/filepath"
    "regexp"
    "sync"
    "fmt"
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

        // Process files in parallel
        wg.Add(1)
        go func() {
            defer wg.Done()

            // Placeholder for file processing logic
            foundTasks, err := processFile(path, config)
            if err != nil {
                fmt.Printf("Error processing file %s: %v\n", path, err)
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

func processFile(filePath string, config *Config) ([]ParsedTask, error) {
    // Implement the logic for processing a single file
    // This would be the logic inside your `file_worker.ts`

    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    fmt.Printf("Processing file %s\n", filePath)

    var tasks []ParsedTask
    scanner := bufio.NewScanner(file)
    lineNumber := 0

    for scanner.Scan() {
        lineNumber++
        scanner.Text()

        // Placeholder for logic to check line against config and populate tasks

    }

    return tasks, scanner.Err()
}
