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
	"reflect"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/labstack/echo"
)

// Gateways

func configuratorCreateGatewayConfig(c echo.Context, networkID string, configType string, configKey string, iConfig interface{}) error {
	cfgInstance := reflect.New(reflect.TypeOf(iConfig).Elem()).Interface().(ConvertibleUserModel)
	if err := c.Bind(cfgInstance); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := cfgInstance.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	// the migrated handler will create the entity and associate the magmad
	// access gateway to it
	// note that this operation is not atomic, so there is a very slim but
	// nonzero chance that the entity is created without the proper assoc
	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:   configType,
		Key:    configKey,
		Config: cfgInstance,
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

func configuratorGetGatewayConfig(c echo.Context, networkID string, configType string, configKey string) error {
	ent, err := configurator.LoadEntity(networkID, configType, configKey, configurator.EntityLoadCriteria{LoadConfig: true})
	if err == errors.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	if ent.Config == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, ent.Config)
}

func configuratorUpdateGatewayConfig(c echo.Context, networkID string, configType string, configKey string, iConfig interface{}) error {
	cfgInstance := reflect.New(reflect.TypeOf(iConfig).Elem()).Interface().(ConvertibleUserModel)
	if err := c.Bind(cfgInstance); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := cfgInstance.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	err := configurator.UpdateEntityConfig(networkID, configType, configKey, cfgInstance)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func configuratorDeleteGatewayConfig(c echo.Context, networkID string, configType string, configKey string) error {
	err := configurator.DeleteEntity(networkID, configType, configKey)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
