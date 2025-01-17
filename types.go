package main

type Config struct {
	Tags   []Tag    `json:"tags"`
	Ignore []string `json:"ignore,omitempty"`
}

type Tag struct {
	Name  string   `json:"name"`
	Match []string `json:"match"`
}

type ParsedTask struct {
	ID      string   `json:"id"`
	File    string   `json:"file"`
	Line    int      `json:"line"`
	Tag     string   `json:"tag"`
	Text    string   `json:"text"`
	Context []string `json:"context"` // assuming context is a slice of lines
}

type DiffResult struct {
	ID   string   `json:"id"`
	Tag  string   `json:"tag"`
	Type string   `json:"type"`
	Data DiffData `json:"data"`
}

type DiffData struct {
	Old *DiffInfo `json:"old"`
	New *DiffInfo `json:"new"`
}

type DiffInfo struct {
	Text    string   `json:"text"`
	Line    int      `json:"line"`
	File    string   `json:"file"`
	Context []string `json:"context"`
}
