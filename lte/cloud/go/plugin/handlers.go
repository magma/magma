/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin

import (
	"fmt"
	"net/http"
	"sort"

	"magma/lte/cloud/go/lte"
	ltemodels "magma/lte/cloud/go/plugin/models"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	orc8rmodels "magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	LteNetworks                  = "ltenetworks"
	ListNetworksPath             = obsidian.V1Root + LteNetworks
	ManageNetworkPath            = ListNetworksPath + "/:network_id"
	ManageNetworkNamePath        = ManageNetworkPath + obsidian.UrlSep + "name"
	ManageNetworkDescriptionPath = ManageNetworkPath + obsidian.UrlSep + "description"
	ManageNetworkFeaturesPath    = ManageNetworkPath + obsidian.UrlSep + "features"
	ManageNetworkDNSPath         = ManageNetworkPath + obsidian.UrlSep + "dns"
)

func GetNetworkHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListNetworksPath, Methods: obsidian.GET, HandlerFunc: ListNetworks},
		{Path: ListNetworksPath, Methods: obsidian.POST, HandlerFunc: CreateNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.GET, HandlerFunc: GetNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.PUT, HandlerFunc: UpdateNetwork},
		{Path: ManageNetworkPath, Methods: obsidian.DELETE, HandlerFunc: DeleteNetwork},
	}
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkNamePath, new(models.NetworkName), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDescriptionPath, new(models.NetworkDescription), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkFeaturesPath, &orc8rmodels.NetworkFeatures{}, orc8r.NetworkFeaturesConfig)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDNSPath, &orc8rmodels.NetworkDNSConfig{}, orc8r.DnsdNetworkType)...)
	return ret
}

func ListNetworks(c echo.Context) error {
	ids, err := configurator.ListNetworksOfType(lte.LteNetworkType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(ids)
	return c.JSON(http.StatusOK, ids)
}

func CreateNetwork(c echo.Context) error {
	payload := &ltemodels.LteNetwork{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	err := configurator.CreateNetwork(payload.ToConfiguratorNetwork())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func GetNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	network, err := configurator.LoadNetwork(nid, true, true)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if network.Type != lte.LteNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not an LTE network", nid))
	}

	ret := (&ltemodels.LteNetwork{}).FromConfiguratorNetwork(network)
	return c.JSON(http.StatusOK, ret)
}

func UpdateNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &ltemodels.LteNetwork{}
	err := c.Bind(payload)
	if err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// check that this is actually an LTE network
	network, err := configurator.LoadNetwork(nid, false, false)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
	}
	if network.Type != lte.LteNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not an LTE network", nid))
	}

	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{payload.ToUpdateCriteria()})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteNetwork(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	// check that this is actually an LTE network
	network, err := configurator.LoadNetwork(nid, false, false)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load network to check type"), http.StatusInternalServerError)
	}
	if network.Type != lte.LteNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not an LTE network", nid))
	}

	err = configurator.DeleteNetwork(nid)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}
