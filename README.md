## Overview
`todo-tracker` is a tool to traverse a directory, identify tasks based on specified tags, and compare the tasks found with an existing dataset to produce a diff. It is implemented in Go for efficient **multi-threaded traversal**.

## Building
To build the application, run `make build` or:
```bash
go build -o todo-tracker
```

## Usage
### Arguments
- `parse <directory> <config.json>`: Traverse the directory and identify tasks based on the configuration => `ParsedTask[]`
- `diff <previous_json> <new_json>`: Compare two JSON outputs and produce a diff => `DiffResult[]`

### Examples
Parsing a directory
```bash
./todo-tracker parse ./resources/codebase ./resources/config.json
```

Parsing a changed directory
```bash
./todo-tracker parse ./resources/codebase-changed ./resources/config.json
```

Generating a diff
```bash
./todo-tracker diff output-base.json output-new.json
```

### Testing
The `make diff` command will run through an example parse of `resources/codebase` and `resources/codebase-changed` and generate an example diff.

## Outputs
```typescript
type ParsedTask = {
    id: string;
    file: string;
    line: number;
    tag: string;
    text: string;
    context: string[];
};

type DiffResult = {
    id: string;
    tag: string;
    type: "SAME" | "MOVE" | "UPDATE" | "NEW" | "DELETE";
    data: {
        old: DiffInfo | null;
        new: DiffInfo | null;
    };
};

type DiffInfo = {
    text: string;
    line: number;
    file: string;
    context: string[];
};
```

## Configuration
The configuration file is a JSON file that specifies the tags to look for and the files or directories to ignore.

### Example Configuration
```jsonc
{
    "tags": [
        {
            "name": "todo",
            "match": [
                "@todo",
                "TODO:"
            ]
        },
        {
            "name": "bug",
            "match": [
                "@bug",
                "BUG:"
            ]
        },
        {
            "name": "note",
            "match": [
                "@note",
                "NOTE:"
            ]
        }
    ],
    "ignore": [
        "node_modules",
		"readme.md"
    ]
}
```
> [resources/config.json](resources/config.json)


### Configuration Schema
```typescript
type Config = {
    tags: Tag[];
    ignore?: string[];
};

type Tag = {
    name: string;
    match: string[];
};
```