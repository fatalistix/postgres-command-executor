# makefile silent mode
#MAKEFLAGS += --silent

# project setup
BINARY_NAME=postgres-command-executor
MAIN_PACKAGE_PATH=./cmd/postgres-command-executor/main.go

# Go environment variables
# 0 or 1
CGO_ENABLED=0
GOOS=linux

# load environment variables
include .env
export $(shell sed 's/=.*//' .env)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	CGO_ENABLED=${CGO_ENABLED} go build -o ${BINARY_DIR}/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the application
.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}

## migrate/up: run migrations (up)
.PHONY: migrate
migrate:
	migrate -path migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

## env: set application environment variables
.PHONY: env
env:


# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## docker/build: build docker image
.PHONY: docker/build
docker/build:
	docker build -t ${BINARY_NAME}:multistage .

## docker/migrate: run migrations
.PHONY: docker/migrate
docker/migrate: