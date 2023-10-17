GOPATH:=$(shell go env GOPATH)

.PHONY: run build build-mac build-linux build-windows test clean

run:
	go run cmd/*.go

build: build-mac build-linux build-windows

build-mac: clean
	GOOS=darwin go build -trimpath -o ./dist/mac/ghost cmd/*.go

build-linux: clean
	GOOS=linux go build -trimpath -o ./dist/linux/ghost cmd/*.go

build-windows: clean
	GOOS=windows go build -trimpath -o ./dist/windows/ghost.exe cmd/*.go

test:
	go test -v -race -timeout 5s ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean:
	@rm -rf dist/mac
	@rm -rf dist/linux
	@rm -rf dist/windows