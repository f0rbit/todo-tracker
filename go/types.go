package main

type Config struct {
    Tags   []Tag `json:"tags"`
    Ignore []string `json:"ignore,omitempty"`
}

type Tag struct {
    Name  string   `json:"name"`
    Match []string `json:"match"`
}
