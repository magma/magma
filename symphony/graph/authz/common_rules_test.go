// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestUserCannotEditWithNoPermission(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), client)
	location, err := client.LocationType.Create().SetName("LocationType").Save(ctx)
	require.NoError(t, err)
	ctx = ent.NewContext(context.Background(), client)
	u := viewer.MustGetOrCreateUser(ctx, viewertest.DefaultUser, user.RoleOWNER)
	v := viewer.NewUser(viewertest.DefaultTenant, u)
	ctx = viewer.NewContext(ctx, v)
	_, err = client.LocationType.Get(ctx, location.ID)
	require.NoError(t, err)
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
	ctx = viewertest.NewContext(ctx, client, viewertest.WithPermissions(authz.EmptyPermissions()))
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
	permissions.CanWrite = true
	ctx = viewertest.NewContext(ctx, client, viewertest.WithPermissions(permissions))
	_, err = client.LocationType.Get(ctx, location.ID)
	require.NoError(t, err)
	_, err = client.LocationType.UpdateOneID(location.ID).SetName("NewLocationType").Save(ctx)
	require.NoError(t, err)
}
