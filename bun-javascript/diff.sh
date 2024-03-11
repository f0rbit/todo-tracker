bun src/parser.js --dir "../resources/codebase" --config "../config.json" > output-base.json
bun src/parser.js --dir "../resources/codebase-changed" --config "../config.json" > output-new.json
bun src/diff.js
