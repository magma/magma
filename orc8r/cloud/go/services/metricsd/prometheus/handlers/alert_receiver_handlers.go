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
	"net/http"
	neturl "net/url"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/receivers"

	"github.com/labstack/echo"
	"github.com/prometheus/alertmanager/config"
)

const (
	ReceiverNamePathParam  = "receiver"
	ReceiverNameQueryParam = "receiver"
)

func GetConfigureAlertReceiverHandler(configManagerURL string) func(c echo.Context) error {
	return getHandlerWithReceiverFunc(configManagerURL, configureAlertReceiver)
}

func GetRetrieveAlertReceiverHandler(configManagerURL string) func(c echo.Context) error {
	return getHandlerWithReceiverFunc(configManagerURL, retrieveAlertReceivers)
}

func GetUpdateAlertReceiverHandler(configManagerURL string) func(c echo.Context) error {
	return getHandlerWithReceiverFunc(configManagerURL, updateAlertReceiver)
}

func GetDeleteAlertReceiverHandler(configManagerURL string) func(c echo.Context) error {
	return getHandlerWithReceiverFunc(configManagerURL, deleteAlertReceiver)
}

// getHandlerWithReceiverFunc returns an echo HandlerFunc that checks the
// networkID and runs the given handlerImplFunc that communicates with the
// alertmanager config service
func getHandlerWithReceiverFunc(configManagerURL string, handlerImplFunc func(echo.Context, string) error) func(echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := makeNetworkReceiverPath(configManagerURL, networkID)
		return handlerImplFunc(c, url)
	}
}

func GetRetrieveAlertRouteHandler(configManagerURL string) func(c echo.Context) error {
	return getHandlerWithRouteFunc(configManagerURL, retrieveAlertRoute)
}

func GetUpdateAlertRouteHandler(configManagerURL string) func(c echo.Context) error {
	return getHandlerWithRouteFunc(configManagerURL, updateAlertRoute)
}

// getHandlerWithRouteFunc returns an echo HandlerFunc that checks the
// networkID and runs the given handlerImplFunc that communicates with the
// alertmanager config service for routing trees
func getHandlerWithRouteFunc(configManagerURL string, handlerImplFunc func(echo.Context, string) error) func(echo.Context) error {
	return func(c echo.Context) error {
		networkID, nerr := obsidian.GetNetworkId(c)
		if nerr != nil {
			return nerr
		}
		url := makeNetworkRoutePath(configManagerURL, networkID)
		return handlerImplFunc(c, url)
	}
}

func configureAlertReceiver(c echo.Context, url string) error {
	receiver, err := buildReceiverFromContext(c)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	sendErr := sendConfig(receiver, url, http.MethodPost)
	if sendErr != nil {
		return obsidian.HttpError(fmt.Errorf("%s", sendErr.Message), sendErr.Code)
	}
	return c.NoContent(http.StatusOK)
}

func retrieveAlertReceivers(c echo.Context, url string) error {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return obsidian.HttpError(fmt.Errorf("error reading receivers: %v", body.Message), resp.StatusCode)
	}
	var recs []receivers.Receiver
	err = json.NewDecoder(resp.Body).Decode(&recs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("error decoding server response %v", err))
	}
	return c.JSON(http.StatusOK, recs)
}

func updateAlertReceiver(c echo.Context, url string) error {
	receiver, err := buildReceiverFromContext(c)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	receiverName := c.Param(ReceiverNamePathParam)
	if receiverName == "" {
		return obsidian.HttpError(fmt.Errorf("receiver name not provided"), http.StatusBadRequest)
	}
	if receiverName != receiver.Name {
		return obsidian.HttpError(fmt.Errorf("new receiver configuration must have same name"), http.StatusBadRequest)
	}
	url += fmt.Sprintf("/%s", neturl.PathEscape(receiverName))

	sendErr := sendConfig(receiver, url, http.MethodPut)
	if sendErr != nil {
		return obsidian.HttpError(sendErr, sendErr.Code)
	}
	return c.NoContent(http.StatusOK)
}

func deleteAlertReceiver(c echo.Context, url string) error {
	receiverName := c.QueryParam(ReceiverNameQueryParam)
	if receiverName == "" {
		return obsidian.HttpError(fmt.Errorf("receiver name not provided"), http.StatusBadRequest)
	}
	url += fmt.Sprintf("?%s=%s", ReceiverNameQueryParam, neturl.QueryEscape(receiverName))

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	resp, err := client.Do(req)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return obsidian.HttpError(fmt.Errorf("error deleting receiver: %v", body.Message), resp.StatusCode)
	}
	return c.NoContent(http.StatusOK)
}

func retrieveAlertRoute(c echo.Context, url string) error {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body echo.HTTPError
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return obsidian.HttpError(fmt.Errorf("error reading alerting route: %v", body.Message), resp.StatusCode)
	}
	var route receivers.RouteJSONWrapper
	err = json.NewDecoder(resp.Body).Decode(&route)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("error decoding server response %v", err), http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, route)
}

func updateAlertRoute(c echo.Context, url string) error {
	route, err := buildRouteFromContext(c)
	if err != nil {
		return obsidian.HttpError(fmt.Errorf("invalid route specification: %v\n", err), http.StatusBadRequest)
	}

	sendErr := sendConfig(route, url, http.MethodPost)
	if sendErr != nil {
		return obsidian.HttpError(fmt.Errorf("error updating alert route: %v", sendErr.Message), sendErr.Code)
	}
	return c.NoContent(http.StatusOK)
}

func buildReceiverFromContext(c echo.Context) (receivers.Receiver, error) {
	wrapper := receivers.Receiver{}
	err := json.NewDecoder(c.Request().Body).Decode(&wrapper)
	if err != nil {
		return receivers.Receiver{}, err
	}
	return wrapper, nil
}

func buildRouteFromContext(c echo.Context) (config.Route, error) {
	jsonRoute := receivers.RouteJSONWrapper{}
	err := json.NewDecoder(c.Request().Body).Decode(&jsonRoute)
	if err != nil {
		return config.Route{}, err
	}
	return jsonRoute.ToPrometheusConfig()
}

func makeNetworkReceiverPath(configManagerURL, networkID string) string {
	return configManagerURL + "/" + networkID + "/receiver"
}

func makeNetworkRoutePath(configManagerURL, networkID string) string {
	return configManagerURL + "/" + networkID + "/receiver/route"
}
