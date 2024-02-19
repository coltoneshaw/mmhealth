.PHONY: build
GO_PACKAGES=$(shell go list ./...)
GO ?= $(shell command -v go 2> /dev/null)
BUILD_HASH ?= $(shell git rev-parse HEAD)
BUILD_VERSION ?= $(shell git ls-remote --tags --refs https://github.com/coltoneshaw/mm-healthcheck.git | tail -n1 | sed 's/.*\///')

LDFLAGS += -X "github.com/coltoneshaw/mm-healthcheck/commands.BuildHash=$(BUILD_HASH)"
LDFLAGS += -X "github.com/coltoneshaw/mm-healthcheck/commands.Version=$(BUILD_VERSION)"
BUILD_COMMAND ?= go build -ldflags '$(LDFLAGS)' -o mmhealth ./healthcli/main.go

build: check-style
	$(BUILD_COMMAND)

buildDocker: build
	docker build -f ./docker/dockerfile -t mm-healthcheck . 

run:
	go run ./healthcli/main.go


package: check-style
	mkdir -p build

	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(BUILD_COMMAND)
	tar cf - mmhealth | gzip -9 > build/linux_amd64.tar.gz


	@echo Build OSX amd64
	env GOOS=darwin GOARCH=amd64 $(BUILD_COMMAND)
	tar cf - mmhealth | gzip -9 > build/darwin_amd64.tar.gz

	@echo Build OSX arm64
	env GOOS=darwin GOARCH=arm64 $(BUILD_COMMAND)
	tar cf - mmhealth | gzip -9 > build/darwin_arm64.tar.gz

	@echo Build Windows amd64
	env GOOS=windows GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o mmhealth.exe ./healthcli/main.go
	zip -9 build/windows_amd64.zip mmhealth.exe

	rm mmhealth mmhealth.exe

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

