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

	"magma/lte/cloud/go/services/cellular/config"
	"magma/lte/cloud/go/services/cellular/obsidian/models"
	"magma/lte/cloud/go/services/cellular/utils"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/config/obsidian"
	magmad_handlers "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"

	"github.com/labstack/echo"
)

const (
	ConfigKey         = "cellular"
	NetworkConfigPath = magmad_handlers.ConfigureNetwork + "/" + ConfigKey
	GatewayConfigPath = magmad_handlers.ConfigureAG + "/" + ConfigKey
	EnodebConfigPath  = magmad_handlers.ConfigureNetwork + "/enodeb/:enodeb_id"
)

// GetObsidianHandlers returns all obsidian handlers for the cellular service
func GetObsidianHandlers() []handlers.Handler {
	defaultUpdateHandler := obsidian.GetUpdateNetworkConfigHandler(NetworkConfigPath, config.CellularNetworkType, &models.NetworkCellularConfigs{})
	ret := []handlers.Handler{
		obsidian.GetReadNetworkConfigHandler(NetworkConfigPath, config.CellularNetworkType, &models.NetworkCellularConfigs{}),
		obsidian.GetCreateNetworkConfigHandler(NetworkConfigPath, config.CellularNetworkType, &models.NetworkCellularConfigs{}),
		obsidian.GetDeleteNetworkConfigHandler(NetworkConfigPath, config.CellularNetworkType),
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
		},
		obsidian.GetReadConfigHandler(EnodebConfigPath, config.CellularEnodebType, getEnodebId, &models.NetworkEnodebConfigs{}),
		obsidian.GetCreateConfigHandler(EnodebConfigPath, config.CellularEnodebType, getEnodebId, &models.NetworkEnodebConfigs{}),
		obsidian.GetUpdateConfigHandler(EnodebConfigPath, config.CellularEnodebType, getEnodebId, &models.NetworkEnodebConfigs{}),
		obsidian.GetDeleteConfigHandler(EnodebConfigPath, config.CellularEnodebType, getEnodebId),
	}
	ret = append(ret, obsidian.GetCRUDGatewayConfigHandlers(GatewayConfigPath, config.CellularGatewayType, &models.GatewayCellularConfigs{})...)
	return ret
}

func getEnodebId(c echo.Context) (string, *echo.HTTPError) {
	operID := c.Param("enodeb_id")
	if operID == "" {
		return operID, handlers.HttpError(
			fmt.Errorf("Invalid/Missing Enodeb ID"),
			http.StatusBadRequest)
	}
	return operID, nil
}

func getNetworkConfigFromRequest(c echo.Context) (echo.Context, error) {
	if c.Request().Body == nil {
		return nil, handlers.HttpError(fmt.Errorf("Network config is nil"), http.StatusBadRequest)
	}
	cfg := &models.NetworkCellularConfigs{}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, handlers.HttpError(err, http.StatusBadRequest)
	}
	err = json.Unmarshal(body, cfg)
	if err != nil {
		return nil, handlers.HttpError(err, http.StatusBadRequest)
	}

	// Config does not have a FDD/TDD sub-config set
	if cfg.Ran.TddConfig == nil && cfg.Ran.FddConfig == nil {
		band, err := utils.GetBand(int32(cfg.Ran.Earfcndl))
		if err != nil {
			return nil, handlers.HttpError(err, http.StatusBadRequest)
		}

		cfg, err = setAppropriateNetworkSubConfig(band, cfg)
		if err != nil {
			return nil, handlers.HttpError(err, http.StatusBadRequest)
		}
	}

	body, err = json.Marshal(cfg)
	if err != nil {
		return nil, handlers.HttpError(fmt.Errorf("Error converting config to TDD/FDD format"), http.StatusBadRequest)
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
