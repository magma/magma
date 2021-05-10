package main

import (
	"flag"

	"magma/fbinternal/cloud/go/services/download/servicers"
	"magma/orc8r/lib/go/registry"
)

func main() {
	flag.Parse()
	registry.MustPopulateServices()

	servicers.Run()
}
