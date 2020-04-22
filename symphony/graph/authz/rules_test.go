package authz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"

	"github.com/facebookincubator/symphony/graph/ent/privacy"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
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

func TestAdminUserCanEditUsers(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	_, err := client.User.Create().SetAuthID("admin_user").SetRole(user.RoleADMIN).Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser("admin_user"))
	_, err = client.UsersGroup.Create().SetName("NewGroup").Save(ctx)
	require.NoError(t, err)
	_, err = client.User.Create().SetAuthID("new_user").Save(ctx)
	require.NoError(t, err)
}

func TestUserCannotEditUsers(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	_, err := client.User.Create().SetAuthID("regular_user").SetRole(user.RoleUSER).Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser("regular_user"))
	_, err = client.UsersGroup.Create().SetName("NewGroup").Save(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
	_, err = client.User.Create().SetAuthID("new_user").Save(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
}

func TestUserIsReadonly(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	_, err := client.User.Create().SetAuthID("simple_user").Save(ctx)
	require.NoError(t, err)
	location, err := client.LocationType.Create().SetName("LocationType").Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser("simple_user"))
	_, err = client.LocationType.Get(ctx, location.ID)
	require.NoError(t, err)
	_, err = client.LocationType.UpdateOneID(location.ID).SetName("NewLocationType").Save(ctx)
	require.Error(t, err)
	require.True(t, errors.Is(err, privacy.Deny))
}

func TestOwnerCanWrite(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	_, err := client.User.Create().SetAuthID("owner_user").SetRole(user.RoleOWNER).Save(ctx)
	require.NoError(t, err)
	location, err := client.LocationType.Create().SetName("LocationType").Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser("owner_user"))
	_, err = client.LocationType.Get(ctx, location.ID)
	require.NoError(t, err)
	_, err = client.LocationType.UpdateOneID(location.ID).SetName("NewLocationType").Save(ctx)
	require.NoError(t, err)
}

func TestUserInGroupCanWrite(t *testing.T) {
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	userInGroup, err := client.User.Create().SetAuthID("user_in_group").Save(ctx)
	require.NoError(t, err)
	_, err = client.UsersGroup.Create().SetName(authz.WritePermissionGroupName).AddMembers(userInGroup).Save(ctx)
	require.NoError(t, err)
	location, err := client.LocationType.Create().SetName("LocationType").Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser("user_in_group"))
	_, err = client.LocationType.Get(ctx, location.ID)
	require.NoError(t, err)
	_, err = client.LocationType.UpdateOneID(location.ID).SetName("NewLocationType").Save(ctx)
	require.NoError(t, err)
}
