/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/services/cellular/obsidian/models"
	"magma/lte/cloud/go/services/cellular/utils"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	cfgObsidian "magma/orc8r/cloud/go/services/config/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	magmad_handlers "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
	"magma/orc8r/cloud/go/storage"

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
	defaultUpdateHandler := cfgObsidian.GetUpdateNetworkConfigHandler(NetworkConfigPath, lte.CellularNetworkType, &models.NetworkCellularConfigs{})
	createGatewayConfigHandler := cfgObsidian.GetCreateGatewayConfigHandler(GatewayConfigPath, lte.CellularGatewayType, &models.GatewayCellularConfigs{})
	updateGatewayConfigHandler := cfgObsidian.GetUpdateGatewayConfigHandler(GatewayConfigPath, lte.CellularGatewayType, &models.GatewayCellularConfigs{})

	// override create and update migrated handler func
	createGatewayConfigHandler.MigratedHandlerFunc = createGatewayConfig
	updateGatewayConfigHandler.MigratedHandlerFunc = updateGatewayConfig

	return []obsidian.Handler{
		cfgObsidian.GetReadNetworkConfigHandler(NetworkConfigPath, lte.CellularNetworkType, &models.NetworkCellularConfigs{}),
		cfgObsidian.GetCreateNetworkConfigHandler(NetworkConfigPath, lte.CellularNetworkType, &models.NetworkCellularConfigs{}),
		cfgObsidian.GetDeleteNetworkConfigHandler(NetworkConfigPath, lte.CellularNetworkType),
		// Patch default config update handler to set TDD/FDD fields in network config
		{
			Path:    defaultUpdateHandler.Path,
			Methods: defaultUpdateHandler.Methods,
			HandlerFunc: func(c echo.Context) error {
				cc, err := getNetworkConfigFromRequest(c)
				if err != nil {
					return err
				}
				return defaultUpdateHandler.HandlerFunc(cc)
			},
			MigratedHandlerFunc: func(c echo.Context) error {
				cc, err := getNetworkConfigFromRequest(c)
				if err != nil {
					return err
				}
				return defaultUpdateHandler.MigratedHandlerFunc(cc)
			},
		},
		cfgObsidian.GetReadConfigHandler(EnodebConfigPath, lte.CellularEnodebType, getEnodebId, &models.NetworkEnodebConfigs{}),
		cfgObsidian.GetCreateConfigHandler(EnodebConfigPath, lte.CellularEnodebType, getEnodebId, &models.NetworkEnodebConfigs{}),
		cfgObsidian.GetUpdateConfigHandler(EnodebConfigPath, lte.CellularEnodebType, getEnodebId, &models.NetworkEnodebConfigs{}),
		cfgObsidian.GetDeleteConfigHandler(EnodebConfigPath, lte.CellularEnodebType, getEnodebId),
		// List all eNodeB devices for a network
		cfgObsidian.GetReadAllKeysConfigHandler(EnodebListPath, lte.CellularEnodebType),
		// Cellular gateway configs
		cfgObsidian.GetReadGatewayConfigHandler(GatewayConfigPath, lte.CellularGatewayType, &models.GatewayCellularConfigs{}),
		cfgObsidian.GetDeleteGatewayConfigHandler(GatewayConfigPath, lte.CellularGatewayType),
		createGatewayConfigHandler,
		updateGatewayConfigHandler,
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

func getNetworkConfigFromRequest(c echo.Context) (echo.Context, error) {
	if c.Request().Body == nil {
		return nil, obsidian.HttpError(fmt.Errorf("Network config is nil"), http.StatusBadRequest)
	}
	cfg := &models.NetworkCellularConfigs{}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}
	err = json.Unmarshal(body, cfg)
	if err != nil {
		return nil, obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Config does not have a FDD/TDD sub-config set
	if cfg.Ran.TddConfig == nil && cfg.Ran.FddConfig == nil {
		band, err := utils.GetBand(cfg.Ran.Earfcndl)
		if err != nil {
			return nil, obsidian.HttpError(err, http.StatusBadRequest)
		}

		cfg, err = setAppropriateNetworkSubConfig(band, cfg)
		if err != nil {
			return nil, obsidian.HttpError(err, http.StatusBadRequest)
		}
	}

	body, err = json.Marshal(cfg)
	if err != nil {
		return nil, obsidian.HttpError(fmt.Errorf("Error converting config to TDD/FDD format"), http.StatusBadRequest)
	}
	// populate request body with the updated config
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return c, nil
}

func setAppropriateNetworkSubConfig(band *utils.LTEBand, config *models.NetworkCellularConfigs) (*models.NetworkCellularConfigs, error) {
	switch band.Mode {
	case utils.TDDMode:
		config.Ran.TddConfig = &models.NetworkRanConfigsTddConfig{
			Earfcndl:               config.Ran.Earfcndl,
			SubframeAssignment:     config.Ran.SubframeAssignment,
			SpecialSubframePattern: config.Ran.SpecialSubframePattern,
		}
		return config, nil
	case utils.FDDMode:
		earfcndl := config.Ran.Earfcndl
		// Use the same math as in validateNetworkRANConfig
		earfcnul := earfcndl - uint32(band.StartEarfcnDl) + uint32(band.StartEarfcnUl)
		config.Ran.FddConfig = &models.NetworkRanConfigsFddConfig{
			Earfcndl: earfcndl,
			Earfcnul: earfcnul,
		}
		return config, nil
	default:
		return nil, fmt.Errorf("Invalid LTE band mode supplied")
	}
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
