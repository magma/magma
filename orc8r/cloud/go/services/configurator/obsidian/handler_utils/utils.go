/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package handler_utils

import (
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/storage"
)

// Create an empty network if it doesn't exist already. If the network already
// exists, return the networkID, otherwise return the created networkID
func CreateNetworkIfNotExists(networkID string) error {
	exists, err := configurator.DoesNetworkExist(networkID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// create a network
	network := configurator.Network{
		ID: networkID,
	}
	_, err = configurator.CreateNetworks([]configurator.Network{network})
	if err != nil {
		return err
	}
	return nil
}

// Create an empty network and/or network entity if it doesn't exist already.
// If the entity already exists, return its networkID and entityID. Otherwise,
// return (networkID, entityID) that point to the newly created entity.
func CreateNetworkEntityIfNotExists(networkID, entityType, entityID string) error {
	err := CreateNetworkIfNotExists(networkID)
	if err != nil {
		return err
	}

	exists, err := configurator.DoesEntityExist(networkID, entityType, entityID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Create an empty entity
	networkEntity := configurator.NetworkEntity{
		Type: entityType,
		Key:  entityID,
	}
	_, err = configurator.CreateEntities(networkID, []configurator.NetworkEntity{networkEntity})
	if err != nil {
		return err
	}
	return nil
}

func CreateInternalNetworkEntityIfNotExists(entityType, entityID string) error {
	return CreateNetworkEntityIfNotExists(storage.InternalNetworkID, entityType, entityID)
}
