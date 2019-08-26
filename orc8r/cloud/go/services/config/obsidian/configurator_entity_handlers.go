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

	magma_errors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

// This set of CRUD handlers are meant for entities that only have one config
// per entity.
func configuratorCreateEntityConfig(c echo.Context, networkID string, entityType string, entityKey string, iConfig interface{}) error {
	config, nerr := GetConfigAndValidate(c, iConfig)
	if nerr != nil {
		return nerr
	}
	entityExists, err := configurator.DoesEntityExist(networkID, entityType, entityKey)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, fmt.Sprintf("Entity %s,%s does not exist in %s", entityType, entityKey, networkID)), http.StatusInternalServerError)
	}
	if !entityExists {
		_, err = configurator.CreateEntity(networkID, configurator.NetworkEntity{
			Key:    entityKey,
			Type:   entityType,
			Config: config,
		})
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "Failed to create entity"), http.StatusInternalServerError)
		}
	} else {
		err := configurator.CreateOrUpdateEntityConfig(networkID, entityType, entityKey, config)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "Failed to create entity config"), http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusCreated, entityKey)
}

func configuratorGetEntityConfig(c echo.Context, networkID string, entityType string, entityKey string) error {
	ent, err := configurator.LoadEntity(networkID, entityType, entityKey, configurator.EntityLoadCriteria{LoadConfig: true})
	if err == magma_errors.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if ent.Config == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, ent.Config)
}

func configuratorUpdateEntityConfig(c echo.Context, networkID string, entityType string, entityKey string, iConfig interface{}) error {
	config, nerr := GetConfigAndValidate(c, iConfig)
	if nerr != nil {
		return nerr
	}
	err := configurator.CreateOrUpdateEntityConfig(networkID, entityType, entityKey, config)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func configuratorDeleteEntityConfig(c echo.Context, networkID string, entityType string, entityKey string) error {
	err := configurator.DeleteEntityConfig(networkID, entityType, entityKey)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func configuratorGetAllKeys(c echo.Context, networkID, entityType string) error {
	keysArr, err := configurator.ListEntityKeys(networkID, entityType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if keysArr == nil {
		return obsidian.HttpError(errors.New("Keys not found"), http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, keysArr)
}
