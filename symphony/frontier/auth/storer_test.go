// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/cloud/log/logtest"
	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/facebookincubator/symphony/frontier/ent/enttest"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/authboss"
)

type storerTestSuite struct {
	suite.Suite
	ctx    context.Context
	client *ent.Client
	storer *UserStorer
}

func (s *storerTestSuite) SetupSuite() {
	client, err := enttest.NewClient()
	s.Require().NoError(err)

	ctx := context.Background()
	tenant, err := client.Tenant.Create().
		SetName("test").
		SetDomains([]string{}).
		SetNetworks([]string{}).
		Save(ctx)
	s.Require().NoError(err)

	s.ctx = context.WithValue(ctx, tenantCtxKey{}, tenant)
	s.client = client
	s.storer = NewUserStorer(client.User, logtest.NewTestLogger(s.T()))
}

func (s *storerTestSuite) TearDownSuite() {
	err := s.client.Close()
	s.Assert().NoError(err)
}

func (s *storerTestSuite) TestLoadUser() {
	email := "tester@example.com"
	password := "testpassword"

	u := s.storer.New(s.ctx)
	user, ok := u.(authboss.AuthableUser)
	s.Require().True(ok)

	user.PutPID(email)
	user.PutPassword(password)
	err := s.storer.Create(s.ctx, user)
	s.Require().NoError(err)

	s.Run("Existing", func() {
		u, err := s.storer.Load(s.ctx, email)
		s.Require().NoError(err)
		s.Assert().Equal(email, u.GetPID())
		s.Assert().Equal(password, u.(authboss.AuthableUser).GetPassword())
	})
	s.Run("NotFound", func() {
		_, err := s.storer.Load(s.ctx, "missing@example.com")
		s.Assert().EqualError(err, authboss.ErrUserNotFound.Error())
	})
}

func (s *storerTestSuite) TestSaveUser() {
	client := s.client.User
	email := "updater@example.com"
	password := "password"

	for _, assertion := range []func(error, ...interface{}){s.Require().NoError, s.Require().Error} {
		u := s.storer.New(s.ctx)
		user, ok := u.(authboss.AuthableUser)
		s.Require().True(ok)

		user.PutPID(email)
		user.PutPassword(password)
		err := s.storer.Create(s.ctx, user)
		assertion(err)
	}

	u, err := s.storer.Load(s.ctx, email)
	s.Require().NoError(err)

	s.Run("OK", func() {
		u.PutPID("new@example.com")
		u.(authboss.AuthableUser).PutPassword("pwd")
		err = s.storer.Save(s.ctx, u)
		s.Require().NoError(err)
		s.Assert().Equal("new@example.com", u.GetPID())
		s.Assert().Equal("pwd", u.(authboss.AuthableUser).GetPassword())
	})
	s.Run("Duplicate", func() {
		email := "root@example.com"
		_, err := client.Create().
			SetEmail("root@example.com").
			SetPassword("root").
			SetTenant(CurrentTenant(s.ctx).Name).
			SetNetworks([]string{}).
			Save(s.ctx)
		s.Require().NoError(err)

		u.PutPID(email)
		err = s.storer.Save(s.ctx, u)
		s.Assert().EqualError(err, authboss.ErrUserFound.Error())
	})
}

func TestStorerTestSuite(t *testing.T) {
	suite.Run(t, &storerTestSuite{})
}
