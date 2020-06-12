/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package sqorc

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	MariaDriver    = "mysql"
	PostgresDriver = "postgres"
	SQLiteDriver   = "sqlite3"
)

// Open is a wrapper for sql.Open which sets the max open connections to 1
// for in memory sqlite3 dbs. In memory sqlite3 creates a new database
// on each connection, so the number of open connections must be limited
// to 1 for thread safety. Otherwise, there is a race condition between
// threads using a cached connection to the original database or opening
// a new connection to a new database.
func Open(driver string, source string) (*sql.DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}
	if driver == SQLiteDriver && strings.Contains(source, ":memory:") {
		db.SetMaxOpenConns(1)
	}
	return db, nil
}
