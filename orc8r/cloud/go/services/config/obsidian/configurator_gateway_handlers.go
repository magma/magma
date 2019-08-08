/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian

import (
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/magmad/obsidian/models"
	upgrade_models "magma/orc8r/cloud/go/services/upgrade/obsidian/models"
	"magma/orc8r/cloud/go/storage"

	"github.com/labstack/echo"
)

// Gateway is a special type of network entity. Since magmad supports multiple
// config types per gateway, each config is modeled as an entity with an
// association to the gateway entity. If the configType is "magmad_gateway" the
// config will simply be stored inside the gateway entity with an association
// to the corresponding tier entity.

func configuratorCreateGatewayConfig(c echo.Context, networkID string, configType string, configKey string, iConfig interface{}) error {
	// If the type is magmad_gateway add the config to the gw entity
	if configType == orc8r.MagmadGatewayType {
		return configuratorCreateMagmadGatewayConfig(c, networkID, configKey, iConfig)
	}
	config, nerr := GetConfigAndValidate(c, iConfig)
	if nerr != nil {
		return nerr
	}

	// the migrated handler will create the entity and associate the magmad
	// access gateway to it
	// note that this operation is not atomic, so there is a very slim but
	// nonzero chance that the entity is created without the proper assoc
	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:   configType,
		Key:    configKey,
		Config: config,
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	_, err = configurator.UpdateEntity(networkID, configurator.EntityUpdateCriteria{
		// hardcoded type to prevent import cycle
		Type:              orc8r.MagmadGatewayType,
		Key:               configKey,
		AssociationsToAdd: []storage.TypeAndKey{{Type: configType, Key: configKey}},
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, configKey)
}

func configuratorUpdateGatewayConfig(c echo.Context, networkID string, configType string, configKey string, iConfig interface{}) error {
	if configType == orc8r.MagmadGatewayType {
		return configuratorUpdateMagmadGatewayConfig(c, networkID, configKey, iConfig)
	}

	// if the config is not for magmad_gateway, this is the same as entity config update
	return configuratorUpdateEntityConfig(c, networkID, configType, configKey, iConfig)
}

func configuratorDeleteGatewayConfig(c echo.Context, networkID string, configType string, configKey string) error {
	if configType == orc8r.MagmadGatewayType {
		return configuratorDeleteMagmadGatewayConfig(c, networkID, configKey)
	}

	err := configurator.DeleteEntity(networkID, configType, configKey)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func configuratorCreateMagmadGatewayConfig(c echo.Context, networkID string, gatewayID string, iConfig interface{}) error {
	iRequestedConfig, nerr := GetConfigAndValidate(c, iConfig)
	if nerr != nil {
		return nerr
	}
	requestedConfig := iRequestedConfig.(*models.MagmadGatewayConfig)

	nerr = createTierEntityIfDoesNotExist(networkID, requestedConfig.Tier)
	if nerr != nil {
		return nerr
	}

	gwUpdate := configurator.EntityUpdateCriteria{
		Type:      orc8r.MagmadGatewayType,
		Key:       gatewayID,
		NewConfig: requestedConfig,
	}
	tierUpdate := configurator.EntityUpdateCriteria{
		Type:              orc8r.UpgradeTierEntityType,
		Key:               requestedConfig.Tier,
		AssociationsToAdd: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
	}
	_, err := configurator.UpdateEntities(networkID, []configurator.EntityUpdateCriteria{gwUpdate, tierUpdate})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, gatewayID)
}

func configuratorUpdateMagmadGatewayConfig(c echo.Context, networkID string, gatewayID string, iConfig interface{}) error {
	iRequestedConfig, nerr := GetConfigAndValidate(c, iConfig)
	if nerr != nil {
		return nerr
	}
	requestedConfig := iRequestedConfig.(*models.MagmadGatewayConfig)

	// fetch previous update tier to compute association change
	iCurrentConfig, err := configurator.LoadEntityConfig(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	currentConfig := iCurrentConfig.(*models.MagmadGatewayConfig)

	updates := []configurator.EntityUpdateCriteria{}
	if currentConfig.Tier != requestedConfig.Tier {
		nerr = createTierEntityIfDoesNotExist(networkID, requestedConfig.Tier)
		if nerr != nil {
			return nerr
		}

		updateToCurrentTier := configurator.EntityUpdateCriteria{
			Type:                 orc8r.UpgradeTierEntityType,
			Key:                  currentConfig.Tier,
			AssociationsToDelete: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
		}
		updateToRequestedTier := configurator.EntityUpdateCriteria{
			Type:              orc8r.UpgradeTierEntityType,
			Key:               requestedConfig.Tier,
			AssociationsToAdd: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gatewayID}},
		}
		updates = []configurator.EntityUpdateCriteria{updateToCurrentTier, updateToRequestedTier}
	}

	updateToGw := configurator.EntityUpdateCriteria{
		Type:      orc8r.MagmadGatewayType,
		Key:       gatewayID,
		NewConfig: requestedConfig,
	}
	updates = append(updates, updateToGw)

	_, err = configurator.UpdateEntities(networkID, updates)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func configuratorDeleteMagmadGatewayConfig(c echo.Context, networkID, gatewayID string) error {
	update := configurator.EntityUpdateCriteria{
		Type:         orc8r.MagmadGatewayType,
		Key:          gatewayID,
		DeleteConfig: true,
	}
	_, err := configurator.UpdateEntity(networkID, update)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func createTierEntityIfDoesNotExist(networkID, tierID string) error {
	exists, err := configurator.DoesEntityExist(networkID, orc8r.UpgradeTierEntityType, tierID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if exists {
		return nil
	}
	entity := configurator.NetworkEntity{
		Type: orc8r.UpgradeTierEntityType,
		Key:  tierID,
		Config: &upgrade_models.Tier{
			ID: tierID,
		},
	}
	_, err = configurator.CreateEntity(networkID, entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return nil
}
