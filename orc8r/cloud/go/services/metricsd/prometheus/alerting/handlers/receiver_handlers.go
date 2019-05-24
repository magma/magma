/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/receivers"

	"github.com/labstack/echo"
	"github.com/prometheus/alertmanager/config"
)

const (
	alertmanagerReloadPath = "/-/reload"
)

// GetReceiverPostHandler returns a handler function that creates a new
// receiver and then reloads alertmanager
func GetReceiverPostHandler(client *receivers.Client, alertmanagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		receiver, err := decodeReceiverPostResponse(c)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		}
		err = client.CreateReceiver(&receiver, getNetworkID(c))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		err = reloadAlertmanager(alertmanagerURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.NoContent(http.StatusOK)
	}
}

// GetGetReceiversHandler returns a handler function to retrieve receivers for
// a network
func GetGetReceiversHandler(client *receivers.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID := getNetworkID(c)
		recs, err := client.GetReceivers(networkID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, recs)
	}
}

func GetGetRouteHandler(client *receivers.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID := getNetworkID(c)
		route, err := client.GetRoute(networkID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, route)
	}
}

func GetUpdateRouteHandler(client *receivers.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		networkID := getNetworkID(c)
		newRoute, err := decodeRoutePostRequest(c)
		if err != nil {
			return err
		}
		err = client.ModifyNetworkRoute(&newRoute, networkID)
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}

func decodeReceiverPostResponse(c echo.Context) (receivers.Receiver, error) {
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
		return config.Route{}, fmt.Errorf("error unmarshalling route: %v", err)
	}
	return route, nil
}

func reloadAlertmanager(url string) error {
	resp, err := http.Post(fmt.Sprintf("http://%s%s", url, alertmanagerReloadPath), "text/plain", &bytes.Buffer{})
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("code: %d error reloading alertmanager: %v", resp.StatusCode, err)
	}
	return nil
}
