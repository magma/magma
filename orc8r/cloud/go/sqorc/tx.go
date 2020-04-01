/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sqorc

import (
	"context"
	"database/sql"

	"github.com/golang/glog"
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
		glog.Errorf("error closing *Rows in %s: %s", callsite, err)
	}
}
