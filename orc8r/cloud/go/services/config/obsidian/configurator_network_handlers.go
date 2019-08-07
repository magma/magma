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
	"net/http"
	"reflect"
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
)

// Since the config service does not differentiate between configs that belong
// to networks vs network entities, this is a bit of a hack that relies on the
// current naming pattern to differentiate between the two.
// The current naming pattern is "*_network" and "*_gateway"
func getConfigTypeForConfigurator(configType string) ConfigType {
	split := strings.Split(configType, "_")
	lastIndex := len(split) - 1
	if len(split) > 1 && split[lastIndex] == "network" {
		return Network
	} else if len(split) == 2 && split[lastIndex] == "gateway" {
		return Gateway
	} else {
		return Entity
	}
}

// Networks

func configuratorCreateNetworkConfig(c echo.Context, networkID string, configType string, iConfig interface{}) error {
	config, nerr := GetConfigAndValidate(c, iConfig)
	if nerr != nil {
		return nerr
	}
	err := configurator.UpdateNetworkConfig(networkID, configType, config)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, networkID)
}

func configuratorGetNetworkConfig(c echo.Context, networkID string, configType string) error {
	cfg, err := configurator.GetNetworkConfigsByType(networkID, configType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if cfg == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, cfg)
}

func configuratorUpdateNetworkConfig(c echo.Context, networkID string, configType string, iConfig interface{}) error {
	config, nerr := GetConfigAndValidate(c, iConfig)
	if nerr != nil {
		return nerr
	}
	err := configurator.UpdateNetworkConfig(networkID, configType, config)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func configuratorDeleteNetworkConfig(c echo.Context, networkID string, configType string) error {
	err := configurator.DeleteNetworkConfig(networkID, configType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func GetConfigAndValidate(c echo.Context, iConfig interface{}) (ConvertibleUserModel, error) {
	cfgInstance := reflect.New(reflect.TypeOf(iConfig).Elem()).Interface().(ConvertibleUserModel)
	if err := c.Bind(cfgInstance); err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := cfgInstance.ValidateModel(); err != nil {
		return nil, obsidian.HttpError(fmt.Errorf("Invalid config: %s", err), http.StatusBadRequest)
	}
	return cfgInstance, nil
}
