bun src/parser.ts --dir "/home/tom/dev/github/todo-tracker/resources/codebase" --config "/home/tom/dev/github/todo-tracker/resources/config.json" > output-base.json
bun src/parser.ts --dir "/home/tom/dev/github/todo-tracker/resources/codebase-changed" --config "/home/tom/dev/github/todo-tracker/resources/config.json" > output-new.json
bun src/diff.ts
