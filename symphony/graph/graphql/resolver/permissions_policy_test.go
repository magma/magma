// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

const (
	policyName        = "my_policy"
	policyDescription = "Same Description"
)

func getInventoryPolicyInput() *models2.InventoryPolicyInput {
	return &models2.InventoryPolicyInput{
		Read: &models2.BasicPermissionRuleInput{IsAllowed: models2.PermissionValueYes},
		Equipment: &models2.BasicCUDInput{
			Create: &models2.BasicPermissionRuleInput{IsAllowed: models2.PermissionValueYes},
			Update: &models2.BasicPermissionRuleInput{IsAllowed: models2.PermissionValueNo},
			Delete: &models2.BasicPermissionRuleInput{IsAllowed: models2.PermissionValueByCondition},
		},
	}
}

func getWorkforcePolicyInput() *models2.WorkforcePolicyInput {
	return &models2.WorkforcePolicyInput{
		Read: &models2.WorkforcePermissionRuleInput{IsAllowed: models2.PermissionValueNo},
		Data: &models2.WorkforceCUDInput{
			Create: &models2.WorkforcePermissionRuleInput{IsAllowed: models2.PermissionValueYes},
			Assign: &models2.WorkforcePermissionRuleInput{IsAllowed: models2.PermissionValueByCondition},
		},
	}
}

func TestQueryInventoryPolicies(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	inventoryPolicyInput := getInventoryPolicyInput()
	_, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           policyName,
		Description:    pointer.ToString(policyDescription),
		InventoryInput: inventoryPolicyInput,
		WorkforceInput: nil,
	})
	require.NoError(t, err)

	ppc, err := qr.PermissionsPolicies(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, ppc.Edges, 1)
}

func TestAddInventoryPolicy(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, ppr := r.Mutation(), r.PermissionsPolicy()

	inventoryPolicyInput := getInventoryPolicyInput()
	policy, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           policyName,
		Description:    pointer.ToString(policyDescription),
		InventoryInput: inventoryPolicyInput,
		WorkforceInput: nil,
	})
	require.NoError(t, err)
	require.Equal(t, policyName, policy.Name)
	require.Equal(t, policyDescription, policy.Description)
	res, err := ppr.Policy(ctx, policy)
	require.NoError(t, err)
	inventoryPolicy, ok := res.(*models.InventoryPolicy)
	require.True(t, ok)

	require.Equal(t, models2.PermissionValueYes, inventoryPolicy.Read.IsAllowed)

	require.Equal(t, models2.PermissionValueNo, inventoryPolicy.Location.Create.IsAllowed)
	require.Equal(t, models2.PermissionValueNo, inventoryPolicy.Location.Update.IsAllowed)
	require.Equal(t, models2.PermissionValueNo, inventoryPolicy.Location.Delete.IsAllowed)

	require.Equal(t, models2.PermissionValueYes, inventoryPolicy.Equipment.Create.IsAllowed)
	require.Equal(t, models2.PermissionValueNo, inventoryPolicy.Equipment.Update.IsAllowed)
	require.Equal(t, models2.PermissionValueByCondition, inventoryPolicy.Equipment.Delete.IsAllowed)
}

func TestAddWorkOrderPolicy(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, ppr := r.Mutation(), r.PermissionsPolicy()

	workforcePolicyInput := getWorkforcePolicyInput()
	policy, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           policyName,
		Description:    pointer.ToString(policyDescription),
		InventoryInput: nil,
		WorkforceInput: workforcePolicyInput,
	})
	require.NoError(t, err)
	require.Equal(t, policyName, policy.Name)
	require.Equal(t, policyDescription, policy.Description)
	res, err := ppr.Policy(ctx, policy)
	require.NoError(t, err)
	workforcePolicy, ok := res.(*models.WorkforcePolicy)
	require.True(t, ok)

	require.Equal(t, models2.PermissionValueNo, workforcePolicy.Read.IsAllowed)

	require.Equal(t, models2.PermissionValueYes, workforcePolicy.Data.Create.IsAllowed)
	require.Equal(t, models2.PermissionValueNo, workforcePolicy.Data.Delete.IsAllowed)
	require.Equal(t, models2.PermissionValueByCondition, workforcePolicy.Data.Assign.IsAllowed)

	require.Equal(t, models2.PermissionValueNo, workforcePolicy.Templates.Create.IsAllowed)
	require.Equal(t, models2.PermissionValueNo, workforcePolicy.Templates.Update.IsAllowed)
	require.Equal(t, models2.PermissionValueNo, workforcePolicy.Templates.Delete.IsAllowed)
}

func TestAddMultipleTypesPermissionsPolicy(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	_, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           policyName,
		Description:    pointer.ToString(policyDescription),
		InventoryInput: getInventoryPolicyInput(),
		WorkforceInput: getWorkforcePolicyInput(),
	})
	require.Error(t, err)
}

func TestDeletePermissionsPolicy(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	_, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           policyName,
		Description:    pointer.ToString(policyDescription),
		InventoryInput: getInventoryPolicyInput(),
		WorkforceInput: nil,
	})
	require.NoError(t, err)

	client := ent.FromContext(ctx)
	pps := client.PermissionsPolicy.Query().AllX(ctx)
	require.Len(t, pps, 1)

	_, err = mr.DeletePermissionsPolicy(ctx, pps[0].ID)
	require.NoError(t, err)
}

func TestAddEmptyPermissionsPolicy(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	_, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           policyName,
		Description:    pointer.ToString(policyDescription),
		InventoryInput: nil,
		WorkforceInput: nil,
	})
	require.Error(t, err)
}

