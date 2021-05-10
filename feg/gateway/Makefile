.PHONY: all build download fmt test clean run install build_only gen precommit

ifndef MAGMA_ROOT
MAGMA_ROOT = /home/$USER/magma
endif
export MAGMA_ROOT

COVER_DIR:=$(MAGMA_ROOT)/feg/gateway/coverage
export COVER_DIR
COVER_FILE:=$(COVER_DIR)/feg.gocov
DOCKER_DIR := ${MAGMA_ROOT}/feg/gateway/docker

all: fmt lint test vet install

build: install

download:
	go mod download

install:
	go install ./... magma/gateway/services/magmad

install_envoy_controller:
	go install ./... magma/feg/gateway/services/envoy_controller

test:
	mkdir -p /tmp/test-results
	gotestsum --junitfile /tmp/test-results/main.xml ./...
	gotestsum --junitfile /tmp/test-results/aaa-eap.xml --packages  magma/feg/gateway/services/eap/... magma/feg/gateway/services/aaa/... -- -tags link_local_service,with_builtin_radius
	go test -tags cli_test magma/feg/gateway/tools/...

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
	gofmt -l -s -w .

gen:
	go generate ./...

vet:
	go vet ./...

lint: lint_tools
	golangci-lint run

lint_tools:
	@if which golangci-lint 2>/dev/null; then \
		echo "-> golangci-lint already installed"; \
	else \
		echo "-> installing golangci " && \
			curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
			| sudo sh -s -- -b /usr/sbin/ v1.35.2; \
	fi

build_only:
	go build ./...

run_local_hss:
	sudo service magma@hss start

precommit: tools fmt lint test vet

cover: tools
	mkdir -p $(COVER_DIR)
	go-acc ./... --covermode count --output $(COVER_FILE)
	# skip some unnecessary files for coverage
	awk '!/\.pb\.go|_swaggergen\.go|\/mocks\/|\/tools\/|\/blobstore\/ent\//' $(COVER_FILE) > $(COVER_FILE).tmp && \
		mv $(COVER_FILE).tmp $(COVER_FILE)

docker_build:
	cd ${DOCKER_DIR}; docker-compose build --parallel

docker_start:
	cd ${DOCKER_DIR}; docker-compose up -d

docker_stop:
	cd  ${DOCKER_DIR}; docker-compose down

docker_healthcheck: docker_start
	cd ${DOCKER_DIR}/tools/; ./docker_ps_healthcheck.sh ${DOCKER_DIR}

# Tool dependencies
TOOL_DEPS:= \
	github.com/ory/go-acc \
	github.com/wadey/gocovmerge \
	gotest.tools/gotestsum

tools:: $(TOOL_DEPS)
$(TOOL_DEPS): %:
	go install $*


