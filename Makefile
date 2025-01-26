RESOURCES_DIR=./resources
OUTPUT_NAME=todo-tracker

all: clean build parse

clean:
	rm -f ${OUTPUT_NAME}
	rm -rf output*.json

build:
	go build -o ${OUTPUT_NAME}

parse: build
	@./${OUTPUT_NAME} parse ${RESOURCES_DIR}/codebase ${RESOURCES_DIR}/config.json

parse-new: build
	@./${OUTPUT_NAME} parse ${RESOURCES_DIR}/codebase-changed ${RESOURCES_DIR}/config.json

diff: build
	@./${OUTPUT_NAME} parse ${RESOURCES_DIR}/codebase ${RESOURCES_DIR}/config.json > output-base.json
	@./${OUTPUT_NAME} parse ${RESOURCES_DIR}/codebase-changed ${RESOURCES_DIR}/config.json > output-new.json
	@./${OUTPUT_NAME} diff output-base.json output-new.json
