.PHONY: all clean 

BIN_NAME=Compiler

build:
	@go build -o $(BIN_NAME) main.go lex.go ds.go
	@chmod +x $(BIN_NAME)

clean:
	@rm -rf $(BIN_NAME)