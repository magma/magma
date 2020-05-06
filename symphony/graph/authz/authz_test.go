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
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/stretchr/testify/require"
)

func testContextGetFullPermissions(ctx context.Context, t *testing.T) {
	h := authz.AuthHandler{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions := authz.FromContext(r.Context())
			require.NotNil(t, permissions)
			require.EqualValues(t, authz.FullPermissions(), permissions)
			w.WriteHeader(http.StatusAccepted)
		}),
		Logger: log.NewNopLogger(),
	}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req.WithContext(ctx))
	require.Equal(t, http.StatusAccepted, rec.Code)
}

func TestAuthHandler(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	u := viewer.MustGetOrCreateUser(
		privacy.DecisionContext(ctx, privacy.Allow),
		viewertest.DefaultUser,
		user.RoleOWNER)
	v := viewer.NewUser(viewertest.DefaultTenant, u)
	ctx = viewer.NewContext(ctx, v)
	testContextGetFullPermissions(ctx, t)
}

func TestAuthHandlerForAutomation(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	v := viewer.NewAutomation(viewertest.DefaultTenant, viewertest.DefaultUser, user.RoleOWNER)
	ctx = viewer.NewContext(ctx, v)
	testContextGetFullPermissions(ctx, t)
}
