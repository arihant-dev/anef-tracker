FROM mcr.microsoft.com/devcontainers/go:1.24-bookworm

RUN apt-get update && apt-get install -y sqlite3 ca-certificates curl git \
    && rm -rf /var/lib/apt/lists/*

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
