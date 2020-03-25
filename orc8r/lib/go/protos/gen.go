/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

//go:generate bash -c "protoc -I /usr/include -I $MAGMA_ROOT/orc8r/protos/prometheus --proto_path=$MAGMA_ROOT --go_out=plugins=grpc,Mgoogle/protobuf/field_mask.proto=google.golang.org/genproto/protobuf/field_mask:$MAGMA_ROOT/.. $MAGMA_ROOT/orc8r/protos/*.proto"
package protos
