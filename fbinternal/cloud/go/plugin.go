package main

import (
	"magma/orc8r/cloud/go/plugin"
	fbinternalp "orc8r/fbinternal/cloud/go/plugin"
)

func main() {}

func GetOrchestratorPlugin() plugin.OrchestratorPlugin {
	return &fbinternalp.FbinternalOrchestratorPlugin{}
}
