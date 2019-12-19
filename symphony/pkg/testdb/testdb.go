// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdb

import (
	"database/sql"
	"errors"
)

type opener interface {
	open() (*sql.DB, error)
}

var openers = map[string]opener{}

func register(name string, opener opener) {
	openers[name] = opener
}

// Open opens a testdb database.
func Open() (*sql.DB, string, error) {
	for name, opener := range openers {
		if db, err := opener.open(); err == nil {
			return db, name, nil
		}
	}
	return nil, "", errors.New("no available testdb")
}
