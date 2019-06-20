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

	"github.com/labstack/echo"
)

// This set of CRUD handlers are meant for entities that only have one config
// per entity.
func configuratorCreateEntityConfig(c echo.Context, networkID string, entityType string, entityKey string, iConfig interface{}) error {
	cfgInstance := reflect.New(reflect.TypeOf(iConfig).Elem()).Interface().(ConvertibleUserModel)
	if err := c.Bind(cfgInstance); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := cfgInstance.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	err := configurator.CreateOrUpdateEntityConfig(networkID, entityType, entityKey, cfgInstance)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, entityKey)
}

func configuratorGetEntityConfig(c echo.Context, networkID string, entityType string, entityKey string) error {
	ent, err := configurator.LoadEntity(networkID, entityType, entityKey, configurator.EntityLoadCriteria{LoadConfig: true})
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

func configuratorUpdateEntityConfig(c echo.Context, networkID string, entityType string, entityKey string, iConfig interface{}) error {
	cfgInstance := reflect.New(reflect.TypeOf(iConfig).Elem()).Interface().(ConvertibleUserModel)
	if err := c.Bind(cfgInstance); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := cfgInstance.ValidateModel(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	err := configurator.CreateOrUpdateEntityConfig(networkID, entityType, entityKey, cfgInstance)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func configuratorDeleteEntityConfig(c echo.Context, networkID string, entityType string, entityKey string) error {
	err := configurator.DeleteEntityConfig(networkID, entityType, entityKey)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
