/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import "magma/orc8r/cloud/go/tools/migrations/m003_configurator/migration"

func main() {}

type plugin struct{}

func (*plugin) GetConfigMigrators() []migration.ConfigMigrator {
	return []migration.ConfigMigrator{}
}

func GetPlugin() migration.ConfiguratorMigrationPlugin {
	return &plugin{}
}
