/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"fmt"
	"net/http"
	"reflect"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
)

// PartialNetworkModels describe models that represents a portion of network
// that can be read, updated, and deleted.
type PartialNetworkModels interface {
	// ValidateModel validates the model to be according to swagger spec, as
	// well as other custom validations
	ValidateModel() error
	// GetFromNetwork grabs the desired model from the configurator network.
	// Returns nil if it is not there.
	GetFromNetwork(network configurator.Network) interface{}
	// ToUpdateCriteria takes in the existing network and applies the change
	// from the model to create a NetworkUpdateCriteria
	ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error)
}

// GetPartialNetworkHandlers returns a set of GET/PUT/DELETE handlers according to the parameters.
// If the configKey is not "", it will add a delete handler for the network config for that key.
func GetPartialNetworkHandlers(path string, model PartialNetworkModels, configKey string) []obsidian.Handler {
	ret := []obsidian.Handler{
		GetPartialReadNetworkHandler(path, model),
		GetPartialUpdateNetworkHandler(path, model),
	}
	if configKey != "" {
		ret = append(ret, GetPartialDeleteNetworkHandler(path, configKey))
	}
	return []obsidian.Handler{
		GetPartialReadNetworkHandler(path, model),
		GetPartialUpdateNetworkHandler(path, model),
		GetPartialDeleteNetworkHandler(path, configKey),
	}
}

// GetPartialReadNetworkHandler returns a GET obsidian handler at the specified path.
// This function loads a network specified by the networkID and returns the
// part of the network that corresponds to the given model.
// Example:
//      (m *NetworkName) GetFromNetwork(network configurator.Network) interface{} {
// 			return string(network.Name)
// 		}
// 		getNameHandler := handlers.GetPartialReadNetworkHandler(URL, &models.NetworkName{})
//
//      would return a GET handler that can read the network name of a network with the specified ID.
func GetPartialReadNetworkHandler(path string, model PartialNetworkModels) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.GET,
		HandlerFunc: func(c echo.Context) error {
			networkID, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			network, err := configurator.LoadNetwork(networkID, true, true)
			if err == errors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			ret := model.GetFromNetwork(network)
			if ret == nil {
				return obsidian.HttpError(fmt.Errorf("Not found"), http.StatusNotFound)
			}
			return c.JSON(http.StatusOK, ret)
		},
	}
}

// GetPartialUpdateNetworkHandler returns a PUT obsidian handler at the specified path.
// The handler will fetch the payload into the configModel and perform validations according to the swagger spec.
// updater will take the model and apply the change into an existing network.
// Example:
//      (m *NetworkName) ToUpdateCriteria(network configurator.Network) interface{} {
// 			return configurator.NetworkUpdateCriteria{
//				ID:   network.ID,
// 				Name: *m,
//			}
//      }
// 		putNameHandler := handlers.GetPartialUpdateNetworkHandler(URL, &models.NetworkName{})
//
//      would return a PUT handler that will intake a NetworkName model and update the corresponding network
func GetPartialUpdateNetworkHandler(path string, model PartialNetworkModels) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.PUT,
		HandlerFunc: func(c echo.Context) error {
			networkID, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			requestedUpdate, nerr := getPayload(c, model)
			if nerr != nil {
				return nerr
			}

			network, err := configurator.LoadNetwork(networkID, true, true)
			if err == errors.ErrNotFound {
				return obsidian.HttpError(err, http.StatusNotFound)
			} else if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}

			updateCriteria, err := requestedUpdate.ToUpdateCriteria(network)
			if err != nil {
				return obsidian.HttpError(err, http.StatusBadRequest)
			}
			err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{updateCriteria})
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

// GetPartialDeleteNetworkHandler returns a DELETE obsidian handler at the specified path.
// The handler will delete a network config specified by the key.
// Example:
// 		deleteNetworkFeaturesHandler := handlers.GetPartialDeleteNetworkHandler(URL, "orc8r_features")
//
//      would return a DELETE handler that will remove the network features config from the corresponding network
func GetPartialDeleteNetworkHandler(path string, key string) obsidian.Handler {
	return obsidian.Handler{
		Path:    path,
		Methods: obsidian.DELETE,
		HandlerFunc: func(c echo.Context) error {
			networkID, nerr := obsidian.GetNetworkId(c)
			if nerr != nil {
				return nerr
			}
			update := configurator.NetworkUpdateCriteria{
				ID:              networkID,
				ConfigsToDelete: []string{key},
			}
			err := configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{update})
			if err != nil {
				return obsidian.HttpError(err, http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusNoContent)
		},
	}
}

func getPayload(c echo.Context, model interface{}) (PartialNetworkModels, *echo.HTTPError) {
	iModel := reflect.New(reflect.TypeOf(model).Elem()).Interface().(PartialNetworkModels)
	if err := c.Bind(iModel); err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}
	// Run validations specified by the swagger spec
	if err := iModel.ValidateModel(); err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}
	return iModel, nil
}
