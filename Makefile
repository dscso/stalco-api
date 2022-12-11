.PHONY: clean build run

APP_NAME = apiserver
BUILD_DIR = $(PWD)/build
GO_BIN_PATH=$(GOPATH)/bin/

clean:
	rm -fr ./build
	rm -fr ./docs

build: swag
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	$(BUILD_DIR)/$(APP_NAME)

swag:
	$(GO_BIN_PATH)swag init -g routes/routes.go --output docs/

install:
	go get
	which $(GO_BIN_PATH)swag || go install github.com/swaggo/swag/cmd/swag@latest
