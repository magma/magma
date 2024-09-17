package sqorc

import (
	"os"
	"testing"

	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"

	"magma/orc8r/lib/go/merrors"
)

const (
	uniqueViolationNum = "23505"
	otherPsqlErrorNum  = "22011"
)

func TestErrorChecker(t *testing.T) {
	suite.Run(t, &ErrorCheckerTestSuite{})
}

type ErrorCheckerTestSuite struct {
	suite.Suite
	prevDialect string
}

func (s *ErrorCheckerTestSuite) SetupSuite() {
	s.prevDialect = os.Getenv(SQLDialectEnv)
}

func (s *ErrorCheckerTestSuite) TearDownTest() {
	_ = os.Setenv(SQLDialectEnv, s.prevDialect)
}

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

func (s *ErrorCheckerTestSuite) TestSQLiteErrorCheckerCreation() {
	_ = os.Setenv(SQLDialectEnv, SQLiteDialect)
	c := GetErrorChecker()
	s.Assert().IsType(SQLiteErrorChecker{}, c)
}

func (s *ErrorCheckerTestSuite) TestPostgresErrorCheckerCreation() {
	_ = os.Setenv(SQLDialectEnv, PostgresDialect)
	c := GetErrorChecker()
	s.Assert().IsType(PostgresErrorChecker{}, c)
}

func (s *ErrorCheckerTestSuite) TestErrorCheckerNotCreatedWithUnknownDialect() {
	_ = os.Setenv(SQLDialectEnv, "someOtherDialect")
	s.Assert().Panics(func() { GetErrorChecker() })
}

func (s *ErrorCheckerTestSuite) TestErrorCheckerDefaultsToPostgres() {
	c := GetErrorChecker()
	s.Assert().IsType(PostgresErrorChecker{}, c)
}

func (s *ErrorCheckerTestSuite) TestSQLiteGetError() {
	testCases := []sqliteGetErrorTestCase{
		{
			name:    "test unique constraint error with SQLiteErrorChecker",
			checker: SQLiteErrorChecker{},
			err: sqlite3.Error{
				ExtendedCode: sqlite3.ErrConstraintUnique,
			},
			expectedError: merrors.ErrAlreadyExists,
		},
		{
			name:    "test pk constraint error with SQLiteErrorChecker",
			checker: SQLiteErrorChecker{},
			err: sqlite3.Error{
				ExtendedCode: sqlite3.ErrConstraintPrimaryKey,
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
		s.Run(tc.name, func() {
			err := tc.checker.GetError(tc.err)
			s.Assert().IsType(tc.expectedError, err)
		})
	}
}

func (s *ErrorCheckerTestSuite) TestPostgresGetError() {
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
		s.Run(tc.name, func() {
			err := tc.checker.GetError(tc.err)
			s.Assert().IsType(tc.expectedError, err)
		})
	}
}
