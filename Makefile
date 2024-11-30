.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	CGO_ENABLED=0​​ go build -o bin/chatapp main.go

clean:
	go clean