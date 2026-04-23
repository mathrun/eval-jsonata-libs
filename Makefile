BIN := bin/eval
TESTDATA ?= ./testdata

.PHONY: all build run clean

all: build

build:
	go build -o $(BIN) ./cmd

run: build
	./$(BIN) $(TESTDATA)

clean:
	rm -rf bin results
