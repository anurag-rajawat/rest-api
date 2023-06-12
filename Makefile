##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	@go fmt ./...

.PHONY: test
test: ## Run tests.
	@go test ./...

.PHONY: cover
cover: ## Run tests with coverage.
	@go test ./... -coverprofile=coverage.out -covermode=atomic

.PHONY: docker-run
docker-run: ## Start API server in docker containers.
	@docker compose up -d

##@ Build

.PHONY: build
build: fmt ## Build API binary.
	@go build -o bin/rest-api cmd/main.go

.PHONY: run
run: fmt ## Start API server.
	@go run ./cmd/main.go
