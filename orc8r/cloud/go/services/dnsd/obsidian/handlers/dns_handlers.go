/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	cfgObsidian "magma/orc8r/cloud/go/services/config/obsidian"
	magmad_handlers "magma/orc8r/cloud/go/services/magmad/obsidian/handlers"
)

const (
	ConfigKey         = "dns"
	NetworkConfigPath = magmad_handlers.ConfigureNetwork + "/" + ConfigKey
)

// GetObsidianHandlers returns all obsidian handlers for dnsd
func GetObsidianHandlers() []obsidian.Handler {
	return cfgObsidian.GetCRUDNetworkConfigHandlers(NetworkConfigPath, orc8r.DnsdNetworkType, &models.NetworkDNSConfig{})
}
