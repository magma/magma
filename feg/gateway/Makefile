.PHONY: all build fmt test clean run install build_only gen precommit

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

all: fmt test vet install

build: install

install:
	go install ./...

test:
	go test ./...

buildenv: stop
	PROTO_LIST="orc8r_protos feg_protos lte_protos" $(MAKE) -C ../../orc8r/gateway/python $@

run: buildenv build
	sudo service magma@magmad start

restart:
	sudo service magma@* stop
	sudo service magma@magmad start

stop:
	sudo service magma@* stop

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
