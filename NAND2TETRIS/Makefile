.PHONY: all clean build submit

BIN_NAME=JackCompiler

build:
	@go build -o $(BIN_NAME) main.go lex.go ds.go codegen.go lltable.go
	@chmod +x $(BIN_NAME)

clean:
	@rm -rf $(BIN_NAME)

submit: clean
	mkdir project11-jamesshima
	cp -r project11-jamesshima main.go ds.go lltable.go lex.go jack.cfg go.mod codegen.go LANGINFO README.md tree-to-graphvis viewAST.sh