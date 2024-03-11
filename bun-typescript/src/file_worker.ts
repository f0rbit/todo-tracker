import { randomUUID } from "node:crypto";
import { EOL } from "node:os";
import { ParsedTask, WorkerMessage } from "./schema";

// @ts-ignore
const me = self;

me.onmessage = async (message: MessageEvent) => {
    const { chunk, config, directory } = message.data as WorkerMessage;
    const found_lines: ParsedTask[] = [];

    for (const f of chunk) {
        const path = `${directory}/${f}` as const;
        const bun_file = Bun.file(path);
        const name = bun_file.name;
        if (!name) continue;

        const is_dir = bun_file.type.includes("octet-stream");
        if (is_dir) continue;
        const lines = (await bun_file.text()).split("\n");
        lines.forEach((l, i) => {
            config.forEach((c) => {
                c.match.forEach((m) => {
                    if (!l.includes(m)) return;
                    const uuid = randomUUID();
                    const file_name = name.substring(directory.length + 1);
                    found_lines.push({ id: uuid, file: file_name, line: i + 1, tag: c.name, text: l, context: lines.slice(i - 3, i + 3).join(EOL) });
                });
            });
        });
    }

    postMessage(found_lines);
}
