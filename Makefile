# butcher — dev tooling
#
#   make dev      run database, server and client together
#   make server   run the Go API only
#   make client   run the Nuxt client only

SHELL := /bin/bash

SERVER_DIR := server
CLIENT_DIR := client
BIN_DIR    := $(SERVER_DIR)/bin

.DEFAULT_GOAL := help

## help: list the available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | awk -F': ' '{printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

## install: install client dependencies and download Go modules
install:
	cd $(CLIENT_DIR) && npm install
	cd $(SERVER_DIR) && go mod download

## db: start the postgres container in the background
db:
	docker compose up -d db

## db-stop: stop the postgres container
db-stop:
	docker compose stop db

## server: run the Go API (localhost:8080)
server:
	cd $(SERVER_DIR) && go run ./cmd/main.go server

## client: run the Nuxt dev server (localhost:3000)
client:
	cd $(CLIENT_DIR) && npm run dev

## dev: run db, server and client together (ctrl-c stops everything)
dev: db
	@trap 'trap - INT TERM; kill 0' INT TERM; \
	$(MAKE) server & \
	$(MAKE) client & \
	wait

## build: compile the server binary and build the client for production
build:
	cd $(SERVER_DIR) && go build -o bin/server ./cmd
	cd $(CLIENT_DIR) && npm run build

## lint: vet the server and lint the client
lint:
	cd $(SERVER_DIR) && go vet ./...
	cd $(CLIENT_DIR) && npm run lint

## clean: remove build artifacts
clean:
	rm -rf $(BIN_DIR) $(CLIENT_DIR)/.nuxt $(CLIENT_DIR)/.output

.PHONY: help install db db-stop server client dev build lint clean
