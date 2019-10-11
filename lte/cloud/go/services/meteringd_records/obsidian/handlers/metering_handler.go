/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers implements all meteringd related REST APIs
package handlers

import (
	"net/http"

	"magma/lte/cloud/go/plugin/handlers"
	"magma/lte/cloud/go/services/meteringd_records"
	"magma/lte/cloud/go/services/meteringd_records/obsidian/models"
	subscriber_handlers "magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian"

	"github.com/labstack/echo"
)

const (
	SubscriberFlowsPath = subscriber_handlers.SubscriberManagePath + "/flow_records"

	SubscriberFlowsPathV1 = handlers.ManageSubscriberPath + "/flow_records"
)

func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{Path: SubscriberFlowsPath, Methods: obsidian.GET, HandlerFunc: listSubscriberFlowRecordsHandler},

		{Path: SubscriberFlowsPathV1, Methods: obsidian.GET, HandlerFunc: listSubscriberFlowRecordsHandler},
	}
}

func listSubscriberFlowRecordsHandler(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubscriberIDs(c)
	if nerr != nil {
		return nerr
	}

	flows, err := meteringd_records.ListSubscriberRecords(networkID, subscriberID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := make([]*models.FlowRecord, 0, len(flows))
	for _, flow := range flows {
		flowRecordModel := (&models.FlowRecord{}).FromProto(flow)
		ret = append(ret, flowRecordModel)
	}
	return c.JSON(http.StatusOK, ret)
}

func getNetworkAndSubscriberIDs(c echo.Context) (string, string, *echo.HTTPError) {
	ids, err := obsidian.GetParamValues(c, "network_id", "subscriber_id")
	if err != nil {
		return "", "", err
	}
	return ids[0], ids[1], nil
}
