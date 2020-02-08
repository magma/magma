PLUGIN_NAME=orc8r
TOOL_DEPS:=golang.org/x/lint/golint github.com/golang/protobuf/protoc-gen-go github.com/go-swagger/go-swagger/cmd/swagger magma/orc8r/cloud/go/tools/combine_swagger github.com/vektra/mockery/cmd/mockery magma/orc8r/cloud/go/tools/swaggergen
include $(MAGMA_ROOT)/orc8r/cloud/go/module.mk

# Add targets for Orc8r submodules
clean::
	$(MAKE) -C ../../lib/go/protos clean
	$(MAKE) -C ../../lib/go clean

fmt::
	$(MAKE) -C ../../lib/go/protos fmt
	$(MAKE) -C ../../lib/go fmt

gen::
	$(MAKE) -C ../../lib/go/protos gen
	$(MAKE) -C ../../lib/go gen

test::
	$(MAKE) -C ../../lib/go/protos test
	$(MAKE) -C ../../lib/go test

vet::
	$(MAKE) -C ../../lib/go/protos vet
	$(MAKE) -C ../../lib/go vet
