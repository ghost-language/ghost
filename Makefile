GOPATH:=$(shell go env GOPATH)

.PHONY: run build-mac build-linux test clean build

run:
	go run ghost.go

build: build-mac build-linux
	@mkdir -p dist/mac
	@mkdir -p dist/linux

build-mac: clean
	@go build -o dist/mac/ghost *.go

build-linux: clean
	CGO_ENABLED=0 GOOS=linux go build -o dist/linux/ghost *.go

test:
	@go test -timeout 1s -v -race ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean:
	@rm -rf dist/mac/ghost
	@rm -rf dist/linux/ghost