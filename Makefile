.PHONY: all clean 

BIN_NAME=LLONE

build:
	@go build -o $(BIN_NAME) main.go
	@chmod +x $(BIN_NAME)

clean:
	@rm -rf $(BIN_NAME)