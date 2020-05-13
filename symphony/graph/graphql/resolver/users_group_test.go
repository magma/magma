// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/usersgroup"
	"github.com/facebookincubator/symphony/graph/graphql/models"
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

func TestAddUsersGroupWithPolicy(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()

	gName := "group_1"
	gDescription := "this is group 1"

	inventoryPolicyInput := getInventoryPolicyInput()

	policy1, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           "policy1",
		Description:    pointer.ToString("policyDescription1"),
		InventoryInput: inventoryPolicyInput,
		WorkforceInput: nil,
	})
	require.NoError(t, err)
	inp := models.AddUsersGroupInput{
		Name:        gName,
		Description: &gDescription,
		Members:     nil,
		Policies:    []int{policy1.ID},
	}
	_, err = mr.AddUsersGroup(ctx, inp)
	require.NoError(t, err)

	client := ent.FromContext(ctx)
	ugs := client.UsersGroup.Query().AllX(ctx)
	require.Len(t, ugs, 1)

	require.Equal(t, ugs[0].Name, gName, "verifying group name")
	require.Equal(t, ugs[0].Status, usersgroup.StatusACTIVE, "verifying group status")
	require.Len(t, ugs[0].QueryPolicies().AllX(ctx), 1)
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

	gUpdatedDescription := "group_description_1_updated"
	updateInput2 := models.EditUsersGroupInput{
		ID:          ug1.ID,
		Description: &gUpdatedDescription,
	}
	ug1update2, err := mr.EditUsersGroup(ctx, updateInput2)
	require.NoError(t, err)
	require.Equal(t, ug1update2.Name, gUpdatedName, "verifying group name stayed the same")
	require.Equal(t, ug1update2.Description, gUpdatedDescription, "verifying group description update")

	gUpdatedStatus := usersgroup.StatusDEACTIVATED
	updateInput3 := models.EditUsersGroupInput{
		ID:     ug1.ID,
		Status: &gUpdatedStatus,
	}
	ug1update3, err := mr.EditUsersGroup(ctx, updateInput3)
	require.NoError(t, err)
	require.Equal(t, ug1update3.Status, gUpdatedStatus, "verifying group status updated")
	require.Len(t, ug1update2.QueryMembers().AllX(ctx), 0)

	uAuthID1 := "user_1@test.ing"
	u1, err := CreateUserEnt(ctx, r.client, uAuthID1)
	require.NoError(t, err)

	userUpdateInput1 := models.EditUsersGroupInput{
		ID:      ug1.ID,
		Members: []int{u1.ID},
	}
	userUpdate1, err := mr.EditUsersGroup(ctx, userUpdateInput1)
	require.NoError(t, err)
	require.Len(t, userUpdate1.QueryMembers().AllX(ctx), 1)

	uAuthID2 := "user_2@test.ing"
	u2, err := CreateUserEnt(ctx, r.client, uAuthID2)
	require.NoError(t, err)

	userUpdateInput2 := models.EditUsersGroupInput{
		ID:      ug1.ID,
		Members: []int{u1.ID, u2.ID},
	}
	userUpdate2, err := mr.EditUsersGroup(ctx, userUpdateInput2)
	require.NoError(t, err)
	require.Len(t, userUpdate2.QueryMembers().AllX(ctx), 2)

	userUpdateInput3 := models.EditUsersGroupInput{
		ID:      ug1.ID,
		Members: []int{u2.ID},
	}
	userUpdate3, err := mr.EditUsersGroup(ctx, userUpdateInput3)
	require.NoError(t, err)
	require.Len(t, userUpdate3.QueryMembers().AllX(ctx), 1)

	inventoryPolicyInput := getInventoryPolicyInput()
	policy1, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           "policy1",
		Description:    pointer.ToString("policyDescription1"),
		InventoryInput: inventoryPolicyInput,
		WorkforceInput: nil,
	})
	require.NoError(t, err)

	policy2, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           "policy2",
		Description:    pointer.ToString("policyDescription2"),
		InventoryInput: inventoryPolicyInput,
		WorkforceInput: nil,
	})
	require.NoError(t, err)

	updatePolicyInput1 := models.EditUsersGroupInput{
		ID:       ug1.ID,
		Policies: []int{policy1.ID},
	}
	ugUpdate1, err := mr.EditUsersGroup(ctx, updatePolicyInput1)
	require.NoError(t, err)
	require.Len(t, ugUpdate1.QueryPolicies().AllX(ctx), 1)

	updatePolicyInput2 := models.EditUsersGroupInput{
		ID:       ug1.ID,
		Policies: []int{policy1.ID, policy2.ID},
	}
	ugUpdate2, err := mr.EditUsersGroup(ctx, updatePolicyInput2)
	require.NoError(t, err)
	require.Len(t, ugUpdate2.QueryPolicies().AllX(ctx), 2)

	updatePolicyInput3 := models.EditUsersGroupInput{
		ID:       ug1.ID,
		Policies: []int{policy2.ID},
	}
	ugUpdate3, err := mr.EditUsersGroup(ctx, updatePolicyInput3)
	require.NoError(t, err)
	require.Len(t, ugUpdate3.QueryPolicies().AllX(ctx), 1)
}

func CreateUserEnt(ctx context.Context, client *ent.Client, userName string) (*ent.User, error) {
	return client.User.Create().SetAuthID(userName).SetEmail(userName).Save(ctx)
}
