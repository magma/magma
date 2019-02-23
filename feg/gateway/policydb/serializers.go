/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package policydb

import (
	"fmt"

	"magma/feg/gateway/object_store"
	"magma/lte/cloud/go/protos"

	"github.com/golang/protobuf/proto"
)

func getProtoSerializer() object_store.Serializer {
	return func(object interface{}) (string, error) {
		msg, ok := object.(proto.Message)
		if !ok {
			return "", fmt.Errorf("Could not cast object to protobuf")
		}
		bytes, err := proto.Marshal(msg)
		if err != nil {
			return "", fmt.Errorf("Could not marshal message")
		}
		return string(bytes[:]), nil
	}
}

func getPolicyDeserializer() object_store.Deserializer {
	return func(serialized string) (interface{}, error) {
		policyPtr := &protos.PolicyRule{}
		bytes := []byte(serialized)
		err := proto.Unmarshal(bytes, policyPtr)
		if err != nil {
			return nil, err
		}
		return policyPtr, nil
	}
}

func getBaseNameDeserializer() object_store.Deserializer {
	return func(serialized string) (interface{}, error) {
		setPtr := &protos.ChargingRuleNameSet{}
		bytes := []byte(serialized)
		err := proto.Unmarshal(bytes, setPtr)
		if err != nil {
			return nil, err
		}
		return setPtr, nil
	}
}
