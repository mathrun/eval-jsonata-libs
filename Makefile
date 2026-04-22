.PHONY: all build run clean deps

all: run

deps:
	go mod tidy
	go mod download

build: deps
	go build -o bin/jsonata-eval .

run: clean build
	./bin/jsonata-eval testdata/test_01.json

clean:
	rm -rf bin/