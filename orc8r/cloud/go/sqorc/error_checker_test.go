package sqorc

import (
	"os"
	"testing"

	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/lib/go/merrors"
)

const (
	uniqueViolationNum = "23505"
	otherPsqlErrorNum  = "22011"
)

type customError struct{}

var _ error = customError{}

func (c customError) Error() string {
	return "custom error message"
}

type sqliteGetErrorTestCase struct {
	name          string
	checker       ErrorChecker
	err           error
	expectedError error
}

type psqlGetErrorTestCase struct {
	name          string
	checker       ErrorChecker
	err           error
	expectedError error
}

func TestSQLiteErrorCheckerCreation(t *testing.T) {
	_ = os.Setenv(SQLDialectEnv, SQLiteDialect)
	c := GetErrorChecker()
	assert.IsType(t, SQLiteErrorChecker{}, c)
}

func TestPostgresErrorCheckerCreation(t *testing.T) {
	_ = os.Setenv(SQLDialectEnv, PostgresDialect)
	c := GetErrorChecker()
	assert.IsType(t, PostgresErrorChecker{}, c)
}

func TestErrorCheckerNotCreatedWithUnknownDialect(t *testing.T) {
	_ = os.Setenv(SQLDialectEnv, "someOtherDialect")
	assert.Panics(t, func() { GetErrorChecker() })
}

func TestErrorCheckerDefaultsToPostgres(t *testing.T) {
	c := GetErrorChecker()
	assert.IsType(t, PostgresErrorChecker{}, c)
}

func TestSQLiteGetError(t *testing.T) {
	testCases := []sqliteGetErrorTestCase{
		{
			name:    "test sqlite constraint error with SQLiteErrorChecker",
			checker: SQLiteErrorChecker{},
			err: sqlite3.Error{
				Code: sqlite3.ErrConstraint,
			},
			expectedError: merrors.ErrAlreadyExists,
		},
		{
			name:    "test other sqlite error with SQLiteErrorChecker",
			checker: SQLiteErrorChecker{},
			err: sqlite3.Error{
				Code: sqlite3.ErrNotFound,
			},
			expectedError: sqlite3.Error{},
		},
		{
			name:          "test any other error with SQLiteErrorChecker",
			checker:       SQLiteErrorChecker{},
			err:           customError{},
			expectedError: customError{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.checker.GetError(tc.err)
			assert.IsType(t, tc.expectedError, err)
		})
	}
}

func TestPostgresGetError(t *testing.T) {
	testCases := []psqlGetErrorTestCase{
		{
			name:    "test postgres constraint error with PostgresErrorChecker",
			checker: PostgresErrorChecker{},
			err: &pq.Error{
				Code: uniqueViolationNum,
			},
			expectedError: merrors.ErrAlreadyExists,
		},
		{
			name:    "test other postgres error with PostgresErrorChecker",
			checker: PostgresErrorChecker{},
			err: &pq.Error{
				Code: otherPsqlErrorNum,
			},
			expectedError: &pq.Error{},
		},
		{
			name:          "test any other error with PostgresErrorChecker",
			checker:       PostgresErrorChecker{},
			err:           customError{},
			expectedError: customError{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.checker.GetError(tc.err)
			assert.IsType(t, tc.expectedError, err)
		})
	}
}
