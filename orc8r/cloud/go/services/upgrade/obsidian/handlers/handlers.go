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

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	upgrade_client "magma/orc8r/cloud/go/services/upgrade"
	"magma/orc8r/cloud/go/services/upgrade/obsidian/models"

	"github.com/labstack/echo"
)

const (
	ReleaseChannelsRootPath   = handlers.CHANNELS_ROOT
	ReleaseChannelsManagePath = ReleaseChannelsRootPath + "/:channel_id"
	TiersRootPath             = handlers.REST_ROOT + "/networks/:network_id/tiers"
	TiersManagePath           = TiersRootPath + "/:tier_id"
)

// GetObsidianHandlers returns the obsidian handlers for upgrade
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		{Path: ReleaseChannelsRootPath, Methods: handlers.GET, HandlerFunc: listReleaseChannelsHandler, MigratedHandlerFunc: listReleaseChannel},
		{Path: ReleaseChannelsRootPath, Methods: handlers.POST, HandlerFunc: createReleaseChannelHandler, MigratedHandlerFunc: createReleaseChannel, MultiplexAfterMigration: true},
		{Path: ReleaseChannelsManagePath, Methods: handlers.GET, HandlerFunc: getReleaseChannelsHandler, MigratedHandlerFunc: getReleaseChannel},
		{Path: ReleaseChannelsManagePath, Methods: handlers.PUT, HandlerFunc: updateReleaseChannelHandler, MigratedHandlerFunc: updateReleaseChannel, MultiplexAfterMigration: true},
		{Path: ReleaseChannelsManagePath, Methods: handlers.DELETE, HandlerFunc: deleteReleaseChannelHandler, MigratedHandlerFunc: deleteReleaseChannel, MultiplexAfterMigration: true},
		{Path: TiersRootPath, Methods: handlers.GET, HandlerFunc: listTiersHandler, MigratedHandlerFunc: listTiers},
		{Path: TiersRootPath, Methods: handlers.POST, HandlerFunc: createTierHandler, MigratedHandlerFunc: createrTier, MultiplexAfterMigration: true},
		{Path: TiersManagePath, Methods: handlers.GET, HandlerFunc: getTierHandler, MigratedHandlerFunc: getTier},
		{Path: TiersManagePath, Methods: handlers.PUT, HandlerFunc: updateTierHandler, MigratedHandlerFunc: updateTier, MultiplexAfterMigration: true},
		{Path: TiersManagePath, Methods: handlers.DELETE, HandlerFunc: deleteTierHandler, MigratedHandlerFunc: deleteTier, MultiplexAfterMigration: true},
	}
}

func noChannelIdError() error {
	return handlers.HttpError(
		errors.New("Missing release channel ID"),
		http.StatusBadRequest,
	)
}

func listReleaseChannel(c echo.Context) error {
	channelNames, err := configurator.ListInternalEntityKeys(upgrade_client.UpgradeReleaseChannelEntityType)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
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
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	entity := configurator.NetworkEntity{
		Type:   upgrade_client.UpgradeReleaseChannelEntityType,
		Key:    channel.Name,
		Name:   channel.Name,
		Config: channel,
	}
	_, err := configurator.CreateInternalEntity(entity)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
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
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	// Release channel name is immutable
	// This could change if release channels are keyed by UUID in their tables
	if channel.Name != channelID {
		return handlers.HttpError(
			errors.New("Release channel name cannot be modified"),
			http.StatusBadRequest)
	}

	update := configurator.EntityUpdateCriteria{
		Key:       channel.Name,
		Type:      upgrade_client.UpgradeReleaseChannelEntityType,
		NewConfig: channel,
	}
	_, err := configurator.UpdateInternalEntity(update)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func deleteReleaseChannel(c echo.Context) error {
	channelID := c.Param("channel_id")
	if channelID == "" {
		return noChannelIdError()
	}
	// the API requires that an error is returned when the channel does not exist
	exists, err := configurator.DoesInternalEntityExist(upgrade_client.UpgradeReleaseChannelEntityType, channelID)
	if err != nil || !exists {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	err = configurator.DeleteInternalEntity(upgrade_client.UpgradeReleaseChannelEntityType, channelID)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func getReleaseChannel(c echo.Context) error {
	channelID := c.Param("channel_id")
	if channelID == "" {
		return noChannelIdError()
	}

	ent, err := configurator.LoadInternalEntity(upgrade_client.UpgradeReleaseChannelEntityType, channelID, configurator.EntityLoadCriteria{LoadConfig: true, LoadMetadata: true})
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}
	releaseChannel := ent.Config.(*models.ReleaseChannel)
	releaseChannel.Name = ent.Name
	return c.JSON(http.StatusOK, releaseChannel)
}

func noTierIdError() error {
	return handlers.HttpError(
		errors.New("Missing tier ID"),
		http.StatusBadRequest,
	)
}

func listTiers(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	tiers, err := configurator.ListEntityKeys(networkID, orc8r.UpgradeTierEntityType)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	// Return a deterministic ordering of tiers
	sort.Strings(tiers)
	return c.JSON(http.StatusOK, tiers)
}

func createrTier(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tier := new(models.Tier)
	if err := c.Bind(tier); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	entity := configurator.NetworkEntity{
		Type:   orc8r.UpgradeTierEntityType,
		Key:    tier.ID,
		Config: tier,
	}
	_, err := configurator.CreateEntity(networkID, entity)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	// Return the ID of the created tier
	return c.JSON(http.StatusCreated, tier.Name)
}

func getTier(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierID := c.Param("tier_id")
	if tierID == "" {
		return noTierIdError()
	}

	tier, err := configurator.LoadEntity(networkID, orc8r.UpgradeTierEntityType, tierID, configurator.EntityLoadCriteria{LoadConfig: true})
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}
	// Return the ID of the created tier
	return c.JSON(http.StatusCreated, tier.Config)
}

func updateTier(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierID := c.Param("tier_id")
	if tierID == "" {
		return noTierIdError()
	}
	tier := new(models.Tier)
	if err := c.Bind(tier); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}
	update := configurator.EntityUpdateCriteria{
		Key:       tier.ID,
		Type:      orc8r.UpgradeTierEntityType,
		NewConfig: tier,
	}
	_, err := configurator.UpdateEntity(networkID, update)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func deleteTier(c echo.Context) error {
	networkID, nerr := handlers.GetNetworkId(c)
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
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	err = configurator.DeleteEntity(networkID, orc8r.UpgradeTierEntityType, tierID)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
