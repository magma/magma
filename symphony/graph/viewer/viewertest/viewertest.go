// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewertest

import (
	"net/http"

	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/facebookincubator/symphony/graph/viewer"
)

const (
	TenantHeader  = viewer.TenantHeader
	DefaultTenant = "test"
	UserHeader    = viewer.UserHeader
	DefaultUser   = "tester@example.com"
	RoleHeader    = viewer.RoleHeader
	DefaultRole   = user.RoleOWNER
)

func SetDefaultViewerHeaders(req *http.Request) {
	req.Header.Set(TenantHeader, DefaultTenant)
	req.Header.Set(UserHeader, DefaultUser)
	req.Header.Set(RoleHeader, string(DefaultRole))
}
