GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
GIT_HASH ?= $(shell git rev-parse --short HEAD)

VERSION = $(GIT_BRANCH)-$(GIT_HASH)
GO_LDFLAGS = -ldflags="-w -s -X 'main.Version=$(VERSION)'"


BINDIR = $(CURDIR)/dist
LINT_VERSION ?= v2.6.2
LINT_BIN := $(BINDIR)/golangci-lint

APP = garikos

ARCHS = arm64 arm

lint: $(LINT_BIN)
	$(LINT_BIN) run ./...

fix: $(LINT_BIN)
	$(LINT_BIN) run --fix ./...

$(LINT_BIN):
	@echo "Installing golangci-lint $(LINT_VERSION)..."
	@mkdir -p $(BINDIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(BINDIR) $(LINT_VERSION)

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: vuln
vuln: tidy
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: test
test: tidy
	@go test -v -race -count=5 ./...

.PHONY: run
run: tidy
	@go run .

.PHONY: build
build: tidy
	@echo "Building..."
	@CGO_ENABLED=0 go build $(GO_LDFLAGS) -o $(BINDIR)/$(APP) .
	@for arch in $(ARCHS); do \
		echo "Building for $$arch..."; \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$arch go build $(GO_LDFLAGS) -o $(BINDIR)/$(APP)-$$arch .; \
	done
	@ls -lh $(BINDIR)


.PHONY: install
install: build
	scp $(BINDIR)/$(APP)-arm64 xakep@raspberrypi.local:/home/xakep/bin/$(APP)
	scp $(APP).service xakep@raspberrypi.local:/home/xakep/bin/

.PHONY: deploy
deploy: build
	ssh xakep@raspberrypi.local "sudo systemctl stop $(APP)"
	scp $(BINDIR)/$(APP)-arm64 xakep@raspberrypi.local:/home/xakep/bin/$(APP)
	ssh xakep@raspberrypi.local "sudo systemctl start $(APP)"
	ssh xakep@raspberrypi.local "sudo systemctl status $(APP)"