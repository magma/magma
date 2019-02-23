/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

//go:generate bash -c "protoc -I . -I /usr/include -I $MAGMA_ROOT/protos --proto_path=$MAGMA_ROOT --go_out=plugins=grpc:. *.proto"
package test_protos

import (
	"testing"

	"magma/orc8r/cloud/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func GetDefaultMconfig(t *testing.T) (*protos.GatewayConfigs, []byte) {
	return GetMconfig(
		t,
		map[string]proto.Message{
			"config1": &Config1{Field: "field"},
			"config2": &Config2{Field1: "field1", Field2: "field2"},
		},
	)
}

func GetMconfig(t *testing.T, configs map[string]proto.Message) (*protos.GatewayConfigs, []byte) {
	anys := map[string]*any.Any{}
	for k, v := range configs {
		anyV, err := ptypes.MarshalAny(v)
		assert.NoError(t, err)
		anys[k] = anyV
	}

	mcfg := &protos.GatewayConfigs{ConfigsByKey: anys}
	marshal, err := protos.MarshalIntern(mcfg)
	assert.NoError(t, err)
	return mcfg, marshal
}
