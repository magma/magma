package main

import (
	"flag"

	"magma/fbinternal/cloud/go/services/download/servicers"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/lib/go/registry"
)

func main() {
	flag.Parse()
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})
	registry.MustPopulateServices()

	servicers.Run()
}
