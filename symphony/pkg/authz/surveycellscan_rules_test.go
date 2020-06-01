// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func TestSurveyCellScanWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "anotherOne", user.RoleUSER)
	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)

	workOrder := c.WorkOrder.Create().
		SetName("WorkOrder").
		SetTypeID(workOrderType.ID).
		SetCreationDate(time.Now()).
		SetOwner(u).
		SaveX(ctx)

	clc := c.CheckListCategory.Create().
		SetTitle("Category").
		SetWorkOrderID(workOrder.ID).
		SaveX(ctx)

	checkListItem := c.CheckListItem.Create().
		SetTitle("Item").
		SetCheckListCategoryID(clc.ID).
		SetType("simple").
		SaveX(ctx)

	surveyCellScan := c.SurveyCellScan.Create().
		SetChecklistItem(checkListItem).
		SetNetworkType("5G").
		SetSignalStrength(10).
		SaveX(ctx)

	createSurveyCellScan := func(ctx context.Context) error {
		_, err := c.SurveyCellScan.Create().
			SetChecklistItem(checkListItem).
			SetNetworkType("5G").
			SetSignalStrength(10).
			Save(ctx)
		return err
	}
	updateSurveyCellScan := func(ctx context.Context) error {
		return c.SurveyCellScan.UpdateOne(surveyCellScan).
			SetSignalStrength(5).
			Exec(ctx)
	}
	deleteSurveyCellScan := func(ctx context.Context) error {
		return c.SurveyCellScan.DeleteOne(surveyCellScan).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		initialPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models.PermissionValueYes
		},
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueByCondition
			p.WorkforcePolicy.Data.Update.WorkOrderTypeIds = []int{workOrderType.ID}
		},
		create: createSurveyCellScan,
		update: updateSurveyCellScan,
		delete: deleteSurveyCellScan,
	})
}

func TestSurveyCellScanReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	_, wo2 := prepareWorkOrderData(ctx, c)

	caetgory1 := c.CheckListCategory.Create().
		SetTitle("Category1").
		SetWorkOrder(wo1).
		SaveX(ctx)
	checkListItem1 := c.CheckListItem.Create().
		SetTitle("Item1").
		SetCheckListCategory(caetgory1).
		SetType("simple").
		SaveX(ctx)
	c.SurveyCellScan.Create().
		SetChecklistItem(checkListItem1).
		SetNetworkType("5G").
		SetSignalStrength(10).
		SaveX(ctx)

	caetgory2 := c.CheckListCategory.Create().
		SetTitle("Category2").
		SetWorkOrder(wo2).
		SaveX(ctx)
	checkListItem2 := c.CheckListItem.Create().
		SetTitle("Item1").
		SetCheckListCategory(caetgory2).
		SetType("simple").
		SaveX(ctx)
	c.SurveyCellScan.Create().
		SetChecklistItem(checkListItem2).
		SetNetworkType("5G").
		SetSignalStrength(10).
		SaveX(ctx)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.SurveyCellScan.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Zero(t, count)
	})
	t.Run("PartialPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models.PermissionValueByCondition
		permissions.WorkforcePolicy.Read.WorkOrderTypeIds = []int{woType1.ID}
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.SurveyCellScan.Query().Count(permissionsContext)
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
		count, err := c.SurveyCellScan.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}
