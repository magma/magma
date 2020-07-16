//go:generate bash -c "protoc -I /usr/include -I $MAGMA_ROOT --proto_path=../../../../protos/mconfig --go_out=plugins=grpc:../../../../../.. ../../../../protos/mconfig/*.proto"
package mconfig
