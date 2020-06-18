// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/facebookincubator/symphony/pkg/ent/usersgroup"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func permissionValue(allowed bool) models.PermissionValue {
	if allowed {
		return models.PermissionValueYes
	}
	return models.PermissionValueNo
}

func newBasicPermissionRuleInput(allowed bool) *models.BasicPermissionRuleInput {
	return &models.BasicPermissionRuleInput{IsAllowed: permissionValue(allowed)}
}

func newLocationPermissionRuleInput(allowed models.PermissionValue, locationTypeIDs []int) *models.LocationPermissionRuleInput {
	return &models.LocationPermissionRuleInput{
		IsAllowed:       allowed,
		LocationTypeIds: locationTypeIDs,
	}
}

func newWorkforcePermissionRuleInput(allowed models.PermissionValue, workOrderTypeIDs []int, projectTypeIDs []int) *models.WorkforcePermissionRuleInput {
	return &models.WorkforcePermissionRuleInput{
		IsAllowed:        allowed,
		WorkOrderTypeIds: workOrderTypeIDs,
		ProjectTypeIds:   projectTypeIDs,
	}
}

type testData struct {
	locationPolicyInput   *models.InventoryPolicyInput
	locationPolicyInput2  *models.InventoryPolicyInput
	catalogInventoryInput *models.InventoryPolicyInput
	workforcePolicyInput1 *models.WorkforcePolicyInput
	workforcePolicyInput2 *models.WorkforcePolicyInput
	locationPolicyID      int
	locationPolicyID2     int
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
		Location: &models.LocationCUDInput{
			Create: newBasicPermissionRuleInput(true),
			Update: newLocationPermissionRuleInput(models.PermissionValueByCondition, []int{rand.Int(), rand.Int()}),
			Delete: newBasicPermissionRuleInput(false),
		},
	}
	data.locationPolicyInput2 = &models.InventoryPolicyInput{
		Location: &models.LocationCUDInput{
			Create: newBasicPermissionRuleInput(true),
			Update: newLocationPermissionRuleInput(models.PermissionValueByCondition, []int{rand.Int(), rand.Int()}),
			Delete: newBasicPermissionRuleInput(true),
		},
	}
	data.catalogInventoryInput = &models.InventoryPolicyInput{
		LocationType: &models.BasicCUDInput{
			Create: newBasicPermissionRuleInput(true),
			Update: newBasicPermissionRuleInput(true),
			Delete: newBasicPermissionRuleInput(true),
		},
		EquipmentType: &models.BasicCUDInput{
			Create: newBasicPermissionRuleInput(true),
			Update: newBasicPermissionRuleInput(true),
			Delete: newBasicPermissionRuleInput(true),
		},
	}
	data.workforcePolicyInput1 = &models.WorkforcePolicyInput{
		Read: newWorkforcePermissionRuleInput(models.PermissionValueYes, nil, nil),
		Data: &models.WorkforceCUDInput{
			Create:            newBasicPermissionRuleInput(false),
			Update:            newBasicPermissionRuleInput(false),
			Delete:            newBasicPermissionRuleInput(true),
			Assign:            newBasicPermissionRuleInput(true),
			TransferOwnership: newBasicPermissionRuleInput(true),
		},
	}
	data.workforcePolicyInput2 = &models.WorkforcePolicyInput{
		Read: newWorkforcePermissionRuleInput(models.PermissionValueByCondition, []int{rand.Int(), rand.Int()}, []int{rand.Int()}),
		Data: &models.WorkforceCUDInput{
			Create: newBasicPermissionRuleInput(false),
			Update: newBasicPermissionRuleInput(true),
			Delete: newBasicPermissionRuleInput(false),
			Assign: newBasicPermissionRuleInput(true),
		},
	}

	p := c.PermissionsPolicy.Create().
		SetName("LocationPolicy").
		SetInventoryPolicy(data.locationPolicyInput).
		SaveX(ctx)
	data.locationPolicyID = p.ID
	p2 := c.PermissionsPolicy.Create().
		SetName("LocationPolicy2").
		SetInventoryPolicy(data.locationPolicyInput2).
		SaveX(ctx)
	data.locationPolicyID2 = p2.ID
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
	require.EqualValues(t, authz.NewInventoryPolicy(false), permissions.InventoryPolicy)
	require.EqualValues(t, authz.NewWorkforcePolicy(false, false), permissions.WorkforcePolicy)
	c.PermissionsPolicy.UpdateOneID(data.locationPolicyID).
		SetIsGlobal(true).
		ExecX(ctx)
	c.PermissionsPolicy.UpdateOneID(data.locationPolicyID2).
		SetIsGlobal(true).
		ExecX(ctx)
	permissions, err = authz.Permissions(ctx)
	require.NoError(t, err)
	require.EqualValues(
		t,
		authz.AppendInventoryPolicies(
			authz.NewInventoryPolicy(false),
			data.locationPolicyInput,
			data.locationPolicyInput2,
		),
		permissions.InventoryPolicy,
	)
	require.EqualValues(t, authz.NewWorkforcePolicy(false, false), permissions.WorkforcePolicy)
	require.Equal(t, models.PermissionValueByCondition, permissions.InventoryPolicy.Location.Create.IsAllowed)
	require.Len(t, permissions.InventoryPolicy.Location.Create.LocationTypeIds, 4)
	require.Equal(t, models.PermissionValueByCondition, permissions.InventoryPolicy.Location.Update.IsAllowed)
	require.Len(t, permissions.InventoryPolicy.Location.Update.LocationTypeIds, 4)
	require.Equal(t, models.PermissionValueByCondition, permissions.InventoryPolicy.Location.Delete.IsAllowed)
	require.Len(t, permissions.InventoryPolicy.Location.Delete.LocationTypeIds, 2)
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
			authz.NewInventoryPolicy(false),
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

