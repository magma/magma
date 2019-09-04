/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"errors"
	"net/http"
	"sort"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
)

const (
	ReleaseChannelsRootPath   = obsidian.RestRoot + obsidian.UrlSep + "channels"
	ReleaseChannelsManagePath = ReleaseChannelsRootPath + "/:channel_id"
	TiersRootPath             = obsidian.RestRoot + "/networks/:network_id/tiers"
	TiersManagePath           = TiersRootPath + "/:tier_id"
)

// GetObsidianHandlers returns the obsidian handlers for upgrade
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{Path: ReleaseChannelsRootPath, Methods: obsidian.GET, HandlerFunc: listReleaseChannel},
		{Path: ReleaseChannelsRootPath, Methods: obsidian.POST, HandlerFunc: createReleaseChannel},
		{Path: ReleaseChannelsManagePath, Methods: obsidian.GET, HandlerFunc: getReleaseChannel},
		{Path: ReleaseChannelsManagePath, Methods: obsidian.PUT, HandlerFunc: updateReleaseChannel},
		{Path: ReleaseChannelsManagePath, Methods: obsidian.DELETE, HandlerFunc: deleteReleaseChannel},
		{Path: TiersRootPath, Methods: obsidian.GET, HandlerFunc: listTiers},
		{Path: TiersRootPath, Methods: obsidian.POST, HandlerFunc: createrTier},
		{Path: TiersManagePath, Methods: obsidian.GET, HandlerFunc: getTier},
		{Path: TiersManagePath, Methods: obsidian.PUT, HandlerFunc: updateTier},
		{Path: TiersManagePath, Methods: obsidian.DELETE, HandlerFunc: deleteTier},
	}
}

func noChannelIdError() error {
	return obsidian.HttpError(
		errors.New("Missing release channel ID"),
		http.StatusBadRequest,
	)
}

func listReleaseChannel(c echo.Context) error {
	channelNames, err := configurator.ListInternalEntityKeys(orc8r.UpgradeReleaseChannelEntityType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	// Return a deterministic ordering of channels
	sort.Strings(channelNames)
	if len(channelNames) == 0 {
		channelNames = nil
	}
	return c.JSON(http.StatusOK, channelNames)
}

func createReleaseChannel(c echo.Context) error {
	channel := new(models.ReleaseChannel)
	if err := c.Bind(channel); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	entity := configurator.NetworkEntity{
		Type:   orc8r.UpgradeReleaseChannelEntityType,
		Key:    string(channel.Name),
		Name:   string(channel.Name),
		Config: channel,
	}
	_, err := configurator.CreateInternalEntity(entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	// Return the ID of the created channel
	return c.JSON(http.StatusCreated, channel.Name)
}

func updateReleaseChannel(c echo.Context) error {
	channelID := c.Param("channel_id")
	if channelID == "" {
		return noChannelIdError()
	}
	channel := new(models.ReleaseChannel)
	if err := c.Bind(channel); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	// Release channel name is immutable
	// This could change if release channels are keyed by UUID in their tables
	if string(channel.Name) != channelID {
		return obsidian.HttpError(
			errors.New("Release channel name cannot be modified"),
			http.StatusBadRequest)
	}

	update := configurator.EntityUpdateCriteria{
		Key:       string(channel.Name),
		Type:      orc8r.UpgradeReleaseChannelEntityType,
		NewConfig: channel,
	}
	_, err := configurator.UpdateInternalEntity(update)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func deleteReleaseChannel(c echo.Context) error {
	channelID := c.Param("channel_id")
	if channelID == "" {
		return noChannelIdError()
	}
	// the API requires that an error is returned when the channel does not exist
	exists, err := configurator.DoesInternalEntityExist(orc8r.UpgradeReleaseChannelEntityType, channelID)
	if err != nil || !exists {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	err = configurator.DeleteInternalEntity(orc8r.UpgradeReleaseChannelEntityType, channelID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func getReleaseChannel(c echo.Context) error {
	channelID := c.Param("channel_id")
	if channelID == "" {
		return noChannelIdError()
	}

	ent, err := configurator.LoadInternalEntity(orc8r.UpgradeReleaseChannelEntityType, channelID, configurator.EntityLoadCriteria{LoadConfig: true, LoadMetadata: true})
	if err != nil {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	releaseChannel := ent.Config.(*models.ReleaseChannel)
	releaseChannel.Name = ent.Name
	return c.JSON(http.StatusOK, releaseChannel)
}

func noTierIdError() error {
	return obsidian.HttpError(
		errors.New("Missing tier ID"),
		http.StatusBadRequest,
	)
}

func listTiers(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	tiers, err := configurator.ListEntityKeys(networkID, orc8r.UpgradeTierEntityType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	// Return a deterministic ordering of tiers
	sort.Strings(tiers)
	return c.JSON(http.StatusOK, tiers)
}

func createrTier(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tier := new(models.Tier)
	if err := c.Bind(tier); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	entity := configurator.NetworkEntity{
		Type:   orc8r.UpgradeTierEntityType,
		Key:    string(tier.ID),
		Name:   string(tier.Name),
		Config: tier,
	}
	_, err := configurator.CreateEntity(networkID, entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	// Return the ID of the created tier
	return c.JSON(http.StatusCreated, tier.Name)
}

func getTier(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierID := c.Param("tier_id")
	if tierID == "" {
		return noTierIdError()
	}

	tier, err := configurator.LoadEntity(networkID, orc8r.UpgradeTierEntityType, tierID, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	// Return the ID of the created tier
	return c.JSON(http.StatusCreated, tier.Config)
}

func updateTier(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierID := c.Param("tier_id")
	if tierID == "" {
		return noTierIdError()
	}
	tier := new(models.Tier)
	if err := c.Bind(tier); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	update := configurator.EntityUpdateCriteria{
		Key:       string(tier.ID),
		Type:      orc8r.UpgradeTierEntityType,
		NewName:   swag.String(string(tier.Name)),
		NewConfig: tier,
	}
	_, err := configurator.UpdateEntity(networkID, update)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func deleteTier(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierID := c.Param("tier_id")
	if tierID == "" {
		return noTierIdError()
	}
	// the API requires that an error is returned when the channel does not exist
	exists, err := configurator.DoesEntityExist(networkID, orc8r.UpgradeTierEntityType, tierID)
	if err != nil || !exists {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	err = configurator.DeleteEntity(networkID, orc8r.UpgradeTierEntityType, tierID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
