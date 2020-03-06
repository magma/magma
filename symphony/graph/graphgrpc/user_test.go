// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/facebookincubator/ent/dialect/sql/schema"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/facebookincubator/ent/dialect/sql"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T) *ent.Client {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	drv := sql.OpenDB(name, db)
	client := ent.NewClient(ent.Driver(drv))
	require.NoError(t, client.Schema.Create(context.Background(), schema.WithGlobalUniqueID(true)))
	return client
}

func TestUserService_Create(t *testing.T) {
	client := newTestClient(t)
	us := NewUserService(func(context.Context, string) (*ent.Client, error) { return client, nil })
	ctx := context.Background()

	u, err := us.Create(ctx, &UserInput{Tenant: "", Id: "XXX"})
	require.Nil(t, u)
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	u, err = us.Create(ctx, &UserInput{Tenant: "XXX", Id: ""})
	require.Nil(t, u)
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	u, err = us.Create(ctx, &UserInput{Tenant: "XXX", Id: "YYY"})
	require.NoError(t, err)
	userObject, err := client.User.Get(ctx, int(u.Id))
	require.NoError(t, err)
	require.Equal(t, user.StatusActive, userObject.Status)
	require.Equal(t, user.RoleUser, userObject.Role)
}

func TestUserService_Delete(t *testing.T) {
	client := newTestClient(t)
	us := NewUserService(func(context.Context, string) (*ent.Client, error) { return client, nil })
	ctx := context.Background()
	u := client.User.Create().SetAuthID("YYY").SaveX(ctx)
	require.Equal(t, user.StatusActive, u.Status)

	_, err := us.Delete(ctx, &UserInput{Tenant: "", Id: "YYY"})
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	_, err = us.Delete(ctx, &UserInput{Tenant: "XXX", Id: ""})
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	_, err = us.Delete(ctx, &UserInput{Tenant: "XXX", Id: "YYY"})
	require.NoError(t, err)
	newU, err := client.User.Get(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, user.StatusDeleted, newU.Status)
}
