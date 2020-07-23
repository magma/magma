/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

//go:generate bash -c "protoc -I /usr/include -I $MAGMA_ROOT --proto_path=../../../../protos/mconfig --go_out=plugins=grpc:../../../../../.. ../../../../protos/mconfig/*.proto"
package mconfig
