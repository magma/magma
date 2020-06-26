// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
)

func prepareWorkOrderData(ctx context.Context, c *ent.Client) (*ent.WorkOrderType, *ent.WorkOrder) {
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleOWNER)
	workOrderTypeName := uuid.New().String()
	workOrderName := uuid.New().String()
	workOrderType := c.WorkOrderType.Create().
		SetName(workOrderTypeName).
		SaveX(ctx)
	workOrder := c.WorkOrder.Create().
		SetName(workOrderName).
		SetType(workOrderType).
		SetOwner(u).
		SetCreationDate(time.Now()).
		SaveX(ctx)
	return workOrderType, workOrder
}

func TestNonUserCannotEditWorkOrder(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	_, workOrder := prepareWorkOrderData(ctx, c)

	v := viewer.NewAutomation(viewertest.DefaultTenant, "BOT", user.RoleUSER)
	ctx = viewer.NewContext(ctx, v)
	ctx = authz.NewContext(ctx, authz.EmptyPermissions())
	err := c.WorkOrder.UpdateOne(workOrder).
		SetName("NewName").
		Exec(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
}

func TestAssignCanEditWOWithOwnerAndDelete(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	_, workOrder := prepareWorkOrderData(ctx, c)
	u := viewer.MustGetOrCreateUser(ctx, "MyAssignee", user.RoleUSER)
	c.WorkOrder.UpdateOne(workOrder).
		SetAssignee(u).
		ExecX(ctx)

	ctx = viewertest.NewContext(ctx, c,
		viewertest.WithUser("MyAssignee"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(authz.EmptyPermissions()))
	err := c.WorkOrder.UpdateOne(workOrder).
		SetName("NewName").
		Exec(ctx)
	require.NoError(t, err)
	err = c.WorkOrder.UpdateOne(workOrder).
		SetOwner(u).
		Exec(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
	err = c.WorkOrder.DeleteOne(workOrder).
		Exec(ctx)
	require.True(t, errors.Is(err, privacy.Deny))
}

func TestOwnerCanEditWO(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	_, workOrder := prepareWorkOrderData(ctx, c)
	u := viewer.MustGetOrCreateUser(ctx, "MyOwner", user.RoleUSER)
	u2 := viewer.MustGetOrCreateUser(ctx, "NewOwner", user.RoleUSER)
	c.WorkOrder.UpdateOne(workOrder).
		SetOwner(u).
		ExecX(ctx)

	ctx = viewertest.NewContext(ctx, c,
		viewertest.WithUser("MyOwner"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(authz.EmptyPermissions()))
	err := c.WorkOrder.UpdateOne(workOrder).
		SetName("NewName").
		Exec(ctx)
	require.NoError(t, err)
	err = c.WorkOrder.UpdateOne(workOrder).
		SetOwner(u2).
		Exec(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(ctx, c,
		viewertest.WithUser("NewOwner"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(authz.EmptyPermissions()))
	err = c.WorkOrder.DeleteOne(workOrder).
		Exec(ctx)
	require.NoError(t, err)
}

func TestWorkOrderWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, workOrder := prepareWorkOrderData(ctx, c)
	workOrderType2, workOrder2 := prepareWorkOrderData(ctx, c)
	createWorkOrder := func(ctx context.Context) error {
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
		_, err := c.WorkOrder.Create().
			SetName("NewWorkOrder").
			SetType(workOrderType).
			SetOwner(u).
			SetCreationDate(time.Now()).
			Save(ctx)
		return err
	}
	createWorkOrder2 := func(ctx context.Context) error {
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
		_, err := c.WorkOrder.Create().
			SetName("NewWorkOrder2").
			SetType(workOrderType2).
			SetOwner(u).
			SetCreationDate(time.Now()).
			Save(ctx)
		return err
	}
	updateWorkOrder := func(ctx context.Context) error {
		return c.WorkOrder.UpdateOne(workOrder).
			SetName("NewName").
			Exec(ctx)
	}
	updateWorkOrder2 := func(ctx context.Context) error {
		return c.WorkOrder.UpdateOne(workOrder2).
			SetName("NewName2").
			Exec(ctx)
	}
	deleteWorkOrder := func(ctx context.Context) error {
		return c.WorkOrder.DeleteOne(workOrder).
			Exec(ctx)
	}
	deleteWorkOrder2 := func(ctx context.Context) error {
		return c.WorkOrder.DeleteOne(workOrder2).
			Exec(ctx)
	}
	initialPermissions := func(p *models.PermissionSettings) {
		p.WorkforcePolicy.Data.TransferOwnership.IsAllowed = models.PermissionValueYes
	}
	tests := []policyTest{
		{
			operationName:      "Create",
			initialPermissions: initialPermissions,
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Create.IsAllowed = models.PermissionValueYes
			},
			operation: createWorkOrder,
		},
		{
			operationName: "CreateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				initialPermissions(p)
				p.WorkforcePolicy.Data.Create.IsAllowed = models.PermissionValueByCondition
				p.WorkforcePolicy.Data.Create.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Create.WorkOrderTypeIds = []int{workOrderType.ID, workOrderType2.ID}
			},
			operation: createWorkOrder2,
		},
		{
			operationName:      "Update",
			initialPermissions: initialPermissions,
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueYes
			},
			operation: updateWorkOrder,
		},
		{
			operationName: "UpdateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				initialPermissions(p)
				p.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueByCondition
				p.WorkforcePolicy.Data.Update.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Update.WorkOrderTypeIds = []int{workOrderType.ID, workOrderType2.ID}
			},
			operation: updateWorkOrder2,
		},
		{
			operationName:      "Delete",
			initialPermissions: initialPermissions,
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Delete.IsAllowed = models.PermissionValueYes
			},
			operation: deleteWorkOrder,
		},
		{
			operationName: "DeleteWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				initialPermissions(p)
				p.WorkforcePolicy.Data.Delete.IsAllowed = models.PermissionValueByCondition
				p.WorkforcePolicy.Data.Delete.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Delete.WorkOrderTypeIds = []int{workOrderType.ID, workOrderType2.ID}
			},
			operation: deleteWorkOrder2,
		},
	}
	runPolicyTest(t, tests)
}

func TestWorkOrderReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, _ := prepareWorkOrderData(ctx, c)
	_, _ = prepareWorkOrderData(ctx, c)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.WorkOrder.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Zero(t, count)
	})
	t.Run("PartialPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models.PermissionValueByCondition
		permissions.WorkforcePolicy.Read.WorkOrderTypeIds = []int{workOrderType.ID}
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.WorkOrder.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
	t.Run("FullPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models.PermissionValueYes
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.WorkOrder.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

func TestWorkOrderTransferOwnershipWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, workOrder := prepareWorkOrderData(ctx, c)
	getCud := func(p *models.PermissionSettings) *models.WorkforceCud {
		return p.WorkforcePolicy.Data
	}
	u := viewer.MustGetOrCreateUser(ctx, "SomeUser", user.RoleUSER)
	u2 := viewer.MustGetOrCreateUser(ctx, "NewUser", user.RoleUSER)
	createWorkOrderWithOwner := func(ctx context.Context) error {
		_, err := c.WorkOrder.Create().
			SetName("NewWorkOrder").
			SetType(workOrderType).
			SetCreationDate(time.Now()).
			SetOwner(u).
			Save(ctx)
		return err
	}
	updateWorkOrderOwner := func(user *ent.User) func(context.Context) error {
		return func(ctx context.Context) error {
			_, err := c.WorkOrder.UpdateOne(workOrder).
				SetOwner(user).
				Save(ctx)
			return err
		}
	}
	tests := []policyTest{
		{
			operationName: "CreateWithOwner",
			appendPermissions: func(p *models.PermissionSettings) {
				getCud(p).Create.IsAllowed = models.PermissionValueYes
			},
			operation: createWorkOrderWithOwner,
		},
		{
			operationName: "UpdateWithOwner",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueYes
			},
			appendPermissions: func(p *models.PermissionSettings) {
				getCud(p).TransferOwnership.IsAllowed = models.PermissionValueYes
			},
			operation: updateWorkOrderOwner(u),
		},
		{
			operationName: "UpdateWithOwnerWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueByCondition
				getCud(p).Update.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				getCud(p).TransferOwnership.IsAllowed = models.PermissionValueByCondition
				getCud(p).TransferOwnership.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			operation: updateWorkOrderOwner(u2),
		},
	}
	runPolicyTest(t, tests)
}

func TestWorkOrderCreateWithViewerAssigneeOwner(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, _ := prepareWorkOrderData(ctx, c)
	ctx = userContextWithPermissions(ctx, "SomeUser", func(p *models.PermissionSettings) {
		p.WorkforcePolicy.Data.Create.IsAllowed = models.PermissionValueYes
	})
	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
	_, err := c.WorkOrder.Create().
		SetName("NewWorkOrder").
		SetType(workOrderType).
		SetCreationDate(time.Now()).
		SetOwner(u).
		SetAssignee(u).
		Save(ctx)
	require.NoError(t, err)
}

