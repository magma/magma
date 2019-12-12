// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqltx

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type testServerInterceptorSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func (s *testServerInterceptorSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	s.db = db
	s.mock = mock
}

func (s *testServerInterceptorSuite) TearDownTest() {
	err := s.mock.ExpectationsWereMet()
	s.Require().NoError(err)
}

func TestServerInterceptor(t *testing.T) {
	suite.Run(t, &testServerInterceptorSuite{})
}

func (s *testServerInterceptorSuite) TestCommit() {
	s.mock.ExpectBegin()
	s.mock.ExpectCommit()
	_, _ = UnaryServerInterceptor(s.db)(
		context.Background(), nil, nil,
		func(ctx context.Context, _ interface{}) (interface{}, error) {
			tx := FromContext(ctx)
			s.Assert().NotNil(tx)
			return nil, nil
		},
	)
}

func (s *testServerInterceptorSuite) TestBadCommit() {
	s.mock.ExpectBegin()
	s.mock.ExpectCommit().WillReturnError(errors.New("bad commit"))
	_, err := UnaryServerInterceptor(s.db)(
		context.Background(), nil, nil,
		func(context.Context, interface{}) (interface{}, error) {
			return nil, nil
		},
	)
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "committing transaction")
}

func (s *testServerInterceptorSuite) TestRollback() {
	s.mock.ExpectBegin()
	s.mock.ExpectRollback()
	_, _ = UnaryServerInterceptor(s.db)(
		context.Background(), nil, nil,
		func(context.Context, interface{}) (interface{}, error) {
			return nil, errors.New("bad request")
		},
	)
}

func (s *testServerInterceptorSuite) TestNoTx() {
	errBadConn := errors.New("bad conn")
	s.mock.ExpectBegin().WillReturnError(errBadConn)
	_, err := UnaryServerInterceptor(s.db)(
		context.Background(), nil, nil,
		func(context.Context, interface{}) (interface{}, error) {
			s.Require().FailNow("invoked handler")
			return nil, nil
		},
	)
	s.Require().Error(err)
	s.Assert().Contains(err.Error(), "beginning transaction")
	s.Assert().True(errors.Is(err, errBadConn))
}

func (s *testServerInterceptorSuite) TestRecovery() {
	s.mock.ExpectBegin()
	s.mock.ExpectRollback()
	s.Assert().Panics(func() {
		_, _ = UnaryServerInterceptor(s.db)(
			context.Background(), nil, nil,
			func(context.Context, interface{}) (interface{}, error) {
				panic("fatal")
			},
		)
	})
}
