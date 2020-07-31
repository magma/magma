/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package mconfig

import (
	"fmt"
	"time"

	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
)

func CreateMconfig(network *storage.Network, graph *storage.EntityGraph, gatewayID string) (*protos.GatewayConfigs, error) {
	builders, err := GetBuilders()
	if err != nil {
		return nil, err
	}

	configs := ConfigsByKey{}
	for _, b := range builders {
		partialConfig, err := b.Build(network, graph, gatewayID)
		if err != nil {
			return nil, errors.Wrapf(err, "mconfig builder %+v error", b)
		}
		for key, config := range partialConfig {
			_, ok := configs[key]
			if ok {
				return nil, fmt.Errorf("received partial config for key %v from multiple mconfig builders", key)
			}
			configs[key] = config
		}
	}

	mconfig := &protos.GatewayConfigs{
		Metadata: &protos.GatewayConfigsMetadata{
			CreatedAt: uint64(time.Now().Unix()),
			Digest:    &protos.GatewayConfigsDigest{},
		},
		ConfigsByKey: configs,
	}
	mconfig.Metadata.Digest.Md5HexDigest, err = mconfig.GetMconfigDigest()
	if err != nil {
		return nil, errors.Wrap(err, "generate mconfig digest")
	}

	return mconfig, nil
}
