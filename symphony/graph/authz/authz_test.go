// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/stretchr/testify/suite"
)

type handlerTestSuite struct {
	suite.Suite
	ctx    context.Context
	client *ent.Client
	viewer viewer.Viewer
}

func (s *handlerTestSuite) SetupTest() {
	s.client = viewertest.NewTestClient(s.T())
	s.ctx = ent.NewContext(context.Background(), s.client)
}

func (s *handlerTestSuite) TearDownTest() {
	s.client.Close()
	s.viewer = nil
}

func (s *handlerTestSuite) testViewerPermissions() {
	h := authz.Handler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions := authz.FromContext(r.Context())
			s.Require().NotNil(permissions)
			s.Require().EqualValues(authz.FullPermissions(), permissions)
			w.WriteHeader(http.StatusAccepted)
		}),
		logtest.NewTestLogger(s.T()),
	)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req.WithContext(viewer.NewContext(s.ctx, s.viewer)))
	s.Require().Equal(http.StatusAccepted, rec.Code)
}

func TestHandler(t *testing.T) {
	suite.Run(t, &handlerTestSuite{})
}

func (s *handlerTestSuite) TestUserHandler() {
	u := viewer.MustGetOrCreateUser(
		privacy.DecisionContext(s.ctx, privacy.Allow),
		viewertest.DefaultUser,
		user.RoleOWNER)
	s.viewer = viewer.NewUser(viewertest.DefaultTenant, u, viewer.WithFeatures(viewer.FeatureUserManagementDev))
	s.testViewerPermissions()
}

func (s *handlerTestSuite) TestAutomationHandler() {
	s.viewer = viewer.NewAutomation(
		viewertest.DefaultTenant,
		viewertest.DefaultUser,
		user.RoleOWNER,
		viewer.WithFeatures(viewer.FeatureUserManagementDev),
	)
	s.testViewerPermissions()
}