func TestAddPermissionsPolicyWithGroup(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	gName1 := "group_1"
	addInp1 := getAddUsersGroupInput(gName1, "this is group 1")
	ug1, err := mr.AddUsersGroup(ctx, addInp1)
	require.NoError(t, err)

	_, err = mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           policyName,
		Description:    pointer.ToString(policyDescription),
		InventoryInput: getInventoryPolicyInput(),
		WorkforceInput: nil,
		Groups:         []int{ug1.ID},
	})
	require.NoError(t, err)
}

func TestEditPermissionsPolicy(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	gName1 := "group_1"
	addInp1 := getAddUsersGroupInput(gName1, "this is group 1")
	ug1, err := mr.AddUsersGroup(ctx, addInp1)
	require.NoError(t, err)

	gName2 := "group_2"
	addInp2 := getAddUsersGroupInput(gName2, "this is group 2")
	ug2, err := mr.AddUsersGroup(ctx, addInp2)
	require.NoError(t, err)

	inventoryPolicyInput := getInventoryPolicyInput()
	workforcePolicyInput := getWorkforcePolicyInput()

	policy, err := mr.AddPermissionsPolicy(ctx, models.AddPermissionsPolicyInput{
		Name:           policyName,
		Description:    pointer.ToString(policyDescription),
		InventoryInput: inventoryPolicyInput,
		WorkforceInput: nil,
	})
	require.NoError(t, err)
	require.Equal(t, policy.InventoryPolicy, inventoryPolicyInput)
	require.Empty(t, policy.WorkforcePolicy)

	newPolicyName := "new_" + policyName
	newDescription := "New " + policyDescription
	newInventoryPolicy := &models2.InventoryPolicyInput{
		Location: &models2.LocationCUDInput{
			Create: &models2.LocationPermissionRuleInput{IsAllowed: models2.PermissionValueYes},
			Update: &models2.LocationPermissionRuleInput{IsAllowed: models2.PermissionValueNo},
			Delete: &models2.LocationPermissionRuleInput{IsAllowed: models2.PermissionValueByCondition},
		},
	}
	fetchedPermissionsPolicy1, err := mr.EditPermissionsPolicy(ctx, models.EditPermissionsPolicyInput{
		ID:             policy.ID,
		Name:           &newPolicyName,
		Description:    pointer.ToString(newDescription),
		InventoryInput: newInventoryPolicy,
		WorkforceInput: nil,
	})
	require.NoError(t, err)
	require.Equal(t, fetchedPermissionsPolicy1.Name, newPolicyName)
	require.Equal(t, fetchedPermissionsPolicy1.Description, newDescription)
	require.Equal(t, fetchedPermissionsPolicy1.InventoryPolicy, newInventoryPolicy)
	require.Equal(t, fetchedPermissionsPolicy1.InventoryPolicy.Location.Create.IsAllowed, models2.PermissionValueYes)
	require.Equal(t, fetchedPermissionsPolicy1.InventoryPolicy.Location.Update.IsAllowed, models2.PermissionValueNo)
	require.Equal(t, fetchedPermissionsPolicy1.InventoryPolicy.Location.Delete.IsAllowed, models2.PermissionValueByCondition)
	require.Empty(t, fetchedPermissionsPolicy1.WorkforcePolicy)

	fetchedPermissionsPolicy2, err := mr.EditPermissionsPolicy(ctx, models.EditPermissionsPolicyInput{
		ID:             policy.ID,
		Name:           nil,
		Description:    nil,
		InventoryInput: nil,
		WorkforceInput: nil,
	})
	require.NoError(t, err)
	require.Equal(t, fetchedPermissionsPolicy2.Name, fetchedPermissionsPolicy1.Name)
	require.Equal(t, fetchedPermissionsPolicy2.Description, fetchedPermissionsPolicy1.Description)
	require.Equal(t, fetchedPermissionsPolicy2.InventoryPolicy, fetchedPermissionsPolicy1.InventoryPolicy)
	require.Empty(t, fetchedPermissionsPolicy2.WorkforcePolicy)

	_, err = mr.EditPermissionsPolicy(ctx, models.EditPermissionsPolicyInput{
		ID:             policy.ID,
		Name:           &newPolicyName,
		Description:    pointer.ToString(newDescription),
		InventoryInput: nil,
		WorkforceInput: workforcePolicyInput,
	})
	require.Error(t, err)

	updateGroupsInput1 := models.EditPermissionsPolicyInput{
		ID:             policy.ID,
		Name:           nil,
		Description:    nil,
		InventoryInput: nil,
		WorkforceInput: nil,
		Groups:         []int{ug1.ID},
	}
	ugUpdate1, err := mr.EditPermissionsPolicy(ctx, updateGroupsInput1)
	require.NoError(t, err)
	require.Len(t, ugUpdate1.QueryGroups().AllX(ctx), 1)

	updateGroupsInput2 := models.EditPermissionsPolicyInput{
		ID:             policy.ID,
		Name:           nil,
		Description:    nil,
		InventoryInput: nil,
		WorkforceInput: nil,
		Groups:         []int{ug1.ID, ug2.ID},
	}
	ugUpdate2, err := mr.EditPermissionsPolicy(ctx, updateGroupsInput2)
	require.NoError(t, err)
	require.Len(t, ugUpdate2.QueryGroups().AllX(ctx), 2)

	updateGroupsInput3 := models.EditPermissionsPolicyInput{
		ID:             policy.ID,
		Name:           nil,
		Description:    nil,
		InventoryInput: nil,
		WorkforceInput: nil,
		Groups:         []int{ug2.ID},
	}
	ugUpdate3, err := mr.EditPermissionsPolicy(ctx, updateGroupsInput3)
	require.NoError(t, err)
	require.Len(t, ugUpdate3.QueryGroups().AllX(ctx), 1)
}
