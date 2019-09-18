/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	cfgObsidian "magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	magmad_handlers "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/storage"
	"orc8r/devmand/cloud/go/devmand"
	"orc8r/devmand/cloud/go/services/devmand/obsidian/models"

	"github.com/labstack/echo"
)

const (
	// ConfigKey is the devmand config key
	ConfigKey = "devmand"
	// GatewayConfigPath is the endpoint path to configure a gateway
	GatewayConfigPath = magmad_handlers.ConfigureAG + "/" + ConfigKey
	// DeviceRootPath is the root path for devices
	DeviceRootPath = magmad_handlers.ConfigureNetwork + "/devices"
	// DeviceConfigPath is the path to configure devices
	DeviceConfigPath = DeviceRootPath + "/:device_id"
)

// GetObsidianHandlers returns all obsidian handlers for Devmand
func GetObsidianHandlers() []obsidian.Handler {
	createGatewayConfigHandler := cfgObsidian.GetCreateGatewayConfigHandler(GatewayConfigPath, devmand.DevmandGatewayType, &models.GatewayDevmandConfigs{})
	updateGatewayConfigHandler := cfgObsidian.GetUpdateGatewayConfigHandler(GatewayConfigPath, devmand.DevmandGatewayType, &models.GatewayDevmandConfigs{})
	readGatewayConfigHandler := cfgObsidian.GetReadGatewayConfigHandler(GatewayConfigPath, devmand.DevmandGatewayType, &models.GatewayDevmandConfigs{})

	// overwrite HandlerFunc for devmand gateway
	createGatewayConfigHandler.HandlerFunc = createGatewayConfig
	updateGatewayConfigHandler.HandlerFunc = updateGatewayConfig
	readGatewayConfigHandler.HandlerFunc = readGatewayConfig

	return []obsidian.Handler{
		cfgObsidian.GetCreateConfigHandler(DeviceRootPath, devmand.DeviceType, newDeviceID, &models.ManagedDevice{}),
		cfgObsidian.GetReadConfigHandler(DeviceConfigPath, devmand.DeviceType, getDeviceID, &models.ManagedDevice{}),
		cfgObsidian.GetUpdateConfigHandler(DeviceConfigPath, devmand.DeviceType, getDeviceID, &models.ManagedDevice{}),
		cfgObsidian.GetDeleteConfigHandler(DeviceConfigPath, devmand.DeviceType, getDeviceID),
		cfgObsidian.GetReadAllKeysConfigHandler(DeviceRootPath, devmand.DeviceType),
		createGatewayConfigHandler,
		updateGatewayConfigHandler,
		readGatewayConfigHandler,
		cfgObsidian.GetDeleteGatewayConfigHandler(GatewayConfigPath, devmand.DevmandGatewayType),
	}
}

func createGatewayConfig(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGWID(c)
	if nerr != nil {
		return nerr
	}
	iConfig, nerr := cfgObsidian.GetConfigAndValidate(c, &models.GatewayDevmandConfigs{})
	if nerr != nil {
		return nerr
	}
	config := iConfig.(*models.GatewayDevmandConfigs)

	// create devmand gateway entity
	associations := getDeviceTKs(config.ManagedDevices)
	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:         devmand.DevmandGatewayType,
		Key:          gatewayID,
		Config:       iConfig,
		Associations: associations,
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	_, err = configurator.UpdateEntity(networkID, configurator.EntityUpdateCriteria{
		Type:              orc8r.MagmadGatewayType,
		Key:               gatewayID,
		AssociationsToSet: []storage.TypeAndKey{{Type: devmand.DevmandGatewayType, Key: gatewayID}},
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, gatewayID)
}

func updateGatewayConfig(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGWID(c)
	if nerr != nil {
		return nerr
	}
	iConfig, nerr := cfgObsidian.GetConfigAndValidate(c, &models.GatewayDevmandConfigs{})
	if nerr != nil {
		return nerr
	}
	config := iConfig.(*models.GatewayDevmandConfigs)

	associationsToSet := getDeviceTKs(config.ManagedDevices)
	associationsToDelete := []storage.TypeAndKey{}
	if len(config.ManagedDevices) == 0 {
		entity, err := configurator.LoadEntity(
			networkID,
			devmand.DevmandGatewayType,
			gatewayID,
			configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		associationsToDelete = entity.Associations
	}

	_, err := configurator.UpdateEntity(networkID, configurator.EntityUpdateCriteria{
		Type:                 devmand.DevmandGatewayType,
		Key:                  gatewayID,
		NewConfig:            iConfig,
		AssociationsToSet:    associationsToSet,
		AssociationsToDelete: associationsToDelete,
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func readGatewayConfig(c echo.Context) error {
	networkID, gatewayID, nerr := obsidian.GetNetworkAndGWID(c)
	if nerr != nil {
		return nerr
	}
	entity, err := configurator.LoadEntity(networkID, devmand.DevmandGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	devices := getDeviceKeys(entity.Associations)
	config := &models.GatewayDevmandConfigs{
		ManagedDevices: devices,
	}
	return c.JSON(http.StatusOK, config)
}

func getDeviceTKs(devices []string) []storage.TypeAndKey {
	deviceTKs := []storage.TypeAndKey{}
	for _, device := range devices {
		deviceTKs = append(deviceTKs, storage.TypeAndKey{Type: devmand.DeviceType, Key: device})
	}
	return deviceTKs
}

func getDeviceKeys(deviceTKs []storage.TypeAndKey) []string {
	deviceKeys := []string{}
	for _, deviceTK := range deviceTKs {
		deviceKeys = append(deviceKeys, deviceTK.Key)
	}
	return deviceKeys
}

func getDeviceID(c echo.Context) (string, *echo.HTTPError) {
	devID := c.Param("device_id")
	if devID == "" {
		return devID, obsidian.HttpError(
			fmt.Errorf("Invalid/Missing Device ID"),
			http.StatusBadRequest)
	}
	return devID, nil
}

func newDeviceID(c echo.Context) (string, *echo.HTTPError) {
	requestedID := c.QueryParam("requested_id")
	if len(requestedID) <= 0 {
		return requestedID, obsidian.HttpError(
			fmt.Errorf("Requested device ID cannot be empty"),
			http.StatusBadRequest,
		)
	}
	r, _ := regexp.Compile("^[a-z_][0-9a-z_]+$")
	if !r.MatchString(requestedID) {
		return requestedID, obsidian.HttpError(
			fmt.Errorf("Device ID '%s' is not allowed. Device ID can only contain "+
				"lowercase alphanumeric characters and underscore, and should start with a letter or underscore.", requestedID),
			http.StatusBadRequest,
		)
	}
	return requestedID, nil
}
