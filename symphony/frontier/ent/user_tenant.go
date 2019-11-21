// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ent

import (
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
)

// QueryTenant returns the tenant query bound to user.
func (u *User) QueryTenant() *TenantQuery {
	return (&TenantQuery{config: u.config}).
		Where(tenant.Name(u.Tenant))
}
