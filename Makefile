.PHONY: clean fromat run build
APP_NAME = govas
BUILD_DIR = $(PWD)/build

clean:
	@rm -rf $(BUILD_DIR)
	
format:
	@go fmt ./...

run: 
	@go run .

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go