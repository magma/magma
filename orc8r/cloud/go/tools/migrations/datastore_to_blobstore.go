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

package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	// Old table cols
	genCol     = "generation_number"
	deletedCol = "deleted"

	// Shared table cols
	keyCol = "\"key\"" // escaped for mysql compat
	valCol = "value"

	// New table cols
	nidCol  = "network_id"
	typeCol = "type"
	verCol  = "version"
)

// MigrateNetworkAgnosticServiceToBlobstore migrates a network-agnostic service's data
// from datastore to blobstore formats.
//
// Schema migration:
// 	- Datastore has cols
//		- key
//		- value
//		- generation_number
//		- deleted
// 	- Blobstore has cols
//		- network_id
//		- type
//		- key
//		- value
//		- version
//	- Conversion (blobstore col <- datastore col)
//		- network_id	<- [nid parameter]
//		- type 			<- [typ parameter]
//		- key			<- key
//		- value			<- value
//		- version		<- generation_number
//		- [n/a]			<- deleted
func MigrateNetworkAgnosticServiceToBlobstore(nid, typ, oldTable, newTable string) {
	dbDriver := GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		glog.Fatal(errors.Wrap(err, "could not open db connection"))
	}

	// Set up transaction
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		glog.Fatal(errors.Wrap(err, "error opening tx"))
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

	// Ensure old table exists
	ob := builder.CreateTable(oldTable).
		IfNotExists().
		Column(keyCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(valCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(genCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		Column(deletedCol).Type(sqorc.ColumnTypeBool).NotNull().Default("FALSE").EndColumn()
	_, err = ob.RunWith(tx).Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create old table")
		return
	}

	// Drop new table if exists (migration idempotency)
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", newTable)
	_, err = tx.Exec(query)
	if err != nil {
		err = errors.Wrap(err, "failed to drop new table")
		return
	}

	// Ensure new table exists
	nb := builder.CreateTable(newTable).
		Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(typeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(keyCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(valCol).Type(sqorc.ColumnTypeBytes).EndColumn().
		Column(verCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		PrimaryKey(nidCol, typeCol, keyCol)
	_, err = nb.RunWith(tx).Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create new table")
		return
	}

	// Migrate old table to new table
	nidVal := fmt.Sprintf("'%s'", nid)
	typeVal := fmt.Sprintf("'%s'", typ)
	sb := builder.Select().
		From(oldTable).
		Column(squirrel.Alias(squirrel.Expr(nidVal), nidCol)).
		Column(squirrel.Alias(squirrel.Expr(typeVal), typeCol)).
		Columns(keyCol, valCol, genCol)
	ib := builder.Insert(newTable).
		Columns(nidCol, typeCol, keyCol, valCol, verCol).
		Select(sb)
	sqlStr, _, _ := ib.ToSql()
	glog.Info("[RUN] ", sqlStr)
	_, err = ib.RunWith(tx).Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to insert")
		return
	}
}
