.PHONY: all build download install test clean fmt vet build_only precommit

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

# Path to the test files
UNIT_TESTS=services

all: fmt test vet install

build: install

download:
	go mod download

install:
	go install ./...

test:
	go test ./$(UNIT_TESTS)/...

clean:
	go clean ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

build_only:
	go build ./...

precommit: fmt test vet
