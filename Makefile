.PHONY: build test lint clean

build:
	go build -o bin/logcheck ./cmd/logcheck/

test:
	go test -v -race ./pkg/rules/...
	go test -v -race ./pkg/analyzer/...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/