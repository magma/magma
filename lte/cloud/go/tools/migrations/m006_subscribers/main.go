/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_ "github.com/lib/pq"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"
)

const (
	tableName = "cfg_entities"
	pkCol     = "pk"
	typeCol   = "type"
	confCol   = "config"

	entType = "subscriber"
)

type legacySubscriber struct {
	Lte json.RawMessage `json:"lte"`
}

func main() {
	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		glog.Fatal(fmt.Errorf("could not open db connection: %w", err))
	}

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(fmt.Errorf("error opening tx: %w", err))
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				glog.Errorf("tx failed to rollback: %s", rollbackErr)
			}
			glog.Fatal(err)
		}

		if err = tx.Commit(); err != nil {
			glog.Fatalf("tx failed to commit: %s", err)
		}
		glog.Info("SUCCESS")
	}()

	sc := squirrel.NewStmtCache(tx)
	defer func() { _ = sc.Clear() }()
	builder := sqorc.GetSqlBuilder().RunWith(sc)

	rows, err := builder.Select(pkCol, confCol).
		From(tableName).
		Where(squirrel.Eq{typeCol: entType}).
		Query()
	if err != nil {
		if rows != nil {
			_ = rows.Close()
		}
		return
	}
	defer func() { _ = rows.Close() }()

	legacySubs := map[string]*legacySubscriber{}
	for rows.Next() {
		var pk string
		var conf []byte

		if err = rows.Scan(&pk, &conf); err != nil {
			err = fmt.Errorf("could not scan row: %w", err)
			return
		}

		legacySub := &legacySubscriber{}
		if err = json.Unmarshal(conf, legacySub); err != nil {
			err = fmt.Errorf("could not unmarshal subscriber config %s: %w", pk, err)
			return
		}
		legacySubs[pk] = legacySub
	}

	for pk, legacy := range legacySubs {
		_, err = builder.Update(tableName).
			Set(confCol, legacy.Lte).
			Where(squirrel.Eq{pkCol: pk}).
			Exec()
		if err != nil {
			err = fmt.Errorf("error updating subscriber %s: %w", pk, err)
			return
		}
	}
}