func TestPoliciesAreNotAppendedForDeactivatedGroups(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c, viewertest.WithRole(user.RoleUSER))
	data := prepareData(ctx)
	c.UsersGroup.UpdateOneID(data.group1ID).
		AddPolicyIDs(data.catalogPolicyID).
		ExecX(ctx)
	c.UsersGroup.UpdateOneID(data.group2ID).
		AddPolicyIDs(data.workforcePolicyID).
		SetStatus(usersgroup.StatusDEACTIVATED).
		ExecX(ctx)
	permissions, err := authz.Permissions(ctx)
	require.NoError(t, err)
	require.EqualValues(
		t,
		authz.AppendInventoryPolicies(
			authz.NewInventoryPolicy(false),
			data.catalogInventoryInput,
		),
		permissions.InventoryPolicy,
	)
	require.EqualValues(
		t,
		authz.NewWorkforcePolicy(false, false),
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
		authz.NewInventoryPolicy(false),
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

	require.Equal(t, models.PermissionValueYes, permissions.WorkforcePolicy.Read.IsAllowed)
	require.Nil(t, permissions.WorkforcePolicy.Read.WorkOrderTypeIds)
	require.Nil(t, permissions.WorkforcePolicy.Read.ProjectTypeIds)
	require.Equal(t, models.PermissionValueNo, permissions.WorkforcePolicy.Data.Create.IsAllowed)
	require.Nil(t, permissions.WorkforcePolicy.Data.Create.WorkOrderTypeIds)
	require.Nil(t, permissions.WorkforcePolicy.Data.Create.ProjectTypeIds)
	require.Equal(t, models.PermissionValueByCondition, permissions.WorkforcePolicy.Data.Update.IsAllowed)
	require.Len(t, permissions.WorkforcePolicy.Data.Update.WorkOrderTypeIds, 2)
	require.Len(t, permissions.WorkforcePolicy.Data.Update.ProjectTypeIds, 1)
	require.Equal(t, models.PermissionValueYes, permissions.WorkforcePolicy.Data.Delete.IsAllowed)
	require.Nil(t, permissions.WorkforcePolicy.Data.Delete.WorkOrderTypeIds)
	require.Nil(t, permissions.WorkforcePolicy.Data.Delete.ProjectTypeIds)
	require.Equal(t, models.PermissionValueYes, permissions.WorkforcePolicy.Data.Assign.IsAllowed)
	require.Nil(t, permissions.WorkforcePolicy.Data.Assign.WorkOrderTypeIds)
	require.Nil(t, permissions.WorkforcePolicy.Data.Assign.ProjectTypeIds)
	require.Equal(t, models.PermissionValueYes, permissions.WorkforcePolicy.Data.TransferOwnership.IsAllowed)
	require.Nil(t, permissions.WorkforcePolicy.Data.TransferOwnership.WorkOrderTypeIds)
	require.Nil(t, permissions.WorkforcePolicy.Data.TransferOwnership.ProjectTypeIds)
}

func TestAdminUserHasAdminEditPermissions(t *testing.T) {
	const admin = "admin_user"
	client := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), client)
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
	ctx := viewertest.NewContext(context.Background(), client)
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
	ctx := viewertest.NewContext(context.Background(), client)
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
	ctx := viewertest.NewContext(context.Background(), client)
	u, err := client.User.Create().SetAuthID(userInGroup).Save(ctx)
	require.NoError(t, err)
	_, err = client.UsersGroup.Create().SetName(authz.WritePermissionGroupName).AddMembers(u).Save(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(context.Background(), client, viewertest.WithUser(userInGroup), viewertest.WithFeatures())
	permissions, err := authz.Permissions(ctx)
	expectedPermissions := authz.FullPermissions()
	expectedPermissions.AdminPolicy.Access.IsAllowed = models.PermissionValueNo
	expectedPermissions.CanWrite = true
	require.NoError(t, err)
	require.EqualValues(t, expectedPermissions, permissions)
}
