/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package datastore

import (
	"os"
)

func GetEnvWithDefault(variable string, defaultValue string) string {
	value := os.Getenv(variable)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}

var (
	SQL_DRIVER      = GetEnvWithDefault("SQL_DRIVER", "sqlite3")
	DATABASE_SOURCE = GetEnvWithDefault("DATABASE_SOURCE", ":memory:")
)
