/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers for subscriberdb requests
package handlers

import (
	"fmt"
	"net/http"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian/handlers"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	SubscriberdRootPath  = handlers.NETWORKS_ROOT + "/:network_id/subscribers"
	SubscriberManagePath = SubscriberdRootPath + "/:subscriber_id"
)

// GetObsidianHandlers returns all obsidian handlers for subscriberdb
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		{Path: SubscriberdRootPath, Methods: handlers.POST, HandlerFunc: addSubscriberHandler},
		{Path: SubscriberManagePath, Methods: handlers.POST, HandlerFunc: addSubscriberHandler},
		{Path: SubscriberdRootPath, Methods: handlers.GET, HandlerFunc: listSubscribersHandler},
		{Path: SubscriberManagePath, Methods: handlers.GET, HandlerFunc: getSubscriberHandler},
		{Path: SubscriberdRootPath, Methods: handlers.PUT, HandlerFunc: updateSubscriberHandler},
		{Path: SubscriberManagePath, Methods: handlers.PUT, HandlerFunc: updateSubscriberHandler},
		{Path: SubscriberManagePath, Methods: handlers.DELETE, HandlerFunc: deleteSubscriberHandler},
	}
}

// REST Handler to add a new subscriber, expects subscriber data as payload
func addSubscriberHandler(c echo.Context) error {
	// Get swagger model
	sub := new(models.Subscriber)
	if err := c.Bind(sub); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := sub.Verify(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	// Get networkId and subscriberId from REST context
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberId := getSubscriberId(c)
	if len(subscriberId) != 0 { // SID is in URL
		sub.ID = models.SubscriberID(subscriberId)
	} else {
		subscriberId = string(sub.ID)
	}

	// Convert swagger model to proto
	sd := new(protos.SubscriberData)
	if err := sub.ToMconfig(sd); err != nil {
		return handlers.HttpError(err)
	}

	// Call subscriberdb service
	if err := subscriberdb.AddSubscriber(networkId, sd); err != nil {
		return handlers.HttpError(err, http.StatusConflict)
	}
	return c.JSON(http.StatusCreated, subscriberId)
}

// REST Handler to list all subscribers, no payload expected
func listSubscribersHandler(c echo.Context) error {
	// Get networkId from REST context
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	fields := getFieldsToFetch(c)
	if fields == "all" {
		return listSubscriberDataHandler(c, networkId)
	} else {
		// Default to listing IDs
		return listSubscriberIdsHandler(c, networkId)
	}
}

func listSubscriberIdsHandler(c echo.Context, networkId string) error {
	subs, err := subscriberdb.ListSubscribers(networkId)
	if err != nil {
		return handlers.HttpError(err)
	}
	sidset := make([]models.SubscriberID, len(subs))
	for i := range subs {
		sidset[i] = models.SubscriberID(subs[i])
	}
	return c.JSON(http.StatusOK, sidset)
}

func listSubscriberDataHandler(c echo.Context, networkId string) error {
	subsBySid, err := subscriberdb.GetAllSubscriberData(networkId)
	if err != nil {
		return handlers.HttpError(err)
	}

	ret := make(map[string]*models.Subscriber, len(subsBySid))
	for _, subProto := range subsBySid {
		subModel := &models.Subscriber{}
		if err = subModel.FromMconfig(subProto); err != nil {
			return handlers.HttpError(fmt.Errorf("Error converting subscriber model: %s", err))
		}
		ret[string(subModel.ID)] = subModel
	}
	return c.JSON(http.StatusOK, ret)
}

// REST Handler to get subscriber info, no payload expected
func getSubscriberHandler(c echo.Context) error {
	// Get networkId and subscriberId from REST context
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberId := getSubscriberId(c)
	if subscriberId == "" {
		return subscriberIdHttpErr()
	}

	// Call subscriberdb service
	data, err := subscriberdb.GetSubscriber(networkId, subscriberId)
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}

	// Create swagger model for response
	var sub models.Subscriber
	if err = sub.FromMconfig(data); err != nil {
		glog.Errorf("Error converting subscriber model: %s", err)
		return handlers.HttpError(err)
	}
	return c.JSON(http.StatusOK, sub)
}

// REST Handler to update a subscriber, expects subscriber data as payload
func updateSubscriberHandler(c echo.Context) error {
	// Get swagger model
	sub := new(models.Subscriber)
	if err := c.Bind(sub); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := sub.Verify(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	// Get networkId and subscriberId from REST context
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberId := getSubscriberId(c)
	if len(subscriberId) != 0 { // SID is in URL
		sub.ID = models.SubscriberID(subscriberId)
	} else {
		subscriberId = string(sub.ID)
	}

	// Convert swagger model to proto
	sd := new(protos.SubscriberData)
	if err := sub.ToMconfig(sd); err != nil {
		return handlers.HttpError(err)
	}

	// Call subscriberdb service
	if err := subscriberdb.UpdateSubscriber(networkId, sd); err != nil {
		return handlers.HttpError(err, http.StatusConflict)
	}
	return c.NoContent(http.StatusOK)
}

// REST handler to delete subscriber, no payload expected
func deleteSubscriberHandler(c echo.Context) error {
	// Get networkId and subscriberId from REST context
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberId := getSubscriberId(c)
	if subscriberId == "" {
		return subscriberIdHttpErr()
	}

	// Call subscriberdb service
	if err := subscriberdb.DeleteSubscriber(networkId, subscriberId); err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}
	return c.NoContent(http.StatusNoContent)
}

func getSubscriberId(c echo.Context) string {
	sidstr := c.Param("subscriber_id")
	if len(sidstr) > 0 {
		err := (*models.SubscriberID)(&sidstr).Verify()
		if err != nil {
			glog.Errorf("Invalid subscriber ID parameter: %s", sidstr)
			sidstr = ""
		}
	}
	return sidstr
}

func subscriberIdHttpErr() *echo.HTTPError {
	return handlers.HttpError(
		fmt.Errorf("Invalid/Missing Subscriber ID"),
		http.StatusBadRequest)
}

func getFieldsToFetch(c echo.Context) string {
	fields := c.QueryParam("fields")
	if len(fields) > 0 {
		return fields
	} else {
		// Default to only returning IDs (back-compat with old clients)
		return "id"
	}
}
