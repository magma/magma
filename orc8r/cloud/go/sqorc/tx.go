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

package sqorc

import (
	"context"
	"database/sql"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// ExecInTx executes a callback inside a sql transaction on the provided DB.
// The transaction is rolled back if any error is encountered.
// initFn is a callback to call before the main txFn, commonly used in our
// codebase to execute a CREATE TABLE IF NOT EXISTS.
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

// CloseRowsLogOnError will close the *Rows object and log if an error is
// returned by Rows.Close(). This function will no-op if rows is nil.
func CloseRowsLogOnError(rows *sql.Rows, callsite string) {
	if rows == nil {
		return
	}
	if err := rows.Close(); err != nil {
		glog.Errorf("Error closing sql rows in %s: %+v", callsite, errors.WithStack(err))
	}
}
