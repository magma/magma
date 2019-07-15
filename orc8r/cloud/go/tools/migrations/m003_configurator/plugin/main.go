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
	"magma/orc8r/cloud/go/tools/migrations/m003_configurator/plugin/types"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func main() {}

type orcPlugin struct{}

func (*orcPlugin) GetConfigMigrators() []migration.ConfigMigrator {
	return []migration.ConfigMigrator{
		&networkDnsMigrator{},
		&magmadGatewayMigrator{},
	}
}

func (*orcPlugin) RunCustomMigrations(
	sc *squirrel.StmtCache,
	builder sqorc.StatementBuilder,
	migratedGatewayMetasByNetwork map[string]map[string]migration.MigratedGatewayMeta,
) error {
	if err := migrateReleaseChannels(sc, builder); err != nil {
		return errors.WithStack(err)
	}
	if err := migrateUpgradeTiers(sc, builder, migratedGatewayMetasByNetwork); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func GetPlugin() migration.ConfiguratorMigrationPlugin {
	return &orcPlugin{}
}

// migrators

type networkDnsMigrator struct{}

func (*networkDnsMigrator) GetType() string {
	return "dnsd_network"
}

func (*networkDnsMigrator) ToNewConfig(oldConfig []byte) ([]byte, error) {
	oldMsg := &types.NetworkDNSConfig{}
	err := migration.Unmarshal(oldConfig, oldMsg)
	if err != nil {
		return nil, err
	}

	newModel := &types.NewNetworkDNSConfig{}
	migration.FillIn(oldMsg, newModel)
	return newModel.MarshalBinary()
}

type magmadGatewayMigrator struct{}

func (*magmadGatewayMigrator) GetType() string {
	return "magmad_gateway"
}

func (*magmadGatewayMigrator) ToNewConfig(oldConfig []byte) ([]byte, error) {
	oldMsg := &types.OldMagmadGatewayConfig{}
	err := migration.Unmarshal(oldConfig, oldMsg)
	if err != nil {
		return nil, err
	}

	newModel := &types.MagmadGatewayConfig{}
	migration.FillIn(oldMsg, newModel)
	return newModel.MarshalBinary()
}
