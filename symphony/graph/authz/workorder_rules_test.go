// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/ent/user"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func prepareWorkOrderData(ctx context.Context, c *ent.Client) (*ent.WorkOrderType, *ent.WorkOrder) {
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleOWNER)
	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)
	workOrder := c.WorkOrder.Create().
		SetName("WorkOrder").
		SetType(workOrderType).
		SetOwner(u).
		SetCreationDate(time.Now()).
		SaveX(ctx)
	return workOrderType, workOrder
}
func TestWorkOrderWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	workOrderType, workOrder := prepareWorkOrderData(ctx, c)
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
	updateWorkOrder := func(ctx context.Context) error {
		return c.WorkOrder.UpdateOne(workOrder).
			SetName("NewName").
			Exec(ctx)
	}
	deleteWorkOrder := func(ctx context.Context) error {
		return c.WorkOrder.DeleteOne(workOrder).
			Exec(ctx)
	}
	getCud := func(p *models.PermissionSettings) *models.WorkforceCud {
		return p.WorkforcePolicy.Data
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return &models.Cud{
				Create: getCud(p).Create,
				Update: getCud(p).Update,
				Delete: getCud(p).Delete,
			}
		},
		initialPermissions: func(p *models.PermissionSettings) {
			getCud(p).TransferOwnership.IsAllowed = models2.PermissionValueYes
		},
		create: createWorkOrder,
		update: updateWorkOrder,
		delete: deleteWorkOrder,
	})
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
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
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
	ctx := context.Background()
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
