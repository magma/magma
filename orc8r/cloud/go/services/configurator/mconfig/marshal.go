/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package mconfig

import (
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

func MarshalConfigs(configs map[string]proto.Message) (ConfigsByKey, error) {
	ret := ConfigsByKey{}
	for k, v := range configs {
		anyVal, err := ptypes.MarshalAny(v)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		bytesVal, err := protos.MarshalJSON(anyVal)
		if err != nil {
			return nil, err
		}
		ret[k] = bytesVal
	}
	return ret, nil
}

func UnmarshalConfigs(configs ConfigsByKey) (map[string]proto.Message, error) {
	ret := map[string]proto.Message{}
	for k, v := range configs {
		anyVal := &any.Any{}
		err := protos.Unmarshal(v, anyVal)
		if err != nil {
			return nil, err
		}
		msgVal, err := ptypes.Empty(anyVal)
		err = ptypes.UnmarshalAny(anyVal, msgVal)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ret[k] = msgVal
	}
	return ret, nil
}
