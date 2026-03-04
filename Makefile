.PHONY: build plugin test lint clean
build:
	go build -o bin/logcheck ./cmd/logcheck/
plugin:
	go build -buildmode=plugin -o bin/logcheck.so ./plugin/
test:
	go test -v -race ./pkg/rules/...
	go test -v -race ./pkg/analyzer/...
lint:
	golangci-lint run ./...
clean:
	rm -rf bin/
