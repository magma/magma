//go:generate bash -c "protoc -I /usr/include -I $MAGMA_ROOT -I $MAGMA_ROOT/orc8r/protos/prometheus -I $MAGMA_ROOT --go_out=plugins=grpc,Mgoogle/protobuf/field_mask.proto=google.golang.org/genproto/protobuf/field_mask:$MAGMA_ROOT/.. $MAGMA_ROOT/fbinternal/protos/*.proto"
package protos
