/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package protos is protoc generated GRPC package
package protos

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func Marshal(msg proto.Message) ([]byte, error) {
	var buff bytes.Buffer
	err := (&jsonpb.Marshaler{
		EnumsAsInts:  true,
		EmitDefaults: true,
		OrigName:     true}).Marshal(&buff, msg)

	return buff.Bytes(), err
}

func MarshalIntern(msg proto.Message) ([]byte, error) {
	var buff bytes.Buffer
	err := (&jsonpb.Marshaler{EmitDefaults: true, Indent: " "}).Marshal(
		&buff, msg)

	return buff.Bytes(), err
}

func Unmarshal(bt []byte, msg proto.Message) error {
	return (&jsonpb.Unmarshaler{AllowUnknownFields: true}).Unmarshal(
		bytes.NewBuffer(bt),
		msg)
}

func TestMarshal(msg proto.Message) string {
	res, _ := Marshal(msg)
	return string(res)
}
