.PHONY: all clean build submit

BIN_NAME=WRECK

build:
	@go build -o $(BIN_NAME) main.go ds.go lltable.go nfa.go
	@chmod +x $(BIN_NAME)

clean:
	@rm -rf $(BIN_NAME)

submit: clean
	mkdir wreck-jamesshima
	cp -r main.go ds.go lltable.go Makefile go.mod README.txt tree-to-graphvis ast.png wreck-jamesshima
	tar - project11-jamesshima.zip project11-jamesshima
	echo "ready for canvas"