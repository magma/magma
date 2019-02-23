/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
)

// plugins must implement a main - these are expected to be empty
func main() {}

// GetOrchestratorPlugin is a function that all modules are expected to provide
// which returns an instance of the module's OrchestratorPlugin implementation
func GetOrchestratorPlugin() plugin.OrchestratorPlugin {
	return &pluginimpl.BaseOrchestratorPlugin{}
}
