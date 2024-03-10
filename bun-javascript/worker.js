import { EOL } from "node:os";

console.log("Worker script started");

const DIR = `../resources/codebase`; // Ensure this is correct relative to the worker script

self.onmessage = async (message) => {
    console.log("Message received in worker");
    const { filesChunk, config } = message.data;
    const found_lines = [];

    for (const f of filesChunk) {
        console.log(`Processing file: ${f}`);
        const path = `${DIR}/${f}`;
        const bun_file = Bun.file(path);

        const is_dir = bun_file.type.includes("octet-stream");
        if (!is_dir) {
            const lines = (await bun_file.text()).split("\n");
            console.log(`Lines found in file: ${lines.length}`);
            lines.forEach((l, i) => {
                config.forEach((c) => {
                    c.match.forEach((m) => {
                        if (l.includes(m)) {
                            const file_name = bun_file.name.substring(DIR.length + 1);
                            found_lines.push({ file: file_name, line: i, tag: c.name, text: l, context: lines.slice(i - 3, i + 3).join(EOL) });
                            console.log(`Match found: ${l}`);
                        }
                    });
                });
            });
        }
    }

    console.log(`Worker done, found lines: ${found_lines.length}`);
    postMessage(found_lines);
};
