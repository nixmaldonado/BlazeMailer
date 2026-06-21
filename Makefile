BINARY := blazemailer
PKG := ./...
CMD := ./cmd/blazemailer

.DEFAULT_GOAL := help

## help: show this help
.PHONY: help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | awk -F': ' '{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

## build: compile the binary into ./bin
.PHONY: build
build:
	go build -o bin/$(BINARY) $(CMD)

## run: build and run the server
.PHONY: run
run:
	go run $(CMD)

## fmt: format all Go code
.PHONY: fmt
fmt:
	gofmt -w .

## fmt-check: fail if any file is not gofmt-clean
.PHONY: fmt-check
fmt-check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "These files need gofmt:"; echo "$$unformatted"; exit 1; \
	fi

## vet: run go vet
.PHONY: vet
vet:
	go vet $(PKG)

## test: run all tests
.PHONY: test
test:
	go test $(PKG)

## race: run all tests with the race detector
.PHONY: race
race:
	go test -race $(PKG)

## cover: run tests and open an HTML coverage report
.PHONY: cover
cover:
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out

## lint: run golangci-lint if installed
.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed; skipping (install: https://golangci-lint.run)"; \
	fi

## tidy: sync go.mod/go.sum
.PHONY: tidy
tidy:
	go mod tidy

## check: full local verification (fmt-check, vet, race) — run this before committing
.PHONY: check
check: fmt-check vet race
	@echo "all checks passed"

## clean: remove build and coverage artifacts
.PHONY: clean
clean:
	rm -rf bin coverage.out
