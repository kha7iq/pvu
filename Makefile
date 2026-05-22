BINARY      := kubectl-pvu
MAIN        := ./cmd/pvu
DIST_DIR    := dist
VERSION     ?= dev
COMMIT      ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE        ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = -ldflags "\
  -X github.com/kha7iq/pvu/internal/version.Version=$(VERSION) \
  -X github.com/kha7iq/pvu/internal/version.Commit=$(COMMIT) \
  -X github.com/kha7iq/pvu/internal/version.BuildDate=$(DATE)"

.PHONY: tidy build run tui release clean

tidy:
	go mod tidy

build:
	go build $(LDFLAGS) -o $(BINARY) $(MAIN)

run:
	go run $(MAIN)

tui:
	go run $(MAIN) -n default

release:
	LOCAL_COMMIT=$(COMMIT) goreleaser release --snapshot --clean

clean:
	rm -rf $(BINARY) $(DIST_DIR)