// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/authz"
	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
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

	ctx = viewertest.NewContext(ctx, c, viewertest.WithUser("MyAssignee"), viewertest.WithPermissions(authz.EmptyPermissions()))
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

	ctx = viewertest.NewContext(ctx, c, viewertest.WithUser("MyOwner"), viewertest.WithPermissions(authz.EmptyPermissions()))
	err := c.WorkOrder.UpdateOne(workOrder).
		SetName("NewName").
		Exec(ctx)
	require.NoError(t, err)
	err = c.WorkOrder.UpdateOne(workOrder).
		SetOwner(u2).
		Exec(ctx)
	require.NoError(t, err)
	ctx = viewertest.NewContext(ctx, c, viewertest.WithUser("NewOwner"), viewertest.WithPermissions(authz.EmptyPermissions()))
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
			SetName("NewWorkOrder").
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
			SetName("NewName").
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
		p.WorkforcePolicy.Data.TransferOwnership.IsAllowed = models2.PermissionValueYes
	}
	tests := []policyTest{
		{
			operationName:      "Create",
			initialPermissions: initialPermissions,
			appendPermissions: func(p *models.PermissionSettings) {
				p.WorkforcePolicy.Data.Create.IsAllowed = models2.PermissionValueYes
			},
			operation: createWorkOrder,
		},
		{
			operationName: "CreateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				initialPermissions(p)
				p.WorkforcePolicy.Data.Create.IsAllowed = models2.PermissionValueByCondition
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
				p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
			},
			operation: updateWorkOrder,
		},
		{
			operationName: "UpdateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				initialPermissions(p)
				p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueByCondition
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
				p.WorkforcePolicy.Data.Delete.IsAllowed = models2.PermissionValueYes
			},
			operation: deleteWorkOrder,
		},
		{
			operationName: "DeleteWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				initialPermissions(p)
				p.WorkforcePolicy.Data.Delete.IsAllowed = models2.PermissionValueByCondition
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

func TestWorkOrderTransferOwnershipWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, workOrder := prepareWorkOrderData(ctx, c)
	getCud := func(p *models.PermissionSettings) *models.WorkforceCud {
		return p.WorkforcePolicy.Data
	}
	appendTransferOwnership := func(p *models.PermissionSettings) {
		getCud(p).TransferOwnership.IsAllowed = models2.PermissionValueYes
	}
	createWorkOrderWithOwner := func(ctx context.Context) error {
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
		_, err := c.WorkOrder.Create().
			SetName("NewWorkOrder").
			SetType(workOrderType).
			SetCreationDate(time.Now()).
			SetOwner(u).
			Save(ctx)
		return err
	}
	updateWorkOrderOwner := func(ctx context.Context) error {
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
		_, err := c.WorkOrder.UpdateOne(workOrder).
			SetOwner(u).
			Save(ctx)
		return err
	}
	tests := []policyTest{
		{
			operationName: "CreateWithOwner",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Create.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: appendTransferOwnership,
			operation:         createWorkOrderWithOwner,
		},
		{
			operationName: "UpdateWithOwner",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: appendTransferOwnership,
			operation:         updateWorkOrderOwner,
		},
	}
	runPolicyTest(t, tests)
}

func TestWorkOrderAssignWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, workOrder := prepareWorkOrderData(ctx, c)
	u2 := viewer.MustGetOrCreateUser(ctx, "SomeUser", user.RoleUSER)
	getCud := func(p *models.PermissionSettings) *models.WorkforceCud {
		return p.WorkforcePolicy.Data
	}
	appendAssign := func(p *models.PermissionSettings) {
		getCud(p).Assign.IsAllowed = models2.PermissionValueYes
	}
	createWorkOrderWithAssignee := func(ctx context.Context) error {
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
		_, err := c.WorkOrder.Create().
			SetName("NewWorkOrder").
			SetType(workOrderType).
			SetCreationDate(time.Now()).
			SetOwner(u).
			SetAssignee(u).
			Save(ctx)
		return err
	}
	updateWorkOrderAssignee := func(ctx context.Context) error {
		_, err := c.WorkOrder.UpdateOne(workOrder).
			SetAssignee(u2).
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
			operationName: "CreateWithAssignee",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Create.IsAllowed = models2.PermissionValueYes
				getCud(p).TransferOwnership.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: appendAssign,
			operation:         createWorkOrderWithAssignee,
		},
		{
			operationName: "UpdateWithAssignee",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: appendAssign,
			operation:         updateWorkOrderAssignee,
		},
		{
			operationName: "ClearWorkOrderAssignee",
			initialPermissions: func(p *models.PermissionSettings) {
				getCud(p).Update.IsAllowed = models2.PermissionValueYes
			},
			appendPermissions: appendAssign,
			operation:         clearWorkOrderAssignee,
		},
	}
	runPolicyTest(t, tests)
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
