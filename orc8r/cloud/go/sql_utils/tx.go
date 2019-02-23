/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sql_utils

import "database/sql"

// ExecInTx executes a callback inside a sql transaction on the provided DB.
// The transaction is rolled back if any error is encountered.
// initFn is a callback to call before the main txFn, commonly used in our
// codebase to execute a CREATE TABLE IF NOT EXISTS.
func ExecInTx(
	db *sql.DB,
	initFn func(*sql.Tx) error,
	txFn func(*sql.Tx) (interface{}, error),
) (ret interface{}, err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	err = initFn(tx)
	if err != nil {
		return
	}

	ret, err = txFn(tx)
	return
}
