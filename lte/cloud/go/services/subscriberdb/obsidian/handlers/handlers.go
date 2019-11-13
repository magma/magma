/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"magma/lte/cloud/go/lte"
	models2 "magma/lte/cloud/go/plugin/models"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/thoas/go-funk"
)

const (
	SubscriberdRootPath  = obsidian.NetworksRoot + "/:network_id/subscribers"
	SubscriberManagePath = SubscriberdRootPath + "/:subscriber_id"
)

// GetObsidianHandlers returns all obsidian handlers for subscriberdb
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{Path: SubscriberdRootPath, Methods: obsidian.POST, HandlerFunc: createSubscriber},
		{Path: SubscriberManagePath, Methods: obsidian.POST, HandlerFunc: createSubscriber},
		{Path: SubscriberdRootPath, Methods: obsidian.GET, HandlerFunc: listSubscribers},
		{Path: SubscriberManagePath, Methods: obsidian.GET, HandlerFunc: getSubscriber},
		{Path: SubscriberdRootPath, Methods: obsidian.PUT, HandlerFunc: updateSubscriber},
		{Path: SubscriberManagePath, Methods: obsidian.PUT, HandlerFunc: updateSubscriber},
		{Path: SubscriberManagePath, Methods: obsidian.DELETE, HandlerFunc: deleteSubscriber},
	}
}

func createSubscriber(c echo.Context) error {
	// Get swagger model
	sub := new(models2.Subscriber)
	if err := c.Bind(sub); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := sub.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Get networkId and subscriberId from REST context
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberID := getSubscriberId(c)
	if len(subscriberID) != 0 {
		sub.ID = models2.SubscriberID(subscriberID)
	} else {
		subscriberID = string(models2.SubscriberID(sub.ID))
	}

	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:   lte.SubscriberEntityType,
		Key:    subscriberID,
		Config: sub.Lte,
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, subscriberID)
}

func listSubscribers(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ents, err := configurator.LoadAllEntitiesInNetwork(networkID, lte.SubscriberEntityType, getListSubscribersLoadCriteria(c))
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// if configs were loaded we'll return those, otherwise just the sids
	sids := make([]string, 0, len(ents))
	entConfs := make(map[string]*models2.Subscriber, len(ents))
	for _, ent := range ents {
		sids = append(sids, ent.Key)
		if ent.Config != nil {
			entConfs[ent.Key] = &models2.Subscriber{
				ID:  models2.SubscriberID(ent.Key),
				Lte: ent.Config.(*models2.LteSubscription),
			}
		}
	}

	if !funk.IsEmpty(entConfs) {
		return c.JSON(http.StatusOK, entConfs)
	}
	sort.Strings(sids)
	return c.JSON(http.StatusOK, sids)
}

func getSubscriber(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberID := c.Param("subscriber_id")
	if subscriberID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "subscriber ID is required")
	}

	ent, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, subscriberID, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	ret := &models2.Subscriber{ID: models2.SubscriberID(ent.Key), Lte: ent.Config.(*models2.LteSubscription)}
	return c.JSON(http.StatusOK, ret)
}

func updateSubscriber(c echo.Context) error {
	sub := new(models2.Subscriber)
	if err := c.Bind(sub); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := sub.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberID := getSubscriberId(c)
	if len(subscriberID) != 0 { // SID is in URL
		sub.ID = models2.SubscriberID(subscriberID)
	} else {
		subscriberID = string(sub.ID)
	}

	err := configurator.CreateOrUpdateEntityConfig(networkID, lte.SubscriberEntityType, subscriberID, sub.Lte)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func deleteSubscriber(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subscriberID := getSubscriberId(c)
	if subscriberID == "" {
		return subscriberIdHttpErr()
	}

	err := configurator.DeleteEntity(networkID, lte.SubscriberEntityType, subscriberID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
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
	return obsidian.HttpError(
		fmt.Errorf("Invalid/Missing Subscriber ID"),
		http.StatusBadRequest)
}
