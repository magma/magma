/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"fmt"
	"log"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/tools/migrations"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	networksTable = "cfg_networks"
	nwTypeCol     = "type"
)

func main() {
	dbDriver := migrations.GetEnvWithDefault("SQL_DRIVER", "postgres")
	dbSource := migrations.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_dev user=magma_dev password=magma_dev host=postgres sslmode=disable")

	log.Println("Starting typed networks migrations script...")

	db, err := sqorc.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not open db connection"))
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(errors.Wrap(err, "error opening tx"))
	}

	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s text", networksTable, nwTypeCol))
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(errors.Wrap(err, "failed to add type field to networks table"))
	}

	_, err = tx.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s_idx ON %s (%s)", nwTypeCol, networksTable, nwTypeCol))
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(errors.Wrap(err, "failed to create type index"))
	}

	tx.Commit()
	log.Println("Success!")
}
