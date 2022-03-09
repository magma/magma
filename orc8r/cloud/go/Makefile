MODULE_NAME=orc8r
TOOL_DEPS:= \
	magma/orc8r/cloud/go/tools/swaggergen \
	github.com/go-swagger/go-swagger/cmd/swagger \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/ory/go-acc \
	github.com/vektra/mockery/cmd/mockery \
	github.com/wadey/gocovmerge \
	gotest.tools/gotestsum \
	golang.org/x/tools/cmd/goimports

#	magma/orc8r/cloud/go/tools/swaggergen \
 #	github.com/go-swagger/go-swagger/cmd/swagger@v0.21.0 \
 #	github.com/golang/protobuf/protoc-gen-go@v1.3.3 \
 #	github.com/ory/go-acc@v0.2.6 \
 #	github.com/vektra/mockery/cmd/mockery@v0.0.0-20181123154057-e78b021dcbb5 \
 #	github.com/wadey/gocovmerge@v0.0.0-20160331181800-b5bfa59ec0ad \
 #	gotest.tools/gotestsum@v1.7.0 \
 #	golang.org/x/tools/cmd/goimports@v0.1.9


include $(MAGMA_ROOT)/orc8r/cloud/go/module.mk

LIB_ROOT=$(MAGMA_ROOT)/orc8r/lib/go

# Orc8r module should also forward to Orc8r lib
# Directly mirrors the base Orc8r/cloud/Makefile

build::
	$(MAKE) -C $(LIB_ROOT) build
	$(MAKE) -C $(LIB_ROOT)/protos build

clean::
	$(MAKE) -C $(LIB_ROOT) clean
	$(MAKE) -C $(LIB_ROOT)/protos clean

clean_gen::
	$(MAKE) -C $(LIB_ROOT) clean_gen
	$(MAKE) -C $(LIB_ROOT)/protos clean_gen

download::
	$(MAKE) -C $(LIB_ROOT) download
	$(MAKE) -C $(LIB_ROOT)/protos download

fmt::
	$(MAKE) -C $(LIB_ROOT) fmt
	$(MAKE) -C $(LIB_ROOT)/protos fmt

gen::
	$(MAKE) -C $(LIB_ROOT) gen
	$(MAKE) -C $(LIB_ROOT)/protos gen

gen_protos::
	$(MAKE) -C $(LIB_ROOT) lib_gen_protos

test::
	$(MAKE) -C $(LIB_ROOT) test
	$(MAKE) -C $(LIB_ROOT)/protos test

tidy::
	$(MAKE) -C $(LIB_ROOT) tidy
	$(MAKE) -C $(LIB_ROOT)/protos tidy
