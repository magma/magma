/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
)

func (m *Network) ToConfiguratorNetwork() configurator.Network {
	return configurator.Network{
		ID:          string(m.ID),
		Type:        string(m.Type),
		Name:        string(m.Name),
		Description: string(m.Description),
		Configs: map[string]interface{}{
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}

func FromConfiguratorNetwork(n configurator.Network) *Network {
	m := &Network{}
	m.ID = models.NetworkID(n.ID)
	m.Type = n.Type
	m.Name = models.NetworkName(n.Name)
	m.Description = models.NetworkDescription(n.Description)
	if cfg, exists := n.Configs[orc8r.DnsdNetworkType]; exists {
		m.DNS = cfg.(*NetworkDNSConfig)
	}
	if cfg, exists := n.Configs[orc8r.NetworkFeaturesConfig]; exists {
		m.Features = cfg.(*NetworkFeatures)
	}
	return m
}

func (m *Network) ToConfiguratorNetworkUpdateCriteria() configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             string(m.ID),
		NewType:        swag.String(m.Type),
		NewName:        swag.String(string(m.Name)),
		NewDescription: swag.String(string(m.Description)),
		ConfigsToAddOrUpdate: map[string]interface{}{
			orc8r.DnsdNetworkType:       m.DNS,
			orc8r.NetworkFeaturesConfig: m.Features,
		},
	}
}
