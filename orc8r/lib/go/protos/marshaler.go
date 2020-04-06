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
	"reflect"
	"strings"

	"github.com/golang/glog"
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

func MarshalJSON(msg proto.Message) ([]byte, error) {
	var buff bytes.Buffer
	err := (&jsonpb.Marshaler{Indent: " "}).Marshal(&buff, msg)
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

type mconfigAnyResolver struct{}

// Resolve - AnyResolver interface implementation, it'll resolve any unregistered Any types to Void instead of
// returning error
func (_ mconfigAnyResolver) Resolve(typeUrl string) (proto.Message, error) {
	// Only the part of typeUrl after the last slash is relevant.
	mname := typeUrl
	if slash := strings.LastIndex(mname, "/"); slash >= 0 {
		mname = mname[slash+1:]
	}
	mt := proto.MessageType(mname)
	var res proto.Message
	if mt == nil {
		glog.V(4).Infof("mconfigAnyResolver: unknown message type %q", mname)
		res = new(Void)
	} else {
		res = reflect.New(mt.Elem()).Interface().(proto.Message)
	}
	return res, nil
}

// MarshalMconfig is a special mconfig marshaler tolerant to unregistered Any types
func MarshalMconfig(msg proto.Message) ([]byte, error) {
	var buff bytes.Buffer
	err := (&jsonpb.Marshaler{AnyResolver: mconfigAnyResolver{}, EmitDefaults: true, Indent: " "}).Marshal(&buff, msg)
	return buff.Bytes(), err
}

// UnmarshalMconfig is a special mconfig Unmarshaler tolerant to unregistered Any types
func UnmarshalMconfig(bt []byte, msg proto.Message) error {
	return (&jsonpb.Unmarshaler{AllowUnknownFields: true, AnyResolver: mconfigAnyResolver{}}).Unmarshal(
		bytes.NewBuffer(bt), msg)
}
