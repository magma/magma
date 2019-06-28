/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"net/http"
	"sort"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/labstack/echo"
	"github.com/thoas/go-funk"
)

const (
	SubscriberdRootPath  = handlers.NETWORKS_ROOT + "/:network_id/subscribers"
	SubscriberManagePath = SubscriberdRootPath + "/:subscriber_id"
)

// GetObsidianHandlers returns all obsidian handlers for subscriberdb
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		{Path: SubscriberdRootPath, Methods: handlers.POST, HandlerFunc: addSubscriberHandler, MigratedHandlerFunc: createSubscriber, MultiplexAfterMigration: true},
		{Path: SubscriberManagePath, Methods: handlers.POST, HandlerFunc: addSubscriberHandler, MigratedHandlerFunc: createSubscriber, MultiplexAfterMigration: true},
		{Path: SubscriberdRootPath, Methods: handlers.GET, HandlerFunc: listSubscribersHandler, MigratedHandlerFunc: listSubscribers},
		{Path: SubscriberManagePath, Methods: handlers.GET, HandlerFunc: getSubscriberHandler, MigratedHandlerFunc: getSubscriber},
		{Path: SubscriberdRootPath, Methods: handlers.PUT, HandlerFunc: updateSubscriberHandler, MigratedHandlerFunc: updateSubscriber, MultiplexAfterMigration: true},
		{Path: SubscriberManagePath, Methods: handlers.PUT, HandlerFunc: updateSubscriberHandler, MigratedHandlerFunc: updateSubscriber, MultiplexAfterMigration: true},
		{Path: SubscriberManagePath, Methods: handlers.DELETE, HandlerFunc: deleteSubscriberHandler, MigratedHandlerFunc: deleteSubscriber, MultiplexAfterMigration: true},
	}
}

func createSubscriber(c echo.Context) error {
	// Get swagger model
	sub := new(models.Subscriber)
	if err := c.Bind(sub); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := sub.Verify(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	// Get networkId and subscriberId from REST context
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberID := getSubscriberId(c)
	if len(subscriberID) != 0 {
		sub.ID = models.SubscriberID(subscriberID)
	} else {
		subscriberID = string(sub.ID)
	}

	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:   lte.SubscriberEntityType,
		Key:    subscriberID,
		Config: sub,
	})
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, subscriberID)
}

func listSubscribers(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ents, err := configurator.LoadAllEntitiesInNetwork(networkID, lte.SubscriberEntityType, getListSubscribersLoadCriteria(c))
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	// if configs were loaded we'll return those, otherwise just the sids
	sids := make([]string, 0, len(ents))
	entConfs := make(map[string]*models.Subscriber, len(ents))
	for _, ent := range ents {
		sids = append(sids, ent.Key)
		if ent.Config != nil {
			entConfs[ent.Key] = ent.Config.(*models.Subscriber)
		}
	}

	if !funk.IsEmpty(entConfs) {
		return c.JSON(http.StatusOK, entConfs)
	}
	sort.Strings(sids)
	return c.JSON(http.StatusOK, sids)
}

func getSubscriber(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberID := c.Param("subscriber_id")
	if subscriberID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "subscriber ID is required")
	}

	ent, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, subscriberID, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, ent.Config.(*models.Subscriber))
}

func updateSubscriber(c echo.Context) error {
	sub := new(models.Subscriber)
	if err := c.Bind(sub); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := sub.Verify(); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberID := getSubscriberId(c)
	if len(subscriberID) != 0 { // SID is in URL
		sub.ID = models.SubscriberID(subscriberID)
	} else {
		subscriberID = string(sub.ID)
	}

	err := configurator.CreateOrUpdateEntityConfig(networkID, lte.SubscriberEntityType, subscriberID, sub)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func deleteSubscriber(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberID := getSubscriberId(c)
	if subscriberID == "" {
		return subscriberIdHttpErr()
	}

	err := configurator.DeleteEntity(networkID, lte.SubscriberEntityType, subscriberID)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getListSubscribersLoadCriteria(c echo.Context) configurator.EntityLoadCriteria {
	fields := c.QueryParam("fields")
	if fields == "all" {
		return configurator.EntityLoadCriteria{LoadConfig: true}
	}
	return configurator.EntityLoadCriteria{}
}
