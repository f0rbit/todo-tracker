bun src/parser.ts --dir "/Users/tom/Documents/dev/todo-tracker/resources/codebase" --config "/Users/tom/Documents/dev/todo-tracker/resources/config.json" > output-base.json
bun src/parser.ts --dir "/Users/tom/Documents/dev/todo-tracker/resources/codebase-changed" --config "/Users/tom/Documents/dev/todo-tracker/resources/config.json" > output-new.json
bun src/diff.ts
