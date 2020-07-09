/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package factory

import (
	"github.com/golang/protobuf/proto"
)

type MconfigBuilder interface {
	// Build a partial mconfig for a gateway given a network and gateway ID.
	// For the mconfig streamer, the keys of the map returned by this builder
	// must be unique across the entire system or else the streamer policy will
	// error out, as we cannot coordinate "merging" configs across different
	// modules.
	Build(networkId string, gatewayId string) (map[string]proto.Message, error)
}
