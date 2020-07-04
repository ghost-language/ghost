GOPATH:=$(shell go env GOPATH)

.PHONY: build test clean

build: build-mac
	@./ghost

build-mac: clean
	@go build -o ghost *.go

test:
	@go test -v -race ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean:
	@rm -rf ghost