APP := cicost

.PHONY: build test fmt vet run-help release-check build-gh

build:
	go build -o bin/$(APP) .

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

run-help:
	go run . help

build-gh:
	go build -o bin/gh-cicost ./cmd/gh-cicost

release-check:
	goreleaser check
