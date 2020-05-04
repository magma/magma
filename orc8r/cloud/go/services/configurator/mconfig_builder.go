/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package configurator

import (
	"sync"
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

// MconfigBuilder fills in a partial mconfig for a given gateway within a
// network.
type MconfigBuilder interface {
	// Build fills the mconfigOut parameter with parts of the gateway mconfig
	// that this mconfig builder is responsible for.
	// The whole entity graph associated with the gateway to build the mconfig
	// for is provided in the graph parameter, as well as the network that
	// the gateway belongs to. Both the graph and network will be loaded with
	// all fields.
	// mconfigOut is an output parameter - all updates should be made in-place
	// to this parameter.
	Build(networkID string, gatewayID string, graph EntityGraph, network Network, mconfigOut map[string]proto.Message) error
}

type builderRegistry struct {
	sync.RWMutex
	builders []MconfigBuilder
}

var builderRegistryInstance = &builderRegistry{builders: []MconfigBuilder{}}

// RegisterMconfigBuilders registers a collection of MconfigBuilders to make
// available to the southbound configurator servicer to call when gateways
// request mconfigs.
func RegisterMconfigBuilders(builders ...MconfigBuilder) {
	builderRegistryInstance.Lock()
	defer builderRegistryInstance.Unlock()
	builderRegistryInstance.builders = append(builderRegistryInstance.builders, builders...)
}

func CreateMconfig(networkID string, gatewayID string, graph *storage.EntityGraph, network *storage.Network) (*protos.GatewayConfigs, error) {
	nativeGraph, err := (EntityGraph{}).FromStorageProto(graph)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert storage graph to native type")
	}
	nativeNW, err := (Network{}).fromStorageProto(network)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert storage network to native type")
	}

	builderRegistryInstance.RLock()
	defer builderRegistryInstance.RUnlock()

	messages := map[string]proto.Message{}
	for _, builder := range builderRegistryInstance.builders {
		err := builder.Build(networkID, gatewayID, nativeGraph, nativeNW, messages)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build mconfig, builder %T errored", builder)
		}
	}

	ret := &protos.GatewayConfigs{
		Metadata: &protos.GatewayConfigsMetadata{
			CreatedAt: uint64(time.Now().Unix()),
			Digest:    &protos.GatewayConfigsDigest{},
		},
		ConfigsByKey: map[string]*any.Any{},
	}
	for k, msg := range messages {
		a, err := ptypes.MarshalAny(msg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal mconfig key %s to Any", k)
		}
		ret.ConfigsByKey[k] = a
	}
	digest, err := ret.GetMconfigDigest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate digest of ConfigsByKey")
	}
	ret.Metadata.Digest.Md5HexDigest = digest

	return ret, nil
}

// ClearMconfigBuilders exists ONLY for testing - this the required but unused
// *testing.T parameter.
// DO NOT USE IN ANYTHING BUT TESTS
func ClearMconfigBuilders(_ *testing.T) {
	builderRegistryInstance.Lock()
	defer builderRegistryInstance.Unlock()
	builderRegistryInstance.builders = []MconfigBuilder{}
}
