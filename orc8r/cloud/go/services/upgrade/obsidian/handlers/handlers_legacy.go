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
	upgrade_client "magma/orc8r/cloud/go/services/upgrade"
	"magma/orc8r/cloud/go/services/upgrade/obsidian/models"
	"magma/orc8r/cloud/go/services/upgrade/protos"

	"github.com/labstack/echo"
)

// List all release channels by ID
func listReleaseChannelsHandler(c echo.Context) error {
	channels, err := upgrade_client.ListReleaseChannels()
	if err != nil {
		return obsidian.HttpError(err)
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
		return obsidian.HttpError(err, http.StatusNotFound)
	}

	swaggerChannel := models.ReleaseChannel{}
	swaggerChannel.Name = channelId
	swaggerChannel.SupportedVersions = channel.GetSupportedVersions()
	return c.JSON(http.StatusOK, swaggerChannel)
}

func createReleaseChannelHandler(c echo.Context) error {
	restChannel := new(models.ReleaseChannel)
	if err := c.Bind(restChannel); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Construct proto model and persist
	channelProto := &protos.ReleaseChannel{
		SupportedVersions: restChannel.SupportedVersions,
	}
	err := upgrade_client.CreateReleaseChannel(restChannel.Name, channelProto)
	if err != nil {
		return obsidian.HttpError(err)
	}

	// Return the ID of the created channel
	return c.JSON(http.StatusCreated, restChannel.Name)
}

func updateReleaseChannelHandler(c echo.Context) error {
	channelId := c.Param("channel_id")
	if channelId == "" {
		return noChannelIdError()
	}
	restChannel := new(models.ReleaseChannel)
	if err := c.Bind(restChannel); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Release channel name is immutable
	// This could change if release channels are keyed by UUID in their tables
	if restChannel.Name != channelId {
		return obsidian.HttpError(
			errors.New("Release channel name cannot be modified"),
			http.StatusBadRequest)
	}
	updatedChannelProto := &protos.ReleaseChannel{
		SupportedVersions: restChannel.SupportedVersions,
	}

	err := upgrade_client.UpdateReleaseChannel(channelId, updatedChannelProto)
	if err != nil {
		return obsidian.HttpError(err)
	}

	return c.NoContent(http.StatusOK)
}

func deleteReleaseChannelHandler(c echo.Context) error {
	channelId := c.Param("channel_id")
	if channelId == "" {
		return noChannelIdError()
	}

	err := upgrade_client.DeleteReleaseChannel(channelId)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func listTiersHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	tiers, err := upgrade_client.GetTiers(networkId, []string{})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
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
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	restTier := new(models.Tier)
	if err := c.Bind(restTier); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	tierProto := tierInfoModelToProto(restTier)

	err := upgrade_client.CreateTier(networkId, restTier.ID, tierProto)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	// Return the ID of the created tier
	return c.JSON(http.StatusCreated, restTier.ID)
}

func getTierHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierId := c.Param("tier_id")
	if tierId == "" {
		return noTierIdError()
	}

	tiers, err := upgrade_client.GetTiers(networkId, []string{tierId})
	if err != nil {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	tierProto, ok := tiers[tierId]
	if !ok {
		return obsidian.HttpError(
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
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierId := c.Param("tier_id")
	if tierId == "" {
		return noTierIdError()
	}
	restTier := new(models.Tier)
	if err := c.Bind(restTier); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	tierProto := tierInfoModelToProto(restTier)

	err := upgrade_client.UpdateTier(networkId, tierId, tierProto)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func deleteTierHandler(c echo.Context) error {
	networkId, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tierId := c.Param("tier_id")
	if tierId == "" {
		return noTierIdError()
	}

	err := upgrade_client.DeleteTier(networkId, tierId)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}
