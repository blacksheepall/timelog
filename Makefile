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

.PHONY: all build build-lite build-linux build-linux-lite buildx buildx-linux docker run clean web mcp migrate fmt install-deps swagger gen-model build-temp-password build-temp-password-lite

all: build

build:
	go build -trimpath -o $(BIN_NAME)

build-lite:
	go build -trimpath -ldflags="-s -w" -o $(BIN_NAME)

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -tags prod -o $(BIN_NAME_LINUX)

# recude ≈22% size with -ldflags
build-linux-lite:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -tags prod -ldflags="-s -w" -o $(BIN_NAME_LINUX)

buildx: web build

buildx-linux: web build-linux

docker: build-linux
	docker build -t $(DOCKER_TAG) .

run: docker
	docker run --rm -e ENV=$(env) -p $(PORT):$(DOCKER_PORT) $(DOCKER_TAG)

clean:
	rm -f $(BIN_NAME) $(BIN_NAME_LINUX) bin/temp-password mcp/timelog-mcp-server
	rm -rf web/dist web/node_modules

# Web frontend targets
web:
	cd web && pnpm install && pnpm run build

# MCP Server target
mcp:
	cd mcp && go build -o timelog-mcp-server .

build-temp-password:
	go build -trimpath -o bin/temp-password ./cmd/temp-password

build-temp-password-lite:
	go build -trimpath -ldflags="-s -w" -o bin/temp-password ./cmd/temp-password

# Migrate target
migrate:
	migrate -database "sqlite3://$(MIGRATE_DB_FILE)" --path model/migrations/ up

# fmt
fmt:
	go fmt ./...
	cd mcp && go fmt ./...
	cd web && npx prettier --write src/ || true

install-deps:
	go install gorm.io/gen/tools/gentool@latest

# Swagger docs target
swagger:
	swag init

gen-model:
	@go run model/gentool/gormgen.go
	@rm -f model/gen/schema_migrations.gen.go model/gen/sqlite_sequence.gen.go 2>/dev/null || true
	@sed -i '' 's/gorm\.DeletedAt `gorm:"column:deleted_at;type:DATETIME" json:"deleted_at"`/gorm.DeletedAt `gorm:"column:deleted_at;type:DATETIME" json:"deleted_at" swaggertype:"string"`/g' model/gen/*.go 2>/dev/null || true
	@echo "Models generated. Check compilation with: go build ./model/..."
