/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */

package serde

import (
	"encoding/json"

	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/upgrade"
	"magma/orc8r/cloud/go/services/upgrade/obsidian/models"
)

type ReleaseChannelVersionsManager struct{}

func (*ReleaseChannelVersionsManager) GetDomain() string {
	return configurator.SerdeDomain
}

func (*ReleaseChannelVersionsManager) GetType() string {
	return upgrade.ReleaseChannelType
}

func (*ReleaseChannelVersionsManager) Serialize(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func (*ReleaseChannelVersionsManager) Deserialize(in []byte) (interface{}, error) {
	ret := []string{}
	err := json.Unmarshal(in, &ret)
	return ret, err
}

type NetworkTierConfigManager struct{}

func (*NetworkTierConfigManager) GetDomain() string {
	return configurator.SerdeDomain
}

func (*NetworkTierConfigManager) GetType() string {
	return upgrade.NetworkTierType
}

func (*NetworkTierConfigManager) Serialize(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func (*NetworkTierConfigManager) Deserialize(in []byte) (interface{}, error) {
	ret := &models.Tier{}
	err := json.Unmarshal(in, &ret)
	return ret, err
}
