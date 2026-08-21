GOPATH:=$(shell go env GOPATH)

# The test target pipes through sed for colour. A pipeline reports the exit
# status of its last command, so without bash and pipefail a failing test run
# is reported as success - which is what CI was doing.
SHELL:=/bin/bash

.PHONY: run build build-mac build-linux build-windows test bench check fmt vet clean

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
	set -o pipefail; go test -v -race -timeout 120s ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

# Benchmark bodies are not compiled into a plain `go test` run, and the Ghost
# programs they hold are only parsed when they execute. Running each one once is
# what catches a language change that leaves the suite unable to parse.
bench:
	go test -run '^$$' -bench=. -benchtime=1x ./...

fmt:
	@test -z "$$(gofmt -l . | tee /dev/stderr)" || (echo "run gofmt -w ." && exit 1)

vet:
	go vet ./...

check: fmt vet test bench

clean:
	@rm -rf dist/mac
	@rm -rf dist/linux
	@rm -rf dist/windows