.PHONY: all build test clean run-server run-agent

VERSION ?= 2.3.0
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")

LDFLAGS = -X github.com/Kentralo/kenpanel/core/config.Version=$(VERSION) \
          -X github.com/Kentralo/kenpanel/core/config.GitCommit=$(GIT_COMMIT) \
          -X github.com/Kentralo/kenpanel/core/config.BuildDate=$(BUILD_DATE)

all: build

build: build-server build-agent build-cli build-standalone

build-server:
	@echo "Building KenPanel Control Plane..."
	go build -ldflags "$(LDFLAGS)" -o bin/kenpanel ./cmd/kenpanel

build-agent:
	@echo "Building KenPanel Node Agent..."
	go build -ldflags "$(LDFLAGS)" -o bin/kenpanel-agent ./cmd/kenpanel-agent

build-cli:
	@echo "Building KenPanel CLI..."
	go build -ldflags "$(LDFLAGS)" -o bin/kenpanel-cli ./cmd/kenpanel-cli

build-standalone:
	@echo "Building Standalone Engines..."
	go build -ldflags "$(LDFLAGS)" -o bin/kenpanel-transplant ./cmd/kenpanel-transplant
	go build -ldflags "$(LDFLAGS)" -o bin/kenpanel-mailbridge ./cmd/kenpanel-mailbridge
	go build -ldflags "$(LDFLAGS)" -o bin/kenpanel-capsule ./cmd/kenpanel-capsule
	go build -ldflags "$(LDFLAGS)" -o bin/kenpanel-recovery ./cmd/kenpanel-recovery

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
