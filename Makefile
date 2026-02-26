APP := cicost

.PHONY: build test fmt vet run-help

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

