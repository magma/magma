/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"magma/feg/cloud/go/tools/migrations/m003_configurator/plugin/types"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations/m003_configurator/migration"
	"magma/orc8r/lib/go/protos"

	"github.com/Masterminds/squirrel"
)

func main() {}

type plugin struct{}

func (*plugin) GetConfigMigrators() []migration.ConfigMigrator {
	return []migration.ConfigMigrator{
		&fegNetworkMigrator{},
		&fegGatewayMigrator{},
	}
}

func (*plugin) RunCustomMigrations(
	sc *squirrel.StmtCache,
	builder sqorc.StatementBuilder,
	migratedGatewayMetasByNetwork map[string]map[string]migration.MigratedGatewayMeta,
) error {
	return nil
}

func GetPlugin() migration.ConfiguratorMigrationPlugin {
	return &plugin{}
}

// migrators

type fegNetworkMigrator struct{}

func (f *fegNetworkMigrator) GetType() string {
	return "federation_network"
}

func (f *fegNetworkMigrator) ToNewConfig(oldConfig []byte) ([]byte, error) {
	oldMsg := &types.Config{}
	err := migration.Unmarshal(oldConfig, oldMsg)
	if err != nil {
		return nil, err
	}

	newModel := &types.NetworkFederationConfigs{}
	migration.FillIn(oldMsg, newModel)
	if newModel.S6a == nil {
		newModel.S6a = &types.NetworkFederationConfigsS6a{Server: &types.DiameterClientConfigs{}}
	} else if newModel.S6a.Server == nil {
		newModel.S6a.Server = &types.DiameterClientConfigs{}
	}
	if newModel.Hss == nil {
		newModel.Hss = &types.NetworkFederationConfigsHss{
			DefaultSubProfile: &types.SubscriptionProfile{},
			Server:            &types.DiameterServerConfigs{},
			SubProfiles:       make(map[string]types.SubscriptionProfile),
		}
	} else {
		if newModel.Hss.DefaultSubProfile == nil {
			newModel.Hss.DefaultSubProfile = &types.SubscriptionProfile{}
		}
		if newModel.Hss.Server == nil {
			newModel.Hss.Server = &types.DiameterServerConfigs{}
		}
		if newModel.Hss.SubProfiles == nil {
			newModel.Hss.SubProfiles = make(map[string]types.SubscriptionProfile)
		}
	}
	if newModel.Gx == nil {
		newModel.Gx = &types.NetworkFederationConfigsGx{Server: &types.DiameterClientConfigs{}}
	} else if newModel.Gx.Server == nil {
		newModel.Gx.Server = &types.DiameterClientConfigs{}
	}
	if newModel.Gy == nil {
		newModel.Gy = &types.NetworkFederationConfigsGy{Server: &types.DiameterClientConfigs{}}
	} else if newModel.Gy.Server == nil {
		newModel.Gy.Server = &types.DiameterClientConfigs{}
	}
	if newModel.Swx == nil {
		newModel.Swx = &types.NetworkFederationConfigsSwx{Server: &types.DiameterClientConfigs{}}
	} else if newModel.Swx.Server == nil {
		newModel.Swx.Server = &types.DiameterClientConfigs{}
	}
	if newModel.Health == nil {
		newModel.Health = &types.NetworkFederationConfigsHealth{}
	}
	if newModel.EapAka == nil {
		newModel.EapAka = &types.NetworkFederationConfigsEapAka{}
	}
	if newModel.AaaServer == nil {
		newModel.AaaServer = &types.NetworkFederationConfigsAaaServer{}
	}
	protos.FillIn(oldMsg.S6A, newModel.S6a)
	protos.FillIn(oldMsg.Hss, newModel.Hss)
	protos.FillIn(oldMsg.Gx, newModel.Gx)
	protos.FillIn(oldMsg.Gy, newModel.Gy)
	protos.FillIn(oldMsg.Swx, newModel.Swx)
	protos.FillIn(oldMsg.Health, newModel.Health)
	protos.FillIn(oldMsg.EapAka, newModel.EapAka)
	protos.FillIn(oldMsg.AaaServer, newModel.AaaServer)
	if newModel.ServedNetworkIds == nil {
		newModel.ServedNetworkIds = []string{}
	}

	return newModel.MarshalBinary()
}

type fegGatewayMigrator struct{}

func (*fegGatewayMigrator) GetType() string {
	return "federation_gateway"
}

func (*fegGatewayMigrator) ToNewConfig(oldConfig []byte) ([]byte, error) {
	oldMsg := &types.Config{}
	err := migration.Unmarshal(oldConfig, oldMsg)
	if err != nil {
		return nil, err
	}

	newModel := &types.GatewayFegConfigs{}
	migration.FillIn(oldMsg, newModel)
	if newModel.S6a == nil {
		newModel.S6a = &types.NetworkFederationConfigsS6a{Server: &types.DiameterClientConfigs{}}
	} else if newModel.S6a.Server == nil {
		newModel.S6a.Server = &types.DiameterClientConfigs{}
	}
	if newModel.Hss == nil {
		newModel.Hss = &types.NetworkFederationConfigsHss{
			DefaultSubProfile: &types.SubscriptionProfile{},
			Server:            &types.DiameterServerConfigs{},
			SubProfiles:       make(map[string]types.SubscriptionProfile),
		}
	} else {
		if newModel.Hss.DefaultSubProfile == nil {
			newModel.Hss.DefaultSubProfile = &types.SubscriptionProfile{}
		}
		if newModel.Hss.Server == nil {
			newModel.Hss.Server = &types.DiameterServerConfigs{}
		}
		if newModel.Hss.SubProfiles == nil {
			newModel.Hss.SubProfiles = make(map[string]types.SubscriptionProfile)
		}
	}
	if newModel.Gx == nil {
		newModel.Gx = &types.NetworkFederationConfigsGx{Server: &types.DiameterClientConfigs{}}
	} else if newModel.Gx.Server == nil {
		newModel.Gx.Server = &types.DiameterClientConfigs{}
	}
	if newModel.Gy == nil {
		newModel.Gy = &types.NetworkFederationConfigsGy{Server: &types.DiameterClientConfigs{}}
	} else if newModel.Gy.Server == nil {
		newModel.Gy.Server = &types.DiameterClientConfigs{}
	}
	if newModel.Swx == nil {
		newModel.Swx = &types.NetworkFederationConfigsSwx{Server: &types.DiameterClientConfigs{}}
	} else if newModel.Swx.Server == nil {
		newModel.Swx.Server = &types.DiameterClientConfigs{}
	}
	if newModel.Health == nil {
		newModel.Health = &types.NetworkFederationConfigsHealth{}
	}
	if newModel.EapAka == nil {
		newModel.EapAka = &types.NetworkFederationConfigsEapAka{}
	}
	if newModel.AaaServer == nil {
		newModel.AaaServer = &types.NetworkFederationConfigsAaaServer{}
	}
	protos.FillIn(oldMsg.S6A, newModel.S6a)
	protos.FillIn(oldMsg.Hss, newModel.Hss)
	protos.FillIn(oldMsg.Gx, newModel.Gx)
	protos.FillIn(oldMsg.Gy, newModel.Gy)
	protos.FillIn(oldMsg.Swx, newModel.Swx)
	protos.FillIn(oldMsg.Health, newModel.Health)
	protos.FillIn(oldMsg.EapAka, newModel.EapAka)
	protos.FillIn(oldMsg.AaaServer, newModel.AaaServer)
	if newModel.ServedNetworkIds == nil {
		newModel.ServedNetworkIds = []string{}
	}

	return newModel.MarshalBinary()
}
