// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func NewBasicPermissionRuleInput(allowed bool) *models.BasicPermissionRuleInput {
	rule := models.PermissionValueNo
	if allowed {
		rule = models.PermissionValueYes
	}
	return &models.BasicPermissionRuleInput{IsAllowed: rule}
}

type testData struct {
	locationPolicyInput   *models.InventoryPolicyInput
	catalogInventoryInput *models.InventoryPolicyInput
	workforcePolicyInput1 *models.WorkforcePolicyInput
	workforcePolicyInput2 *models.WorkforcePolicyInput
	locationPolicyID      int
	catalogPolicyID       int
	workforcePolicyID     int
	workforcePolicy2ID    int
	group1ID              int
	group2ID              int
	group3ID              int
}

func prepareData(ctx context.Context) (data testData) {
	c := ent.FromContext(ctx)
	v := viewer.FromContext(ctx).(*viewer.UserViewer)
	data.locationPolicyInput = &models.InventoryPolicyInput{
		Location: &models.BasicCUDInput{
			Create: NewBasicPermissionRuleInput(true),
			Update: NewBasicPermissionRuleInput(true),
			Delete: NewBasicPermissionRuleInput(false),
		},
	}
	data.catalogInventoryInput = &models.InventoryPolicyInput{
		LocationType: &models.BasicCUDInput{
			Create: NewBasicPermissionRuleInput(true),
			Update: NewBasicPermissionRuleInput(true),
			Delete: NewBasicPermissionRuleInput(true),
		},
		EquipmentType: &models.BasicCUDInput{
			Create: NewBasicPermissionRuleInput(true),
			Update: NewBasicPermissionRuleInput(true),
			Delete: NewBasicPermissionRuleInput(true),
		},
	}
	data.workforcePolicyInput1 = &models.WorkforcePolicyInput{
		Data: &models.BasicWorkforceCUDInput{
			Create: NewBasicPermissionRuleInput(false),
			Update: NewBasicPermissionRuleInput(false),
			Delete: NewBasicPermissionRuleInput(true),
			Assign: NewBasicPermissionRuleInput(true),
		},
	}
	data.workforcePolicyInput2 = &models.WorkforcePolicyInput{
		Data: &models.BasicWorkforceCUDInput{
			Create:            NewBasicPermissionRuleInput(false),
			Update:            NewBasicPermissionRuleInput(true),
			Delete:            NewBasicPermissionRuleInput(false),
			Assign:            NewBasicPermissionRuleInput(true),
			TransferOwnership: NewBasicPermissionRuleInput(true),
		},
	}

	p := c.PermissionsPolicy.Create().
		SetName("LocationPolicy").
		SetInventoryPolicy(data.locationPolicyInput).
		SaveX(ctx)
	data.locationPolicyID = p.ID
	p = c.PermissionsPolicy.Create().
		SetName("CatalogPolicy").
		SetInventoryPolicy(data.catalogInventoryInput).
		SaveX(ctx)
	data.catalogPolicyID = p.ID
	p = c.PermissionsPolicy.Create().
		SetName("WorkforcePolicy").
		SetWorkforcePolicy(data.workforcePolicyInput1).
		SaveX(ctx)
	data.workforcePolicyID = p.ID
	p = c.PermissionsPolicy.Create().
		SetName("WorkforcePolicy2").
		SetWorkforcePolicy(data.workforcePolicyInput2).
		SaveX(ctx)
	data.workforcePolicy2ID = p.ID

	g := c.UsersGroup.Create().
		SetName("Group1").
		AddMembers(v.User()).
		SaveX(ctx)
	data.group1ID = g.ID
	g = c.UsersGroup.Create().
		SetName("Group2").
		AddMembers(v.User()).
		SaveX(ctx)
	data.group2ID = g.ID
	g = c.UsersGroup.Create().
		SetName("Group3").
		AddMembers(v.User()).
		SaveX(ctx)
	data.group3ID = g.ID
	return
}

func TestGlobalPolicyIsAppliedForUsers(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c, viewertest.WithRole(user.RoleUSER))
	data := prepareData(ctx)
	permissions, err := authz.Permissions(ctx)
	require.NoError(t, err)
	require.EqualValues(t, authz.NewInventoryPolicy(false, false), permissions.InventoryPolicy)
	require.EqualValues(t, authz.NewWorkforcePolicy(false, false), permissions.WorkforcePolicy)
	c.PermissionsPolicy.UpdateOneID(data.locationPolicyID).
		SetIsGlobal(true).
		ExecX(ctx)
	permissions, err = authz.Permissions(ctx)
	require.NoError(t, err)
	require.EqualValues(
		t,
		authz.AppendInventoryPolicies(
			authz.NewInventoryPolicy(false, false),
			data.locationPolicyInput,
		),
		permissions.InventoryPolicy,
	)
	require.EqualValues(t, authz.NewWorkforcePolicy(false, false), permissions.WorkforcePolicy)
}

