// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !buckbuild

package testdb

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	// registers sqlite3 driver with sql package
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	register("sqlite3", sqlite3{})
}

type sqlite3 struct{}

func (sqlite3) open() (*sql.DB, error) {
	var dbid [10]byte
	if _, err := rand.Read(dbid[:]); err != nil {
		return nil, fmt.Errorf("generating random bytes: %w", err)
	}
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1",
		hex.EncodeToString(dbid[:]))
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	return db, nil
}
