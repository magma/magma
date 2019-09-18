/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package main

import (
	"magma/orc8r/cloud/go/plugin"
	devmandp "orc8r/devmand/cloud/go/plugin"
)

func main() {}

// GetOrchestratorPlugin gets the orchestrator plugin for devmand
func GetOrchestratorPlugin() plugin.OrchestratorPlugin {
	return &devmandp.DevmandOrchestratorPlugin{}
}
