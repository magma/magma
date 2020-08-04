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
	"context"
	"database/sql"
	"flag"
	"fmt"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	// Directoryd-specific.
	oldDirectorydTable = "HWID_TO_HOSTNAME"
)

// main runs the directoryd_to_blobstore migration.
// Migration without SQL error will result in `SUCCESS` printed to stderr.
//
// This migration just deletes the old hwid_to_hostname table.
func main() {
	flag.Parse()
	_ = flag.Set("alsologtostderr", "true") // enable printing to console
	defer glog.Flush()

	glog.Info("BEGIN MIGRATION")

	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")
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

	// Drop new table if exists (migration idempotency)
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", oldDirectorydTable)
	glog.Infof("[RUN] %s", query)
	_, err = tx.Exec(query)
	if err != nil {
		err = errors.Wrap(err, "failed to drop new table")
		return
	}

	glog.Info("END MIGRATION")
}
