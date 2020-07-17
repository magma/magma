package main

import (
	"flag"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/lib/go/registry"
	"orc8r/fbinternal/cloud/go/services/download/servicers"
)

func main() {
	flag.Parse()
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})
	registry.MustPopulateServices()

	servicers.Run()
}
