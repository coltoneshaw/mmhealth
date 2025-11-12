.PHONY: build check-style test lint fmt verify-gomod install-go-tools
GO_PACKAGES=$(shell go list ./...)
GO ?= $(shell command -v go 2> /dev/null)
export GOBIN ?= $(PWD)/bin
GOLANGCI_LINT_VERSION=latest

BUILD_HASH ?= $(shell git rev-parse HEAD)
BUILD_VERSION ?= $(shell git ls-remote --tags --refs https://github.com/coltoneshaw/mmhealth.git | tail -n1 | sed 's/.*\///')

DOCKER_IMAGE_PROD ?= ghcr.io/coltoneshaw/mmhealth
DOCKER_IMAGE_DEV ?= mmhealth

BUILD_ENV ?= dev

ifeq ($(BUILD_ENV),prod)  
	LDFLAGS += -X "github.com/coltoneshaw/mmhealth/mmhealth.GitCommit=$(BUILD_HASH)"
	LDFLAGS += -X "github.com/coltoneshaw/mmhealth/mmhealth.GitVersion=$(BUILD_VERSION)"
	else 
endif


BUILD_COMMAND ?= go build -ldflags '$(LDFLAGS)' -o ./bin/mmhealth 

build: test
	mkdir -p bin
	$(BUILD_COMMAND)

buildDocker: build
	docker build --platform=linux/amd64 -f ./docker/Dockerfile -t $(DOCKER_IMAGE_DEV) .

run:
	go run ./main.go

package: test
	mkdir -p build bin 

	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(BUILD_COMMAND)
	tar cf - -C bin mmhealth | gzip -9 > build/linux_amd64.tar.gz


	@echo Build OSX amd64
	env GOOS=darwin GOARCH=amd64 $(BUILD_COMMAND)
	tar cf - -C bin mmhealth | gzip -9 > build/darwin_amd64.tar.gz

	@echo Build OSX arm64
	env GOOS=darwin GOARCH=arm64 $(BUILD_COMMAND)
	tar cf - -C bin mmhealth | gzip -9 > build/darwin_arm64.tar.gz

	@echo Build Windows amd64
	env GOOS=windows GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o ./bin/mmhealth.exe 
	zip -9 build/windows_amd64.zip ./bin/mmhealth.exe

	rm ./bin/mmhealth ./bin/mmhealth.exe

check-style: install-go-tools
	@echo Checking gofmt
	@test -z "$$(gofmt -l . | tee /dev/stderr)" || (echo "Files not formatted with gofmt"; exit 1)
	@echo Running golangci-lint
	$(GOBIN)/golangci-lint run --timeout 5m ./...

test: check-style
	@echo Running tests
	$(GO) test -race -cover -v $(GO_PACKAGES)

# Quick lint check (without tests)
lint: install-go-tools
	@echo Running go vet
	$(GO) vet ./...
	@echo Running golangci-lint
	$(GOBIN)/golangci-lint run --timeout 5m ./...

# Format code
fmt: install-go-tools
	@echo Formatting code with gofmt
	$(GO) fmt ./...
	@echo Formatting code with gofumpt
	$(GOBIN)/gofumpt -l -w .

verify-gomod:
	$(GO) mod download
	$(GO) mod verify

# Install go tools locally in bin/
install-go-tools:
	@mkdir -p $(GOBIN)
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION) to $(GOBIN)..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION)
	@echo "Installing gofumpt to $(GOBIN)..."
	@GOBIN=$(GOBIN) $(GO) install mvdan.cc/gofumpt@latest
	@echo "✓ Tools installed successfully!"

make plugins:
	bash ./scripts/update_plugins.sh