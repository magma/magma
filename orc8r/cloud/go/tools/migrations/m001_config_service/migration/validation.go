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

// DB validation script to run as a sanity check after a completed migration.
// This loads all existing configs from old magmad and mesh service tables
// and verifies that they all exist in the new config service tables.
package migration

import (
	"bytes"
	"database/sql"
	"fmt"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"

	"github.com/golang/glog"
)

type TypeAndKey struct {
	Type, Key string
}

func Validate(dbDriver string, dbSource string) error {
	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		return fmt.Errorf("Could not open DB connection: %s", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error opening transaction: %s", err)
	}

	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE READ ONLY DEFERRABLE")
	if err != nil {
		return fmt.Errorf("Error setting transaction mode to serializable: %s", err)
	}

	allNetworkConfigs, err := getAllConfigs(tx, NetworkConfigTable)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Error loading networks: %s", err)
	}

	sortedNetworkIds := getSortedConfigIds(allNetworkConfigs)
	for _, networkId := range sortedNetworkIds {
		glog.V(2).Infof("Validating network %s", networkId)
		if err := ValidateNetwork(tx, networkId, allNetworkConfigs[networkId]); err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error committing read-only transaction: %s", err)
	}
	return nil
}

func ValidateNetwork(tx *sql.Tx, networkId string, oldNetworkConfig *Config) error {
	gatewayConfigTable := migrations.GetTableName(networkId, GatewayConfigTable)
	meshConfigTable := migrations.GetTableName(networkId, MeshConfigTable)
	oldGatewayConfigs, err := getAllConfigs(tx, gatewayConfigTable)
	if err != nil {
		return err
	}
	oldMeshConfigs, err := migrations.GetAllValuesFromTable(tx, meshConfigTable)
	if err != nil {
		return err
	}
	newConfigs, err := loadAllNewConfigs(tx, networkId)
	if err != nil {
		return err
	}

	if err := ValidateNetworkConfigs(networkId, oldNetworkConfig, newConfigs); err != nil {
		return fmt.Errorf("Failed validation for network %s: %s", networkId, err)
	}
	if err := ValidateGatewayConfigs(oldGatewayConfigs, newConfigs); err != nil {
		return fmt.Errorf("Failed validation for gateways for network %s: %s", networkId, err)
	}
	for meshId, oldVal := range oldMeshConfigs {
		if err := ValidateMeshConfigs(meshId, oldVal, newConfigs); err != nil {
			return fmt.Errorf("Failed vaidation for mesh %s on network %s: %s", networkId, meshId, err)
		}
	}
	return nil
}

func ValidateNetworkConfigs(networkId string, oldNetworkConfig *Config, newConfigs map[TypeAndKey][]byte) error {
	return validateMagmadConfigs(networkId, oldNetworkConfig, newConfigs, newNetworkTypesByOldKey)
}

func ValidateGatewayConfigs(oldGatewayConfigs map[string]*Config, newConfigs map[TypeAndKey][]byte) error {
	for gatewayId, gatewayConfig := range oldGatewayConfigs {
		if err := validateMagmadConfigs(gatewayId, gatewayConfig, newConfigs, newGatewayTypesByOldKey); err != nil {
			return err
		}
	}
	return nil
}

func ValidateMeshConfigs(meshId string, oldMeshConfig []byte, newConfigs map[TypeAndKey][]byte) error {
	newVal, exists := newConfigs[TypeAndKey{Type: MeshType, Key: meshId}]
	if !exists {
		return fmt.Errorf("No corresponding new config for mesh %s", meshId)
	}
	if !bytes.Equal(oldMeshConfig, newVal) {
		return fmt.Errorf("New config for mesh %s not equal to old", meshId)
	}
	return nil
}

func validateMagmadConfigs(id string, oldConfig *Config, newConfigs map[TypeAndKey][]byte, newTypesByOldKey map[string]string) error {
	for oldKey, oldVal := range oldConfig.ConfigsByKey {
		newType, ok := newTypesByOldKey[oldKey]
		if !ok {
			return fmt.Errorf("Unrecognized old config key %s for ID %s", oldKey, id)
		}

		// Skip empty config values
		newVal, exists := newConfigs[TypeAndKey{Type: newType, Key: id}]
		if newVal == nil || len(newVal) == 0 {
			continue
		}

		if !exists {
			return fmt.Errorf("Old config key %s not migrated for id %s", oldKey, id)
		}
		if !bytes.Equal(oldVal, newVal) {
			return fmt.Errorf("Old config value for key %s not equal to new for id %s", oldKey, id)
		}
	}
	return nil
}

func loadAllNewConfigs(tx *sql.Tx, networkId string) (map[TypeAndKey][]byte, error) {
	table := migrations.GetTableName(networkId, NewConfigTable)
	query := fmt.Sprintf("SELECT type, key, value FROM %s", table)
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	ret := map[TypeAndKey][]byte{}
	for rows.Next() {
		var configType, key string
		var value []byte

		err = rows.Scan(&configType, &key, &value)
		if err != nil {
			return nil, err
		}
		ret[TypeAndKey{Type: configType, Key: key}] = value
	}
	return ret, nil
}
