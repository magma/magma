// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux
// +build !nolibfb

package testdb

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"libfb/go/fbmysql/testdb"
)

type fbmysql struct {
	id uint64
}

func init() { register("mysql", &fbmysql{}) }

func (m *fbmysql) open() (*sql.DB, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("resolving hostname: %w")
	}
	prefix := strings.SplitN(hostname, ".", 2)[0] +
		"_" + strconv.FormatUint(atomic.AddUint64(&m.id, 1), 10)

	db, err := testdb.CreateWithOpts(testdb.Options{
		Prefix: prefix,
		DSNQueryParameters: map[string]string{
			"loc":       "Local",
			"parseTime": "true",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	return db, nil
}
