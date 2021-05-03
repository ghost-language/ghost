GOPATH:=$(shell go env GOPATH)

.PHONY: run build build-mac build-linux build-wasm build-windows test clean

run:
	go run ghost.go

build: build-mac build-linux build-wasm build-windows

build-mac: clean
	GOOS=darwin go build -o dist/mac/ghost ghost.go

build-linux: clean
	GOOS=linux go build -o dist/linux/ghost ghost.go

build-wasm: clean
	GOOS=js GOARCH=wasm go build -o dist/wasm/ghost.wasm wasm/wasm.go

build-windows: clean
	GOOS=windows go build -o dist/windows/ghost.exe ghost.go

test:
	@go test -v -race `go list ./... | grep -v "/wasm"` | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean:
	@rm -rf dist/mac
	@rm -rf dist/linux
	@rm -rf dist/wasm
	@rm -rf dist/windows