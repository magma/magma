PHONY: all build download fmt test clean install build_only gen precommit

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT


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

MAGMA_BUILD_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
MAGMA_BUILD_TAG ?= $(shell git describe --tags --abbrev=0)
MAGMA_BUILD_COMMIT_HASH ?= $(shell git rev-parse HEAD)
MAGMA_BUILD_COMMIT_DATE ?= $(shell git log -1 --format=%cd)

LDFLAGS := -ldflags="-X 'magma/orc8r/lib/go/build_info.buildBranch=$(MAGMA_BUILD_BRANCH)' \
	-X 'magma/orc8r/lib/go/build_info.buildCommitHash=$(MAGMA_BUILD_COMMIT_HASH)' \
	-X 'magma/orc8r/lib/go/build_info.buildCommitDate=$(MAGMA_BUILD_COMMIT_DATE)' \
	-X 'magma/orc8r/lib/go/build_info.buildTag=$(MAGMA_BUILD_TAG)' \
	-X 'magma/orc8r/lib/go/build_info.buildDate=$(shell date)'"

all: fmt test vet install

build: install

download:
	go mod download

install:
	mkdir -p $(GOBIN)
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o $(BUILD_OUT) $(GOFLAGS) $(LDFLAGS) magma/gateway/services/magmad
	$(MAKE) -C $(MAGMA_ROOT)/orc8r/cloud/go/tools GOARCH=$(GOARCH) GOOS=$(GOOS) GOFLAGS=$(GOFLAGS) BUILD_OUT=$(BUILD_OUT) gateway_tools

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

lint:
	golint -min_confidence 1. ./...

build_only:
	go build ./...

precommit: fmt build_only test vet

cover:
	go test ./... -coverprofile ./cover.tmp >/dev/null
	go tool cover -func=./cover.tmp | tail -n 1
	find . -name '*.go' | xargs wc -l | tail -n 1
	rm ./cover.tmp

