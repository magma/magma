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
	"fmt"
	"net/http"
	"reflect"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config"

	"github.com/labstack/echo"
)

// ConvertibleUserModel defines a configuration object exposed via obsidian
// which can be converted to and from a corresponding configuration object
// from the config service.
type ConvertibleUserModel interface {
	serde.BinaryConvertible

	ValidateModel() error

	// DEPRECATED
	ToServiceModel() (interface{}, error)

	// DEPRECATED
	FromServiceModel(serviceModel interface{}) error
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

// instantiateNewConvertibleUserModel creates a new, empty instance of the
// provided userModel struct. The parameter is expected to be a pointer to a
// struct which implements ConvertibleUserModel.
//
// We need this because the userModel parameters in the functions below are
// all only instantiated once upon creation of a Handler object, so they
// will be re-used across obsidian calls unless we reflectively instantiate a
// zeroed-out copy during each call.
//
// If you refactor any of the handler factory functions below, make sure that
// this function is in each handler's call chain.
func instantiateNewConvertibleUserModel(userModel ConvertibleUserModel) ConvertibleUserModel {
	userModelType := reflect.TypeOf(userModel).Elem()
	return reflect.New(userModelType).Interface().(ConvertibleUserModel)
}

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
) []handlers.Handler {
	return []handlers.Handler{
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
) []handlers.Handler {
	return []handlers.Handler{
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
) []handlers.Handler {
	return []handlers.Handler{
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
) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, nerr := handlers.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			return handleGetAllKeys(c, networkID, configType)
		},
		MigratedHandlerFunc: func(c echo.Context) error {
			networkID, nerr := handlers.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			return configuratorGetAllKeys(c, networkID, configType)
		},
	}
}

func handleGetAllKeys(c echo.Context, networkID string, configType string) error {
	keysArr, err := config.ListKeysForType(networkID, configType)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	if keysArr == nil {
		return handlers.HttpError(errors.New("Keys not found"), http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, keysArr)
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
) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.GET,
		HandlerFunc: func(c echo.Context) error {
			networkId, nerr := handlers.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}

			configKey, cerr := configKeyGetter(c)
			if cerr != nil {
				return cerr
			}
			return handleGetConfig(c, networkId, configType, configKey, userModel)
		},
		MigratedHandlerFunc: func(c echo.Context) error {
			networkId, nerr := handlers.GetNetworkId(c)
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
				return handlers.HttpError(errors.New("not implemented"), http.StatusNotImplemented)
			}
		},
	}
}

// GetReadNetworkConfigHandler returns an obsidian handler for getting a
// network config from the config service.
// See GetReadConfigHandler for additional documentation.
func GetReadNetworkConfigHandler(path string, configType string, userModel ConvertibleUserModel) handlers.Handler {
	return GetReadConfigHandler(path, configType, handlers.GetNetworkId, userModel)
}

// GetReadGatewayConfigHandler returns an obsidian handler for getting a
// gateway config from the config service.
// See GetReadConfigHandler for additional documentation.
func GetReadGatewayConfigHandler(path string, configType string, userModel ConvertibleUserModel) handlers.Handler {
	return GetReadConfigHandler(path, configType, handlers.GetLogicalGwId, userModel)
}

func handleGetConfig(c echo.Context, networkId string, configType string, configKey string, userModel ConvertibleUserModel) error {
	iConfig, err := config.GetConfig(networkId, configType, configKey)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	if iConfig == nil {
		return handlers.HttpError(errors.New("Config not found"), http.StatusNotFound)
	}

	userModel = instantiateNewConvertibleUserModel(userModel)
	if err := userModel.FromServiceModel(iConfig); err != nil {
		return handlers.HttpError(fmt.Errorf("Could not fill user config model from config service: %s", err), http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, userModel)
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
) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.POST,
		HandlerFunc: func(c echo.Context) error {
			networkId, err := handlers.GetNetworkId(c)
			if err != nil {
				return err
			}

			configKey, err := configKeyGetter(c)
			if err != nil {
				return err
			}
			return handleCreateConfig(c, networkId, configType, configKey, userModel)
		},
		MigratedHandlerFunc: func(c echo.Context) error {
			networkId, nerr := handlers.GetNetworkId(c)
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
				return handlers.HttpError(errors.New("not implemented"), http.StatusNotImplemented)
			}
		},
		MultiplexAfterMigration: true,
	}
}

// GetCreateNetworkConfigHandler returns an obsidian handler for creating a
// network config using the config service.
// See GetCreateConfigHandler for additional documentation.
func GetCreateNetworkConfigHandler(path string, configType string, userModel ConvertibleUserModel) handlers.Handler {
	return GetCreateConfigHandler(path, configType, handlers.GetNetworkId, userModel)
}

