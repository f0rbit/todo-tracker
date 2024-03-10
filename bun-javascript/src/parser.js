import { readdir } from "node:fs/promises";
import config from "./config";
import args from "./args";

const DIR = `${args.dir}`;

const files = await readdir(DIR, { recursive: true });

// Split files array into chunks for each worker
const num_workers = 4; // Or any number you find appropriate
const chunk_size = Math.ceil(files.length / num_workers);

const chunks = Array.from({ length: num_workers }, (_, i) =>
    files.slice(i * chunk_size, (i + 1) * chunk_size)
);

const workers = [];
const promises = [];

for (const chunk of chunks) {
    const worker = new Worker(new URL('./file_worker.js', import.meta.url));
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
    worker.postMessage({ chunk, config, directory: DIR });
}

// Wait for all workers to complete
const results = (await Promise.all(promises)).flat();
// Clean up
workers.forEach(worker => worker.terminate());

console.log(JSON.stringify(results, null, 2));
