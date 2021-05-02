GOPATH:=$(shell go env GOPATH)

.PHONY: run build-mac build-linux build-wasm test clean

run:
	go run ghost.go

build-mac: clean
	GOOS=darwin go build -o dist/mac/ghost ghost.go

build-linux: clean
	GOOS=linux go build -o dist/linux/ghost ghost.go

build-wasm: clean
	GOOS=js GOARCH=wasm go build -o dist/wasm/ghost.wasm wasm/wasm.go

test:
	@go test -v -race ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean:
	@rm -rf dist/mac
	@rm -rf dist/linux
	@rm -rf dist/wasm