ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Makefile for timelog project

# Default values
env ?= prod
BIN_NAME := main
BIN_NAME_LINUX := main.linux

ifeq ($(env),prod)
	PORT := 8083
	DOCKER_TAG := timelog-app
	DOCKER_PORT := 8083
else ifeq ($(env),dev)
	PORT := 18083
	DOCKER_TAG := timelog-app-dev
	DOCKER_PORT := 8083
else
	$(error Unknown env: $(env))
endif

# Set DB file for migrate target
ifeq ($(env),prod)
	MIGRATE_DB_FILE := prod.db
else
	MIGRATE_DB_FILE := dev.db
endif

.PHONY: all build build-linux buildx buildx-linux docker run clean web mcp passkey-temp migrate fmt install-deps gen-model gen-api check-api test

all: build

build: gen-api
	go build -trimpath -o $(BIN_NAME)

build-linux: gen-api
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -tags prod -o $(BIN_NAME_LINUX)

buildx: web build

buildx-linux: web build-linux

docker: build-linux
	docker build -t $(DOCKER_TAG) .

run: docker
	docker run --rm -e ENV=$(env) -p $(PORT):$(DOCKER_PORT) $(DOCKER_TAG)

clean:
	rm -f $(BIN_NAME) $(BIN_NAME_LINUX) mcp/timelog-mcp-server
	rm -rf web/dist web/node_modules

# Web frontend targets
web:
	cd web && pnpm install
	$(MAKE) gen-api
	cd web && pnpm run build

# MCP Server target
mcp:
	cd mcp && go build -o timelog-mcp-server .

# Passkey temp password management
# Usage:
#   make passkey-temp                    # Create with default TTL (from config.yml)
#   make passkey-temp PASSKEY_TEMP_TTL=900  # Create with custom TTL (seconds)
PASSKEY_TEMP_TTL ?=
passkey-temp:
	@go run ./cmd/passkey-temp-admin $(PASSKEY_TEMP_TTL)

# Migrate target
# Usage:
#   make migrate                                 # Postgres (default), reads MIGRATE_DSN from .env
#   make migrate DATABASE_TYPE=sqlite MIGRATE_DB_FILE=dev.db
# MIGRATE_DSN is intentionally not hardcoded here to avoid leaking credentials.
# Put it in .env (gitignored), e.g.:
#   MIGRATE_DSN=postgres://timelog:timelog@localhost:5432/timelog?sslmode=disable
DATABASE_TYPE ?= postgres

migrate:
ifeq ($(DATABASE_TYPE),sqlite)
	migrate -database "sqlite3://$(MIGRATE_DB_FILE)" --path model/migrations/sqlite up
else ifeq ($(DATABASE_TYPE),postgres)
	@test -n "$(MIGRATE_DSN)" || (echo "MIGRATE_DSN is not set. Add it to .env, e.g.:"; \
		echo "  MIGRATE_DSN=postgres://user:pass@localhost:5432/timelog?sslmode=disable"; exit 1)
	migrate -database "$(MIGRATE_DSN)" --path model/migrations/postgres up
else
	$(error Unsupported DATABASE_TYPE: $(DATABASE_TYPE))
endif

# fmt
fmt:
	go fmt ./...
	cd mcp && go fmt ./...
	cd web && npx prettier --write src/ || true

install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
	go install github.com/chrusty/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema@1.4.1
	go install gorm.io/gen/tools/gentool@latest
	go install -tags 'sqlite3 postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	cd web && pnpm install

# Swagger docs target removed; API docs are generated via make gen-api (OpenAPI + Redoc).

gen-model:
	@go run model/gentool/gormgen.go
	@rm -f model/gen/schema_migrations.gen.go model/gen/sqlite_sequence.gen.go 2>/dev/null || true
	@echo "Models generated. Check compilation with: go build ./model/..."

gen-api:
	@test -x ./web/node_modules/.bin/buf || (cd web && pnpm install)
	@PATH="$(or $(shell go env GOBIN),$(shell go env GOPATH)/bin):./web/node_modules/.bin:$(PATH)" ./web/node_modules/.bin/buf generate
	go run ./cmd/merge-openapi

check-api:
	./web/node_modules/.bin/buf lint

test: gen-api
	go test ./...

cover: gen-api
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

cover-html: cover
	go tool cover -html=coverage.out -o coverage.html
	@echo "打开 coverage.html 可浏览详细覆盖率"
