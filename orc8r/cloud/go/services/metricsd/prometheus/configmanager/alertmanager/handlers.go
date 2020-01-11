/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

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
	rootPath     = "/:file_prefix"
	ReceiverPath = rootPath + "/receiver"
	RoutePath    = ReceiverPath + "/route"

	ReceiverNamePathParam  = "receiver"
	ReceiverNameQueryParam = "receiver"
)

// GetReceiverPostHandler returns a handler function that creates a new
// receiver and then reloads alertmanager
func GetReceiverPostHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		receiver, err := decodeReceiverPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		err = client.CreateReceiver(getFilePrefix(c), receiver)
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
		getFilePrefix := getFilePrefix(c)
		recs, err := client.GetReceivers(getFilePrefix)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, recs)
	}
}

// GetUpdateReceiverHandler returns a handler function to update a receivers
func GetUpdateReceiverHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		getFilePrefix := getFilePrefix(c)
		newReceiver, err := decodeReceiverPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = client.UpdateReceiver(getFilePrefix, &newReceiver)
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
		getFilePrefix := getFilePrefix(c)
		receiverName := c.QueryParam(ReceiverNameQueryParam)

		err := client.DeleteReceiver(getFilePrefix, receiverName)
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
		getFilePrefix := getFilePrefix(c)
		route, err := client.GetRoute(getFilePrefix)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, receivers.NewRouteJSONWrapper(*route))
	}
}

func GetUpdateRouteHandler(client receivers.AlertmanagerClient) func(c echo.Context) error {
	return func(c echo.Context) error {
		getFilePrefix := getFilePrefix(c)
		newRoute, err := decodeRoutePostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = client.ModifyNetworkRoute(getFilePrefix, &newRoute)
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

func getFilePrefix(c echo.Context) string {
	return c.Param("file_prefix")
}
