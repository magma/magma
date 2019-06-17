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
	"strings"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
)

// Since the config service does not differentiate between configs that belong
// to networks vs network entities, this is a bit of a hack that relies on the
// current naming pattern to differentiate between the two.
func getConfigTypeForConfigurator(configType string) ConfigType {
	splittedConfigType := strings.Split(configType, "_")
	if len(splittedConfigType) > 1 && splittedConfigType[1] == "network" {
		return Network
	} else if len(splittedConfigType) == 1 || splittedConfigType[1] == "gateway" {
		return Entity
	} else {
		return Unrecognized
	}
}

// Networks

func configuratorCreateNetworkConfig(c echo.Context, networkID string, configType string, iConfig interface{}) error {
	cfgInstance := reflect.New(reflect.TypeOf(iConfig).Elem()).Interface().(ConvertibleUserModel)
	if err := c.Bind(cfgInstance); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := cfgInstance.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	err := configurator.UpdateNetworkConfig(networkID, configType, cfgInstance)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, networkID)
}

func configuratorGetNetworkConfig(c echo.Context, networkID string, configType string) error {
	cfg, err := configurator.GetNetworkConfigsByType(networkID, configType)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	if cfg == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, cfg)
}

func configuratorUpdateNetworkConfig(c echo.Context, networkID string, configType string, iConfig interface{}) error {
	cfgInstance := reflect.New(reflect.TypeOf(iConfig).Elem()).Interface().(ConvertibleUserModel)
	if err := c.Bind(cfgInstance); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := cfgInstance.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	err := configurator.UpdateNetworkConfig(networkID, configType, cfgInstance)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func configuratorDeleteNetworkConfig(c echo.Context, networkID string, configType string) error {
	err := configurator.DeleteNetworkConfig(networkID, configType)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
