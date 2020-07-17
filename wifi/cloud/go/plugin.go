package main

import (
	"magma/orc8r/cloud/go/plugin"
	wifip "magma/wifi/cloud/go/plugin"
)

func main() {}

func GetOrchestratorPlugin() plugin.OrchestratorPlugin {
	return &wifip.WifiOrchestratorPlugin{}
}
