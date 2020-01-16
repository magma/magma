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

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/receivers"

	"github.com/labstack/echo"
	"github.com/prometheus/alertmanager/config"
)

const (
	v0rootPath   = "/:tenant_id"
	v0receiverPath = "/receiver"
	v0RoutePath = "/receiver/route"

	v1rootPath   = "/v1"
	v1receiverPath = "/receiver"
	v1routePath   = "/route"

	receiverNameParam  = "receiver"

	tenantIDParam = "tenant_id"
)

func RegisterBaseHandlers(e *echo.Echo) {
	e.GET("/", statusHandler)
}

func RegisterV0Handlers(e *echo.Echo, client receivers.AlertmanagerClient) {
	v0 := e.Group(v0rootPath)
	v0.Use(tenancyMiddlewareProvider(client, pathTenantProvider))

	v0.POST(v0receiverPath, GetReceiverPostHandler(client))
	v0.GET(v0receiverPath, GetGetReceiversHandler(client))
	v0.DELETE(v0receiverPath, GetDeleteReceiverHandler(client))
	v0.PUT(v0receiverPath+"/:"+receiverNameParam, GetUpdateReceiverHandler(client))

	v0.POST(v0RoutePath, GetUpdateRouteHandler(client))
	v0.GET(v0RoutePath, GetGetRouteHandler(client))
}

func RegisterV1Handlers(e *echo.Echo, client receivers.AlertmanagerClient) {
	v1 := e.Group(v1rootPath)
	v1.Use(tenancyMiddlewareProvider(client, queryTenantProvider))

	v1.POST(v1receiverPath, GetReceiverPostHandler(client))
	v1.GET(v1receiverPath, GetGetReceiversHandler(client))
	v1.DELETE(v1receiverPath, GetDeleteReceiverHandler(client))
	v1.PUT(v1receiverPath, GetUpdateReceiverHandler(client))

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

// V1 tenantID is a query parameter
var queryTenantProvider = func(c echo.Context) string {
	return c.QueryParam(tenantIDParam)
}

// Returns middleware func to check for tenant_id dependent on tenancy of the client
func tenancyMiddlewareProvider(client receivers.AlertmanagerClient, getTenantID paramProvider) echo.MiddlewareFunc {
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
func GetReceiverPostHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
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
func GetGetReceiversHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)

		recs, err := client.GetReceivers(tenantID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, recs)
	}
}

// GetUpdateReceiverHandler returns a handler function to update a receivers
func GetUpdateReceiverHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)

		newReceiver, err := decodeReceiverPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.UpdateReceiver(tenantID, &newReceiver)
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

func GetDeleteReceiverHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)

		err := client.DeleteReceiver(tenantID, c.QueryParam(receiverNameParam))
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

func GetGetRouteHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		tenantID := c.Get(tenantIDParam).(string)

		route, err := client.GetRoute(tenantID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, receivers.NewRouteJSONWrapper(*route))
	}
}

func GetUpdateRouteHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
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

func decodeRoutePostRequest(c echo.Context) (config.Route, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return config.Route{}, fmt.Errorf("error reading request body: %v", err)
	}
	route := config.Route{}
	err = json.Unmarshal(body, &route)
	if err != nil {
		// Try decoding into a JSON-compatible struct
		wrappedRoute := receivers.RouteJSONWrapper{}
		err = json.Unmarshal(body, &wrappedRoute)
		if err != nil {
			return config.Route{}, fmt.Errorf("error unmarshalling route: %v", err)
		}
		unwrappedRoute, err := wrappedRoute.ToPrometheusConfig()
		if err != nil {
			return config.Route{}, fmt.Errorf("error handling route: %v", err)
		}
		return unwrappedRoute, nil
	}
	return route, nil
}
