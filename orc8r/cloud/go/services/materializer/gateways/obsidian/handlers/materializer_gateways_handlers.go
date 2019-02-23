/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/services/materializer"
	"magma/orc8r/cloud/go/services/materializer/gateways/obsidian/models"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/dynamo"
	storage_sql "magma/orc8r/cloud/go/services/materializer/gateways/storage/sql"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/xservice"

	"github.com/labstack/echo"
)

func GetStorage() (storage.GatewayViewStorage, error) {
	// Loading the config manually because the service is initialized separately
	// Hardcoding in orc8r to avoid cyclical import
	configMap, err := config.GetServiceConfig("orc8r", materializer.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving config map from materializer %v\n", err)
	}
	v := configMap.GetRequiredStringParam("obsidian_read_storage")
	if strings.ToLower(v) == "sql" {
		db, err := sql.Open(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
		if err != nil {
			return nil, fmt.Errorf("Could not initialize SQL connection: %s", err)
		}
		return storage_sql.NewSqlGatewayViewStorage(db), nil
	} else if strings.ToLower(v) == "xservice" {
		return xservice.NewCrossServiceGatewayViewsStorage(), nil
	}

	// Default to dynamo for back-compat
	return dynamo.GetInitializedDynamoStorage()
}

func ListGatewayViews(c echo.Context, store storage.GatewayViewStorage) error {
	networkID, httpErr := handlers.GetNetworkId(c)
	if httpErr != nil {
		return httpErr
	}
	gatewayIDs := getGatewayIDs(c.QueryParams())
	gatewayStates, err := getGatewayStates(networkID, gatewayIDs, store)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	modelStates, err := models.GatewayStateMapToModelList(gatewayStates)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, modelStates)
}

func getGatewayIDs(queryParams url.Values) []string {
	gatewayIDs := []string{}
	format2Regex := regexp.MustCompile("^gateway_ids\\[[0-9]+\\]$")
	for queryKey, values := range queryParams {
		if queryKey == "gateway_ids" && len(values) > 0 && len(values[0]) > 0 {
			// Format 1: gateway_ids=gw1,gw2,gw3
			gatewayIDs = append(gatewayIDs, strings.Split(values[0], ",")...)
		} else if format2Regex.MatchString(queryKey) {
			// Format 2: gateway_ids[0]=gw1&gateway_ids[1]=gw2&gateway_ids[2]=gw3
			gatewayIDs = append(gatewayIDs, values...)
		}
	}
	return gatewayIDs
}

func getGatewayStates(
	networkID string,
	gatewayIDs []string,
	store storage.GatewayViewStorage,
) (map[string]*storage.GatewayState, error) {
	if len(gatewayIDs) > 0 {
		return store.GetGatewayViews(networkID, gatewayIDs)
	}
	return store.GetGatewayViewsForNetwork(networkID)
}
