//go:generate bash -c "protoc -I . -I /usr/include -I $MAGMA_ROOT/protos --proto_path=$MAGMA_ROOT --go_out=plugins=grpc:. *.proto"
package protos
