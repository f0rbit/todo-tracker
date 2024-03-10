import { EOL } from "node:os";

const DIR = `../resources/codebase`; // Ensure this is correct relative to the worker script

self.onmessage = async (message) => {
    const { chunk, config } = message.data;
    const found_lines = [];

    for (const f of chunk) {
        const path = `${DIR}/${f}`;
        const bun_file = Bun.file(path);

        const is_dir = bun_file.type.includes("octet-stream");
        if (!is_dir) {
            const lines = (await bun_file.text()).split("\n");
            lines.forEach((l, i) => {
                config.forEach((c) => {
                    c.match.forEach((m) => {
                        if (l.includes(m)) {
                            const file_name = bun_file.name.substring(DIR.length + 1);
                            found_lines.push({ file: file_name, line: i, tag: c.name, text: l, context: lines.slice(i - 3, i + 3).join(EOL) });
                        }
                    });
                });
            });
        }
    }

    postMessage(found_lines);
};
