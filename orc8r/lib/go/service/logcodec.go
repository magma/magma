/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package service

import (
	"bytes"
	"fmt"
	"log"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
	grpc_proto "google.golang.org/grpc/encoding/proto"
)

// logCodec is a debugging Codec implementation for protobuf.
// It'll be used if debug GRPC printout is enabled
type logCodec struct {
	protoCodec encoding.Codec
}

func printMessage(prefix string, v interface{}) {
	var payload string
	if pm, ok := v.(proto.Message); ok {
		var buff bytes.Buffer
		err := (&jsonpb.Marshaler{EmitDefaults: true, Indent: "\t", OrigName: true}).Marshal(&buff, pm)
		if err == nil {
			payload = string(buff.Bytes())
		} else {
			payload = fmt.Sprintf("\n\t JSON encoding error: %v; %s", err, string(buff.Bytes()))
		}
	} else {
		payload = fmt.Sprintf("\n\t %T is not proto.Message; %+v", v, v)
	}
	log.Printf("%s%T: %s", prefix, v, payload)
}

// Marshal of GRPC Codec interface
func (lc logCodec) Marshal(v interface{}) ([]byte, error) {
	printMessage("Sending: ", v)
	return lc.protoCodec.Marshal(v)
}

// Unmarshal of GRPC Codec interface
func (lc logCodec) Unmarshal(data []byte, v interface{}) error {
	err := lc.protoCodec.Unmarshal(data, v)
	printMessage("Received: ", v)
	return err
}

// Name of GRPC Codec interface
func (logCodec) Name() string {
	return grpc_proto.Name
}
