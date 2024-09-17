package sqorc

import (
	"fmt"
	"os"
	"strings"

	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"

	"magma/orc8r/lib/go/merrors"
)

const (
	uniqueViolation = "unique_violation"
)

type ErrorChecker interface {
	GetError(error) error
}

type SQLiteErrorChecker struct{}

type PostgresErrorChecker struct{}

// GetErrorChecker returns a squirrel Builder for the configured SQL dialect as
// found in the SQL_DIALECT env var.
func GetErrorChecker() ErrorChecker {
	dialect, envFound := os.LookupEnv(SQLDialectEnv)
	// Default to postgresql
	if !envFound {
		return PostgresErrorChecker{}
	}

	switch strings.ToLower(dialect) {
	case PostgresDialect:
		return PostgresErrorChecker{}
	case SQLiteDialect:
		return SQLiteErrorChecker{}
	default:
		panic(fmt.Sprintf("unsupported sql dialect %s", dialect))
	}
}

func (c SQLiteErrorChecker) GetError(err error) error {
	if e, ok := err.(sqlite3.Error); ok {
		switch e.ExtendedCode {
		case sqlite3.ErrConstraintUnique, sqlite3.ErrConstraintPrimaryKey:
			return merrors.ErrAlreadyExists
		}
	}
	return err
}

func (c PostgresErrorChecker) GetError(err error) error {
	if e, ok := err.(*pq.Error); ok {
		switch e.Code.Name() {
		case uniqueViolation:
			return merrors.ErrAlreadyExists
		}
	}
	return err
}
