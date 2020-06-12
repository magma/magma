// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
)

func TestUsersGroupWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleUSER)
	g := c.UsersGroup.Create().
		SetName("Group").
		AddMembers(u).
		SaveX(ctx)
	createGroup := func(ctx context.Context) error {
		_, err := c.UsersGroup.Create().
			SetName("Group2").
			AddMembers(u).
			Save(ctx)
		return err
	}
	updateGroup := func(ctx context.Context) error {
		return c.UsersGroup.UpdateOne(g).
			SetName("NewName").
			Exec(ctx)
	}
	deleteGroup := func(ctx context.Context) error {
		return c.UsersGroup.DeleteOne(g).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.AdminPolicy.Access.IsAllowed = models.PermissionValueYes
		},
		create: createGroup,
		update: updateGroup,
		delete: deleteGroup,
	})
}
