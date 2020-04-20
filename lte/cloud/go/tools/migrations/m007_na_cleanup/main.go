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
	"flag"
	"fmt"
	"log"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	tableName = "cfg_entities"
	pkCol     = "pk"
	keyCol    = "\"key\""
	typeCol   = "type"

	cellType = "cellular_gateway"
	mdType   = "magmad_gateway"
)

// This migration cleans up cellular gateways on the backend which are
// "hanging" - i.e. have no associated magmad gateway.
// See https://github.com/facebookincubator/magma/issues/1071
func main() {
	flag.Parse()
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

	/*
		WITH
			cell_gw AS (SELECT pk, key FROM cfg_entities WHERE type = 'cellular_gateway'),
			md_gw AS (SELECT key FROM cfg_entities WHERE type = 'magmad_gateway')
		DELETE FROM cfg_entities
		WHERE pk IN (
			SELECT cell_gw.pk as pk from cell_gw
			LEFT OUTER JOIN md_gw ON cell_gw.key = md_gw.key
			WHERE md_gw.key IS NULL
		)
	*/

	builder := sqorc.GetSqlBuilder().RunWith(tx)
	mainSelect, _, _ := builder.Select(fmt.Sprintf("cell_gw.%s as %s", pkCol, pkCol)).
		From("cell_gw").
		JoinClause(fmt.Sprintf("LEFT OUTER JOIN md_gw ON cell_gw.%s = md_gw.%s", keyCol, keyCol)).
		Where(sq.Eq{fmt.Sprintf("md_gw.%s", keyCol): nil}).
		ToSql()
	cellGWSelect, _, _ := builder.Select(pkCol, keyCol).
		From(tableName).
		Where(fmt.Sprintf("%s = '%s'", typeCol, cellType)).
		ToSql()
	mdGWSelect, _, _ := builder.Select(keyCol).
		From(tableName).
		Where(fmt.Sprintf("%s = '%s'", typeCol, mdType)).
		ToSql()

	_, err = builder.Delete(tableName).
		Where(sq.Expr(fmt.Sprintf("%s IN (%s)", pkCol, mainSelect))).
		Prefix(fmt.Sprintf("WITH cell_gw AS (%s), md_gw AS (%s)", cellGWSelect, mdGWSelect)).
		Exec()
}
