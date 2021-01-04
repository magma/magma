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

//go:generate bash -c "protoc -I . -I /usr/include -I $MAGMA_ROOT/.. --go_out=plugins=grpc:. *.proto"

// DB migration script to clean up old magmad config tables.
// Entries in the network table are migrated from the old network configuration
// type to the new network record by pulling the "magmad" entry from the
// old namespaced config.
// As a cleanup, this will also delete all magmad network configs from the new
// config service tables.
// All old gateway and mesh config tables are deleted as a cleanup if
// `shouldDropTables` is true.
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
)

// Redeclare table names
const NetworkTable = "networks"
const GatewayConfigTable = "configs"
const MeshConfigTable = "mesh_config"
const NewConfigTable = "configurations"

const MagmadNetworkType = "magmad_network"

func Migrate(dbDriver string, dbSource string, shouldDropTables bool) error {
	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		return fmt.Errorf("Could not open DB connection: %s", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error opening transaction: %s", err)
	}

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	if err != nil {
		return fmt.Errorf("Error setting transaction mode to serializable: %s", err)
	}

	glog.Error("Migrating network records...")
	if err := MigrateNetworkConfigsToRecords(tx); err != nil {
		tx.Rollback()
		return err
	}

	glog.Error("Deleting magmad network configs...")
	if err := DeleteMagmadNetworkConfigs(tx); err != nil {
		tx.Rollback()
		return err
	}

	if shouldDropTables {
		glog.Error("Deleting old gateway configs...")
		if err := DeleteOldGatewayConfigTables(tx); err != nil {
			tx.Rollback()
			return err
		}

		glog.Errorf("Deleting old mesh configs...")
		if err := DeleteOldMeshConfigTables(tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error committing transaction: %s", err)
	}
	return nil
}

func MigrateNetworkConfigsToRecords(tx *sql.Tx) error {
	marshaledLegacyConfigs, err := migrations.GetAllValuesFromTable(tx, NetworkTable)
	if err != nil {
		return fmt.Errorf("Could not load existing network configs: %s", err)
	}

	valuesToWrite, err := getNewNetworkRecordsToWrite(marshaledLegacyConfigs)
	if err != nil {
		return err
	}
	queryFormat := "UPDATE %s SET value = $1 WHERE key = $2"
	stmt, err := tx.Prepare(fmt.Sprintf(queryFormat, NetworkTable))
	if err != nil {
		return fmt.Errorf("Error preparing upsert statement: %s", err)
	}
	defer stmt.Close()

	sortedNetworkIds := getSortedKeys(valuesToWrite)
	for _, networkId := range sortedNetworkIds {
		_, err = stmt.Exec(valuesToWrite[networkId], networkId)
		if err != nil {
			return fmt.Errorf("Error updating network %s: %s", networkId, err)
		}
	}
	return nil
}

func DeleteMagmadNetworkConfigs(tx *sql.Tx) error {
	networkIds, err := migrations.GetAllKeysFromTable(tx, NetworkTable)
	if err != nil {
		return fmt.Errorf("Could not load network IDs: %s", err)
	}

	queryFormat := "DELETE FROM %s WHERE type = $1 AND key = $2"
	for _, networkId := range networkIds {
		tableName := migrations.GetTableName(networkId, NewConfigTable)
		_, err := tx.Exec(fmt.Sprintf(queryFormat, tableName), MagmadNetworkType, networkId)
		if err != nil {
			return fmt.Errorf("Failed to delete magmad network configs for network %s: %s", networkId, err)
		}
	}
	return nil
}

func DeleteOldGatewayConfigTables(tx *sql.Tx) error {
	return dropTableForAllNetworks(tx, GatewayConfigTable)
}

func DeleteOldMeshConfigTables(tx *sql.Tx) error {
	return dropTableForAllNetworks(tx, MeshConfigTable)
}

func getNewNetworkRecordsToWrite(oldNetworkConfigsOrRecords map[string][]byte) (map[string][]byte, error) {
	ret := make(map[string][]byte, len(oldNetworkConfigsOrRecords))
	for networkId, marshaledConfigOrRecord := range oldNetworkConfigsOrRecords {
		legacyConfig, _, err := deserializeToConfigOrRecord(networkId, marshaledConfigOrRecord)
		if err != nil {
			return nil, err
		}
		// Network is already migrated
		if legacyConfig == nil {
			continue
		}

		newRecord, err := getNewRecord(networkId, legacyConfig)
		if err != nil {
			return nil, err
		}
		ret[networkId] = newRecord
	}
	return ret, nil
}

// Only 1 on *Config or *Record return types will be non-nil
func deserializeToConfigOrRecord(networkId string, bt []byte) (*Config, *Record, error) {
	legacyConfig, err := deserializeToConfig(bt)
	if err == nil {
		return legacyConfig, nil, nil
	}
	record, err := deserializeToRecord(bt)
	if err == nil {
		glog.Errorf("Network %s is already migrated", networkId)
		return nil, record, nil
	} else {
		return nil, nil, fmt.Errorf("Could not deserialize network %s to legacy config or new record - manual intervention required.", networkId)
	}
}

func deserializeToConfig(bt []byte) (*Config, error) {
	legacyConfig := &Config{}
	err := (&jsonpb.Unmarshaler{AllowUnknownFields: false}).Unmarshal(bytes.NewBuffer(bt), legacyConfig)
	return legacyConfig, err
}

func deserializeToRecord(bt []byte) (*Record, error) {
	record := &Record{}
	err := (&jsonpb.Unmarshaler{AllowUnknownFields: false}).Unmarshal(bytes.NewBuffer(bt), record)
	return record, err
}

func getNewRecord(networkId string, legacyConfig *Config) ([]byte, error) {
	if legacyConfig.ConfigsByKey == nil || legacyConfig.ConfigsByKey["magmad"] == nil {
		glog.Errorf("Empty config entry for magmad on network %s, writing %s as network name", networkId, networkId)

		record := &Record{Name: networkId}
		var buff bytes.Buffer
		err := (&jsonpb.Marshaler{EmitDefaults: true, Indent: " "}).Marshal(&buff, record)
		return buff.Bytes(), err
	} else {
		// Check that magmad network config can be deserialized to a record
		_, err := deserializeToRecord(legacyConfig.ConfigsByKey["magmad"])
		if err != nil {
			return nil, fmt.Errorf("Could not deserialize magmad network config to network record on network %s, manual intervention required: %s", networkId, err)
		}
		return legacyConfig.ConfigsByKey["magmad"], nil
	}
}

func dropTableForAllNetworks(tx *sql.Tx, tableName string) error {
	networkIds, err := migrations.GetAllKeysFromTable(tx, NetworkTable)
	if err != nil {
		return fmt.Errorf("Could not load network IDs: %s", err)
	}

	for _, networkId := range networkIds {
		query := fmt.Sprintf("DROP TABLE %s IF EXISTS", migrations.GetTableName(networkId, tableName))
		_, err = tx.Exec(query)
		if err != nil {
			return fmt.Errorf("Error dropping table %s for network %s: %s", tableName, networkId, err)
		}
	}
	return nil
}

// ---
// For deterministic test results
// ---

func getSortedKeys(in map[string][]byte) []string {
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
