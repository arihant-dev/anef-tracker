VERSION ?= v0.9.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u +'%Y-%m-%d')
LDFLAGS = -ldflags "-X github.com/arihant-dev/anef-tracker/internal/version.Version=$(VERSION) -X github.com/arihant-dev/anef-tracker/internal/version.Commit=$(COMMIT) -X github.com/arihant-dev/anef-tracker/internal/version.BuildDate=$(BUILD_DATE)"

.PHONY: all build test race lint coverage coverage-pkg benchmark fuzz sbom docker release clean docs-check install

all: build test

build:
	@echo "Building anef binary ($(VERSION) - $(COMMIT))..."
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/anef ./cmd/anef

test:
	@echo "Running unit tests across all packages..."
	go test ./...

race:
	@echo "Running unit tests with Go race detector..."
	go test -race ./...

examples-test:
	@echo "Exercising all executable documentation examples..."
	go run ./examples/basic
	go run ./examples/watch
	go run ./examples/notifications
	go run ./examples/reports
	go run ./examples/bundle
	@echo "✓ All 5 documentation examples executed cleanly."

lint:
	@echo "Running code vetting & static analysis..."
	go vet ./...

coverage:
	@echo "Generating test coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

coverage-pkg:
	@echo "=== PACKAGE COVERAGE BREAKDOWN ==="
	@go test -cover ./...

fuzz:
	@echo "Running native Go fuzz tests..."
	go test -fuzz=FuzzScanPayload -fuzztime=5s ./pkg/privacy/...
	go test -fuzz=FuzzCompareSnapshots -fuzztime=5s ./pkg/diff/...

sbom:
	@echo "Generating Software Bill of Materials (SBOM)..."
	@syft dir:. -o spdx-json=sbom.spdx.json 2>/dev/null || echo "Note: syft tool not installed. Install syft for automated SBOM generation."

benchmark:
	@echo "Running performance benchmark suite..."
	go test -bench=. ./pkg/benchmark/...

docker:
	@echo "Building Docker container image..."
	docker build -t anef-tracker:$(VERSION) .

release:
	@echo "Building cross-platform release archives (v0.9.0-beta)..."
	@mkdir -p bin/dist bin/build
	GOOS=darwin GOARCH=arm64 go build -trimpath -buildvcs=false $(LDFLAGS) -o bin/build/anef ./cmd/anef && tar -czf bin/dist/anef-tracker_0.9.0-beta_darwin_arm64.tar.gz -C bin/build anef
	GOOS=darwin GOARCH=amd64 go build -trimpath -buildvcs=false $(LDFLAGS) -o bin/build/anef ./cmd/anef && tar -czf bin/dist/anef-tracker_0.9.0-beta_darwin_amd64.tar.gz -C bin/build anef
	GOOS=linux GOARCH=amd64 go build -trimpath -buildvcs=false $(LDFLAGS) -o bin/build/anef ./cmd/anef && tar -czf bin/dist/anef-tracker_0.9.0-beta_linux_amd64.tar.gz -C bin/build anef
	GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false $(LDFLAGS) -o bin/build/anef ./cmd/anef && tar -czf bin/dist/anef-tracker_0.9.0-beta_linux_arm64.tar.gz -C bin/build anef
	GOOS=windows GOARCH=amd64 go build -trimpath -buildvcs=false $(LDFLAGS) -o bin/build/anef.exe ./cmd/anef && zip -j bin/dist/anef-tracker_0.9.0-beta_windows_amd64.zip bin/build/anef.exe
	@rm -rf bin/build
	@cd bin/dist && shasum -a 256 * > checksums.txt
	@echo "✓ Cross-platform release archives and checksums.txt generated in bin/dist/"

docs-check:
	@echo "Validating technical documentation completeness..."
	@test -f docs/architecture.md || exit 1
	@test -f docs/database.md || exit 1
	@test -f docs/evidence-model.md || exit 1
	@test -f docs/api-observability.md || exit 1
	@test -f docs/tui-guide.md || exit 1
	@test -f docs/security.md || exit 1
	@test -f docs/troubleshooting.md || exit 1
	@test -f docs/configuration.md || exit 1
	@test -f docs/compatibility.md || exit 1
	@test -f docs/upgrading.md || exit 1
	@test -f docs/development.md || exit 1
	@test -f docs/release-checklist.md || exit 1
	@test -f mkdocs.yml || exit 1
	@echo "✓ All technical documentation files present."

install: build
	@echo "Installing anef to $(GOPATH)/bin..."
	cp bin/anef $(GOPATH)/bin/anef

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/ exports/ backups/ coverage.out coverage.html sbom.spdx.json
