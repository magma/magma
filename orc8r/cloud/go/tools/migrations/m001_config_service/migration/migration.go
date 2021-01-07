/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//go:generate bash -c "protoc -I . -I /usr/include -I $MAGMA_ROOT --go_out=plugins=grpc:. *.proto"

// DB migration script for the config service refactor. This migration moves
// network and gateway configuration management from magmad to the config
// service, and mesh configuration management from mesh service to config
// service.
// Note this script does not pull in imports from any other magma code, instead
// opting for some code duplication (mainly of constants). This is intentional,
// as the migration should do the exact same thing regardless of how the rest
// of the code has changed.
package migration

import (
	"bytes"
	"database/sql"
	"fmt"
	"sort"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"

	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// Redeclare table names
const NetworkConfigTable = "networks"
const GatewayConfigTable = "configs"
const MeshConfigTable = "mesh_config"

const NewConfigTable = "configurations"

// Redeclare config types and old config keys here for the same reason
const CellularNetworkType = "cellular_network"
const CellularGatewayType = "cellular_gateway"
const DnsdNetworkType = "dnsd_network"
const DnsdGatewayType = "dnsd_gateway" // Technically this is unused
const MagmadGatewayType = "magmad_gateway"
const MagmadNetworkType = "magmad_network"
const MeshType = "mesh"
const WifiNetworkType = "wifi_network"
const WifiGatewayType = "wifi_gateway"

const CellularConfigKey = "cellular"
const DnsConfigKey = "dns"
const MagmadConfigKey = "magmad"
const WifiConfigKey = "wifi"

var newNetworkTypesByOldKey = map[string]string{
	CellularConfigKey: CellularNetworkType,
	DnsConfigKey:      DnsdNetworkType,
	MagmadConfigKey:   MagmadNetworkType,
	WifiConfigKey:     WifiNetworkType,
}

var newGatewayTypesByOldKey = map[string]string{
	CellularConfigKey: CellularGatewayType,
	DnsConfigKey:      DnsdGatewayType,
	MagmadConfigKey:   MagmadGatewayType,
	WifiConfigKey:     WifiGatewayType,
}

// Entry point for the migration. Everything runs in 1 serializable postgres
// transaction so we don't end up in some weird half-migrated state.
func Migrate(dbDriver string, dbSource string) error {
	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		return fmt.Errorf("Could not open DB connection: %s", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error opening transaction: %s", err)
	}
	// Setting isolation level to serializable should result in the transaction
	// being rolled back if it detects any non-serializable concurrent changes.
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		return fmt.Errorf("Error setting transaction mode to serializable: %s", err)
	}

	glog.Error("Migrating network configs...")
	if err := MigrateNetworkConfigs(tx); err != nil {
		tx.Rollback()
		return err
	}

	glog.Error("Migrating gateway configs...")
	if err := MigrateGatewayConfigs(tx); err != nil {
		tx.Rollback()
		return err
	}

	glog.Error("Migrating mesh configs...")
	if err := MigrateMeshConfigs(tx); err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error committing transaction: %s", err)
	}
	return nil
}

func MigrateNetworkConfigs(tx *sql.Tx) error {
	networkConfigs, err := getAllConfigs(tx, NetworkConfigTable)
	if err != nil {
		return fmt.Errorf("Error getting existing network configs: %s", err)
	}

	sortedNetworkIds := getSortedConfigIds(networkConfigs)
	for _, networkId := range sortedNetworkIds {
		if err := migrateMagmadConfigs(tx, networkId, newNetworkTypesByOldKey, map[string]*Config{networkId: networkConfigs[networkId]}); err != nil {
			return fmt.Errorf("Error migrating network config for network %s: %s", networkId, err)
		}
	}
	return nil
}

func MigrateGatewayConfigs(tx *sql.Tx) error {
	networkIds, err := getAllNetworkIds(tx)
	if err != nil {
		return fmt.Errorf("Error getting existing network ids: %s", err)
	}

	for _, networkId := range networkIds {
		if err := migrateGatewayConfigsForNetwork(tx, networkId); err != nil {
			return fmt.Errorf("Error migrating gateway configs for network %s: %s", networkId, err)
		}
	}
	return nil
}

func MigrateMeshConfigs(tx *sql.Tx) error {
	networkIds, err := getAllNetworkIds(tx)
	if err != nil {
		return fmt.Errorf("Error getting existing network ids: %s", err)
	}

	for _, networkId := range networkIds {
		if err := migrateMeshConfigsForNetwork(tx, networkId); err != nil {
			return fmt.Errorf("Error migrating mesh configs for network %s: %s", networkId, err)
		}
	}
	return nil
}

func migrateGatewayConfigsForNetwork(tx *sql.Tx, networkId string) error {
	gatewayConfigTable := migrations.GetTableName(networkId, GatewayConfigTable)
	gatewayConfigs, err := getAllConfigs(tx, gatewayConfigTable)
	if err != nil {
		return err
	}

	return migrateMagmadConfigs(tx, networkId, newGatewayTypesByOldKey, gatewayConfigs)
}

