/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin

import (
	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/services/carrier_wifi/config"
	"magma/cwf/cloud/go/services/carrier_wifi/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type Builder struct{}

func (*Builder) Build(
	networkID string,
	gatewayID string,
	graph configurator.EntityGraph,
	network configurator.Network,
	mconfigOut map[string]proto.Message,
) error {
	// we only build an mconfig if carrier_wifi network configs exist
	inwConfig, found := network.Configs[cwf.CwfNetworkType]
	if !found || inwConfig == nil {
		return nil
	}
	nwConfig := inwConfig.(*models.NetworkCarrierWifiConfigs)

	vals, err := config.BuildFromNetworkConfig(nwConfig)
	if err != nil {
		return errors.WithStack(err)
	}
	for k, v := range vals {
		mconfigOut[k] = v
	}
	return nil
}
