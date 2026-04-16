BINARY      := express2
CMD         := ./...
GOFLAGS     ?=
COVEROUT    := coverage.out

# Version injection — override from CI or command line.
VERSION ?= v0.0.1-dev
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILT   ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.built=$(BUILT)

# Build the binary for the current platform (static, no CGO).
.PHONY: build
build:
	CGO_ENABLED=0 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

# Build a static Linux/amd64 binary (for deployment).
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

# Build a static Linux/amd64 binary with version ldflags injected (used by CI).
.PHONY: build-release
build-release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

# Run all unit tests (static, no CGO).
.PHONY: test
test:
	CGO_ENABLED=0 go test ./...

# Run unit tests with verbose output (static, no CGO).
.PHONY: test-v
test-v:
	CGO_ENABLED=0 go test -v ./...

# Run unit tests with race detector (CGO required by the race runtime).
.PHONY: test-race
test-race:
	go test -race ./...

# Produce a coverage report (static, no CGO).
.PHONY: cover
cover:
	CGO_ENABLED=0 go test -coverpkg=./... -coverprofile=$(COVEROUT) ./...
	go tool cover -func=$(COVEROUT)

# Open the HTML coverage report in the browser.
.PHONY: cover-html
cover-html: cover
	go tool cover -html=$(COVEROUT) -o coverage.html
	xdg-open coverage.html 2>/dev/null || open coverage.html 2>/dev/null || echo "report: coverage.html"

# Run golangci-lint.
.PHONY: lint
lint:
	golangci-lint run ./...

# Run all static checks (lint + test).
.PHONY: check
check: lint test

# Remove build artifacts.
.PHONY: clean
clean:
	rm -f $(BINARY) $(COVEROUT) coverage.html

.PHONY: help
help:
	@echo "Targets:"
	@echo "  build              Build binary for current platform"
	@echo "  build-linux        Build static Linux/amd64 binary"
	@echo "  build-release      Build static Linux/amd64 binary with version ldflags (CI)"
	@echo "  test               Run unit tests"
	@echo "  test-v             Run unit tests (verbose)"
	@echo "  test-race          Run unit tests with race detector"
	@echo "  cover              Coverage report (func summary)"
	@echo "  cover-html         Coverage report (HTML, opens browser)"
	@echo "  lint               Run golangci-lint"
	@echo "  check              lint + test"
	@echo "  clean              Remove build artifacts"
