GOPATH:=$(shell go env GOPATH)

.PHONY: run build-mac build-linux test clean

run:
	go run ghost.go

build-mac: clean
	@go build -o dist/ghost *.go

build-linux: clean
	CGO_ENABLED=0 GOOS=linux go build -o dist/ghost *.go

test:
	@go test -v -race ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean:
	@rm -rf dist/ghost