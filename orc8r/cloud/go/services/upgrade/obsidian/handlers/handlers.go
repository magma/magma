/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_utils "magma/orc8r/cloud/go/services/configurator/obsidian/handler_utils"
	upgrade_client "magma/orc8r/cloud/go/services/upgrade"
	"magma/orc8r/cloud/go/services/upgrade/obsidian/models"
	"magma/orc8r/cloud/go/services/upgrade/protos"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
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
		{Path: ReleaseChannelsRootPath, Methods: handlers.GET, HandlerFunc: listReleaseChannelsHandler},
		{Path: ReleaseChannelsRootPath, Methods: handlers.POST, HandlerFunc: createReleaseChannelHandler},
		{Path: ReleaseChannelsManagePath, Methods: handlers.GET, HandlerFunc: getReleaseChannelsHandler},
		{Path: ReleaseChannelsManagePath, Methods: handlers.PUT, HandlerFunc: updateReleaseChannelHandler},
		{Path: ReleaseChannelsManagePath, Methods: handlers.DELETE, HandlerFunc: deleteReleaseChannelHandler},
		{Path: TiersRootPath, Methods: handlers.GET, HandlerFunc: listTiersHandler},
		{Path: TiersRootPath, Methods: handlers.POST, HandlerFunc: createTierHandler},
		{Path: TiersManagePath, Methods: handlers.GET, HandlerFunc: getTierHandler},
		{Path: TiersManagePath, Methods: handlers.PUT, HandlerFunc: updateTierHandler},
		{Path: TiersManagePath, Methods: handlers.DELETE, HandlerFunc: deleteTierHandler},
	}
}

// List all release channels by ID
func listReleaseChannelsHandler(c echo.Context) error {
	channels, err := upgrade_client.ListReleaseChannels()
	if err != nil {
		return handlers.HttpError(err)
	}
	// Return a deterministic ordering of channels
	sort.Strings(channels)
	return c.JSON(http.StatusOK, channels)
}

func getReleaseChannelsHandler(c echo.Context) error {
	channelId := c.Param("channel_id")
	if channelId == "" {
		return noChannelIdError()
	}

	channel, err := upgrade_client.GetReleaseChannel(channelId)
	if err != nil {
		return handlers.HttpError(err, http.StatusNotFound)
	}

	swaggerChannel := models.ReleaseChannel{}
	swaggerChannel.Name = channelId
	swaggerChannel.SupportedVersions = channel.GetSupportedVersions()
	return c.JSON(http.StatusOK, swaggerChannel)
}

func createReleaseChannelHandler(c echo.Context) error {
	restChannel := new(models.ReleaseChannel)
	if err := c.Bind(restChannel); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	// Construct proto model and persist
	channelProto := &protos.ReleaseChannel{
		SupportedVersions: restChannel.SupportedVersions,
	}
	err := upgrade_client.CreateReleaseChannel(restChannel.Name, channelProto)
	if err != nil {
		return handlers.HttpError(err)
	}

	err = multiplexCreateReleaseChannelIntoConfigurator(restChannel)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex create into configurator : %v", err))
	}

	// Return the ID of the created channel
	return c.JSON(http.StatusCreated, restChannel.Name)
}

func multiplexCreateReleaseChannelIntoConfigurator(channel *models.ReleaseChannel) error {
	entity := configurator.NetworkEntity{
		Type:   upgrade_client.ReleaseChannelType,
		Key:    channel.Name,
		Config: channel,
	}
	_, err := configurator.CreateInternalEntities([]configurator.NetworkEntity{entity})
	return err
}

