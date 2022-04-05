package sqorc

import (
	"magma/orc8r/lib/go/merrors"

	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

const (
	uniqueViolation = "unique_violation"
)

type ErrorChecker interface {
	GetError(error) error
}

type SQLiteErrorChecker struct{}

type PostgresErrorChecker struct{}

func (c SQLiteErrorChecker) GetError(err error) error {
	if e, ok := err.(sqlite3.Error); ok {
		switch e.Code {
		case sqlite3.ErrConstraint:
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
