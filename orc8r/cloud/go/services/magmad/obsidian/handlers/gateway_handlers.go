/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/labstack/echo"
)

const (
	RegisterAG  = ManageNetwork + "/gateways"
	ManageAG    = RegisterAG + "/:logical_ag_id"
	ConfigureAG = ManageAG + "/configs"

	CommandRoot           = ManageAG + "/command"
	RebootGateway         = CommandRoot + "/reboot"
	RestartServices       = CommandRoot + "/restart_services"
	GatewayPing           = CommandRoot + "/ping"
	GatewayGenericCommand = CommandRoot + "/generic"
)

func getListGatewaysHandler(factory view_factory.FullGatewayViewFactory) func(echo.Context) error {
	return func(c echo.Context) error {
		fields := c.QueryParam("view")
		if fields == "full" {
			return ListFullGatewayViews(c, factory)
		}
		return listGateways(c)
	}
}

func listGateways(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayIds, err := magmad.ListGateways(networkId)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	// Return a deterministic ordering of IDs
	sort.Strings(gatewayIds)
	return c.JSON(http.StatusOK, gatewayIds)
}

func registerGateway(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	swaggerRecord := &magmad_models.AccessGatewayRecord{}
	if err := c.Bind(swaggerRecord); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	if err := swaggerRecord.Verify(); err != nil {
		return handlers.HttpError(
			fmt.Errorf("Invalid Gateway Record, Error: %s", err),
			http.StatusBadRequest)
	}
	record, err := swaggerRecord.ToMconfig()
	if err != nil {
		return handlers.HttpError(err, http.StatusUnsupportedMediaType)
	}

	var gatewayId string
	requestedId := c.QueryParam("requested_id")
	if len(requestedId) > 0 {
		r, _ := regexp.Compile("^[a-zA-Z_][0-9a-zA-Z_-]+$")
		if !r.MatchString(requestedId) {
			return handlers.HttpError(
				fmt.Errorf("Gateway ID '%s' is not allowed. Gateway ID can only contain "+
					"alphanumeric characters and underscore, and should start with a letter or underscore.", requestedId),
				http.StatusBadRequest,
			)
		}
		gatewayId, err = magmad.RegisterGatewayWithId(networkId, record, requestedId)
	} else {
		gatewayId, err = magmad.RegisterGateway(networkId, record)
	}

	if err != nil {
		return handlers.HttpError(err, http.StatusConflict)
	}
	return c.JSON(http.StatusCreated, gatewayId)
}

func getGateway(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	lid, gerr := handlers.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	record, err := magmad.FindGatewayRecord(networkId, lid)
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}
	swaggerRecord := magmad_models.AccessGatewayRecord{}
	err = swaggerRecord.FromMconfig(record)
	if err != nil {
		return handlers.HttpError(err, http.StatusUnsupportedMediaType)
	}
	return c.JSON(http.StatusOK, &swaggerRecord)
}

func updateGateway(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	lid, gerr := handlers.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	swaggerRecord := magmad_models.MutableGatewayRecord{}
	if berr := c.Bind(&swaggerRecord); berr != nil {
		return handlers.HttpError(berr, http.StatusBadRequest)
	}
	if err := swaggerRecord.Verify(); err != nil {
		return handlers.HttpError(
			fmt.Errorf("Invalid Gateway Record, Error: %s", err),
			http.StatusBadRequest)
	}
	record := magmadprotos.AccessGatewayRecord{}
	berr := swaggerRecord.ToMconfig(&record)
	if berr != nil {
		return handlers.HttpError(berr, http.StatusUnsupportedMediaType)
	}
	err := magmad.UpdateGatewayRecord(networkId, lid, &record)
	if err != nil {
		return handlers.HttpError(err, http.StatusConflict)
	}
	return c.NoContent(http.StatusOK)
}

func deleteGateway(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	lid, gerr := handlers.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	err := magmad.RemoveGateway(networkId, lid)
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}

func rebootGateway(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := handlers.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	err := magmad.GatewayReboot(networkId, gatewayId)
	if err != nil {
		if datastore.IsErrNotFound(err) {
			return handlers.HttpError(err, http.StatusNotFound)
		}
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func restartServices(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := handlers.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	var services []string
	err := c.Bind(&services)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	err = magmad.GatewayRestartServices(networkId, gatewayId, services)
	if err != nil {
		if datastore.IsErrNotFound(err) {
			return handlers.HttpError(err, http.StatusNotFound)
		}
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func gatewayPing(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := handlers.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	pingRequest := magmad_models.PingRequest{}
	err := c.Bind(&pingRequest)
	response, err := magmad.GatewayPing(networkId, gatewayId, pingRequest.Packets, pingRequest.Hosts)
	if err != nil {
		if datastore.IsErrNotFound(err) {
			return handlers.HttpError(err, http.StatusNotFound)
		}
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	var pingResponse magmad_models.PingResponse
	for _, ping := range response.Pings {
		pingResult := &magmad_models.PingResult{
			HostOrIP:           &ping.HostOrIp,
			NumPackets:         &ping.NumPackets,
			PacketsTransmitted: ping.PacketsTransmitted,
			PacketsReceived:    ping.PacketsReceived,
			AvgResponseMs:      ping.AvgResponseMs,
			Error:              ping.Error,
		}
		pingResponse.Pings = append(pingResponse.Pings, pingResult)
	}
	return c.JSON(http.StatusOK, &pingResponse)
}

func gatewayGenericCommand(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := handlers.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	request := magmad_models.GenericCommandParams{}
	err := c.Bind(&request)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	params, err := magmad_models.JSONMapToProtobufStruct(request.Params)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	genericCommandParams := protos.GenericCommandParams{
		Command: *request.Command,
		Params:  params,
	}

	response, err := magmad.GatewayGenericCommand(networkId, gatewayId, &genericCommandParams)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	resp, err := magmad_models.ProtobufStructToJSONMap(response.Response)
	genericCommandResponse := magmad_models.GenericCommandResponse{
		Response: resp,
	}
	return c.JSON(http.StatusOK, genericCommandResponse)
}
