/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"

	"magma/orc8r/cloud/go/datastore"
	merrors "magma/orc8r/cloud/go/errors"
	models2 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"

	"github.com/go-openapi/swag"
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
	TailGatewayLogs       = CommandRoot + "/tail_logs"

	CommandRootV1           = handlers.ManageGatewayPath + "/command"
	RebootGatewayV1         = CommandRootV1 + "/reboot"
	RestartServicesV1       = CommandRootV1 + "/restart_services"
	GatewayPingV1           = CommandRootV1 + "/ping"
	GatewayGenericCommandV1 = CommandRootV1 + "/generic"
	TailGatewayLogsV1       = CommandRootV1 + "/tail_logs"
)

func getListGateways(factory view_factory.FullGatewayViewFactory) func(echo.Context) error {
	return func(c echo.Context) error {
		fields := c.QueryParam("view")
		if fields == "full" {
			return ListFullGatewayViews(c, factory)
		}
		return listGateways(c)
	}
}

func listGateways(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayIDs, err := configurator.ListEntityKeys(networkID, orc8r.MagmadGatewayType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// Return a deterministic ordering of IDs
	sort.Strings(gatewayIDs)
	return c.JSON(http.StatusOK, gatewayIDs)
}

func createGateway(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	record := &models.GatewayDevice{}
	if err := c.Bind(record); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := record.ValidateModel(); err != nil {
		return obsidian.HttpError(
			fmt.Errorf("Invalid Gateway Record, Error: %s", err),
			http.StatusBadRequest)
	}

	gatewayID := c.QueryParam("requested_id")
	if len(gatewayID) > 0 {
		r, _ := regexp.Compile("^[a-zA-Z_][0-9a-zA-Z_-]+$")
		if !r.MatchString(gatewayID) {
			return obsidian.HttpError(
				fmt.Errorf("Gateway ID '%s' is not allowed. Gateway ID can only contain "+
					"alphanumeric characters and underscore, and should start with a letter or underscore.", gatewayID),
				http.StatusBadRequest,
			)
		}
	} else {
		gatewayID = record.HardwareID
	}

	if device.DoesDeviceExist(networkID, orc8r.AccessGatewayRecordType, record.HardwareID) {
		return fmt.Errorf("Hwid is already registered %s", record.HardwareID)
	}
	// write into device
	err := device.RegisterDevice(networkID, orc8r.AccessGatewayRecordType, record.HardwareID, record)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	// write into configurator
	gwEntity := configurator.NetworkEntity{
		Type:       orc8r.MagmadGatewayType,
		Key:        gatewayID,
		PhysicalID: record.HardwareID,
	}
	_, err = configurator.CreateEntity(networkID, gwEntity)
	if err != nil {
		derr := device.DeleteDevice(networkID, orc8r.AccessGatewayRecordType, record.HardwareID)
		if derr != nil {
			return obsidian.HttpError(
				fmt.Errorf("Failed to create gateway entity: %v, failed to delete device entity: %v", err, derr),
				http.StatusInternalServerError)
		}
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, gatewayID)
}

func getGateway(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayID, nerr := obsidian.GetLogicalGwId(c)
	if nerr != nil {
		return nerr
	}

	gatewayEntity, err := configurator.LoadEntity(networkID, orc8r.MagmadGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadMetadata: true})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	deviceEntity, err := device.GetDevice(networkID, orc8r.AccessGatewayRecordType, gatewayEntity.PhysicalID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	record := &models.GatewayDevice{}
	record.HardwareID = gatewayEntity.PhysicalID
	record.Key = deviceEntity.(*models.GatewayDevice).Key
	return c.JSON(http.StatusOK, record)
}

func updateGateway(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayID, gerr := obsidian.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	record := &models.GatewayDevice{}
	if berr := c.Bind(&record); berr != nil {
		return obsidian.HttpError(berr, http.StatusBadRequest)
	}
	if err := record.ValidateModel(); err != nil {
		return obsidian.HttpError(
			fmt.Errorf("Invalid Gateway Record, Error: %s", err),
			http.StatusBadRequest)
	}

	err := updateChallengeKey(networkID, gatewayID, record.Key)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func updateGatewayNameHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayID, gerr := obsidian.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}
	payload := models2.GatewayName("")
	if err := c.Bind(&payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	updateRequest := configurator.EntityUpdateCriteria{
		Key:     gatewayID,
		Type:    orc8r.MagmadGatewayType,
		NewName: swag.String(string(payload)),
	}
	_, err := configurator.UpdateEntities(networkID, []configurator.EntityUpdateCriteria{updateRequest})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func updateChallengeKey(networkID, gatewayID string, challengeKey *models.ChallengeKey) error {
	deviceID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return err
	}
	iRecord, err := device.GetDevice(networkID, orc8r.AccessGatewayRecordType, deviceID)
	if err != nil {
		return err
	}
	record, ok := iRecord.(*models.GatewayDevice)
	if !ok {
		return fmt.Errorf("Info stored in deviceID %s is not of type GatewayDevice", deviceID)
	}
	record.Key = challengeKey
	return device.UpdateDevice(networkID, orc8r.AccessGatewayRecordType, deviceID, record)
}

func deleteGateway(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayID, gerr := obsidian.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	err = device.DeleteDevice(networkID, orc8r.AccessGatewayRecordType, physicalID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	err = configurator.DeleteEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getLegacyOrV1GatewayID(c echo.Context) (string, *echo.HTTPError) {
	v1GwID := c.Param("gateway_id")
	if v1GwID != "" {
		return v1GwID, nil
	}

	return obsidian.GetLogicalGwId(c)
}

func rebootGateway(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := getLegacyOrV1GatewayID(c)
	if gerr != nil {
		return gerr
	}

	err := magmad.GatewayReboot(networkId, gatewayId)
	if err != nil {
		if datastore.IsErrNotFound(err) {
			return obsidian.HttpError(err, http.StatusNotFound)
		}
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func restartServices(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := getLegacyOrV1GatewayID(c)
	if gerr != nil {
		return gerr
	}

	var services []string
	err := c.Bind(&services)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	err = magmad.GatewayRestartServices(networkId, gatewayId, services)
	if err != nil {
		if datastore.IsErrNotFound(err) {
			return obsidian.HttpError(err, http.StatusNotFound)
		}
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func gatewayPing(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := getLegacyOrV1GatewayID(c)
	if gerr != nil {
		return gerr
	}

	pingRequest := magmad_models.PingRequest{}
	err := c.Bind(&pingRequest)
	response, err := magmad.GatewayPing(networkId, gatewayId, pingRequest.Packets, pingRequest.Hosts)
	if err != nil {
		if datastore.IsErrNotFound(err) {
			return obsidian.HttpError(err, http.StatusNotFound)
		}
		return obsidian.HttpError(err, http.StatusInternalServerError)
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
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := getLegacyOrV1GatewayID(c)
	if gerr != nil {
		return gerr
	}

	request := magmad_models.GenericCommandParams{}
	err := c.Bind(&request)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	params, err := models2.JSONMapToProtobufStruct(request.Params)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	genericCommandParams := protos.GenericCommandParams{
		Command: *request.Command,
		Params:  params,
	}

	response, err := magmad.GatewayGenericCommand(networkId, gatewayId, &genericCommandParams)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	resp, err := models2.ProtobufStructToJSONMap(response.Response)
	genericCommandResponse := magmad_models.GenericCommandResponse{
		Response: resp,
	}
	return c.JSON(http.StatusOK, genericCommandResponse)
}

func tailGatewayLogs(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := getLegacyOrV1GatewayID(c)
	if gerr != nil {
		return gerr
	}

	request := magmad_models.TailLogsRequest{}
	err := c.Bind(&request)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	stream, err := magmad.TailGatewayLogs(networkId, gatewayId, request.Service)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	go func() {
		<-c.Request().Context().Done()
	}()
	// https://echo.labstack.com/cookbook/streaming-response
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	c.Response().Header().Set(echo.HeaderXContentTypeOptions, "nosniff")
	c.Response().WriteHeader(http.StatusOK)
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		if _, err := c.Response().Write([]byte(line.Line)); err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		c.Response().Flush()
	}

	return c.NoContent(http.StatusNoContent)
}

// we need to fill in tbe tier ID of the legacy config struct with what the
// configurator client API returned as the parent assocs of the gateway
func getGatewayConfig(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := obsidian.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	ent, err := configurator.LoadEntity(networkId, orc8r.MagmadGatewayType, gatewayId, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsToThis: true})
	if err == merrors.ErrNotFound {
		return echo.ErrNotFound
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// ignore the error since we'll just be setting tier ID to the empty
	// string if no such assoc exists
	tierTK, _ := ent.GetFirstParentOfType(orc8r.UpgradeTierEntityType)
	cfg := ent.Config.(*models.MagmadGatewayConfigs)

	retConfig := &magmad_models.MagmadGatewayConfig{
		AutoupgradeEnabled:      swag.BoolValue(cfg.AutoupgradeEnabled),
		AutoupgradePollInterval: cfg.AutoupgradePollInterval,
		CheckinInterval:         int32(cfg.CheckinInterval),
		CheckinTimeout:          int32(cfg.CheckinTimeout),
		DynamicServices:         cfg.DynamicServices,
		FeatureFlags:            cfg.FeatureFlags,
		Tier:                    tierTK.Key,
	}
	return c.JSON(http.StatusOK, retConfig)
}
