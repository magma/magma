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

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/golang/glog"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"
)

const (
	tableName = "cfg_entities"
	pkCol     = "pk"
	typeCol   = "type"
	configCol = "config"

	enodebEntityType = "cellular_enodeb"
)

type EnodebConfig struct {
	ConfigType    string          `json:"config_type"`
	ManagedConfig json.RawMessage `json:"managed_config"`
}

// This migration updates the config of all eNodeb cellular entities to the
// EnodebConfig struct instead of the old Config struct.
func main() {
	glog.Info("BEGIN MIGRATION")

	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
	glog.Infof("SQL_DRIVER: %s", dbDriver)
	glog.Infof("DATABASE_SOURCE: %s", dbSource)

	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		glog.Fatal(fmt.Errorf("could not open db connection: %w", err))
	}

	_, err = migrations.ExecInTx(db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, doMigration)
	if err != nil {
		glog.Fatalf("unexpected error occurred during migration: %s", err)
	}

	glog.Info("SUCCESS")
	glog.Info("END MIGRATION")
}

func doMigration(tx *sql.Tx) (interface{}, error) {
	sc := squirrel.NewStmtCache(tx)
	defer func() { _ = sc.Clear() }()
	builder := sqorc.GetSqlBuilder().RunWith(sc)

	b := builder.Select(pkCol, configCol).
		From(tableName).
		Where(squirrel.Eq{typeCol: enodebEntityType})
	sqlStr, args, _ := b.ToSql()
	glog.Infof("[RUN] %s %v", sqlStr, args)
	rows, err := b.RunWith(tx).Query()
	if err != nil {
		return nil, fmt.Errorf("error loading enodeb configs: %w", err)
	}

	defer sqorc.CloseRowsLogOnError(rows, "m014_enodeb_config")
	newConfsByPk := map[string][]byte{}
	for rows.Next() {
		var pk string
		var oldConf []byte

		if err = rows.Scan(&pk, &oldConf); err != nil {
			return nil, fmt.Errorf("error scanning enodeb config row: %w", err)
		}
		shouldMigrate, err := shouldMigrateConf(oldConf)
		if err != nil {
			return nil, err
		}
		if !shouldMigrate {
			continue
		}

		newConf := EnodebConfig{ConfigType: "MANAGED", ManagedConfig: oldConf}
		newConfBytes, err := json.Marshal(newConf)
		if err != nil {
			return nil, fmt.Errorf("error marshaling new enodeb config: %w", err)
		}
		newConfsByPk[pk] = newConfBytes
	}

	for pk, newConf := range newConfsByPk {
		bu := builder.Update(tableName).
			Set(configCol, newConf).
			Where(squirrel.Eq{pkCol: pk})
		sqlStr, args, _ := bu.ToSql()
		glog.Infof("[RUN] %s %v", sqlStr, args)
		_, err = bu.RunWith(tx).Exec()
		if err != nil {
			return nil, fmt.Errorf("error updating enodeb config %s: %w", pk, err)
		}
	}
	return nil, nil
}

// To keep the migration idempotent, we will check to see if it can deserialize
// into the new config type already.
func shouldMigrateConf(oldConf []byte) (bool, error) {
	parsedMessage := map[string]interface{}{}
	err := json.Unmarshal(oldConf, &parsedMessage)
	if err != nil {
		return false, fmt.Errorf("could not unmarshal legacy config: %w", err)
	}

	_, alreadyMigrated := parsedMessage["config_type"]
	return !alreadyMigrated, nil
}
