.PHONY: all clean run

BIN_NAME=LUTHOR

build:
	@go build -o $(BIN_NAME) main.go
	@chmod +x $(BIN_NAME)

run: build
	@./$(BIN_NAME) test.lga

clean:
	@rm -rf $(BIN_NAME)