/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	fegplugin "magma/feg/cloud/go/plugin"
	"magma/orc8r/cloud/go/plugin"
)

func main() {}

func GetOrchestratorPlugin() plugin.OrchestratorPlugin {
	return &fegplugin.FegOrchestratorPlugin{}
}
