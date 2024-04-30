.PHONY: all clean build test

BIN_NAME=JackCompiler

build:
	@go build -o $(BIN_NAME) main.go lex.go ds.go codegen.go lltable.go
	@chmod +x $(BIN_NAME)

clean:
	@rm -rf $(BIN_NAME)

test:

submit: clean
	mkdir project11-jamesshima
	cp -r project11-jamesshima main.go ds.go lltable.go lex.go jack.cfg go.mod 