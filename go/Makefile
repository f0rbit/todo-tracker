RESOURCES_DIR=/home/tom/dev/github/todo-tracker/resources

all: clean build parse

clean:
	rm -f todo-tracker

build:
	go build -o todo-tracker

parse: build
	@./todo-tracker parse ${RESOURCES_DIR}/codebase ${RESOURCES_DIR}/config.json
