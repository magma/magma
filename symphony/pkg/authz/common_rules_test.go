// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/facebookincubator/symphony/pkg/authz/models"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestUserCannotEditOrViewWithNoPermission(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), client)
	location, err := client.LocationType.Create().SetName("LocationType").Save(ctx)
	require.NoError(t, err)
	u := viewer.MustGetOrCreateUser(ctx, viewertest.DefaultUser, user.RoleOWNER)
	v := viewer.NewUser(viewertest.DefaultTenant, u)
	ctx = ent.NewContext(context.Background(), client)
	ctx = viewer.NewContext(ctx, v)
	_, err = client.LocationType.Get(ctx, location.ID)
	require.True(t, errors.Is(err, privacy.Deny))
	_, err = client.LocationType.UpdateOneID(location.ID).SetName("NewLocationType").Save(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
	ctx = authz.NewContext(ctx, authz.FullPermissions())
	_, err = client.UsersGroup.Create().SetName("NewGroup").Save(ctx)
	require.NoError(t, err)
}

func TestUserCannotEditWithEmptyPermission(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), client)
	location, err := client.LocationType.Create().SetName("LocationType").Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(ctx,
		client,
		viewertest.WithUser("user"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(authz.EmptyPermissions()))
	_, err = client.UsersGroup.Create().SetName("NewGroup").Save(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
	_, err = client.User.Create().SetAuthID("new_user").Save(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
	_, err = client.LocationType.Get(ctx, location.ID)
	require.NoError(t, err)
	_, err = client.LocationType.UpdateOneID(location.ID).SetName("NewLocationType").Save(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
}

func TestUserCanWrite(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), client)
	location, err := client.LocationType.Create().SetName("LocationType").Save(ctx)
	require.NoError(t, err)
	permissions := authz.EmptyPermissions()
	permissions.InventoryPolicy.LocationType.Update.IsAllowed = models.PermissionValueYes
	ctx = viewertest.NewContext(ctx,
		client,
		viewertest.WithUser("user"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(permissions))
	_, err = client.LocationType.Get(ctx, location.ID)
	require.NoError(t, err)
	_, err = client.LocationType.UpdateOneID(location.ID).SetName("NewLocationType").Save(ctx)
	require.NoError(t, err)
}
