/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"net/http"
	"path"
	"reflect"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
)

func getSerdeModelFromRegistry(configType string) interface{} {
	model, _ := serde.GetSerdeByDomainAndType(configurator.NetworkConfigSerdeDomain, configType)
	modelType := reflect.TypeOf(model).Elem()
	return reflect.New(modelType).Interface().(serde.BinaryConvertible)
}

// GetCRUDNetworkConfigHandlers returns 4 Handlers which implement CRUD for
// a network config.
// Path should look like '.../networks/:network_id/{config_type}'
// Serde is used to serialize/deserialize the config stored
func GetCRUDNetworkConfigHandlers(path string) []handlers.Handler {
	return []handlers.Handler{
		GetCreateNetworkConfigHandler(path),
		GetReadNetworkConfigHandler(path),
		GetUpdateNetworkConfigHandler(path),
		GetDeleteNetworkConfigHandler(path),
	}
}

func getNetworkIDAndConfigType(c echo.Context, url string) (string, string, *echo.HTTPError) {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return "", "", nerr
	}
	configType := path.Base(url)
	return networkID, configType, nil
}

func GetReadNetworkConfigHandler(path string) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, configType, nerr := getNetworkIDAndConfigType(c, path)
			if nerr != nil {
				return nerr
			}
			iConfig, err := configurator.GetNetworkConfigsByType(networkID, configType)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, iConfig)
		},
	}
}

func GetCreateNetworkConfigHandler(path string) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.POST,
		HandlerFunc: func(c echo.Context) error {
			return createOrUpdateNetworkConfig(c, path)
		},
	}
}

func GetUpdateNetworkConfigHandler(path string) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.PUT,
		HandlerFunc: func(c echo.Context) error {
			return createOrUpdateNetworkConfig(c, path)
		},
	}
}

func GetDeleteNetworkConfigHandler(path string) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.DELETE,
		HandlerFunc: func(c echo.Context) error {
			networkID, configType, nerr := getNetworkIDAndConfigType(c, path)
			if nerr != nil {
				return nerr
			}
			err := configurator.DeleteNetworkConfig(networkID, configType)
			if err != nil {
				return handlers.HttpError(err, http.StatusBadRequest)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

func createOrUpdateNetworkConfig(c echo.Context, path string) error {
	networkID, configType, nerr := getNetworkIDAndConfigType(c, path)
	if nerr != nil {
		return nerr
	}
	config, err := getConfigFromContext(c, configType)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	updateCriteria := configurator.NetworkUpdateCriteria{
		ID:                   networkID,
		ConfigsToAddOrUpdate: map[string]interface{}{configType: config},
	}
	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{updateCriteria})
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, "Created Config")
}

func getConfigFromContext(c echo.Context, configType string) (interface{}, error) {
	model := getSerdeModelFromRegistry(configType)
	err := c.Bind(&model)
	if err != nil {
		return nil, handlers.HttpError(err, http.StatusBadRequest)
	}
	return model, nil
}
