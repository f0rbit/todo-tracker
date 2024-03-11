import fs from 'node:fs';

// Function to read and parse a JSON file
function readJsonFile(filePath) {
    const data = fs.readFileSync(filePath, 'utf8');
    return JSON.parse(data);
}

// Determine if two tasks are equivalent
function areTasksEquivalent(task1, task2) {
    return (task1.line === task2.line || task1.text.trim() === task2.text.trim()) && task1.tag === task2.tag;
}

// Determine if a task has been moved
function isTaskMoved(task, otherTask) {
    return areTasksEquivalent(task, otherTask) && (task.file !== otherTask.file || task.line !== otherTask.line);
}

// Determine if a task has been updated
function isTaskUpdated(task, otherTask) {
    return task.file === otherTask.file && task.line === otherTask.line && (task.text.trim() !== otherTask.text.trim() || task.tag !== otherTask.tag);
}

// Read the base and new task data
const baseTasks = readJsonFile('./output-base.json');
const newTasks = readJsonFile('./output-new.json');

/** @typedef {{ tag: string, type: "NEW" | "UPDATE" | "MOVE" | "DELETE", data: { old: { text: string, line: number, file: string }, new: { text: string, line: number, file: string } }}} DiffResult */

/** @type {DiffResult[]} */
const diffs = [];

const extract_task = (task) => {
    const { text, line, file } = task;
    return { text, line, file };
}


// process the tasks
newTasks.forEach(new_task => {
    const base_task = baseTasks.find(bt => areTasksEquivalent(bt, new_task));
    const _old = base_task ? extract_task(base_task) : null;
    const _new = extract_task(new_task);
    const data = { old: _old, new: _new };

    if (!base_task) {
        diffs.push({ tag: new_task.tag, type: 'NEW', data });
        return;
    }

    if (isTaskMoved(new_task, base_task)) {
        diffs.push({ tag: new_task.tag, type: 'MOVE', data });
    } else if (isTaskUpdated(new_task, base_task)) {
        diffs.push({ tag: new_task.tag, type: 'UPDATE', data });
    }
});

baseTasks.forEach(base_task => {
    const new_task = newTasks.find(nt => areTasksEquivalent(nt, base_task));
    if (new_task) return;
    const data = { old: extract_task(base_task), new: null };
    diffs.push({ tag: base_task.tag, type: 'DELETE', data });
});

console.log(JSON.stringify(diffs, null, 2));
