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
