/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
)

func (m *NetworkName) ToUpdateCriteria(network configurator.Network) configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:      network.ID,
		NewName: swag.String(string(*m)),
	}
}

func (m *NetworkName) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkName(network.Name)
}

func (m *NetworkType) ToUpdateCriteria(network configurator.Network) configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:      network.ID,
		NewType: swag.String(string(*m)),
	}
}

func (m *NetworkType) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkType(network.Type)
}
func (m *NetworkDescription) ToUpdateCriteria(network configurator.Network) configurator.NetworkUpdateCriteria {
	return configurator.NetworkUpdateCriteria{
		ID:             network.ID,
		NewDescription: swag.String(string(*m)),
	}
}

func (m *NetworkDescription) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkDescription(network.Description)
}
