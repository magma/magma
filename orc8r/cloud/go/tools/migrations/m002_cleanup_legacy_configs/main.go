/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"
	"log"

	"magma/orc8r/cloud/go/tools/migrations"
	"magma/orc8r/cloud/go/tools/migrations/m002_cleanup_legacy_configs/migration"

	_ "github.com/lib/pq"
)

func main() {
	shouldDropTables := flag.Bool("dropTables", false, "Set this flag to drop the old gateway and mesh config tables as part of the migration")
	flag.Parse()

	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma user=magma password=magma host=192.168.80.20")

	err := migration.Migrate(dbDriver, dbSource, *shouldDropTables)
	if err != nil {
		log.Fatal(err)
	}
}
