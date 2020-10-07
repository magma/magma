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

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"

	"github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	tableName = "cfg_entities"
	pkCol     = "pk"
	typeCol   = "type"
	configCol = "config"

	subscriberType = "subscriber"
)

type SubscriberConfig struct {
	Lte json.RawMessage `json:"lte"`
}

// This migration updates the config of all subscriber entities to the
// SubscriberConfig struct instead of the old LteSubscription struct.
func main() {
	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		glog.Fatal(errors.Wrap(err, "could not open db connection"))
	}

	_, err = migrations.ExecInTx(db, &sql.TxOptions{Isolation: sql.LevelSerializable}, nil, doMigration)
	if err != nil {
		glog.Fatalf("unexpected error occurred during migration: %s", err)
	}

	glog.Info("Subscriber migration successfully completed")
}

func doMigration(tx *sql.Tx) (interface{}, error) {
	sc := squirrel.NewStmtCache(tx)
	defer func() { _ = sc.Clear() }()
	builder := sqorc.GetSqlBuilder().RunWith(sc)

	rows, err := builder.Select(pkCol, configCol).
		From(tableName).
		Where(squirrel.Eq{typeCol: subscriberType}).
		Query()
	if err != nil {
		return nil, errors.Wrap(err, "error loading subscriber configs")
	}

	newConfsByPk := map[string][]byte{}
	for rows.Next() {
		var pk string
		var oldConf []byte

		if err = rows.Scan(&pk, &oldConf); err != nil {
			return nil, errors.Wrap(err, "error scanning subscriber row")
		}
		shouldMigrate, err := shouldMigrateConf(oldConf)
		if err != nil {
			return nil, err
		}
		if !shouldMigrate {
			continue
		}

		newConf := SubscriberConfig{Lte: oldConf}
		newConfBytes, err := json.Marshal(newConf)
		if err != nil {
			return nil, errors.Wrap(err, "error marshalling new subscriber config")
		}
		newConfsByPk[pk] = newConfBytes
	}
	defer sqorc.CloseRowsLogOnError(rows, "m011_subscriber_config")

	for pk, newConf := range newConfsByPk {
		_, err = builder.Update(tableName).
			Set(configCol, newConf).
			Where(squirrel.Eq{pkCol: pk}).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "error updating subscriber %s", pk)
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
		return false, errors.Wrap(err, "could not unmarshal legacy config")
	}

	_, alreadyMigrated := parsedMessage["lte"]
	return !alreadyMigrated, nil
}
