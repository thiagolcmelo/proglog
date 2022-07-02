.DEFAULT_GOAL := run

fmt:
	go fmt ./...
.PHONY:fmt

vet: fmt
	go vet ./...
.PHONY:vet

lint: vet
	golint ./...
.PHONY:lint

build: lint
	go build -o main.out cmd/server/main.go
.PHONY:build

run: build
	./main.out
.PHONY:run