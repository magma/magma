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
	"sort"
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

func GetTableName(networkId string, baseName string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(networkId), baseName)
}

// If the table DNE, log and return empty map
func GetAllValuesFromTable(tx *sql.Tx, table string) (map[string][]byte, error) {
	// Not every network may have gateways or meshes, in which case the
	// corresponding tables won't exist. Check and return early if so.
	exists, err := doesTableExist(tx, table)
	if err != nil {
		return nil, fmt.Errorf("Error checking if table %s exists: %s", table, err)
	}
	if !exists {
		glog.Errorf("Table %s does not exist, returning empty result from get", table)
		return map[string][]byte{}, nil
	}

	query := fmt.Sprintf("SELECT key, value FROM %s", table)
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	ret := map[string][]byte{}
	for rows.Next() {
		var key string
		var val []byte

		err = rows.Scan(&key, &val)
		if err != nil {
			return nil, err
		}
		ret[key] = val
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "sql rows err")
	}
	return ret, nil
}

// IMPORTANT: This is NOT portable, and ONLY works on postgres!
func doesTableExist(tx *sql.Tx, table string) (bool, error) {
	row := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name=$1)", table)
	ret := false
	err := row.Scan(&ret)
	return ret, err
}

func GetAllKeysFromTable(tx *sql.Tx, table string) ([]string, error) {
	query := fmt.Sprintf("SELECT key FROM %s", table)
	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var ret []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, err
		}
		ret = append(ret, key)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "sql rows err")
	}

	sort.Strings(ret)
	return ret, nil
}

// ExecInTx executes a callback inside a sql transaction on the provided DB.
// The transaction is rolled back if any error is encountered.
// initFn is a callback to call before the main txFn, commonly used in our
// codebase to execute a CREATE TABLE IF NOT EXISTS.
// Copied from orc8r/cloud/go/sqorc/tx.go.
func ExecInTx(
	db *sql.DB,
	opts *sql.TxOptions,
	initFn func(*sql.Tx) error,
	txFn func(*sql.Tx) (interface{}, error),
) (ret interface{}, err error) {
	tx, err := db.BeginTx(context.Background(), opts)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			// Rollback the tx immediately so the DB engine doesn't have to
			// wait for the conn to close
			_ = tx.Rollback()
			glog.Fatalf("recovered from panic: %v", r)
		}

		switch err {
		case nil:
			err = tx.Commit()
		default:
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				glog.Errorf("error rolling back tx: %s", rollbackErr)
			}
		}
	}()

	if initFn != nil {
		err = initFn(tx)
		if err != nil {
			return
		}
	}

	ret, err = txFn(tx)
	return
}
