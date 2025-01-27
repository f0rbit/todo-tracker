## Overview
`todo-tracker` is a tool to traverse a directory, identify tasks based on specified tags, and compare the tasks found with an existing dataset to produce a diff. It is implemented in Go for efficient **multi-threaded traversal**.

## Building
To build the application, run `make build` or:
```sh
go build -o todo-tracker
```

## Usage
### Arguments
- `parse <directory> <config.json>`: Traverse the directory and identify tasks based on the configuration.
- `diff <previous_json> <new_json>`: Compare two JSON outputs and produce a diff.

### Examples
#### Parsing a directory
```sh
./todo-tracker parse ./resources/codebase ./resources/config.json
```

#### Parsing a changed directory
```sh
./todo-tracker parse ./resources/codebase-changed ./resources/config.json
```

#### Generating a diff
```sh
./todo-tracker diff output-base.json output-new.json
```

## Outputs
### ParsedTask
```typescript
type ParsedTask = {
    id: string;
    file: string;
    line: number;
    tag: string;
    text: string;
    context: string[];
};
```

### DiffResult
```typescript
type DiffResult = {
    id: string;
    tag: string;
    type: "SAME" | "MOVE" | "UPDATE" | "NEW" | "DELETE";
    data: {
        old: DiffInfo | null;
        new: DiffInfo | null;
    };
};
```

### DiffInfo
```typescript
type DiffInfo = {
    text: string;
    line: number;
    file: string;
    context: string[];
};
```

## Configuration