func TestWorkOrderUpdateAssignee(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, workOrder := prepareWorkOrderData(ctx, c)
	getCud := func(p *models.PermissionSettings) *models.WorkforceCud {
		return p.WorkforcePolicy.Data
	}
	u := viewer.MustGetOrCreateUser(ctx, "SomeUser", user.RoleUSER)
	appendAssign := func(p *models.PermissionSettings) {
		getCud(p).Assign.IsAllowed = models.PermissionValueYes
	}
	updateWorkOrderAssignee := func(ctx context.Context) error {
		_, err := c.WorkOrder.UpdateOne(workOrder).
			SetAssignee(u).
			Save(ctx)
		return err
	}
	clearWorkOrderAssignee := func(ctx context.Context) error {
		_, err := c.WorkOrder.UpdateOne(workOrder).
			ClearAssignee().
			Save(ctx)
		return err
	}
	tests := []policyTest{
		{
			operationName: "UpdateWithAssignee",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueYes
			},
			appendPermissions: appendAssign,
			operation:         updateWorkOrderAssignee,
		},
		{
			operationName: "ClearWorkOrderAssignee",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueYes
			},
			appendPermissions: appendAssign,
			operation:         clearWorkOrderAssignee,
		},
		{
			operationName: "UpdateWithAssigneeWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueByCondition
				getCud(p).Update.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				getCud(p).Assign.IsAllowed = models.PermissionValueByCondition
				getCud(p).Assign.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			operation: updateWorkOrderAssignee,
		},
		{
			operationName: "ClearWorkOrderAssigneeWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models.PermissionValueByCondition
				getCud(p).Update.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				getCud(p).Assign.IsAllowed = models.PermissionValueByCondition
				getCud(p).Assign.WorkOrderTypeIds = []int{workOrderType.ID}
			},
			operation: clearWorkOrderAssignee,
		},
	}
	runPolicyTest(t, tests)
}

func TestWorkOrderAssigneeUnchangedWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	_, workOrder := prepareWorkOrderData(ctx, c)
	_, workOrder2 := prepareWorkOrderData(ctx, c)
	u := viewer.MustGetOrCreateUser(ctx, "SomeUser", user.RoleUSER)
	c.WorkOrder.UpdateOne(workOrder).
		SetAssigneeID(u.ID).
		ExecX(ctx)
	permissions := authz.EmptyPermissions()
	permissions.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueYes
	ctx = viewertest.NewContext(
		context.Background(),
		c,
		viewertest.WithUser("user"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(permissions))
	err := c.WorkOrder.UpdateOne(workOrder).
		SetAssigneeID(u.ID).
		Exec(ctx)
	require.NoError(t, err)
	err = c.WorkOrder.UpdateOne(workOrder2).
		ClearAssignee().
		Exec(ctx)
	require.NoError(t, err)
	err = c.WorkOrder.UpdateOne(workOrder).
		ClearAssignee().
		Exec(ctx)
	require.Error(t, err)
}

func TestWorkOrderOwnerUnchangedWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	_, workOrder := prepareWorkOrderData(ctx, c)
	permissions := authz.EmptyPermissions()
	permissions.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueYes
	ctx = viewertest.NewContext(
		context.Background(),
		c,
		viewertest.WithUser("user"),
		viewertest.WithRole(user.RoleUSER),
		viewertest.WithPermissions(permissions))
	ownerID := workOrder.QueryOwner().OnlyXID(ctx)
	err := c.WorkOrder.UpdateOne(workOrder).
		SetOwnerID(ownerID).
		Exec(ctx)
	require.NoError(t, err)
}

func TestWorkorderTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workorderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)
	createWorkOrderType := func(ctx context.Context) error {
		_, err := c.WorkOrderType.Create().
			SetName("NewWorkOrderType").
			Save(ctx)
		return err
	}
	updateWorkOrderType := func(ctx context.Context) error {
		return c.WorkOrderType.UpdateOne(workorderType).
			SetName("NewName").
			Exec(ctx)
	}
	deleteWorkOrderType := func(ctx context.Context) error {
		return c.WorkOrderType.DeleteOne(workorderType).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.WorkforcePolicy.Templates
		},
		create: createWorkOrderType,
		update: updateWorkOrderType,
		delete: deleteWorkOrderType,
	})
}
