package storage

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"magma/orc8r/cloud/go/sqorc"
)

type dpManager struct {
	db           *sql.DB
	builder      sqorc.StatementBuilder
	cache        *enumCache
	errorChecker sqorc.ErrorChecker
	locker       sqorc.Locker
}

type queryRunner struct {
	builder sq.StatementBuilderType
	cache   *enumCache
	locker  sqorc.Locker
}

func (m *dpManager) getQueryRunner(tx sq.BaseRunner) *queryRunner {
	return &queryRunner{
		builder: m.builder.RunWith(tx),
		cache:   m.cache,
		locker:  m.locker,
	}
}