func TestPoliciesAreAppendedForGroups(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c, viewertest.WithRole(user.RoleUSER))
	data := prepareData(ctx)
	c.PermissionsPolicy.UpdateOneID(data.locationPolicyID).
		SetIsGlobal(true).
		ExecX(ctx)
	c.UsersGroup.UpdateOneID(data.group1ID).
		AddPolicyIDs(data.catalogPolicyID).
		ExecX(ctx)
	c.UsersGroup.UpdateOneID(data.group2ID).
		AddPolicyIDs(data.catalogPolicyID).
		ExecX(ctx)
	c.UsersGroup.UpdateOneID(data.group3ID).
		AddPolicyIDs(data.workforcePolicyID).
		ExecX(ctx)
	permissions, err := authz.Permissions(ctx)
	require.NoError(t, err)
	require.EqualValues(
		t,
		authz.AppendInventoryPolicies(
			authz.NewInventoryPolicy(false, false),
			data.locationPolicyInput,
			data.catalogInventoryInput,
		),
		permissions.InventoryPolicy,
	)
	require.EqualValues(
		t,
		authz.AppendWorkforcePolicies(
			authz.NewWorkforcePolicy(false, false),
			data.workforcePolicyInput1,
		),
		permissions.WorkforcePolicy,
	)
}

func TestPoliciesAppendingOutput(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c, viewertest.WithRole(user.RoleUSER))
	data := prepareData(ctx)
	c.UsersGroup.UpdateOneID(data.group1ID).
		AddPolicyIDs(data.workforcePolicyID, data.workforcePolicy2ID).
		ExecX(ctx)
	permissions, err := authz.Permissions(ctx)
	require.NoError(t, err)
	require.EqualValues(
		t,
		authz.NewInventoryPolicy(false, false),
		permissions.InventoryPolicy,
	)
	require.EqualValues(
		t,
		authz.AppendWorkforcePolicies(
			authz.NewWorkforcePolicy(false, false),
			data.workforcePolicyInput1,
			data.workforcePolicyInput2,
		),
		permissions.WorkforcePolicy,
	)
	require.Equal(t, models.PermissionValueNo, permissions.WorkforcePolicy.Data.Create.IsAllowed)
	require.Equal(t, models.PermissionValueYes, permissions.WorkforcePolicy.Data.Update.IsAllowed)
	require.Equal(t, models.PermissionValueYes, permissions.WorkforcePolicy.Data.Delete.IsAllowed)
	require.Equal(t, models.PermissionValueYes, permissions.WorkforcePolicy.Data.Assign.IsAllowed)
	require.Equal(t, models.PermissionValueYes, permissions.WorkforcePolicy.Data.TransferOwnership.IsAllowed)
}

func TestAdminUserHasAdminEditPermissions(t *testing.T) {
	const admin = "admin_user"
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	_, err := client.User.Create().
		SetAuthID(admin).
		SetRole(user.RoleADMIN).
		Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(
		context.Background(), client,
		viewertest.WithUser(admin))
	permissions, err := authz.Permissions(ctx)
	require.NoError(t, err)
	require.Equal(t, models.PermissionValueYes, permissions.AdminPolicy.Access.IsAllowed)
	require.Equal(t, false, permissions.CanWrite)
}

func TestUserHasNoReadonlyPermissions(t *testing.T) {
	const regular = "regular_user"
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	_, err := client.User.Create().SetAuthID(regular).SetRole(user.RoleUSER).Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser(regular))
	permissions, err := authz.Permissions(ctx)
	require.NoError(t, err)
	expectedPermissions := authz.EmptyPermissions()
	require.EqualValues(t, expectedPermissions, permissions)
}

func TestOwnerHasWritePermissions(t *testing.T) {
	const owner = "owner_user"
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	_, err := client.User.Create().SetAuthID(owner).SetRole(user.RoleOWNER).Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser(owner))
	permissions, err := authz.Permissions(ctx)
	require.NoError(t, err)
	require.EqualValues(t, authz.FullPermissions(), permissions)
}

func TestUserInGroupHasWritePermissionsButNoAdmin(t *testing.T) {
	const userInGroup = "user_in_group"
	client := viewertest.NewTestClient(t)
	ctx := ent.NewContext(context.Background(), client)
	u, err := client.User.Create().SetAuthID(userInGroup).Save(ctx)
	require.NoError(t, err)
	_, err = client.UsersGroup.Create().SetName(authz.WritePermissionGroupName).AddMembers(u).Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser(userInGroup))
	permissions, err := authz.Permissions(ctx)
	expectedPermissions := authz.FullPermissions()
	expectedPermissions.AdminPolicy.Access.IsAllowed = models.PermissionValueNo
	require.NoError(t, err)
	require.EqualValues(t, expectedPermissions, permissions)
}
