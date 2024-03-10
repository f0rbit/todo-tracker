import { readdir } from "node:fs/promises";
import config from "./config";

console.log("Starting main script");

const DIR = `../resources/codebase`;

console.log(`Reading files from directory: ${DIR}`);
const files = await readdir(DIR, { recursive: true });
console.log(`Total files read: ${files.length}`);

// Split files array into chunks for each worker
const numWorkers = 4; // Or any number you find appropriate
const chunkSize = Math.ceil(files.length / numWorkers);
console.log(`Number of workers: ${numWorkers}, Chunk size: ${chunkSize}`);

const filesChunks = Array.from({ length: numWorkers }, (_, i) =>
  files.slice(i * chunkSize, (i + 1) * chunkSize)
);

const workers = [];
const promises = [];

for (const chunk of filesChunks) {
    console.log(`Sending chunk to worker, chunk size: ${chunk.length}`);
    const worker = new Worker(new URL('./worker.js', import.meta.url));
    workers.push(worker);
    const promise = new Promise((resolve) => {
        worker.onmessage = (message) => {
            console.log("Message received from worker");
            resolve(message.data);
        };
        worker.onerror = (error) => {
            console.log(`Worker error: ${error.message}`);
        };
    });
    promises.push(promise);
    worker.postMessage({ filesChunk: chunk, config });
    console.log("Post message sent to worker");
}

// Wait for all workers to complete
console.log("Waiting for all workers to complete");
const results = (await Promise.all(promises)).flat();
console.log("All workers have completed");

// Clean up
console.log("Terminating all workers");
workers.forEach(worker => worker.terminate());

console.log("Results:", results);
