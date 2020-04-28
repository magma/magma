package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/authz/models"
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

type policyInputs struct {
	locationPolicyInput   *models.InventoryPolicyInput
	catalogInventoryInput *models.InventoryPolicyInput
	workforcePolicyInput1 *models.WorkforcePolicyInput
	workforcePolicyInput2 *models.WorkforcePolicyInput
}

func preparePolicyInputs() (inputs policyInputs) {
	inputs.locationPolicyInput = &models.InventoryPolicyInput{
		Location: &models.BasicCUDInput{
			Create: NewBasicPermissionRuleInput(true),
			Update: NewBasicPermissionRuleInput(true),
			Delete: NewBasicPermissionRuleInput(false),
		},
	}
	inputs.catalogInventoryInput = &models.InventoryPolicyInput{
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
	inputs.workforcePolicyInput1 = &models.WorkforcePolicyInput{
		Data: &models.BasicWorkforceCUDInput{
			Create: NewBasicPermissionRuleInput(false),
			Update: NewBasicPermissionRuleInput(false),
			Delete: NewBasicPermissionRuleInput(true),
			Assign: NewBasicPermissionRuleInput(true),
		},
	}
	inputs.workforcePolicyInput2 = &models.WorkforcePolicyInput{
		Data: &models.BasicWorkforceCUDInput{
			Create:            NewBasicPermissionRuleInput(false),
			Update:            NewBasicPermissionRuleInput(true),
			Delete:            NewBasicPermissionRuleInput(false),
			Assign:            NewBasicPermissionRuleInput(true),
			TransferOwnership: NewBasicPermissionRuleInput(true),
		},
	}
	return
}

func TestGlobalPolicyIsAppliedForUsers(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	inputs := preparePolicyInputs()
	p := c.PermissionsPolicy.Create().
		SetName("LocationPolicy").
		SetInventoryPolicy(inputs.locationPolicyInput).
		SaveX(ctx)
	inventoryPolicy, workforcePolicy, err := authz.PermissionPolicies(ctx)
	require.NoError(t, err)
	require.EqualValues(t, authz.NewInventoryPolicy(false, false), inventoryPolicy)
	require.EqualValues(t, authz.NewWorkforcePolicy(false, false), workforcePolicy)
	c.PermissionsPolicy.UpdateOne(p).SetIsGlobal(true).ExecX(ctx)
	inventoryPolicy, workforcePolicy, err = authz.PermissionPolicies(ctx)
	require.NoError(t, err)
	require.EqualValues(
		t,
		authz.AppendInventoryPolicies(
			authz.NewInventoryPolicy(false, false),
			inputs.locationPolicyInput,
		),
		inventoryPolicy,
	)
	require.EqualValues(t, authz.NewWorkforcePolicy(false, false), workforcePolicy)
}

func TestPoliciesAreAppendedForGroups(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	inputs := preparePolicyInputs()
	_ = c.PermissionsPolicy.Create().
		SetName("LocationPolicy").
		SetInventoryPolicy(inputs.locationPolicyInput).
		SetIsGlobal(true).
		SaveX(ctx)
	p := c.PermissionsPolicy.Create().
		SetName("CatalogPolicy").
		SetInventoryPolicy(inputs.catalogInventoryInput).
		SaveX(ctx)
	p2 := c.PermissionsPolicy.Create().
		SetName("WorkforcePolicy").
		SetWorkforcePolicy(inputs.workforcePolicyInput1).
		SaveX(ctx)
	_ = c.UsersGroup.Create().
		SetName("Group1").
		AddMembers(viewer.FromContext(ctx).User()).
		AddPolicies(p).
		SaveX(ctx)
	_ = c.UsersGroup.Create().
		SetName("Group2").
		AddMembers(viewer.FromContext(ctx).User()).
		AddPolicies(p).
		SaveX(ctx)
	_ = c.UsersGroup.Create().
		SetName("Group3").
		AddMembers(viewer.FromContext(ctx).User()).
		AddPolicies(p2).
		SaveX(ctx)
	inventoryPolicy, workforcePolicy, err := authz.PermissionPolicies(ctx)
	require.NoError(t, err)
	require.EqualValues(
		t,
		authz.AppendInventoryPolicies(
			authz.NewInventoryPolicy(false, false),
			inputs.locationPolicyInput,
			inputs.catalogInventoryInput,
		),
		inventoryPolicy,
	)
	require.EqualValues(
		t,
		authz.AppendWorkforcePolicies(
			authz.NewWorkforcePolicy(false, false),
			inputs.workforcePolicyInput1,
		),
		workforcePolicy,
	)
}

func TestPoliciesAppendingOutput(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	inputs := preparePolicyInputs()
	p := c.PermissionsPolicy.Create().
		SetName("WorkforcePolicy").
		SetWorkforcePolicy(inputs.workforcePolicyInput1).
		SaveX(ctx)
	p2 := c.PermissionsPolicy.Create().
		SetName("WorkforcePolicy2").
		SetWorkforcePolicy(inputs.workforcePolicyInput2).
		SaveX(ctx)
	_ = c.UsersGroup.Create().
		SetName("Group1").
		AddMembers(viewer.FromContext(ctx).User()).
		AddPolicies(p, p2).
		SaveX(ctx)
	inventoryPolicy, workforcePolicy, err := authz.PermissionPolicies(ctx)
	require.NoError(t, err)
	require.EqualValues(
		t,
		authz.NewInventoryPolicy(false, false),
		inventoryPolicy,
	)
	require.EqualValues(
		t,
		authz.AppendWorkforcePolicies(
			authz.NewWorkforcePolicy(false, false),
			inputs.workforcePolicyInput1,
			inputs.workforcePolicyInput2,
		),
		workforcePolicy,
	)
	require.Equal(t, models.PermissionValueNo, workforcePolicy.Data.Create.IsAllowed)
	require.Equal(t, models.PermissionValueYes, workforcePolicy.Data.Update.IsAllowed)
	require.Equal(t, models.PermissionValueYes, workforcePolicy.Data.Delete.IsAllowed)
	require.Equal(t, models.PermissionValueYes, workforcePolicy.Data.Assign.IsAllowed)
	require.Equal(t, models.PermissionValueYes, workforcePolicy.Data.TransferOwnership.IsAllowed)
}
