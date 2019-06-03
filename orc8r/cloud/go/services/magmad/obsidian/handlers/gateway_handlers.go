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
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_utils "magma/orc8r/cloud/go/services/configurator/obsidian/handler_utils"
	configuratorprotos "magma/orc8r/cloud/go/services/configurator/protos"
	"magma/orc8r/cloud/go/services/device"
	deviceprotos "magma/orc8r/cloud/go/services/device/protos"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	magmad_models "magma/orc8r/cloud/go/services/magmad/obsidian/models"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/golang/glog"
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

	err = multiplexGatewayCreateIntoDeviceAndConfigurator(networkId, gatewayId, swaggerRecord)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex write into configurator/device %v", err), http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, gatewayId)
}

func multiplexGatewayCreateIntoDeviceAndConfigurator(networkID, gatewayID string, gwRecord *magmad_models.AccessGatewayRecord) error {
	err := configurator_utils.CreateNetworkIfNotExists(networkID)
	if err != nil {
		return err
	}

	if device.DoesDeviceExist(networkID, device.GatewayInfoType, gwRecord.HwID.ID) {
		return fmt.Errorf("Hwid is already registered %s", gwRecord.HwID.ID)
	}
	// write into configurator
	gwEntity := &configuratorprotos.NetworkEntity{
		Name:       gwRecord.Name,
		Type:       configurator.GatewayEntityType,
		Id:         gatewayID,
		PhysicalId: gwRecord.HwID.ID,
	}
	_, err = configurator.CreateEntities(networkID, []*configuratorprotos.NetworkEntity{gwEntity})
	if err != nil {
		return err
	}

	// write into device
	return device.CreateOrUpdate(networkID, device.GatewayInfoType, gwRecord.HwID.ID, gwRecord)
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
	swaggerRecord, err := getSwaggerGWRecordFromMagmad(networkId, lid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, swaggerRecord)
}

func getSwaggerGWRecordFromMagmad(networkID, logicalID string) (*magmad_models.AccessGatewayRecord, error) {
	record, err := magmad.FindGatewayRecord(networkID, logicalID)
	if err != nil {
		return nil, handlers.HttpError(err, http.StatusNotFound)
	}
	swaggerRecord := magmad_models.AccessGatewayRecord{}
	err = swaggerRecord.FromMconfig(record)
	if err != nil {
		return nil, handlers.HttpError(err, http.StatusUnsupportedMediaType)
	}
	return &swaggerRecord, nil
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

	err = multiplexGatewayUpdateIntoDeviceAndConfigurator(networkId, lid, &swaggerRecord)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex update into configurator/device %v", err), http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func multiplexGatewayUpdateIntoDeviceAndConfigurator(networkID, gatewayID string, updateRecord *magmad_models.MutableGatewayRecord) error {
	entityExists, err := configurator.DoesEntityExist(networkID, configurator.GatewayEntityType, gatewayID)
	if err != nil {
		return err
	}
	if !entityExists {
		// fetch the existing gw record from magmad to get the HWID since it is needed for the device service
		storedRecord, err := getSwaggerGWRecordFromMagmad(networkID, gatewayID)
		if err != nil {
			return err
		}
		storedRecord.Name = updateRecord.Name
		storedRecord.Key = updateRecord.Key
		return multiplexGatewayCreateIntoDeviceAndConfigurator(networkID, gatewayID, storedRecord)
	}
	err = updateChallengeKey(networkID, gatewayID, updateRecord.Key)
	if err != nil {
		return err
	}
	return updateGatewayName(networkID, gatewayID, updateRecord.Name)
}

func updateChallengeKey(networkID, gatewayID string, challengeKey *magmad_models.ChallengeKey) error {
	deviceID, err := configurator.GetPhysicalIDOfEntity(networkID, configurator.GatewayEntityType, gatewayID)
	if err != nil {
		return err
	}
	iRecord, err := device.GetDevice(networkID, device.GatewayInfoType, deviceID)
	if err != nil {
		return err
	}
	record, ok := iRecord.(*magmad_models.AccessGatewayRecord)
	if !ok {
		return fmt.Errorf("Info stored in deviceID %s is not of type AccessGatewayRecord", deviceID)
	}
	record.Key = challengeKey
	return device.CreateOrUpdate(networkID, device.GatewayInfoType, deviceID, record)
}

func updateGatewayName(networkID, gatewayID, name string) error {
	updateRequest := &configuratorprotos.EntityUpdateCriteria{
		Key:     gatewayID,
		Type:    configurator.GatewayEntityType,
		NewName: configuratorprotos.GetStringWrapper(&name),
	}
	_, err := configurator.UpdateEntities(networkID, []*configuratorprotos.EntityUpdateCriteria{updateRequest})
	return err
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

	err = multiplexGatewayDeleteIntoDeviceAndConfigurator(networkId, lid)
	if err != nil {
		glog.Errorf("Failed to multiplex delete into configurator/device %v", err)
	}

	return c.NoContent(http.StatusNoContent)
}

func multiplexGatewayDeleteIntoDeviceAndConfigurator(networkID, gatewayID string) error {
	physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, configurator.GatewayEntityType, gatewayID)
	if err != nil {
		return err
	}
	err = device.DeleteDevices(networkID, []*deviceprotos.DeviceID{{DeviceID: physicalID, Type: device.GatewayInfoType}})
	if err != nil {
		return err
	}
	return configurator.DeleteEntities(networkID, []*configuratorprotos.EntityID{{Id: gatewayID, Type: configurator.GatewayEntityType}})
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

func tailGatewayLogs(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	gatewayId, gerr := handlers.GetLogicalGwId(c)
	if gerr != nil {
		return gerr
	}

	request := magmad_models.TailLogsRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	stream, err := magmad.TailGatewayLogs(networkId, gatewayId, request.Service)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
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
			return handlers.HttpError(err, http.StatusInternalServerError)
		}

		if _, err := c.Response().Write([]byte(line.Line)); err != nil {
			return handlers.HttpError(err, http.StatusInternalServerError)
		}
		c.Response().Flush()
	}

	return c.NoContent(http.StatusNoContent)
}
