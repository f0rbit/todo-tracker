import { z } from "zod";


export const CONFIG_SCHEMA = z.object({
    tags: z.array(z.object({
        name: z.string(),
        match: z.array(z.string())
    })),
    ignore: z.array(z.string()).optional()
});

export type Config = z.infer<typeof CONFIG_SCHEMA>;

export type WorkerMessage = {
    chunk: string[];
    config: Config["tags"];
    directory: string;
};

export const TASK_SCHEMA = z.object({
    id: z.string(),
    file: z.string(),
    line: z.number(),
    tag: z.string(),
    text: z.string(),
    context: z.string()
});

export const TASK_FILE = z.array(TASK_SCHEMA);

export type ParsedTask = z.infer<typeof TASK_SCHEMA>;

export type DiffInfo = {
    text: string;
    line: number;
    file: string;
}

export type DiffResult = {
    id: string;
    tag: string;
    type: "NEW" | "UPDATE" | "MOVE" | "DELETE" | "SAME";
    data: {
        old: DiffInfo | null,
        new: DiffInfo | null
    };
};

