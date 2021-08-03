MODULE_NAME=orc8r
include $(MAGMA_ROOT)/orc8r/cloud/go/module.mk

lib_gen_protos:
	cd $(MAGMA_ROOT) ; \
	for x in $$(find $(MAGMA_ROOT)/orc8r/protos $(MAGMA_ROOT)/orc8r/protos/mconfig -maxdepth 1 -name '*.proto') ; do \
		protoc \
			--proto_path $(MAGMA_ROOT) \
			--proto_path $(MAGMA_ROOT)/orc8r/protos/prometheus \
			--proto_path $(PROTO_INCLUDES) \
			--go_out=plugins=grpc,Mgoogle/protobuf/field_mask.proto=google.golang.org/genproto/protobuf/field_mask:$(MAGMA_ROOT)/.. \
			$${x} ; \
	done
