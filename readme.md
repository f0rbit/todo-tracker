## Implementation
Arguments: `tracker <directory> <config.json> <existing_data.json>`
Program should traverse the directory using multi-threaded technologies and build a list of all detected 'tasks'. The task identification would be based on the specified `config.json`. For example user could specify multiple different tags for tasks like `todo: '@todo', 'TODO:'` and `bug: '@bug', 'BUG:'`.
It would then compare the produced list of tasks found with an existing data set and produce a diff of the two datasets with intelligent handling. There are 3 main cases to consider
- Line changes of tasks
- If line is same but text has changed of the task
- If existing task was moved to a different file but kept same text

## Languages
In order to develop my own understanding and demonstrate knowledge of various programs, I plan to implement the same program in various different languages. I will blog my experiences with each language and the pros and cons of each language compared to the others, and then at the end I will have broad understanding of the strengths of each language.
### Initial Build
- Bun (JavaScript, for rapid and agile development & prototyping)
- Bun (TypeScript, to get a sense of the types and structure of program)
### Language Development
- Go (make use of GoRoutines for multi-threaded traversal)
- Kotlin (new language)
- C++ (refresh skills in a popular language)
- Rust (in-depth understanding of the language)
### Experimental / Optional Languages
- OCaml
- Elixir
- Common Lisp / Clojure

## Roadmap
This project will be a pivotal part in my devpad rescope and will be integrated in that project as soon as there is a working project. The roadmap for the first build is
1. Traverse Directory (using multi-threaded capabilities)
2. Load config from .json file
3. Produce a .json output of the codebase's tasks
4. Parse existing output for a diff
