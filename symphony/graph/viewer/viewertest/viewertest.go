// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewertest

import (
	"net/http"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/pkg/log/logtest"

	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/facebookincubator/symphony/graph/viewer"
)

const (
	TenantHeader  = viewer.TenantHeader
	DefaultTenant = "test"
	UserHeader    = viewer.UserHeader
	DefaultUser   = "tester@example.com"
	RoleHeader    = viewer.RoleHeader
	DefaultRole   = user.RoleUSER
)

func SetDefaultViewerHeaders(req *http.Request) {
	req.Header.Set(TenantHeader, DefaultTenant)
	req.Header.Set(UserHeader, DefaultUser)
	req.Header.Set(RoleHeader, string(DefaultRole))
}

func TestHandler(t *testing.T, h http.Handler, client *ent.Client) http.Handler {
	auth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = authz.NewContext(ctx, authz.FullPermissions())
		h.ServeHTTP(w, r.WithContext(ctx))
	})
	return viewer.TenancyHandler(auth,
		viewer.NewFixedTenancy(client),
		logtest.NewTestLogger(t),
	)
}
