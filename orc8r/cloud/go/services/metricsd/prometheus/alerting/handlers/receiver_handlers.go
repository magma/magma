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

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/prometheus/alertmanager/config"
)

const (
	rootPath     = "/:file_prefix"
	ReceiverPath = rootPath + "/receiver"
	RoutePath    = ReceiverPath + "/route"

	ReceiverNamePathParam  = "receiver"
	ReceiverNameQueryParam = "receiver"

	alertmanagerReloadPath = "/-/reload"
)

// GetReceiverPostHandler returns a handler function that creates a new
// receiver and then reloads alertmanager
func GetReceiverPostHandler(client receivers.AlertmanagerClient, alertmanagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		receiver, err := decodeReceiverPostRequest(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		err = client.CreateReceiver(getFilePrefix(c), receiver)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = reloadAlertmanager(alertmanagerURL)
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
func GetUpdateReceiverHandler(client receivers.AlertmanagerClient, alertmanagerURL string) func(c echo.Context) error {
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

		err = reloadAlertmanager(alertmanagerURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func GetDeleteReceiverHandler(client receivers.AlertmanagerClient, alertmanagerURL string) func(c echo.Context) error {
	return func(c echo.Context) error {
		getFilePrefix := getFilePrefix(c)
		receiverName := c.QueryParam(ReceiverNameQueryParam)

		err := client.DeleteReceiver(getFilePrefix, receiverName)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = reloadAlertmanager(alertmanagerURL)
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
		return c.JSON(http.StatusOK, route)
	}
}

func GetUpdateRouteHandler(client receivers.AlertmanagerClient, alertmanagerURL string) func(c echo.Context) error {
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

		err = reloadAlertmanager(alertmanagerURL)
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
		return config.Route{}, fmt.Errorf("error unmarshalling route: %v", err)
	}
	return route, nil
}

func reloadAlertmanager(url string) error {
	if url == "" {
		glog.Info("Not reloading alertmanager: No url given")
		return nil
	}
	resp, err := http.Post(fmt.Sprintf("http://%s%s", url, alertmanagerReloadPath), "text/plain", &bytes.Buffer{})
	if err != nil {
		return fmt.Errorf("error reloading alertmanager: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("code: %d error reloading alertmanager: %s", resp.StatusCode, msg)
	}
	return nil
}
