/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian

import (
	"fmt"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_utils "magma/orc8r/cloud/go/services/configurator/obsidian/handler_utils"
	"magma/orc8r/cloud/go/services/configurator/storage"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
)

// This is a duplicated constant from magmad
const accessGatewayEntityType = "magmad_gateway"

// case on configType and propagate create/update into configurator
func multiplexCreateOrUpdateConfigIntoConfigurator(networkID, configType string, configKey string, iConfig interface{}) error {
	switch getConfigTypeForConfigurator(configType) {
	case NETWORK:
		return multiplexCreateOrUpdateNetworkConfig(networkID, configType, iConfig)
	case NETWORK_ENTITY:
		return multiplexCreateOrUpdateEntityConfig(networkID, configType, configKey, iConfig)
	default:
		return fmt.Errorf("Unexpected config type : %s", configType)
	}
}

func multiplexCreateOrUpdateNetworkConfig(networkID, configType string, config interface{}) error {
	// Create an empty network if it doesn't exist already
	err := configurator_utils.CreateNetworkIfNotExists(networkID)
	if err != nil {
		return err
	}
	err = configurator.UpdateNetworkConfig(networkID, configType, config)
	if err != nil {
		return fmt.Errorf(
			"Failed to multiplex create network config %s:%s into configurator: %v", networkID, configType, err)
	}
	return nil
}

func multiplexCreateOrUpdateEntityConfig(networkID, entityType, entityKey string, config interface{}) error {
	serializedConfig, err := serde.Serialize(configurator.NetworkEntitySerdeDomain, entityType, config)
	if err != nil {
		return err
	}

	err = configurator_utils.CreateNetworkEntityIfNotExists(networkID, accessGatewayEntityType, entityKey)
	if err != nil {
		return err
	}

	err = configurator_utils.CreateNetworkEntityIfNotExists(networkID, entityType, entityKey)
	if err != nil {
		return err
	}

	if entityType != accessGatewayEntityType {
		// send 2 updates - first, associate the magmad AG to this entity
		// second, update this entity's config
		_, err = configurator.UpdateEntities(
			networkID,
			[]*storage.EntityUpdateCriteria{
				{
					Type: accessGatewayEntityType,
					Key:  entityKey,
					AssociationsToAdd: []*storage.EntityID{
						{Type: entityType, Key: entityKey},
					},
				},
				{
					Type:      entityType,
					Key:       entityKey,
					NewConfig: &wrappers.BytesValue{Value: serializedConfig},
				},
			},
		)
	} else {
		// otherwise, we're just setting the magmad gateway config
		err = configurator.UpdateEntityConfig(networkID, entityType, entityKey, config)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to multiplex network entity config write %s:%s:%s into configurator", networkID, entityType, entityKey)
	}
	return nil
}

// case on configType and propagate delete into configurator
func multiplexDeleteConfigIntoConfigurator(networkID, configType, configKey string) error {
	switch getConfigTypeForConfigurator(configType) {
	case NETWORK:
		return multiplexDeleteNetworkConfig(networkID, configType)
	case NETWORK_ENTITY:
		return multiplexDeleteEntityConfig(networkID, configType, configKey)
	default:
		return fmt.Errorf("Unexpected config type : %s", configType)
	}
}

func multiplexDeleteNetworkConfig(networkID, configType string) error {
	// Create an empty network if it doesn't exist already
	err := configurator_utils.CreateNetworkIfNotExists(networkID)
	if err != nil {
		return err
	}
	err = configurator.DeleteNetworkConfig(networkID, configType)
	if err != nil {
		return fmt.Errorf(
			"Failed to multiplex delete network config %s:%s into configurator: %v", networkID, configType, err)
	}
	return nil
}

func multiplexDeleteEntityConfig(networkID, configType, configKey string) error {
	err := configurator_utils.CreateNetworkEntityIfNotExists(networkID, configType, configKey)
	if err != nil {
		return err
	}
	err = configurator.DeleteEntityConfig(networkID, configType, configKey)
	if err != nil {
		return fmt.Errorf(
			"Failed to multiplex delete network entity config %s:%s:%s into configurator: %v", networkID, configType, configKey, err)
	}
	return nil
}
