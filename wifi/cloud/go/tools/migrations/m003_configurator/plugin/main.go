/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations/m003_configurator/migration"
	"orc8r/wifi/cloud/go/tools/migrations/m003_configurator/plugin/types"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func main() {}

type wifiPlugin struct{}

func (*wifiPlugin) GetConfigMigrators() []migration.ConfigMigrator {
	return []migration.ConfigMigrator{
		&wifiNetworkMigrator{},
		&wifiGatewayMigrator{},
	}
}

func (*wifiPlugin) RunCustomMigrations(
	sc *squirrel.StmtCache,
	builder sqorc.StatementBuilder,
	migratedGatewayMetasByNetwork map[string]map[string]migration.MigratedGatewayMeta,
) error {
	return migrateMeshes(sc, builder, migratedGatewayMetasByNetwork)
}

func GetPlugin() migration.ConfiguratorMigrationPlugin {
	return &wifiPlugin{}
}

// migrators

type wifiNetworkMigrator struct{}
type wifiGatewayMigrator struct{}

func (*wifiNetworkMigrator) GetType() string {
	return "wifi_network"
}

func (*wifiNetworkMigrator) ToNewConfig(oldConfig []byte) ([]byte, error) {
	oldMsg := &types.WifiNetworkConfig{}
	err := migration.Unmarshal(oldConfig, oldMsg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newModel := &types.NetworkWifiConfigs{}
	migration.FillIn(oldMsg, newModel)
	return newModel.MarshalBinary()
}

func (*wifiGatewayMigrator) GetType() string {
	return "wifi_gateway"
}

func (*wifiGatewayMigrator) ToNewConfig(oldConfig []byte) ([]byte, error) {
	oldMsg := &types.WifiGatewayConfig{}
	err := migration.Unmarshal(oldConfig, oldMsg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newModel := &types.GatewayWifiConfigs{}
	migration.FillIn(oldMsg, newModel)
	return newModel.MarshalBinary()
}
