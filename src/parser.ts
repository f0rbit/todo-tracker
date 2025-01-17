import { readdir } from "node:fs/promises";
import config from "./config";
import args from "./args";
import { WorkerMessage } from "./schema";

const DIR = `${args.dir}`;

const crawl_files = await readdir(DIR, { recursive: true });

const ignore_regexes = config.ignore.map((i) => new RegExp(i));

const files = crawl_files.filter((f) => {
    return !(ignore_regexes.some((r) => r.test(f)));
});

// Split files array into chunks for each worker
const num_workers = 4; // Or any number you find appropriate
const chunk_size = Math.ceil(files.length / num_workers);

const chunks = Array.from({ length: num_workers }, (_, i) =>
    files.slice(i * chunk_size, (i + 1) * chunk_size)
)

const workers: Worker[] = [];
const promises: Promise<any>[] = [];

const url = new URL("./file_worker.js", import.meta.url);
if (!url) throw new Error("URL of file_worker.js not found");

for (const chunk of chunks) {
    const worker = new Worker(url.toString());
    workers.push(worker);

    const promise = new Promise((resolve) => {
        worker.onmessage = (message) => {
            resolve(message.data);
        };
        worker.onerror = (error) => {
            console.error(`Worker error: ${error.message}`);
        };
    });
    promises.push(promise);
    const worker_message: WorkerMessage = { chunk, config: config.tags, directory: DIR };
    worker.postMessage(worker_message);
}

// Wait for all workers to complete
const results = (await Promise.all(promises)).flat();
// Clean up
workers.forEach(worker => worker.terminate());

console.log(JSON.stringify(results, null, 2));
