//go:generate bash -c "protoc -I /usr/include --proto_path=$MAGMA_ROOT --go_out=plugins=grpc:$MAGMA_ROOT/.. $MAGMA_ROOT/wifi/protos/mconfig/*.proto"
package mconfig