func updateReleaseChannelHandler(c echo.Context) error {
	channelId := c.Param("channel_id")
	if channelId == "" {
		return noChannelIdError()
	}
	restChannel := new(models.ReleaseChannel)
	if err := c.Bind(restChannel); err != nil {
		return handlers.HttpError(err, http.StatusBadRequest)
	}

	// Release channel name is immutable
	// This could change if release channels are keyed by UUID in their tables
	if restChannel.Name != channelId {
		return handlers.HttpError(
			errors.New("Release channel name cannot be modified"),
			http.StatusBadRequest)
	}
	updatedChannelProto := &protos.ReleaseChannel{
		SupportedVersions: restChannel.SupportedVersions,
	}

	err := upgrade_client.UpdateReleaseChannel(channelId, updatedChannelProto)
	if err != nil {
		return handlers.HttpError(err)
	}

	err = multiplexUpdateReleaseChannelIntoConfigurator(restChannel)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex update into configurator : %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func multiplexUpdateReleaseChannelIntoConfigurator(channel *models.ReleaseChannel) error {
	err := configurator_utils.CreateInternalNetworkEntityIfNotExists(upgrade_client.ReleaseChannelType, channel.Name)
	if err != nil {
		return err
	}
	update := configurator.EntityUpdateCriteria{
		Key:       channel.Name,
		Type:      upgrade_client.ReleaseChannelType,
		NewConfig: channel,
	}
	_, err = configurator.UpdateInternalEntity([]configurator.EntityUpdateCriteria{update})
	return err
}

func deleteReleaseChannelHandler(c echo.Context) error {
	channelId := c.Param("channel_id")
	if channelId == "" {
		return noChannelIdError()
	}

	err := upgrade_client.DeleteReleaseChannel(channelId)
	if err != nil {
		return handlers.HttpError(err, http.StatusInternalServerError)
	}

	// multiplex delete into configurator
	err = configurator.DeleteInternalEntities([]storage.TypeAndKey{{Key: channelId, Type: upgrade_client.ReleaseChannelType}})
	if err != nil {
		glog.Errorf("Failed to multiplex delete into configurator: %v", err)
	}

	return c.NoContent(http.StatusNoContent)
}

func noChannelIdError() error {
	return handlers.HttpError(
		errors.New("Missing release channel ID"),
		http.StatusBadRequest,
	)
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

	err = multiplexCreateTierIntoConfigurator(networkId, restTier)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex create into configurator : %v", err))
	}
	// Return the ID of the created tier
	return c.JSON(http.StatusCreated, restTier.ID)
}

func multiplexCreateTierIntoConfigurator(networkID string, tier *models.Tier) error {
	entity := configurator.NetworkEntity{
		Type:   upgrade_client.NetworkTierType,
		Key:    tier.ID,
		Config: tier,
	}
	_, err := configurator.CreateEntities(networkID, []configurator.NetworkEntity{entity})
	return err
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

	err = multiplexUpdateTierIntoConfigurator(networkId, tierId, restTier)
	if err != nil {
		return handlers.HttpError(fmt.Errorf("Failed to multiplex update into configurator : %v", err))
	}

	return c.NoContent(http.StatusOK)
}

func multiplexUpdateTierIntoConfigurator(networkID, tierID string, tier *models.Tier) error {
	err := configurator_utils.CreateNetworkEntityIfNotExists(networkID, upgrade_client.NetworkTierType, tierID)
	if err != nil {
		return err
	}
	entity := configurator.EntityUpdateCriteria{
		Type:      upgrade_client.NetworkTierType,
		Key:       tierID,
		NewConfig: tier,
	}
	_, err = configurator.UpdateEntities(networkID, []configurator.EntityUpdateCriteria{entity})
	return err
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

	err = configurator.DeleteEntities(networkId, []storage.TypeAndKey{{Type: upgrade_client.NetworkTierType, Key: tierId}})
	if err != nil {
		glog.Errorf("Failed to multiplex delete into configurator: %v", err)
	}
	return c.NoContent(http.StatusNoContent)
}

func noTierIdError() error {
	return handlers.HttpError(
		errors.New("Missing tier ID"),
		http.StatusBadRequest,
	)
}
