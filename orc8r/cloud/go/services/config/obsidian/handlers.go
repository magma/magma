/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Generic obsidian handlers for configuration management
package obsidian

import (
	"errors"
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/serde"

	"github.com/labstack/echo"
)

// ConvertibleUserModel defines a configuration object exposed via obsidian
// which can be converted to and from a corresponding configuration object
// from the config service.
type ConvertibleUserModel interface {
	serde.BinaryConvertible

	ValidateModel() error
}

type ConfigType int

const (
	Network ConfigType = 1
	Gateway ConfigType = 2
	// Entity is a network entity in configurator that is not a magmad_gateway.
	// This distinction is important since the current setup allows a gateway
	// to have multiple configs. A gateway config is treated as a separate
	// entity with an association to the gateway entity. With a simple entity,
	// its config is stored inside the entity.
	Entity ConfigType = 3
)

// ConfigKeyGetter is a function which returns a config key from an
// echo.Context.
type ConfigKeyGetter func(echo.Context) (string, *echo.HTTPError)

// GetCRUDConfigHandlers returns 4 Handlers which implement GET/POST/PUT/DELETE
// for a given config type.
// configKeyGetter is a function which returns the desired config key value
// from a request's echo context.
// userModel is a pointer to an instance of the config struct that these
// handlers will manage.
func GetCRUDConfigHandlers(
	path string,
	configType string,
	configKeyGetter ConfigKeyGetter,
	userModel ConvertibleUserModel,
) []obsidian.Handler {
	return []obsidian.Handler{
		GetReadConfigHandler(path, configType, configKeyGetter, userModel),
		GetCreateConfigHandler(path, configType, configKeyGetter, userModel),
		GetUpdateConfigHandler(path, configType, configKeyGetter, userModel),
		GetDeleteConfigHandler(path, configType, configKeyGetter),
	}
}

// GetCRUDNetworkConfigHandlers returns 4 Handlers which implement CRUD for
// a network config. See GetCRUDConfigHandlers for additional documentation.
func GetCRUDNetworkConfigHandlers(
	path string,
	configType string,
	userModel ConvertibleUserModel,
) []obsidian.Handler {
	return []obsidian.Handler{
		GetReadNetworkConfigHandler(path, configType, userModel),
		GetCreateNetworkConfigHandler(path, configType, userModel),
		GetUpdateNetworkConfigHandler(path, configType, userModel),
		GetDeleteNetworkConfigHandler(path, configType),
	}
}

// GetCRUDGatewayConfigHandlers returns 4 Handlers which implement CRUD for
// a gateway config. See GetCRUDConfigHandlers for additional documentation.
func GetCRUDGatewayConfigHandlers(
	path string,
	configType string,
	userModel ConvertibleUserModel,
) []obsidian.Handler {
	return []obsidian.Handler{
		GetReadGatewayConfigHandler(path, configType, userModel),
		GetCreateGatewayConfigHandler(path, configType, userModel),
		GetUpdateGatewayConfigHandler(path, configType, userModel),
		GetDeleteGatewayConfigHandler(path, configType),
	}
}

// GetReadAllKeysConfigHandler returns an obsidian handler for reading all
// keys of a type using the config service.
// The returned Handler will have Methods set to GET.
//
// path is the URI for the handler to serve.
func GetReadAllKeysConfigHandler(
	path string,
	configType string,
) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			return configuratorGetAllKeys(c, networkID, configType)
		},
	}
}

// GetReadConfigHandler returns an obsidian handler for getting a config
// from the config service. The returned Handler will have Methods set to GET.
//
// path is the URI for the handler to serve.
// configKeyGetter is a function which returns the desired config key value
// from the request's echo context.
// userModel is a pointer to an instance of the config struct that this handler
// is for.
func GetReadConfigHandler(
	path string,
	configType string,
	configKeyGetter ConfigKeyGetter,
	userModel ConvertibleUserModel,
) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			networkId, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}

			switch getConfigTypeForConfigurator(configType) {
			case Network:
				return configuratorGetNetworkConfig(c, networkId, configType)
			case Gateway, Entity:
				configKey, cerr := configKeyGetter(c)
				if cerr != nil {
					return cerr
				}
				return configuratorGetEntityConfig(c, networkId, configType, configKey)
			default:
				return obsidian.HttpError(errors.New("not implemented"), http.StatusNotImplemented)
			}
		},
	}
}

// GetReadNetworkConfigHandler returns an obsidian handler for getting a
// network config from the config service.
// See GetReadConfigHandler for additional documentation.
func GetReadNetworkConfigHandler(path string, configType string, userModel ConvertibleUserModel) obsidian.Handler {
	return GetReadConfigHandler(path, configType, obsidian.GetNetworkId, userModel)
}

// GetReadGatewayConfigHandler returns an obsidian handler for getting a
// gateway config from the config service.
// See GetReadConfigHandler for additional documentation.
func GetReadGatewayConfigHandler(path string, configType string, userModel ConvertibleUserModel) obsidian.Handler {
	return GetReadConfigHandler(path, configType, obsidian.GetLogicalGwId, userModel)
}