func migrateMeshConfigsForNetwork(tx *sql.Tx, networkId string) error {
	err := initConfigTable(tx, networkId)
	if err != nil {
		return fmt.Errorf("Error initializing new config table: %s", err)
	}

	oldTable := migrations.GetTableName(networkId, MeshConfigTable)
	newTable := migrations.GetTableName(networkId, NewConfigTable)
	queryFormat := "INSERT INTO %s (type, key, value) VALUES($1, $2, $3) ON CONFLICT (type, key) DO UPDATE SET value=$4"
	stmt, err := tx.Prepare(fmt.Sprintf(queryFormat, newTable))
	if err != nil {
		return fmt.Errorf("Error preparing upsert statement: %s", err)
	}

	meshConfigs, err := migrations.GetAllValuesFromTable(tx, oldTable)
	if err != nil {
		return err
	}
	sortedMeshIds := getSortedMeshIds(meshConfigs)
	for _, meshId := range sortedMeshIds {
		_, err = stmt.Exec(MeshType, meshId, meshConfigs[meshId], meshConfigs[meshId])
		if err != nil {
			return err
		}
	}
	return nil
}

func migrateMagmadConfigs(tx *sql.Tx, networkId string, newTypesByOldKey map[string]string, configsByKey map[string]*Config) error {
	err := initConfigTable(tx, networkId)
	if err != nil {
		return fmt.Errorf("Error initializing new config table: %s", err)
	}

	table := migrations.GetTableName(networkId, NewConfigTable)
	queryFormat := "INSERT INTO %s (type, key, value) VALUES ($1, $2, $3) ON CONFLICT (type, key) DO UPDATE SET value=$4"
	stmt, err := tx.Prepare(fmt.Sprintf(queryFormat, table))
	if err != nil {
		return fmt.Errorf("Error preparing upsert statement: %s", err)
	}
	defer stmt.Close()

	sortedKeys := getSortedConfigIds(configsByKey)
	for _, newKey := range sortedKeys {
		err = migrateMagmadConfig(stmt, newKey, newTypesByOldKey, configsByKey[newKey])
		if err != nil {
			return fmt.Errorf("Error migrating magmad config for %s in network %s: %s", newKey, networkId, err)
		}
	}
	return nil
}

func migrateMagmadConfig(stmt *sql.Stmt, newKey string, newTypesByOldKey map[string]string, config *Config) error {
	sortedConfigKeys := getSortedConfigKeys(config)
	for _, configKey := range sortedConfigKeys {
		newConfigType, ok := newTypesByOldKey[configKey]
		if !ok {
			return fmt.Errorf("No new config type defined for magmad key %s", configKey)
		}

		// Skip empty config values
		val := config.ConfigsByKey[configKey]
		if val == nil || len(val) == 0 {
			continue
		}

		_, err := stmt.Exec(newConfigType, newKey, val, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAllNetworkIds(tx *sql.Tx) ([]string, error) {
	return migrations.GetAllKeysFromTable(tx, NetworkConfigTable)
}

func getAllConfigs(tx *sql.Tx, table string) (map[string]*Config, error) {
	marshaledConfigs, err := migrations.GetAllValuesFromTable(tx, table)
	if err != nil {
		return nil, err
	}
	return unmarshalConfigs(marshaledConfigs)
}

func unmarshalConfigs(marshaledConfigs map[string][]byte) (map[string]*Config, error) {
	ret := make(map[string]*Config, len(marshaledConfigs))
	for k, v := range marshaledConfigs {
		msg := &Config{}
		err := unmarshal(v, msg)
		if err != nil {
			return nil, err
		}
		ret[k] = msg
	}
	return ret, nil
}

func unmarshal(bt []byte, msg proto.Message) error {
	return (&jsonpb.Unmarshaler{AllowUnknownFields: true}).Unmarshal(bytes.NewBuffer(bt), msg)
}

func initConfigTable(tx *sql.Tx, networkId string) error {
	queryFormat := `
		CREATE TABLE IF NOT EXISTS %s
		(
			type text NOT NULL,
			key text NOT NULL,
			value bytea,
			version INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (type, key)
		)
	`
	table := migrations.GetTableName(networkId, NewConfigTable)
	_, err := tx.Exec(fmt.Sprintf(queryFormat, table))
	return err
}

// ---
// For deterministic test results
// ---

func getSortedConfigKeys(cfg *Config) []string {
	keys := make([]string, 0, len(cfg.ConfigsByKey))
	for k := range cfg.ConfigsByKey {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getSortedConfigIds(networkConfigs map[string]*Config) []string {
	keys := make([]string, 0, len(networkConfigs))
	for k := range networkConfigs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getSortedMeshIds(meshConfigs map[string][]byte) []string {
	keys := make([]string, 0, len(meshConfigs))
	for k := range meshConfigs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
