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
	"log"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	tableName = "device"

	nidCol  = "network_id"
	typeCol = "type"
	keyCol  = "\"key\""
	valCol  = "value"
	verCol  = "version"
)

const (
	gwType = "access_gateway_record"
)

type blobRow struct {
	nid, t, k string
	val       []byte
	version   uint64
}

type legacyRecord struct {
	HwID legacyHWID      `json:"hw_id"`
	Key  json.RawMessage `json:"key"`
}

type legacyHWID struct {
	ID string `json:"id"`
}

type gatewayDevice struct {
	HardwareID string          `json:"hardware_id"`
	Key        json.RawMessage `json:"key"`
}

func main() {
	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		glog.Fatal(errors.Wrap(err, "could not open db connection"))
	}

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(errors.Wrap(err, "error opening tx"))
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				glog.Errorf("tx failed to rollback: %s", err)
			}
			glog.Fatal(err)
		}

		if err = tx.Commit(); err != nil {
			glog.Fatalf("tx failed to commit: %s", err)
		}
		glog.Info("SUCCESS")
	}()

	builder := sqorc.GetSqlBuilder()
	// since the DDL to check for existence of a table is not portable, we'll
	// just start by creating the table if it doesn't exist
	_, err = builder.CreateTable(tableName).
		IfNotExists().
		Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(typeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(keyCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(valCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(verCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		PrimaryKey(nidCol, typeCol, keyCol).
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create table")
	}

	rows, err := builder.Select(nidCol, typeCol, keyCol, valCol, verCol).
		From(tableName).
		Where(squirrel.Eq{typeCol: gwType}).
		RunWith(tx).
		Query()
	if err != nil {
		if rows != nil {
			_ = rows.Close()
		}
		return
	}
	defer func() { _ = rows.Close() }()

	migratedRows := []*blobRow{}
	for rows.Next() {
		var nid, t, k string
		var val []byte
		var version uint64

		if err = rows.Scan(&nid, &t, &k, &val, &version); err != nil {
			err = errors.Wrap(err, "could not scan row")
			return
		}

		oldVal := &legacyRecord{}
		if err = json.Unmarshal(val, oldVal); err != nil {
			err = errors.Wrapf(err, "could not unmarshal AGW record for (%s, %s)", nid, k)
			return
		}

		newVal := &gatewayDevice{HardwareID: oldVal.HwID.ID, Key: oldVal.Key}
		var newBytes []byte
		newBytes, err = json.Marshal(newVal)
		if err != nil {
			err = errors.Wrapf(err, "could not remarshal migrated gateway device for (%s, %s)", nid, k)
		}

		migratedRows = append(migratedRows, &blobRow{nid: nid, t: t, k: k, val: newBytes, version: version})
	}

	sc := squirrel.NewStmtCache(tx)
	defer func() { _ = sc.Clear() }()
	for _, row := range migratedRows {
		_, err = builder.Update(tableName).
			Set(valCol, row.val).
			Where(squirrel.Eq{nidCol: row.nid, typeCol: row.t, keyCol: row.k}).
			RunWith(sc).
			Exec()
		if err != nil {
			err = errors.Wrapf(err, "error updating gateway device (%s, %s)", row.nid, row.k)
			return
		}
	}
}
