/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
)

func (m *NetworkName) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID:      network.ID,
		NewName: swag.String(string(*m)),
	}, nil
}

func (m *NetworkName) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkName(network.Name)
}

func (m *NetworkType) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID:      network.ID,
		NewType: swag.String(string(*m)),
	}, nil
}

func (m *NetworkType) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkType(network.Type)
}

func (m *NetworkDescription) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID:             network.ID,
		NewDescription: swag.String(string(*m)),
	}, nil
}

func (m *NetworkDescription) GetFromNetwork(network configurator.Network) interface{} {
	return NetworkDescription(network.Description)
}

func (m *GatewayName) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Key:     gatewayID,
			Type:    orc8r.MagmadGatewayType,
			NewName: swag.String(string(*m)),
		},
	}, nil
}

func (m *GatewayName) FromBackendModels(networkID string, gatewayID string) error {
	entity, err := configurator.LoadEntity(networkID, orc8r.MagmadGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadMetadata: true})
	if err != nil {
		return err
	}
	*m = GatewayName(entity.Name)
	return nil
}

func (m *GatewayDescription) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	return []configurator.EntityUpdateCriteria{
		{
			Key:            gatewayID,
			Type:           orc8r.MagmadGatewayType,
			NewDescription: swag.String(string(*m)),
		},
	}, nil
}

func (m *GatewayDescription) FromBackendModels(networkID string, gatewayID string) error {
	entity, err := configurator.LoadEntity(networkID, orc8r.MagmadGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadMetadata: true})
	if err != nil {
		return err
	}
	*m = GatewayDescription(entity.Description)
	return nil
}
