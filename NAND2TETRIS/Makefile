.PHONY: all clean build test

BIN_NAME=Compiler

build:
	@go build -o $(BIN_NAME) main.go lex.go ds.go codegen.go lltable.go
	@chmod +x $(BIN_NAME)

clean:
	@rm -rf $(BIN_NAME)