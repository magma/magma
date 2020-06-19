PHONY: all build download fmt test clean install build_only gen precommit

ifndef BUILD_OUT
BUILD_OUT := $(shell go env GOBIN)
ifdef GOOS
BUILD_OUT := $(BUILD_OUT)/$(GOOS)
endif
ifdef GOARCH
BUILD_OUT := $(BUILD_OUT)/$(GOARCH)
endif
endif

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

GW_TOOL_LIST = ./gateway_cli ./service303_cli

all: fmt test vet install

test:
	go test ./...

clean:
	go clean ./...

fmt:
	go fmt ./...

gen:
	go generate ./...

vet:
	go vet ./...

gateway_tools:
	mkdir -p $(GOBIN)
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o $(GOBIN) $(GW_TOOL_LIST)

