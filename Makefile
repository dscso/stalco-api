.PHONY: clean build run

APP_NAME = apiserver
BUILD_DIR = $(PWD)/build
GO_BIN_PATH=$(GOPATH)/bin/

clean:
	rm -rf ./build

build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

swag:
	$(GO_BIN_PATH)swag init -g routes/routes.go --output docs/

install:
	which $(GO_BIN_PATH)swag || go install github.com/swaggo/swag/cmd/swag@latest