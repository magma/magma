// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func TestUserViewer(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	vr := r.Viewer()

	v := viewer.FromContext(ctx).(*viewer.UserViewer)
	r.client.User.UpdateOne(v.User()).SetRole(user.RoleUSER).ExecX(ctx)
	permissions, err := vr.Permissions(ctx, v)
	require.NoError(t, err)
	require.Equal(t, &models.BasicPermissionRule{IsAllowed: models.PermissionValueNo}, permissions.AdminPolicy.Access)
	require.False(t, permissions.CanWrite)
}

func TestUserViewerInWriteGroup(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client, viewertest.WithRole(user.RoleUSER), viewertest.WithFeatures())
	vr := r.Viewer()

	v := viewer.FromContext(ctx).(*viewer.UserViewer)
	_ = r.client.UsersGroup.Create().
		SetName(authz.WritePermissionGroupName).
		AddMembers(v.User()).
		SaveX(ctx)
	r.client.User.UpdateOne(v.User()).SetRole(user.RoleUSER).ExecX(ctx)

	permissions, err := vr.Permissions(ctx, v)
	require.NoError(t, err)
	require.Equal(t, &models.BasicPermissionRule{IsAllowed: models.PermissionValueNo}, permissions.AdminPolicy.Access)
	require.True(t, permissions.CanWrite)
}

func TestAdminViewer(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	vr := r.Viewer()

	v := viewer.FromContext(ctx).(*viewer.UserViewer)
	r.client.User.UpdateOne(v.User()).SetRole(user.RoleADMIN).ExecX(ctx)
	permissions, err := vr.Permissions(ctx, v)
	require.NoError(t, err)
	require.Equal(t, &models.BasicPermissionRule{IsAllowed: models.PermissionValueYes}, permissions.AdminPolicy.Access)
	require.False(t, permissions.CanWrite)
}

func TestOwnerViewer(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	vr := r.Viewer()

	v := viewer.FromContext(ctx).(*viewer.UserViewer)
	r.client.User.UpdateOne(v.User()).SetRole(user.RoleOWNER).ExecX(ctx)
	permissions, err := vr.Permissions(ctx, v)
	require.NoError(t, err)
	require.EqualValues(t, authz.FullPermissions(), permissions)
}
