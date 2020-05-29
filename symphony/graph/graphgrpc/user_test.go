// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/graphgrpc/schema"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/enttest"
	"github.com/facebookincubator/symphony/pkg/ent/migrate"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/testdb"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTestClient(t *testing.T) *ent.Client {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	drv := sql.OpenDB(name, db)
	return enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
}

func TestUserService_Create(t *testing.T) {
	client := newTestClient(t)
	us := NewUserService(func(context.Context, string) (*ent.Client, error) { return client, nil })
	ctx, err := CreateServiceContext(context.Background(), viewertest.DefaultTenant, UserServiceName, user.RoleADMIN)
	require.NoError(t, err)

	u, err := us.Create(ctx, &schema.AddUserInput{Tenant: "", Id: "XXX", IsOwner: false})
	require.Nil(t, u)
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	u, err = us.Create(ctx, &schema.AddUserInput{Tenant: "XXX", Id: "", IsOwner: false})
	require.Nil(t, u)
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	u, err = us.Create(ctx, &schema.AddUserInput{Tenant: "XXX", Id: "YYY", IsOwner: false})
	require.NoError(t, err)
	userObject, err := client.User.Get(ctx, int(u.Id))
	require.NoError(t, err)
	require.Equal(t, user.StatusACTIVE, userObject.Status)
	require.Equal(t, user.RoleUSER, userObject.Role)
}

func TestUserService_Delete(t *testing.T) {
	client := newTestClient(t)
	us := NewUserService(func(context.Context, string) (*ent.Client, error) { return client, nil })
	ctx, err := CreateServiceContext(context.Background(), viewertest.DefaultTenant, UserServiceName, user.RoleADMIN)
	require.NoError(t, err)
	u := client.User.Create().SetAuthID("YYY").SaveX(ctx)
	require.Equal(t, user.StatusACTIVE, u.Status)

	_, err = us.Delete(ctx, &schema.UserInput{Tenant: "", Id: "YYY"})
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	_, err = us.Delete(ctx, &schema.UserInput{Tenant: "XXX", Id: ""})
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	_, err = us.Delete(ctx, &schema.UserInput{Tenant: "XXX", Id: "YYY"})
	require.NoError(t, err)
	newU, err := client.User.Get(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, user.StatusDEACTIVATED, newU.Status)
}

func TestUserService_CreateAfterDelete(t *testing.T) {
	client := newTestClient(t)
	us := NewUserService(func(context.Context, string) (*ent.Client, error) { return client, nil })
	ctx, err := CreateServiceContext(context.Background(), viewertest.DefaultTenant, UserServiceName, user.RoleADMIN)
	require.NoError(t, err)
	u := client.User.Create().SetAuthID("YYY").SaveX(ctx)
	require.Equal(t, user.StatusACTIVE, u.Status)

	_, err = us.Delete(ctx, &schema.UserInput{Tenant: "XXX", Id: "YYY"})
	require.NoError(t, err)

	_, err = us.Create(ctx, &schema.AddUserInput{Tenant: "XXX", Id: "YYY", IsOwner: true})
	require.NoError(t, err)
	userObject, err := client.User.Get(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, user.StatusACTIVE, userObject.Status)
	require.Equal(t, user.RoleOWNER, userObject.Role)
}

func TestUserService_CreateGroup(t *testing.T) {
	client := newTestClient(t)
	us := NewUserService(func(context.Context, string) (*ent.Client, error) { return client, nil })
	ctx, err := CreateServiceContext(context.Background(), viewertest.DefaultTenant, UserServiceName, user.RoleADMIN)
	require.NoError(t, err)
	exist, err := client.UsersGroup.Query().Exist(ctx)
	require.NoError(t, err)
	require.False(t, exist)
	_, err = us.Create(ctx, &schema.AddUserInput{Tenant: "XXX", Id: "YYY", IsOwner: false})
	require.NoError(t, err)
	count, err := client.UsersGroup.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, count)
	_, err = us.Create(ctx, &schema.AddUserInput{Tenant: "XXX", Id: "YYY2", IsOwner: false})
	require.NoError(t, err)
	count, err = client.UsersGroup.Query().Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}
