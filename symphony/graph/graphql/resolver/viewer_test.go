// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
)

func TestUserViewer(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	vr := r.Viewer()

	v := viewer.FromContext(ctx)
	r.client.User.UpdateOne(v.User()).SetRole(user.RoleUSER).ExecX(ctx)
	permissions, err := vr.Permissions(ctx, v)
	require.NoError(t, err)
	require.Equal(t, &models.BasicPermissionRule{IsAllowed: models2.PermissionValueNo}, permissions.AdminPolicy.Access)
	require.False(t, permissions.CanWrite)
}

func TestUserViewerInWriteGroup(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	vr := r.Viewer()

	v := viewer.FromContext(ctx)
	r.client.User.UpdateOne(v.User()).SetRole(user.RoleUSER).ExecX(ctx)
	_ = r.client.UsersGroup.Create().SetName(viewer.WritePermissionGroupName).AddMembers(v.User()).SaveX(ctx)

	permissions, err := vr.Permissions(ctx, v)
	require.NoError(t, err)
	require.Equal(t, &models.BasicPermissionRule{IsAllowed: models2.PermissionValueNo}, permissions.AdminPolicy.Access)
	require.True(t, permissions.CanWrite)
}

func TestAdminViewer(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	vr := r.Viewer()

	v := viewer.FromContext(ctx)
	r.client.User.UpdateOne(v.User()).SetRole(user.RoleADMIN).ExecX(ctx)
	permissions, err := vr.Permissions(ctx, v)
	require.NoError(t, err)
	require.Equal(t, &models.BasicPermissionRule{IsAllowed: models2.PermissionValueYes}, permissions.AdminPolicy.Access)
	require.False(t, permissions.CanWrite)
}

func TestOwnerViewer(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	vr := r.Viewer()

	v := viewer.FromContext(ctx)
	r.client.User.UpdateOne(v.User()).SetRole(user.RoleOWNER).ExecX(ctx)
	permissions, err := vr.Permissions(ctx, v)
	require.NoError(t, err)
	require.Equal(t, &models.BasicPermissionRule{IsAllowed: models2.PermissionValueYes}, permissions.AdminPolicy.Access)
	require.True(t, permissions.CanWrite)
}