// GetCreateGatewayConfigHandler returns an obsidian handler for creating a
// gateway config using the config service.
// See GetCreateConfigHandler for additional documentation.
func GetCreateGatewayConfigHandler(path string, configType string, userModel ConvertibleUserModel) handlers.Handler {
	return GetCreateConfigHandler(path, configType, handlers.GetLogicalGwId, userModel)
}

func handleCreateConfig(c echo.Context, networkId string, configType string, configKey string, userModel ConvertibleUserModel) error {
	userModel = instantiateNewConvertibleUserModel(userModel)
	if err := c.Bind(userModel); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := userModel.ValidateModel(); err != nil {
		return handlers.HttpError(fmt.Errorf("Invalid config: %s", err), http.StatusBadRequest)
	}

	iConfig, err := userModel.ToServiceModel()
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Error converting config model: %s", err), http.StatusBadRequest)
	}
	if err := config.CreateConfig(networkId, configType, configKey, iConfig); err != nil {
		return handlers.HttpError(fmt.Errorf("Error creating config: %s", err), http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, configKey)
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
) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.PUT,
		HandlerFunc: func(c echo.Context) error {
			networkId, err := handlers.GetNetworkId(c)
			if err != nil {
				return err
			}

			configKey, err := configKeyGetter(c)
			if err != nil {
				return err
			}
			return handleConfigUpdate(c, networkId, configType, configKey, userModel)
		},
		MigratedHandlerFunc: func(c echo.Context) error {
			networkId, err := handlers.GetNetworkId(c)
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
				return handlers.HttpError(errors.New("not implemented"), http.StatusNotImplemented)
			}
		},
		MultiplexAfterMigration: true,
	}
}

// GetUpdateNetworkConfigHandler returns an obsidian handler for updating a
// network config using the config service.
// See GetUpdateConfigHandler for additional documentation
func GetUpdateNetworkConfigHandler(path string, configType string, userModel ConvertibleUserModel) handlers.Handler {
	return GetUpdateConfigHandler(path, configType, handlers.GetNetworkId, userModel)
}

// GetUpdateGatewayConfigHandler returns an obsidian handler for updating a
// gateway config using the config service.
// See GetUpdateConfigHandler for additional documentation
func GetUpdateGatewayConfigHandler(path string, configType string, userModel ConvertibleUserModel) handlers.Handler {
	return GetUpdateConfigHandler(path, configType, handlers.GetLogicalGwId, userModel)
}

func handleConfigUpdate(c echo.Context, networkId string, configType string, configKey string, userModel ConvertibleUserModel) error {
	userModel = instantiateNewConvertibleUserModel(userModel)
	if err := c.Bind(userModel); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := userModel.ValidateModel(); err != nil {
		return handlers.HttpError(fmt.Errorf("Invalid config: %s", err), http.StatusBadRequest)
	}

	iConfig, err := userModel.ToServiceModel()
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Error converting config model: %s", err), http.StatusBadRequest)
	}
	if err := config.UpdateConfig(networkId, configType, configKey, iConfig); err != nil {
		return handlers.HttpError(fmt.Errorf("Error updating config: %s", err), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

// GetDeleteConfigHandler returns an obsidian handler for deleting a config
// from the config service. The returned Handler will have Methods set to DELETE.
//
// path is the URI for the handler to serve.
// configKeyGetter is a function which returns the desired config key value
// from the request's echo context.
func GetDeleteConfigHandler(path string, configType string, configKeyGetter ConfigKeyGetter) handlers.Handler {
	return handlers.Handler{
		Path:    path,
		Methods: handlers.DELETE,
		HandlerFunc: func(c echo.Context) error {
			networkId, err := handlers.GetNetworkId(c)
			if err != nil {
				return err
			}

			configKey, err := configKeyGetter(c)
			if err != nil {
				return err
			}
			return handleConfigDelete(c, networkId, configType, configKey)
		},
		MigratedHandlerFunc: func(c echo.Context) error {
			networkId, err := handlers.GetNetworkId(c)
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
				return handlers.HttpError(errors.New("not implemented"), http.StatusNotImplemented)
			}
		},
		MultiplexAfterMigration: true,
	}
}

// GetDeleteNetworkConfigHandler returns an obsidian handler for deleting
// a network config using the config service.
// See GetDeleteConfigHandler for additional documentation.
func GetDeleteNetworkConfigHandler(path string, configType string) handlers.Handler {
	return GetDeleteConfigHandler(path, configType, handlers.GetNetworkId)
}

// GetDeleteGatewayConfigHandler returns an obsidian handler for deleting
// a gateway config using the config service.
// See GetDeleteConfigHandler for additional documentation.
func GetDeleteGatewayConfigHandler(path string, configType string) handlers.Handler {
	return GetDeleteConfigHandler(path, configType, handlers.GetLogicalGwId)
}

func handleConfigDelete(c echo.Context, networkId string, configType string, configKey string) error {
	if err := config.DeleteConfig(networkId, configType, configKey); err != nil {
		return handlers.HttpError(fmt.Errorf("Error deleting config: %s", err), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
