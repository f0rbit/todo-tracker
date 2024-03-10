import fs from 'node:fs';

// Function to read and parse a JSON file
function readJsonFile(filePath) {
    const data = fs.readFileSync(filePath, 'utf8');
    return JSON.parse(data);
}

// Determine if two tasks are equivalent
function areTasksEquivalent(task1, task2) {
    return task1.text.trim() === task2.text.trim() && task1.tag === task2.tag;
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

const diffs = [];

// Track moved and updated tasks
newTasks.forEach(newTask => {
    const baseTask = baseTasks.find(bt => areTasksEquivalent(bt, newTask));
    if (baseTask) {
        if (isTaskMoved(newTask, baseTask)) {
            diffs.push({ ...newTask, type: 'MOVED' });
        } else if (isTaskUpdated(newTask, baseTask)) {
            diffs.push({ ...newTask, type: 'UPDATED' });
        }
    } else {
        diffs.push({ ...newTask, type: 'NEW' }); // Task does not exist in baseTasks
    }
});

// Track deleted tasks
baseTasks.forEach(baseTask => {
    const newTask = newTasks.find(nt => areTasksEquivalent(nt, baseTask));
    if (!newTask) {
        diffs.push({ ...baseTask, type: 'DELETE' }); // Task does not exist in newTasks
    }
});

console.log(JSON.stringify(diffs, null, 2));
