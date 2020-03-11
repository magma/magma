/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package protos is protoc generated GRPC package and related continence functions
package protos

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/golang/protobuf/proto"
)

// GetMconfigDigest generates a representative hash of the configs (sans metadata).
func (cfg *GatewayConfigs) GetMconfigDigest() (string, error) {
	configsWithoutMetadata := &GatewayConfigs{ConfigsByKey: cfg.GetConfigsByKey()}
	serializedConfig, err := encodePbDeterministic(configsWithoutMetadata)
	if err != nil {
		return "", err
	}

	sum := md5.Sum(serializedConfig)
	digest := hex.EncodeToString(sum[:])
	return digest, nil
}

// encodePbDeterministic encodes protobuf while enforcing deterministic serialization.
// NOTE: deterministic != canonical, so do not expect this encoding to be
// equal across languages or even versions of golang/protobuf/proto.
// For further reading, see below.
// 	- https://developers.google.com/protocol-buffers/docs/encoding#implications
//	- https://gist.github.com/kchristidis/39c8b310fd9da43d515c4394c3cd9510
func encodePbDeterministic(pb proto.Message) ([]byte, error) {
	buf := &proto.Buffer{}
	buf.SetDeterministic(true)

	err := buf.Marshal(pb)
	return buf.Bytes(), err
}
