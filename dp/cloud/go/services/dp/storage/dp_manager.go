package storage

import (
	"database/sql"

	"magma/orc8r/cloud/go/sqorc"
)

type dpManager struct {
	db           *sql.DB
	builder      sqorc.StatementBuilder
	cache        *enumCache
	errorChecker sqorc.ErrorChecker
	locker       sqorc.Locker
}
