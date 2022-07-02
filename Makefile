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

test: build
	go test ./...
.PHONY:test

cover: test
	go test -v -cover -coverprofile=c.out ./...
	go tool cover -html=c.out
.PHONY:cover

compile-protobuf:
	protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.
.PHONY:compile-protobuf