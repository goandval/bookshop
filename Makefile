# Makefile для bookshop

APP_NAME=bookshop
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build run test lint mocks migrate deps

all: build

deps:
	go mod tidy

build: deps
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

run: deps
	go run ./cmd/$(APP_NAME)

test: deps
	go test -v -race -cover ./...

lint: deps
	golangci-lint run ./...

mocks: deps
	mockery --all --output ./internal/mocks --case underscore

migrate: deps
	go run ./cmd/migrate 