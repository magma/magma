/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sql_utils

import (
	"database/sql"
	"fmt"

	"github.com/golang/glog"
)

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
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				glog.Errorf("error rolling back tx: %s", err)
			}
		}
	}()

	err = initFn(tx)
	if err != nil {
		return
	}

	ret, err = txFn(tx)
	return
}

// PrepareStatements prepares a list of SQL query strings and returns the
// *sql.Stmt statements in the order that the query strings were provided.
// Responsibility is on the caller to close the statements.
func PrepareStatements(tx *sql.Tx, stmtStrings []string) ([]*sql.Stmt, error) {
	ret := make([]*sql.Stmt, 0, len(stmtStrings))
	for _, stmtStr := range stmtStrings {
		stmt, err := tx.Prepare(stmtStr)
		if err != nil {
			GetCloseStatementsDeferFunc(ret, "PrepareStatements")()
			return nil, fmt.Errorf("error preparing DB statement: %s", err)
		}
		ret = append(ret, stmt)
	}
	return ret, nil
}

// GetCloseStatementsDeferFunc returns a function which closes all provided
// sql statements to defer. Any error encountered while closing a statement
// will be logged.
// The callsite argument will be included in the log message for context.
//
// IMPORTANT: don't forget to call the returned func() in the defer clause.
// In other words, this should be used like:
// 		defer GetCloseStatementsDeferFunc(stmts, "foo")()
// If you forget the last set of parens, your statements will NOT close.
func GetCloseStatementsDeferFunc(stmts []*sql.Stmt, callsite string) func() {
	return func() {
		for _, stmt := range stmts {
			if stmt != nil {
				if err := stmt.Close(); err != nil {
					glog.Errorf("error closing statement in %s: %s", callsite, err)
				}
			}
		}
	}
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
