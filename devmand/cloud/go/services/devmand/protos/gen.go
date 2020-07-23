/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

//go:generate bash -c "protoc -I . -I /usr/include -I $MAGMA_ROOT/protos --proto_path=$MAGMA_ROOT --go_out=plugins=grpc:. *.proto"
package protos
