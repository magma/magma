MODULE_NAME=orc8r
PLUGIN_NAME=orc8r
TOOL_DEPS:= \
	magma/orc8r/cloud/go/tools/swaggergen \
	github.com/facebookincubator/ent/cmd/entc \
	github.com/go-swagger/go-swagger/cmd/swagger \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/ory/go-acc \
	github.com/vektra/mockery/cmd/mockery \
	github.com/wadey/gocovmerge \
	gotest.tools/gotestsum
include $(MAGMA_ROOT)/orc8r/cloud/go/module.mk

LIB_ROOT=$(MAGMA_ROOT)/orc8r/lib/go

# Add targets for orc8r submodules
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
