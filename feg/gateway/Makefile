.PHONY: all build download fmt test clean run install build_only gen precommit

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

COVER_DIR := $(MAGMA_ROOT)/feg/gateway/coverage
export COVER_DIR


all: fmt test vet install

build: install

download:
	go mod download

install:
	go install ./... magma/gateway/services/magmad

install_envoy_controller:
	go install ./... magma/feg/gateway/services/envoy_controller

test:
	go test ./...
	go test -tags link_local_service,with_builtin_radius magma/feg/gateway/services/eap/... magma/feg/gateway/services/aaa/...

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

precommit: fmt test vet

TOOL_DEPS:= \
	github.com/ory/go-acc \
	github.com/wadey/gocovmerge
tools:: $(TOOL_DEPS)
$(TOOL_DEPS): %:
	go install $*

cover: tools cover_pre

COVER_FILE=$(COVER_DIR)/feg.gocov
cover: tools cover_pre
	go-acc ./... --covermode count --output $(COVER_FILE)
	# skip some unnecessary files for coverage
	awk '!/\.pb\.go|_swaggergen\.go|\/mocks\/|\/tools\/|\/blobstore\/ent\//' $(COVER_FILE) > $(COVER_FILE).tmp && mv $(COVER_FILE).tmp $(COVER_FILE)
cover_pre:
	mkdir -p $(COVER_DIR)
