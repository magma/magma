/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/client"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/config"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/receivers"

	"github.com/labstack/echo"
	amconfig "github.com/prometheus/alertmanager/config"
)

const (
	v0rootPath               = "/:tenant_id"
	v0receiverPath           = "/receiver"
	v0RoutePath              = "/receiver/route"
	v0receiverNameQueryParam = "receiver"

	v1rootPath         = "/v1"
	v1receiverPath     = "/receiver"
	v1receiverNamePath = v1receiverPath + "/:" + receiverNameParam
	v1routePath        = "/route"
	v1GlobalPath       = "/global"

	receiverNameParam = "receiver_name"
	tenantIDParam     = "tenant_id"
)

func RegisterBaseHandlers(e *echo.Echo) {
	e.GET("/", statusHandler)
}

func RegisterV0Handlers(e *echo.Echo, client client.AlertmanagerClient) {
	v0 := e.Group(v0rootPath)
	v0.Use(tenancyMiddlewareProvider(client, pathTenantProvider))

	v0.POST(v0receiverPath, GetReceiverPostHandler(client))
	v0.GET(v0receiverPath, GetGetReceiversHandler(client))
	v0.DELETE(v0receiverPath, GetDeleteReceiverHandler(client, v0receiverNameQueryProvider))
	v0.PUT(v0receiverPath+"/:"+receiverNameParam, GetUpdateReceiverHandler(client, receiverNamePathProvider))

	v0.POST(v0RoutePath, GetUpdateRouteHandler(client))
	v0.GET(v0RoutePath, GetGetRouteHandler(client))
}

func RegisterV1Handlers(e *echo.Echo, client client.AlertmanagerClient) {
	v1 := e.Group(v1rootPath)

	// these don't require tenancy so register before middleware
	v1.POST(v1GlobalPath, GetUpdateGlobalConfigHandler(client))
	v1.GET(v1GlobalPath, GetGetGlobalConfigHandler(client))

	v1.Use(tenancyMiddlewareProvider(client, pathTenantProvider))

	v1.POST(v1receiverPath, GetReceiverPostHandler(client))
	v1.GET(v1receiverPath, GetGetReceiversHandler(client))

	v1.DELETE(v1receiverNamePath, GetDeleteReceiverHandler(client, receiverNamePathProvider))
	v1.PUT(v1receiverNamePath, GetUpdateReceiverHandler(client, receiverNamePathProvider))
	v1.GET(v1receiverNamePath, GetGetReceiversHandler(client))

	v1.POST(v1routePath, GetUpdateRouteHandler(client))
	v1.GET(v1routePath, GetGetRouteHandler(client))

}

func statusHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Alertmanager Config server")
}

type paramProvider func(c echo.Context) string

// For v0 tenant_id field in path
var pathTenantProvider = func(c echo.Context) string {
	return c.Param(tenantIDParam)
}

var v0receiverNameQueryProvider = func(c echo.Context) string {
	return c.QueryParam(v0receiverNameQueryParam)
}

var receiverNamePathProvider = func(c echo.Context) string {
	return c.Param(receiverNameParam)
}

// Returns middleware func to check for tenant_id dependent on tenancy of the client
func tenancyMiddlewareProvider(client client.AlertmanagerClient, getTenantID paramProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			providedTenantID := getTenantID(c)
			if client.Tenancy() != nil && providedTenantID == "" {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Must provide %s parameter", tenantIDParam))
			}
			c.Set(tenantIDParam, providedTenantID)
			return next(c)
		}
	}
}

// GetReceiverPostHandler returns a handler function that creates a new
// receiver and then reloads alertmanager
func GetReceiverPostHandler(client client.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		receiver, err := decodeReceiverPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		tenantID := c.Get(tenantIDParam).(string)

		err = client.CreateReceiver(tenantID, receiver)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.ReloadAlertmanager()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

// GetGetReceiversHandler returns a handler function to retrieve receivers for
// a filePrefix
func GetGetReceiversHandler(client client.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)
		receiverName := c.Param(receiverNameParam)

		recs, err := client.GetReceivers(tenantID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if receiverName != "" {
			for _, rec := range recs {
				if rec.Name == receiverName {
					return c.JSON(http.StatusOK, rec)
				}
			}
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Receiver %s not found", receiverName))
		}
		return c.JSON(http.StatusOK, recs)
	}
}

// GetUpdateReceiverHandler returns a handler function to update a receivers
func GetUpdateReceiverHandler(client client.AlertmanagerClient, getReceiverName paramProvider) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)
		receiverName := getReceiverName(c)

		newReceiver, err := decodeReceiverPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.UpdateReceiver(tenantID, receiverName, &newReceiver)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.ReloadAlertmanager()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetDeleteReceiverHandler(client client.AlertmanagerClient, getReceiverName paramProvider) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)

		err := client.DeleteReceiver(tenantID, getReceiverName(c))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.ReloadAlertmanager()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetGetRouteHandler(client client.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)

		route, err := client.GetRoute(tenantID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, receivers.NewRouteJSONWrapper(*route))
	}
}

func GetUpdateRouteHandler(client client.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)

		newRoute, err := decodeRoutePostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = client.ModifyTenantRoute(tenantID, &newRoute)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.ReloadAlertmanager()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetUpdateGlobalConfigHandler(client client.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		newGlobalConfig, err := decodeGlobalConfigPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = client.SetGlobalConfig(newGlobalConfig)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.ReloadAlertmanager()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetGetGlobalConfigHandler(client client.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		globalConf, err := client.GetGlobalConfig()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, globalConf)
	}
}

func decodeGlobalConfigPostRequest(c echo.Context) (config.GlobalConfig, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return config.GlobalConfig{}, fmt.Errorf("error reading request body: %v", err)
	}
	globalConfig := config.GlobalConfig{}
	err = json.Unmarshal(body, &globalConfig)
	if err != nil {
		return config.GlobalConfig{}, fmt.Errorf("error unmarshalling payload: %v", err)
	}
	return globalConfig, nil
}

func decodeReceiverPostRequest(c echo.Context) (receivers.Receiver, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return receivers.Receiver{}, fmt.Errorf("error reading request body: %v", err)
	}
	receiver := receivers.Receiver{}
	err = json.Unmarshal(body, &receiver)
	if err != nil {
		return receivers.Receiver{}, fmt.Errorf("error unmarshalling payload: %v", err)
	}
	return receiver, nil
}

func decodeRoutePostRequest(c echo.Context) (amconfig.Route, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return amconfig.Route{}, fmt.Errorf("error reading request body: %v", err)
	}
	route := amconfig.Route{}
	err = json.Unmarshal(body, &route)
	if err != nil {
		// Try decoding into a JSON-compatible struct
		wrappedRoute := receivers.RouteJSONWrapper{}
		err = json.Unmarshal(body, &wrappedRoute)
		if err != nil {
			return amconfig.Route{}, fmt.Errorf("error unmarshalling route: %v", err)
		}
		unwrappedRoute, err := wrappedRoute.ToPrometheusConfig()
		if err != nil {
			return amconfig.Route{}, fmt.Errorf("error handling route: %v", err)
		}
		return unwrappedRoute, nil
	}
	return route, nil
}
