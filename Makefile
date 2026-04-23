BIN := bin/eval

.PHONY: all build run clean

all: build

build:
	go build -o $(BIN) ./cmd

run: build
	./$(BIN) ./testdata

clean:
	rm -rf bin results
