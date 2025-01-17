bun src/parser.js --dir "./resources/codebase" --config "./resources/config.json" > output-base.json
bun src/parser.js --dir "./resources/codebase-changed" --config "./resources/config.json" > output-new.json
bun src/diff.js
