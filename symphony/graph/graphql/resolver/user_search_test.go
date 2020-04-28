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
		SetLastName("Rheed").
		SetRole(user.RoleUSER).
		SaveX(ctx)

	client.Create().
		SetAuthID("user4").
		SetEmail("danit@workspace.com").
		SetFirstName("Dana").
		SetLastName("Cohen").
		SetRole(user.RoleUSER).
		SaveX(ctx)

	client.Create().
		SetAuthID("user5").
		SetEmail("user5@test.ing").
		SetFirstName("Raul").
		SetLastName("Himemes").
		SetRole(user.RoleUSER).
		SetStatus(user.StatusDEACTIVATED).
		SaveX(ctx)
}

func searchByStatus(
	ctx context.Context,
	t *testing.T,
	qr generated.QueryResolver,
	status user.Status) *models.UserSearchResult {
	limit := 100
	f1 := models.UserFilterInput{
		FilterType:  models.UserFilterTypeUserStatus,
		Operator:    models.FilterOperatorIs,
		StatusValue: &status,
	}
	res, err := qr.UserSearch(ctx, []*models.UserFilterInput{&f1}, &limit)
	require.NoError(t, err)
	return res
}

func searchByName(
	ctx context.Context,
	t *testing.T,
	qr generated.QueryResolver,
	searchTerm string) *models.UserSearchResult {
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

	search1 := searchByName(ctx, t, qr, "Cohen")
	require.Len(t, search1.Users, 2)

	search2 := searchByName(ctx, t, qr, "monster")
	require.Len(t, search2.Users, 1)

	search3 := searchByName(ctx, t, qr, "willis")
	require.Len(t, search3.Users, 2)

	search4 := searchByName(ctx, t, qr, "sam")
	require.Len(t, search4.Users, 1)

	search5 := searchByName(ctx, t, qr, "li")
	require.Len(t, search5.Users, 3)

	search6 := searchByName(ctx, t, qr, "ra hi")
	require.Len(t, search6.Users, 1)

	search7 := searchByName(ctx, t, qr, "li he")
	require.Len(t, search7.Users, 2)

	search8 := searchByStatus(ctx, t, qr, user.StatusACTIVE)
	require.Len(t, search8.Users, 5) // including 'tester@example.com'

	search9 := searchByStatus(ctx, t, qr, user.StatusDEACTIVATED)
	require.Len(t, search9.Users, 1)
}
