# Makefile для bookshop

APP_NAME=bookshop
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build run test lint mocks migrate deps up swag

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

up:
	docker-compose up -d
	sleep 10
	docker cp keycloak/init-users.sh gh-keycloak-1:/opt/keycloak/init-users.sh
	docker exec gh-keycloak-1 bash /opt/keycloak/init-users.sh 

swag:
	swag init -g cmd/bookshop/main.go -o ./docs 