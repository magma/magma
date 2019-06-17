/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

//go:generate bash -c "protoc -I . -I /usr/include --proto_path=$MAGMA_ROOT --go_out=plugins=grpc:. *.proto"
package main

import (
	"fmt"

	"magma/lte/cloud/go/tools/migrations/m003_configurator/plugin/types"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations/m003_configurator/migration"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const policyTable = "policydb"
const baseNameTable = "base_names"

func main() {}

type plugin struct{}

func GetPlugin() migration.ConfiguratorMigrationPlugin {
	return &plugin{}
}

func (*plugin) GetConfigMigrators() []migration.ConfigMigrator {
	return []migration.ConfigMigrator{
		&cellularNetworkMigrator{},
		&cellularGatewayMigrator{},
	}
}

func (*plugin) RunCustomMigrations(
	sc *squirrel.StmtCache,
	builder sqorc.StatementBuilder,
	migratedGatewayMetasByNetwork map[string]map[string]migration.MigratedGatewayMeta,
) error {
	return migratePolicydb(sc, builder, funk.Keys(migratedGatewayMetasByNetwork).([]string))
}

// migrators

type cellularNetworkMigrator struct{}
type cellularGatewayMigrator struct{}

func (*cellularNetworkMigrator) GetType() string {
	return "cellular_network"
}

func (*cellularNetworkMigrator) ToNewConfig(oldConfig []byte) ([]byte, error) {
	oldMsg := &types.CellularNetworkConfig{}
	err := migration.Unmarshal(oldConfig, oldMsg)
	if err != nil {
		return nil, err
	}

	newModel := &types.NetworkCellularConfigs{}
	migration.FillIn(oldMsg, newModel)
	newModel.FegNetworkID = oldMsg.FegNetworkId
	newModel.Epc.RelayEnabled = oldMsg.Epc.RelayEnabled

	if oldMsg.Epc.NetworkServices != nil {
		for _, serviceEnum := range oldMsg.Epc.NetworkServices {
			serviceName, ok := networkServiceEnumToNameMap[serviceEnum]
			if !ok {
				return nil, fmt.Errorf("Unknown network service enum: %s", serviceEnum)
			}
			newModel.Epc.NetworkServices = append(newModel.Epc.NetworkServices, serviceName)
		}
	}

	return newModel.MarshalBinary()
}

var networkServiceEnumToNameMap = map[types.NetworkEPCConfig_NetworkServices]string{
	types.NetworkEPCConfig_METERING:    "metering",
	types.NetworkEPCConfig_DPI:         "dpi",
	types.NetworkEPCConfig_ENFORCEMENT: "policy_enforcement",
}

func (*cellularGatewayMigrator) GetType() string {
	return "cellular_gateway"
}

func (*cellularGatewayMigrator) ToNewConfig(oldConfig []byte) ([]byte, error) {
	oldMsg := &types.CellularGatewayConfig{}
	err := migration.Unmarshal(oldConfig, oldMsg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newModel := &types.GatewayCellularConfigs{}
	migration.FillIn(oldMsg, newModel)
	return newModel.MarshalBinary()

}
