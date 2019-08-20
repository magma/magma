/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers implements all meteringd related REST APIs
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"magma/lte/cloud/go/services/meteringd_records"
	"magma/lte/cloud/go/services/meteringd_records/obsidian/models"
	subscriber_handlers "magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	subscriber_models "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	SubscriberFlowsPath = subscriber_handlers.SubscriberManagePath + "/flow_records"
	ListFlowsPath       = obsidian.NetworksRoot + "/:network_id/flow_records"
	FlowDetailsPath     = ListFlowsPath + "/:flow_record_id"
)

func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{Path: SubscriberFlowsPath, Methods: obsidian.GET, HandlerFunc: listSubscriberFlowRecordsHandler},
		{Path: FlowDetailsPath, Methods: obsidian.GET, HandlerFunc: getFlowRecordHandler},
	}
}

// REST Handler to list subscriber flow records, no payload expected
func listSubscriberFlowRecordsHandler(c echo.Context) error {
	// Get networkid from context
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return obsidian.HttpError(errors.New("Invalid/missing network ID"), http.StatusBadRequest)
	}

	// Get subscriberId from context
	subscriberId, serr := getSubscriberId(c)
	if serr != nil {
		return obsidian.HttpError(errors.New("Invalid/missing subscriber ID"), http.StatusBadRequest)
	}

	flows, err := meteringd_records.ListSubscriberRecords(networkId, subscriberId)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// Convert to swagger model
	ret := make([]*models.FlowRecord, 0, len(flows))
	for _, flow := range flows {
		flowRecordModel := &models.FlowRecord{}
		err := flowRecordModel.FromProto(flow)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		ret = append(ret, flowRecordModel)
	}
	return c.JSON(http.StatusOK, ret)
}

// REST Handler to get flow record info, no payload expected
func getFlowRecordHandler(c echo.Context) error {
	// Get networkId and flow record id from REST context
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	flowRecordId := c.Param("flow_record_id")
	if flowRecordId == "" {
		return obsidian.HttpError(fmt.Errorf("Invalid/Missing Flow Record ID"), http.StatusBadRequest)
	}

	// Call meteringd_records service
	data, err := meteringd_records.GetRecord(networkId, flowRecordId)
	if err != nil {
		return obsidian.HttpError(err, http.StatusNotFound)
	}

	// Create swagger model for response
	var flowRecord models.FlowRecord
	if err = flowRecord.FromProto(data); err != nil {
		glog.Errorf("Error converting flow record model: %s", err)
		return obsidian.HttpError(err)
	}
	return c.JSON(http.StatusOK, flowRecord)
}

func getSubscriberId(c echo.Context) (string, *echo.HTTPError) {
	sidstr := c.Param("subscriber_id")
	err := (*subscriber_models.SubscriberID)(&sidstr).Verify()
	if err != nil {
		return sidstr, obsidian.HttpError(
			fmt.Errorf("Invalid subscriber ID %s: %s", sidstr, err),
			http.StatusBadRequest,
		)
	}
	return sidstr, nil
}
