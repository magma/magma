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
	upgrade_client "magma/orc8r/cloud/go/services/upgrade"
	"magma/orc8r/cloud/go/services/upgrade/obsidian/models"
	"magma/orc8r/cloud/go/services/upgrade/protos"

	"github.com/labstack/echo"
)

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

	return c.NoContent(http.StatusOK)
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
	return c.NoContent(http.StatusNoContent)
}
