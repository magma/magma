// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/authboss"
	"go.opencensus.io/trace"
)

type traceStorerSuite struct {
	suite.Suite
	ctx      context.Context
	storer   *testStorer
	tracer   *traceStorer
	exporter testExporter
}

func (s *traceStorerSuite) SetupSuite() {
	s.ctx = context.Background()
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
}

func (s *traceStorerSuite) SetupTest() {
	s.storer = &testStorer{}
	s.tracer = TraceStorer(s.storer)
	s.exporter = testExporter{}
	trace.RegisterExporter(&s.exporter)
}

func (s *traceStorerSuite) TearDownTest() {
	s.storer.AssertExpectations(s.T())
	trace.UnregisterExporter(&s.exporter)
}

func (s *traceStorerSuite) GetSpanData(name string) *trace.SpanData {
	for _, data := range s.exporter {
		if data.Name == name {
			return data
		}
	}
	return nil
}

func (s *traceStorerSuite) TestTraceStorer() {
	var (
		user  = &User{User: &ent.User{Email: "tester@example.com"}}
		token = "token"
	)
	tests := []struct {
		name       string
		prepare    func()
		spanName   string
		statusCode int32
	}{
		{
			name: "Load",
			prepare: func() {
				pid := "tester@example.com"
				s.storer.On("Load", mock.Anything, pid).
					Return(nil, authboss.ErrUserNotFound).
					Once()
				_, err := s.tracer.Load(s.ctx, pid)
				s.Require().EqualError(err, authboss.ErrUserNotFound.Error())
			},
			spanName:   "storer.LoadUser",
			statusCode: trace.StatusCodeNotFound,
		},
		{
			name: "Save",
			prepare: func() {
				s.storer.On("Save", mock.Anything, user).
					Return(nil).
					Once()
				err := s.tracer.Save(s.ctx, user)
				s.Assert().NoError(err)
			},
			spanName:   "storer.SaveUser",
			statusCode: trace.StatusCodeOK,
		},
		{
			name: "Create",
			prepare: func() {
				s.storer.On("Create", mock.Anything, user).
					Return(nil).
					Once()
				err := s.tracer.Create(s.ctx, user)
				s.Assert().NoError(err)
			},
			spanName:   "storer.CreateUser",
			statusCode: trace.StatusCodeOK,
		},
		{
			name: "AddRememberToken",
			prepare: func() {
				s.storer.On("AddRememberToken", mock.Anything, user.Email, token).
					Return(nil).
					Once()
				err := s.tracer.AddRememberToken(s.ctx, user.Email, token)
				s.Assert().NoError(err)
			},
			spanName:   "storer.AddRememberToken",
			statusCode: trace.StatusCodeOK,
		},
		{
			name: "DelRememberTokens",
			prepare: func() {
				s.storer.On("DelRememberTokens", mock.Anything, user.Email).
					Return(nil).
					Once()
				err := s.tracer.DelRememberTokens(s.ctx, user.Email)
				s.Assert().NoError(err)
			},
			spanName:   "storer.DelRememberTokens",
			statusCode: trace.StatusCodeOK,
		},
		{
			name: "UseRememberToken",
			prepare: func() {
				s.storer.On("UseRememberToken", mock.Anything, user.Email, token).
					Return(authboss.ErrTokenNotFound).
					Once()
				err := s.tracer.UseRememberToken(s.ctx, user.Email, token)
				s.Assert().EqualError(err, authboss.ErrTokenNotFound.Error())
			},
			spanName:   "storer.UseRememberToken",
			statusCode: trace.StatusCodeNotFound,
		},
	}

	for _, tc := range tests {
		tc.prepare()
	}
	for _, tc := range tests {
		data := s.GetSpanData(tc.spanName)
		s.Require().NotNil(data)
		s.Assert().Equal(tc.statusCode, data.Code)
		s.Assert().Contains(data.Attributes, "user")
	}
	s.Assert().Condition(func() bool {
		var found bool
		for _, data := range s.exporter {
			if value, ok := data.Attributes["token"]; ok {
				found = true
				s.Assert().NotEqual(token, value)
			}
		}
		return found
	})
}

func TestTraceStorerSuite(t *testing.T) {
	suite.Run(t, &traceStorerSuite{})
}

type testExporter []*trace.SpanData

func (e *testExporter) ExportSpan(s *trace.SpanData) {
	*e = append(*e, s)
}

type testStorer struct {
	mock.Mock
}

func (m *testStorer) Load(ctx context.Context, key string) (authboss.User, error) {
	args := m.Called(ctx, key)
	user, _ := args.Get(0).(authboss.User)
	return user, args.Error(1)
}

func (m *testStorer) Save(ctx context.Context, user authboss.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *testStorer) New(ctx context.Context) authboss.User {
	user, _ := m.Called(ctx).Get(0).(authboss.User)
	return user
}

func (m *testStorer) Create(ctx context.Context, user authboss.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *testStorer) AddRememberToken(ctx context.Context, pid, token string) error {
	return m.Called(ctx, pid, token).Error(0)
}

func (m *testStorer) DelRememberTokens(ctx context.Context, pid string) error {
	return m.Called(ctx, pid).Error(0)
}

func (m *testStorer) UseRememberToken(ctx context.Context, pid, token string) error {
	return m.Called(ctx, pid, token).Error(0)
}
