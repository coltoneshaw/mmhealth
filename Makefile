.PHONY: build
GO_PACKAGES=$(shell go list ./...)
GO ?= $(shell command -v go 2> /dev/null)
BUILD_HASH ?= $(shell git rev-parse HEAD)
# BUILD_VERSION ?= $(shell git ls-remote --tags --refs https://github.com/coltoneshaw/mm-healthcheck.git | tail -n1 | sed 's/.*\///')

LDFLAGS += -X "github.com/coltoneshaw/mm-healthcheck/commands.BuildHash=$(BUILD_HASH)"
# LDFLAGS += -X "github.com/coltoneshaw/mm-healthcheck/commands.Version=$(BUILD_VERSION)"
BUILD_COMMAND ?= go build -ldflags '$(LDFLAGS)'

build: check-style
	$(BUILD_COMMAND)

run:
	go run main.go


package: check-style
	mkdir -p build

	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(BUILD_COMMAND)
	env GZIP=-9 tar czf build/linux_amd64.tar.gz healthcheck

	@echo Build OSX amd64
	env GOOS=darwin GOARCH=amd64 $(BUILD_COMMAND)
	GZIP=-9 tar czf build/darwin_amd64.tar.gz healthcheck

	@echo Build OSX arm64
	env GOOS=darwin GOARCH=arm64 $(BUILD_COMMAND)
	GZIP=-9 tar czf build/darwin_arm64.tar.gz healthcheck

	@echo Build Windows amd64
	env GOOS=windows GOARCH=amd64 $(BUILD_COMMAND)
	zip -9 build/windows_amd64.zip healthcheck.exe

	rm healthcheck healthcheck.exe

golangci-lint:
# https://stackoverflow.com/a/677212/1027058 (check if a command exists or not)
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi; \

	@echo Running golangci-lint
	golangci-lint run --skip-dirs-use-default --timeout 5m -E gofmt ./...

test:
	@echo Running tests
	$(GO) test -race -v $(GO_PACKAGES)

check-style: golangci-lint

verify-gomod:
	$(GO) mod download
	$(GO) mod verify

