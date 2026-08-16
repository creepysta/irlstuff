.ONESHELL:

.PHONY: run debug build build-arm64

clean:
	@rm -rf dist/

run:
	@go run main.go

debug:
	@go run main.go -debug

build:
	@go build -o dist/smoke-test

build-arm64:
	@GOOS=linux GOARCH=arm64 go build -o dist/echo-server-linux-arm64
