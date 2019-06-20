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

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/labstack/echo"
)

// Gateway is a special type of network entity. Since magmad supports multiple
// config types per gateway, each config if modeled as an entity with an
// association to the gateway entity. This assumes that the configType passed
// in is different from the entity type of the gateway.

func configuratorCreateGatewayConfig(c echo.Context, networkID string, configType string, configKey string, iConfig interface{}) error {
	config, nerr := getConfigAndValidate(c, iConfig)
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
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	_, err = configurator.UpdateEntity(networkID, configurator.EntityUpdateCriteria{
		// hardcoded type to prevent import cycle
		Type:              "magmad_gateway",
		Key:               configKey,
		AssociationsToAdd: []storage.TypeAndKey{{Type: configType, Key: configKey}},
	})
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, configKey)
}

func configuratorDeleteGatewayConfig(c echo.Context, networkID string, configType string, configKey string) error {
	err := configurator.DeleteEntity(networkID, configType, configKey)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
