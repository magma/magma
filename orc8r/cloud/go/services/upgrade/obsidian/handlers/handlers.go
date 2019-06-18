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
	"magma/orc8r/cloud/go/services/configurator"
	upgrade_client "magma/orc8r/cloud/go/services/upgrade"
	"magma/orc8r/cloud/go/services/upgrade/obsidian/models"
	"magma/orc8r/cloud/go/services/upgrade/protos"

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
		{Path: ReleaseChannelsRootPath, Methods: handlers.POST, HandlerFunc: createReleaseChannelHandler, MigratedHandlerFunc: createReleaseChannel},
		{Path: ReleaseChannelsManagePath, Methods: handlers.GET, HandlerFunc: getReleaseChannelsHandler, MigratedHandlerFunc: getReleaseChannel},
		{Path: ReleaseChannelsManagePath, Methods: handlers.PUT, HandlerFunc: updateReleaseChannelHandler, MigratedHandlerFunc: updateReleaseChannel},
		{Path: ReleaseChannelsManagePath, Methods: handlers.DELETE, HandlerFunc: deleteReleaseChannelHandler, MigratedHandlerFunc: deleteReleaseChannel},
		{Path: TiersRootPath, Methods: handlers.GET, HandlerFunc: listTiersHandler},
		{Path: TiersRootPath, Methods: handlers.POST, HandlerFunc: createTierHandler},
		{Path: TiersManagePath, Methods: handlers.GET, HandlerFunc: getTierHandler},
		{Path: TiersManagePath, Methods: handlers.PUT, HandlerFunc: updateTierHandler},
		{Path: TiersManagePath, Methods: handlers.DELETE, HandlerFunc: deleteTierHandler},
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

func listTiersHandler(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	tiers, err := upgrade_client.GetTiers(networkId, []string{})
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	ret := make([]string, 0, len(tiers))
	for tierId := range tiers {
		ret = append(ret, tierId)
	}
	// Return a deterministic ordering of tiers
	sort.Strings(ret)
	return c.JSON(http.StatusOK, ret)
}

func tierInfoModelToProto(model *models.Tier) *protos.TierInfo {
	// Copy each image spec into a protobuf
	var imageArray []*protos.ImageSpec
	for _, elem := range model.Images {
		imageArray = append(imageArray, &protos.ImageSpec{
			Name:  elem.Name,
			Order: elem.Order})
	}

	return &protos.TierInfo{
		Name:    model.Name,
		Version: model.Version,
		Images:  imageArray,
	}
}

func createTierHandler(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	restTier := new(models.Tier)
	if err := c.Bind(restTier); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	tierProto := tierInfoModelToProto(restTier)

	err := upgrade_client.CreateTier(networkId, restTier.ID, tierProto)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	// Return the ID of the created tier
	return c.JSON(http.StatusCreated, restTier.ID)
}

func getTierHandler(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierId := c.Param("tier_id")
	if tierId == "" {
		return noTierIdError()
	}

	tiers, err := upgrade_client.GetTiers(networkId, []string{tierId})
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}
	tierProto, ok := tiers[tierId]
	if !ok {
		return handlers.HttpError(
			errors.New("Error while loading tier from service"),
			http.StatusNotFound)
	}

	var imgList []*models.TierImagesItems0
	for _, elem := range tierProto.GetImages() {
		imgList = append(imgList, &models.TierImagesItems0{
			Name:  elem.Name,
			Order: elem.Order})
	}

	restTier := models.Tier{
		ID:      tierId,
		Name:    tierProto.GetName(),
		Version: tierProto.GetVersion(),
		Images:  imgList,
	}
	return c.JSON(http.StatusOK, restTier)
}

func updateTierHandler(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierId := c.Param("tier_id")
	if tierId == "" {
		return noTierIdError()
	}
	restTier := new(models.Tier)
	if err := c.Bind(restTier); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	tierProto := tierInfoModelToProto(restTier)

	err := upgrade_client.UpdateTier(networkId, tierId, tierProto)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func deleteTierHandler(c echo.Context) error {
	networkId, nerr := handlers.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierId := c.Param("tier_id")
	if tierId == "" {
		return noTierIdError()
	}

	err := upgrade_client.DeleteTier(networkId, tierId)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func noTierIdError() error {
	return handlers.HttpError(
		errors.New("Missing tier ID"),
		http.StatusBadRequest,
	)
}
