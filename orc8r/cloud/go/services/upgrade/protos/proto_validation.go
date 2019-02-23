/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"errors"
)

func ValidateCreateOrUpdateReleaseChannelReq(req *CreateOrUpdateReleaseChannelRequest) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	if req.GetChannelName() == "" {
		return errors.New("Channel must be specified")
	}
	if req.GetChannel() == nil {
		return errors.New("Channel body must be specified")
	}
	return nil
}

func ValidateGetReleaseChannelRequest(req *GetReleaseChannelRequest) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	if req.GetChannelName() == "" {
		return errors.New("Channel name must be specified")
	}
	return nil
}

func ValidateDeleteReleaseChannelReq(req *DeleteReleaseChannelRequest) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	if req.GetChannelName() == "" {
		return errors.New("Channel name must be specified")
	}
	return nil
}

func ValidateGetTiersReq(req *GetTiersRequest) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	if req.GetNetworkId() == "" {
		return errors.New("NetworkID must be specified")
	}
	return nil
}

func ValidateCreateTierReq(req *CreateTierRequest) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	if req.GetNetworkId() == "" {
		return errors.New("NetworkID must be specified")
	}
	if req.GetTierId() == "" {
		return errors.New("Tier ID must be specified")
	}
	if req.GetTierInfo() == nil {
		return errors.New("Tier info must be specified")
	}
	return ValidateTierInfo(req.GetTierInfo())
}

func ValidateUpdateTierReq(req *UpdateTierRequest) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	if req.GetNetworkId() == "" {
		return errors.New("NetworkID must be specified")
	}
	if req.GetTierId() == "" {
		return errors.New("Tier ID must be specified")
	}
	if req.GetUpdatedTier() == nil {
		return errors.New("Updated tier info must be specified")
	}
	return ValidateTierInfo(req.GetUpdatedTier())
}

func ValidateDeleteTierReq(req *DeleteTierRequest) error {
	if req == nil {
		return errors.New("Request is nil")
	}
	if req.GetNetworkId() == "" {
		return errors.New("NetworkID must be specified")
	}
	if req.GetTierIdToDelete() == "" {
		return errors.New("Tier ID to delete must be specified")
	}
	return nil
}

func ValidateTierInfo(info *TierInfo) error {
	if info == nil {
		return errors.New("Tier info is nil")
	}
	if info.GetName() == "" {
		return errors.New("Name must be specified")
	}
	if info.GetVersion() == "" {
		return errors.New("Version must be specified")
	}
	return nil
}
