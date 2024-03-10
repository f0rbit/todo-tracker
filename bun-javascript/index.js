import { readdir } from "node:fs/promises";
import { EOL } from "node:os";

const DIR = `../resources/codebase`

// read all the files in the current directory, recursively
const files = await readdir(DIR, { recursive: true });

/** @type {import("bun").BunFile[]} */
const search = [];

for (const f of files) {
    const path = `${DIR}/${f}`;
    const bun_file = Bun.file(path);

    const is_dir = bun_file.type.includes("octet-stream");
    if (!is_dir) search.push(bun_file);
}

const config = [{ name: "todo", match: ["@todo", "TODO"] }];

const found_lines = [];

for (const f of search) {
    const lines = (await f.text()).split("\n");
    lines.forEach((l, i) => {
        config.forEach((c) => {
            c.match.forEach((m) => {
                if (l.includes(m)) {
                    const file_name = f.name.substring(DIR.length + 1);
                    found_lines.push({ file: file_name, line: i, tag: c.name, text: l, context: lines.slice(i - 3, i + 3).join(EOL) });
                }
            });
        });
    });
}


console.log(found_lines);

