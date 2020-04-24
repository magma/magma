// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/generated"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func prepareUserData(ctx context.Context) {
	client := ent.FromContext(ctx).User
	client.Create().
		SetAuthID("user1").
		SetEmail("sam@workspace.com").
		SetFirstName("Samuel").
		SetLastName("Willis").
		SetRole(user.RoleUSER).
		SaveX(ctx)
	client.Create().
		SetAuthID("user2").
		SetEmail("the-monster@workspace.com").
		SetFirstName("Eli").
		SetLastName("Cohen").
		SetRole(user.RoleUSER).
		SaveX(ctx)

	client.Create().
		SetAuthID("user3").
		SetEmail("funny@workspace.com").
		SetFirstName("Willis").
		SetLastName("Reed").
		SetRole(user.RoleUSER).
		SaveX(ctx)

	client.Create().
		SetAuthID("user4").
		SetEmail("danit@workspace.com").
		SetFirstName("Dana").
		SetLastName("Cohen").
		SetRole(user.RoleUSER).
		SaveX(ctx)
}

func searchByName(t *testing.T, ctx context.Context, qr generated.QueryResolver, searchTerm string) *models.UserSearchResult {
	limit := 100
	f1 := models.UserFilterInput{
		FilterType:  models.UserFilterTypeUserName,
		Operator:    models.FilterOperatorContains,
		StringValue: &searchTerm,
	}
	res, err := qr.UserSearch(ctx, []*models.UserFilterInput{&f1}, &limit)
	require.NoError(t, err)
	return res
}

func TestSearchUsersByName(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	qr := r.Query()
	ctx := viewertest.NewContext(context.Background(), r.client)
	prepareUserData(ctx)

	search1 := searchByName(t, ctx, qr, "Cohen")
	require.Len(t, search1.Users, 2)
	search2 := searchByName(t, ctx, qr, "monster")
	require.Len(t, search2.Users, 1)
	search3 := searchByName(t, ctx, qr, "willis")
	require.Len(t, search3.Users, 2)
	search4 := searchByName(t, ctx, qr, "sam")
	require.Len(t, search4.Users, 1)
}
