package main

import (
	fbinternalp "magma/fbinternal/cloud/go/plugin"
	"magma/orc8r/cloud/go/plugin"
)

func main() {}

func GetOrchestratorPlugin() plugin.OrchestratorPlugin {
	return &fbinternalp.FbinternalOrchestratorPlugin{}
}
