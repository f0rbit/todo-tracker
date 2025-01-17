import fs from 'node:fs';
import { DiffResult, ParsedTask, TASK_FILE } from './schema';

// Function to read and parse a JSON file
function read(path: string) {
    const data = fs.readFileSync(path, 'utf8');
    const result = TASK_FILE.safeParse(JSON.parse(data));
    if (!result.success) throw new Error('Invalid JSON');
    return result.data;
}



// Read the base and new task data
const base_tasks = read('./output-base.json');
const new_tasks = read('./output-new.json');

const diffs: DiffResult[] = [];

const extract_task = (task: ParsedTask) => {
    const { text, line, file } = task;
    return { text, line, file };
}

const same_text = (task1: ParsedTask, task2: ParsedTask) => {
    // return true if the tasks have the same text
    return task1.text.trim() === task2.text.trim();
}

// process the tasks
new_tasks.forEach((new_task: ParsedTask) => {
    // first we want to check if we can find a task in the base tasks that has the same text but different line number or file
    const same_text_result = base_tasks.find(bt => same_text(bt, new_task));
    if (same_text_result) {
        if (same_text_result.line !== new_task.line || same_text_result.file !== new_task.file) {
            const data = { old: extract_task(same_text_result), new: extract_task(new_task) };
            diffs.push({ id: same_text_result.id, tag: new_task.tag, type: 'MOVE', data });
        } else {
            const data = { old: extract_task(same_text_result), new: extract_task(new_task) };
            diffs.push({ id: same_text_result.id, tag: new_task.tag, type: 'SAME', data });
        }
        return;
    }
    // then if the line number and tag are the same but the text is different
    const same_line_result = base_tasks.find(bt => bt.line === new_task.line && bt.tag === new_task.tag);
    if (same_line_result && !same_text(same_line_result, new_task)) {
        const data = { old: extract_task(same_line_result), new: extract_task(new_task) };
        diffs.push({ id: same_line_result.id, tag: new_task.tag, type: 'UPDATE', data });
        return;
    }

    // if we get here then we have a new task
    const data = { old: null, new: extract_task(new_task) };
    diffs.push({ id: new_task.id, tag: new_task.tag, type: 'NEW', data });    
});

base_tasks.forEach(base_task => {
    // check if the task in the new tasks, and if not add a 'DELETE' diff
    const in_diffs = diffs.find(d => d.id === base_task.id);
    if (!in_diffs) {
        const data = { old: extract_task(base_task), new: null };
        diffs.push({ id: base_task.id, tag: base_task.tag, type: 'DELETE', data });
    }
});

console.log(JSON.stringify(diffs, null, 2));
