// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func TestUserCannotBeDeleted(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u, err := c.User.Create().SetAuthID("someone").Save(ctx)
	require.NoError(t, err)
	err = c.User.DeleteOne(u).Exec(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
}

func TestUserWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleUSER)
	createUser := func(ctx context.Context) error {
		_, err := c.User.Create().
			SetAuthID("AuthID2").
			Save(ctx)
		return err
	}
	updateUser := func(ctx context.Context) error {
		return c.User.UpdateOne(u).
			SetFirstName("NewName").
			Exec(ctx)
	}
	deleteUser := func(ctx context.Context) error {
		if authz.FromContext(ctx).AdminPolicy.Access.IsAllowed == models.PermissionValueYes {
			return nil
		}
		return privacy.Deny
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.AdminPolicy.Access.IsAllowed = models.PermissionValueYes
		},
		create: createUser,
		update: updateUser,
		delete: deleteUser,
	})
}
