package main

import (
	"encoding/json"
	"os"
	"strings"
)

func read(path string) ([]ParsedTask, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tasks []ParsedTask
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func extractTask(task ParsedTask) *DiffInfo {
	return &DiffInfo{
		Text: task.Text,
		Line: task.Line,
		File: task.File,
	}
}

func sameText(task1, task2 ParsedTask) bool {
	return strings.TrimSpace(task1.Text) == strings.TrimSpace(task2.Text)
}

func generate_diff(dir1 string, dir2 string) ([]DiffResult, error) {
	base_tasks, err := read(dir1)
	if err != nil {
		return nil, err
	}

	new_tasks, err := read(dir2)
	if err != nil {
		return nil, err
	}

	var diffs []DiffResult

	for _, new_task := range new_tasks {
		found := false

		// search for the same text but different line or file
		for _, base_task := range base_tasks {
			if sameText(base_task, new_task) {
				diffType := "SAME"
				if base_task.Line != new_task.Line || base_task.File != new_task.File {
					diffType = "MOVE"
				}

				diffs = append(diffs, DiffResult{
					ID:   base_task.ID,
					Tag:  new_task.Tag,
					Type: diffType,
					Data: DiffData{
						Old: extractTask(base_task),
						New: extractTask(new_task),
					},
				})
				found = true
				break
			}
		}
		if found {
			continue
		}

		// then we need to see if line number and tag are same but text is different
		for _, base_task := range base_tasks {
			if base_task.Line == new_task.Line && base_task.Tag == new_task.Tag && !sameText(base_task, new_task) {
				diffs = append(diffs, DiffResult{
					ID:   base_task.ID,
					Tag:  new_task.Tag,
					Type: "UPDATE",
					Data: DiffData{
						Old: extractTask(base_task),
						New: extractTask(new_task),
					},
				})
				found = true
				break
			}
		}
		if found {
			continue
		}

		diffs = append(diffs, DiffResult{
			ID:   new_task.ID,
			Tag:  new_task.Tag,
			Type: "NEW",
			Data: DiffData{
				Old: nil,
				New: extractTask(new_task),
			},
		})
	}

	for _, base_task := range base_tasks {
		found := false
		for _, diff := range diffs {
			if diff.ID == base_task.ID {
				found = true
				break
			}
		}

		if !found {
			diffs = append(diffs, DiffResult{
				ID:   base_task.ID,
				Tag:  base_task.Tag,
				Type: "DELETE",
				Data: DiffData{
					Old: extractTask(base_task),
					New: nil,
				},
			})
		}
	}

	return diffs, nil
}
