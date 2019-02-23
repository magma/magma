/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"log"

	"magma/lte/cloud/go/tools/migrations/m003_fdd_tdd_configs/migration"
	"magma/orc8r/cloud/go/plugin"
)

func main() {
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})
	err := migration.Migrate()
	if err != nil {
		log.Fatal(err)
	}
}