// GetCreateConfigHandler returns an obsidian handler for creating a config
// using the config service. The returned Handler will have Methods set to POST.
//
// path is the URI for the handler to serve.
// configKeyGetter is a function which returns the desired config key value
// from the request's echo context.
// userModel is a pointer to an instance of the config struct that this handler
// is for.
func GetCreateConfigHandler(
	path string,
	configType string,
	configKeyGetter ConfigKeyGetter,
	userModel ConvertibleUserModel,
) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.POST,
		HandlerFunc: func(c echo.Context) error {
			networkId, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}

			switch getConfigTypeForConfigurator(configType) {
			case Network:
				return configuratorCreateNetworkConfig(c, networkId, configType, userModel)
			case Gateway:
				configKey, err := configKeyGetter(c)
				if err != nil {
					return err
				}
				return configuratorCreateGatewayConfig(c, networkId, configType, configKey, userModel)
			case Entity:
				configKey, err := configKeyGetter(c)
				if err != nil {
					return err
				}
				return configuratorCreateEntityConfig(c, networkId, configType, configKey, userModel)
			default:
				return obsidian.HttpError(errors.New("not implemented"), http.StatusNotImplemented)
			}
		},
	}
}

// GetCreateNetworkConfigHandler returns an obsidian handler for creating a
// network config using the config service.
// See GetCreateConfigHandler for additional documentation.
func GetCreateNetworkConfigHandler(path string, configType string, userModel ConvertibleUserModel) obsidian.Handler {
	return GetCreateConfigHandler(path, configType, obsidian.GetNetworkId, userModel)
}

// GetCreateGatewayConfigHandler returns an obsidian handler for creating a
// gateway config using the config service.
// See GetCreateConfigHandler for additional documentation.
func GetCreateGatewayConfigHandler(path string, configType string, userModel ConvertibleUserModel) obsidian.Handler {
	return GetCreateConfigHandler(path, configType, obsidian.GetLogicalGwId, userModel)
}

// GetUpdateConfigHandler returns an obsidian handler for updating a config
// using the config service. The returned Handler will have Methods set to PUT.
//
// path is the URI for the handler to serve.
// configKeyGetter is a function which returns the desired config key value
// from the request's echo context.
// userModel is a pointer to an instance of the config struct that this handler
// is for.
func GetUpdateConfigHandler(
	path string,
	configType string,
	configKeyGetter ConfigKeyGetter,
	userModel ConvertibleUserModel,
) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.PUT,
		HandlerFunc: func(c echo.Context) error {
			networkId, err := obsidian.GetNetworkId(c)
			if err != nil {
				return err
			}

			switch getConfigTypeForConfigurator(configType) {
			case Network:
				return configuratorUpdateNetworkConfig(c, networkId, configType, userModel)
			case Gateway:
				configKey, err := configKeyGetter(c)
				if err != nil {
					return err
				}
				return configuratorUpdateGatewayConfig(c, networkId, configType, configKey, userModel)
			case Entity:
				configKey, err := configKeyGetter(c)
				if err != nil {
					return err
				}
				return configuratorUpdateEntityConfig(c, networkId, configType, configKey, userModel)
			default:
				return obsidian.HttpError(errors.New("not implemented"), http.StatusNotImplemented)
			}
		},
	}
}

// GetUpdateNetworkConfigHandler returns an obsidian handler for updating a
// network config using the config service.
// See GetUpdateConfigHandler for additional documentation
func GetUpdateNetworkConfigHandler(path string, configType string, userModel ConvertibleUserModel) obsidian.Handler {
	return GetUpdateConfigHandler(path, configType, obsidian.GetNetworkId, userModel)
}

// GetUpdateGatewayConfigHandler returns an obsidian handler for updating a
// gateway config using the config service.
// See GetUpdateConfigHandler for additional documentation
func GetUpdateGatewayConfigHandler(path string, configType string, userModel ConvertibleUserModel) obsidian.Handler {
	return GetUpdateConfigHandler(path, configType, obsidian.GetLogicalGwId, userModel)
}

// GetDeleteConfigHandler returns an obsidian handler for deleting a config
// from the config service. The returned Handler will have Methods set to DELETE.
//
// path is the URI for the handler to serve.
// configKeyGetter is a function which returns the desired config key value
// from the request's echo context.
func GetDeleteConfigHandler(path string, configType string, configKeyGetter ConfigKeyGetter) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.DELETE,
		HandlerFunc: func(c echo.Context) error {
			networkId, err := obsidian.GetNetworkId(c)
			if err != nil {
				return err
			}

			switch getConfigTypeForConfigurator(configType) {
			case Network:
				return configuratorDeleteNetworkConfig(c, networkId, configType)
			case Gateway:
				configKey, err := configKeyGetter(c)
				if err != nil {
					return err
				}
				return configuratorDeleteGatewayConfig(c, networkId, configType, configKey)
			case Entity:
				configKey, err := configKeyGetter(c)
				if err != nil {
					return err
				}
				return configuratorDeleteEntityConfig(c, networkId, configType, configKey)
			default:
				return obsidian.HttpError(errors.New("not implemented"), http.StatusNotImplemented)
			}
		},
	}
}

// GetDeleteNetworkConfigHandler returns an obsidian handler for deleting
// a network config using the config service.
// See GetDeleteConfigHandler for additional documentation.
func GetDeleteNetworkConfigHandler(path string, configType string) obsidian.Handler {
	return GetDeleteConfigHandler(path, configType, obsidian.GetNetworkId)
}

// GetDeleteGatewayConfigHandler returns an obsidian handler for deleting
// a gateway config using the config service.
// See GetDeleteConfigHandler for additional documentation.
func GetDeleteGatewayConfigHandler(path string, configType string) obsidian.Handler {
	return GetDeleteConfigHandler(path, configType, obsidian.GetLogicalGwId)
}
