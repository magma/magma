PLUGIN_NAME=orc8r
TOOL_DEPS:=golang.org/x/lint/golint github.com/golang/protobuf/protoc-gen-go github.com/go-swagger/go-swagger/cmd/swagger magma/orc8r/cloud/go/tools/combine_swagger github.com/vektra/mockery/cmd/mockery magma/orc8r/cloud/go/tools/swaggergen
include $(MAGMA_ROOT)/orc8r/cloud/go/module.mk

LIB_ROOT=$(MAGMA_ROOT)/orc8r/lib/go

# Add targets for Orc8r submodules
clean::
	$(MAKE) -C $(LIB_ROOT)/protos clean
	$(MAKE) -C $(LIB_ROOT) clean

fmt::
	$(MAKE) -C $(LIB_ROOT)/protos fmt
	$(MAKE) -C $(LIB_ROOT) fmt

gen::
	$(MAKE) -C $(LIB_ROOT)/protos gen
	$(MAKE) -C $(LIB_ROOT) gen

test::
	$(MAKE) -C $(LIB_ROOT)/protos test
	$(MAKE) -C $(LIB_ROOT) test

vet::
	$(MAKE) -C $(LIB_ROOT)/protos vet
	$(MAKE) -C $(LIB_ROOT) vet
