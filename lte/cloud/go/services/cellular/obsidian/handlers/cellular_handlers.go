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

	"magma/lte/cloud/go/lte"
	models2 "magma/lte/cloud/go/plugin/models"
	"magma/lte/cloud/go/services/cellular/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	cfgObsidian "magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	magmad_handlers "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
)

const (
	ConfigKey         = "cellular"
	NetworkConfigPath = magmad_handlers.ConfigureNetwork + "/" + ConfigKey
	GatewayConfigPath = magmad_handlers.ConfigureAG + "/" + ConfigKey
	EnodebListPath    = magmad_handlers.ConfigureNetwork + "/enodeb"
	EnodebConfigPath  = magmad_handlers.ConfigureNetwork + "/enodeb/:enodeb_id"
)

// GetObsidianHandlers returns all obsidian handlers for the cellular service
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		cfgObsidian.GetReadNetworkConfigHandler(NetworkConfigPath, lte.CellularNetworkType, &models2.NetworkCellularConfigs{}),
		cfgObsidian.GetCreateNetworkConfigHandler(NetworkConfigPath, lte.CellularNetworkType, &models2.NetworkCellularConfigs{}),
		cfgObsidian.GetUpdateNetworkConfigHandler(NetworkConfigPath, lte.CellularNetworkType, &models2.NetworkCellularConfigs{}),
		cfgObsidian.GetDeleteNetworkConfigHandler(NetworkConfigPath, lte.CellularNetworkType),

		cfgObsidian.GetReadConfigHandler(EnodebConfigPath, lte.CellularEnodebType, getEnodebId, &models2.EnodebConfiguration{}),
		cfgObsidian.GetCreateConfigHandler(EnodebConfigPath, lte.CellularEnodebType, getEnodebId, &models2.EnodebConfiguration{}),
		cfgObsidian.GetUpdateConfigHandler(EnodebConfigPath, lte.CellularEnodebType, getEnodebId, &models2.EnodebConfiguration{}),
		cfgObsidian.GetDeleteConfigHandler(EnodebConfigPath, lte.CellularEnodebType, getEnodebId),
		// List all eNodeB devices for a network
		cfgObsidian.GetReadAllKeysConfigHandler(EnodebListPath, lte.CellularEnodebType),

		// Cellular gateway configs
		{Path: GatewayConfigPath, Methods: obsidian.POST, HandlerFunc: createGatewayConfig},
		{Path: GatewayConfigPath, Methods: obsidian.GET, HandlerFunc: getGatewayConfig},
		{Path: GatewayConfigPath, Methods: obsidian.PUT, HandlerFunc: updateGatewayConfig},
		cfgObsidian.GetDeleteGatewayConfigHandler(GatewayConfigPath, lte.CellularGatewayType),
	}
}

func getEnodebId(c echo.Context) (string, *echo.HTTPError) {
	operID := c.Param("enodeb_id")
	if operID == "" {
		return operID, obsidian.HttpError(
			fmt.Errorf("Invalid/Missing Enodeb ID"),
			http.StatusBadRequest)
	}
	return operID, nil
}

func getGatewayConfig(c echo.Context) error {
	networkID, gatewayID, nerr := getIDs(c)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		networkID, lte.CellularGatewayType, gatewayID,
		configurator.EntityLoadCriteria{
			LoadConfig:         true,
			LoadAssocsFromThis: true,
		},
	)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	loadedConfig := ent.Config.(*models2.GatewayCellularConfigs)
	retConfig := &models.GatewayCellularConfigs{
		Epc: &models.GatewayEpcConfigs{
			NatEnabled: swag.BoolValue(loadedConfig.Epc.NatEnabled),
			IPBlock:    loadedConfig.Epc.IPBlock,
		},
		Ran: &models.GatewayRanConfigs{
			TransmitEnabled: swag.BoolValue(loadedConfig.Ran.TransmitEnabled),
			Pci:             loadedConfig.Ran.Pci,
		},
		NonEpsService: &models.GatewayNonEpsServiceConfigs{
			Arfcn2g:              loadedConfig.NonEpsService.Arfcn2g,
			CsfbMcc:              loadedConfig.NonEpsService.CsfbMcc,
			CsfbMnc:              loadedConfig.NonEpsService.CsfbMnc,
			CsfbRat:              swag.Uint32Value(loadedConfig.NonEpsService.CsfbRat),
			Lac:                  swag.Uint32Value(loadedConfig.NonEpsService.Lac),
			NonEpsServiceControl: swag.Uint32Value(loadedConfig.NonEpsService.NonEpsServiceControl),
		},
	}
	for _, tk := range ent.Associations {
		if tk.Type == lte.CellularEnodebType {
			retConfig.AttachedEnodebSerials = append(retConfig.AttachedEnodebSerials, tk.Key)
		}
	}
	return c.JSON(http.StatusOK, retConfig)
}

func createGatewayConfig(c echo.Context) error {
	networkID, gatewayID, nerr := getIDs(c)
	if nerr != nil {
		return nerr
	}
	iConfig, nerr := cfgObsidian.GetConfigAndValidate(c, &models.GatewayCellularConfigs{})
	if nerr != nil {
		return nerr
	}
	config := iConfig.(*models.GatewayCellularConfigs)

	associationsToAdd := getEnodebTKs(config.AttachedEnodebSerials)

	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:         lte.CellularGatewayType,
		Key:          gatewayID,
		Config:       config,
		Associations: associationsToAdd,
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	_, err = configurator.UpdateEntity(networkID, configurator.EntityUpdateCriteria{
		Type:              orc8r.MagmadGatewayType,
		Key:               gatewayID,
		AssociationsToSet: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: gatewayID}},
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, gatewayID)
}

func updateGatewayConfig(c echo.Context) error {
	networkID, gatewayID, nerr := getIDs(c)
	if nerr != nil {
		return nerr
	}
	iConfig, nerr := cfgObsidian.GetConfigAndValidate(c, &models.GatewayCellularConfigs{})
	if nerr != nil {
		return nerr
	}
	config := iConfig.(*models.GatewayCellularConfigs)

	associationsToDelete := []storage.TypeAndKey{}
	associationsToSet := getEnodebTKs(config.AttachedEnodebSerials)

	if len(config.AttachedEnodebSerials) == 0 {
		// due to the way protobuf serialize/deserializes,
		// associationsToSet = [] does not delete all associations, so here we
		// look up the entity's association to pass in as associationsToDelete.
		entity, err := configurator.LoadEntity(networkID, lte.CellularGatewayType, gatewayID, configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		associationsToDelete = entity.Associations
	}

	_, err := configurator.UpdateEntity(networkID, configurator.EntityUpdateCriteria{
		Type:                 lte.CellularGatewayType,
		Key:                  gatewayID,
		NewConfig:            config,
		AssociationsToSet:    associationsToSet,
		AssociationsToDelete: associationsToDelete,
	})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func getEnodebTKs(enodbSerials []string) []storage.TypeAndKey {
	enodebTKs := []storage.TypeAndKey{}
	for _, enodebSerial := range enodbSerials {
		enodebTKs = append(enodebTKs, storage.TypeAndKey{Key: enodebSerial, Type: lte.CellularEnodebType})
	}
	return enodebTKs
}

func getIDs(c echo.Context) (string, string, error) {
	networkID, err := obsidian.GetNetworkId(c)
	if err != nil {
		return "", "", err
	}
	gatewayID, err := obsidian.GetLogicalGwId(c)
	if err != nil {
		return "", "", err
	}
	return networkID, gatewayID, nil
}
