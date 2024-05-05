.PHONY: all clean build submit

BIN_NAME=Zobos

build:
	@go build -o $(BIN_NAME) main.go ds.go grammar.go
	@chmod +x $(BIN_NAME)

clean:
	@rm -rf $(BIN_NAME)

submit: clean
	mkdir project11-jamesshima
	cp -r main.go ds.go lltable.go lex.go codegen.go Makefile jack.cfg go.mod LANGINFO README.md tree-to-graphvis viewAST.sh ast.png project11-jamesshima
	zip -r project11-jamesshima.zip project11-jamesshima
	echo "ready for canvas"