bun src/parser.ts --dir "../resources/codebase" --config "../resources/config.json" > output-base.json
bun src/parser.ts --dir "../resources/codebase-changed" --config "../resources/config.json" > output-new.json
bun src/diff.ts
