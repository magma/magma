.PHONY: all build download fmt test clean run install build_only gen precommit

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

all: fmt test vet install

build: install

download:
	go mod download

install:
	go install ./...

test:
	go test ./...

buildenv: stop
	PROTO_LIST="orc8r_protos feg_protos lte_protos" $(MAKE) -C $(MAGMA_ROOT)/orc8r/gateway/python $@

run: buildenv build
	sudo service magma@magmad start

restart:
	sudo service magma@* stop
	sudo service magma@magmad start

stop:
	sudo service magma@* stop

clean:
	$(MAKE) -C $(MAGMA_ROOT)/lte/gateway/python $@
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

run_local_hss:
	sudo service magma@hss start

precommit: fmt gen build_only test vet
