.PHONY: build test clean

all: build
	@./ghost

build: clean
	@go build .

test:
	@go test -v \
		-cover -coverprofile=coverage.txt -covermode=atomic \
		./...

clean:
	@rm -rf ghost