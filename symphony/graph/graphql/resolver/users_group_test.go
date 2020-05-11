// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func getAddUsersGroupInput(name, description string) models.AddUsersGroupInput {
	return models.AddUsersGroupInput{
		Name:        name,
		Description: &description,
	}
}

func TestAddUsersGroup(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()

	u1 := viewer.MustGetOrCreateUser(ctx, "user_1@test.ing", user.RoleUSER)
	memberIds := []int{u1.ID}

	gName := "group_1"
	gDescription := "this is group 1"
	inp := models.AddUsersGroupInput{
		Name:        gName,
		Description: &gDescription,
		Members:     memberIds,
	}

	_, err := mr.AddUsersGroup(ctx, inp)
	require.NoError(t, err)

	client := ent.FromContext(ctx)
	ugs := client.UsersGroup.Query().AllX(ctx)
	require.Len(t, ugs, 1)

	require.Equal(t, ugs[0].Name, gName, "verifying group name")
	require.Equal(t, ugs[0].Status, usersgroup.StatusACTIVE, "verifying group status")
	require.Len(t, ugs[0].QueryMembers().AllX(ctx), 1)
}

func TestDeleteUsersGroup(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()

	gName := "group_1"
	inp := getAddUsersGroupInput(gName, "this is group 1")
	_, err := mr.AddUsersGroup(ctx, inp)
	require.NoError(t, err)

	client := ent.FromContext(ctx)
	ugs := client.UsersGroup.Query().AllX(ctx)
	require.Len(t, ugs, 1)

	_, err = mr.DeleteUsersGroup(ctx, ugs[0].ID)
	require.NoError(t, err)
}

func TestEditUsersGroup(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()

	gName := "group_1"
	addInp1 := getAddUsersGroupInput(gName, "this is group 1")
	ug1, err := mr.AddUsersGroup(ctx, addInp1)
	require.NoError(t, err)

	gUpdatedName := "group_1_updated"
	updateInput1 := models.EditUsersGroupInput{
		ID:   ug1.ID,
		Name: &gUpdatedName,
	}
	ug1update1, err := mr.EditUsersGroup(ctx, updateInput1)
	require.NoError(t, err)
	require.Equal(t, ug1update1.Name, gUpdatedName, "verifying group name update")

	uAuthID1 := "user_1@test.ing"
	u1, err := CreateUserEnt(ctx, r.client, uAuthID1)
	require.NoError(t, err)
	memberIds := []int{u1.ID}

	gUpdatedDescription := "group_description_1_updated"
	updateInput2 := models.EditUsersGroupInput{
		ID:          ug1.ID,
		Description: &gUpdatedDescription,
		Members:     memberIds,
	}
	ug1update2, err := mr.EditUsersGroup(ctx, updateInput2)
	require.NoError(t, err)
	require.Equal(t, ug1update2.Name, gUpdatedName, "verifying group name stayed the same")
	require.Equal(t, ug1update2.Description, gUpdatedDescription, "verifying group description update")
	require.Len(t, ug1update2.QueryMembers().AllX(ctx), 1)

	gUpdatedStatus := usersgroup.StatusDEACTIVATED
	noMemberIds := []int{}
	updateInput3 := models.EditUsersGroupInput{
		ID:      ug1.ID,
		Status:  &gUpdatedStatus,
		Members: noMemberIds,
	}
	ug1update3, err := mr.EditUsersGroup(ctx, updateInput3)
	require.NoError(t, err)
	require.Equal(t, ug1update3.Status, gUpdatedStatus, "verifying group status updated")
	require.Len(t, ug1update2.QueryMembers().AllX(ctx), 0)
}

func CreateUserEnt(ctx context.Context, client *ent.Client, userName string) (*ent.User, error) {
	return client.User.Create().SetAuthID(userName).SetEmail(userName).Save(ctx)
}
