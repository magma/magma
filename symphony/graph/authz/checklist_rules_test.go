// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/stretchr/testify/require"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestCheckListCategoryDefinitionWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)

	clcDef := c.CheckListCategoryDefinition.Create().
		SetTitle("CategoryDefinition").
		SetWorkOrderTypeID(workOrderType.ID).
		SaveX(ctx)

	createItem := func(ctx context.Context) error {
		_, err := c.CheckListCategoryDefinition.Create().
			SetTitle("CategoryDefinition").
			SetWorkOrderTypeID(workOrderType.ID).
			Save(ctx)
		return err
	}
	updateItem := func(ctx context.Context) error {
		return c.CheckListCategoryDefinition.UpdateOne(clcDef).
			SetTitle("NewTitle").
			Exec(ctx)
	}
	deleteItem := func(ctx context.Context) error {
		return c.CheckListCategoryDefinition.DeleteOne(clcDef).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Templates.Update.IsAllowed = models2.PermissionValueYes
		},
		create: createItem,
		update: updateItem,
		delete: deleteItem,
	})
}

func TestCheckListCategoryWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)

	workOrder := c.WorkOrder.Create().
		SetName("WorkOrder").
		SetTypeID(workOrderType.ID).
		SetCreationDate(time.Now()).
		SetOwner(viewer.FromContext(ctx).(*viewer.UserViewer).User()).
		SaveX(ctx)

	clc := c.CheckListCategory.Create().
		SetTitle("Category").
		SetWorkOrderID(workOrder.ID).
		SaveX(ctx)

	createItem := func(ctx context.Context) error {
		_, err := c.CheckListCategory.Create().
			SetTitle("Item").
			SetWorkOrderID(workOrder.ID).
			Save(ctx)
		return err
	}
	updateItem := func(ctx context.Context) error {
		return c.CheckListCategory.UpdateOne(clc).
			SetTitle("NewTitle").
			Exec(ctx)
	}
	deleteItem := func(ctx context.Context) error {
		return c.CheckListCategory.DeleteOne(clc).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueByCondition
			p.WorkforcePolicy.Data.Update.WorkOrderTypeIds = []int{workOrderType.ID}
		},
		create: createItem,
		update: updateItem,
		delete: deleteItem,
	})
}

func TestCheckListItemDefinitionWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)

	clcDef := c.CheckListCategoryDefinition.Create().
		SetTitle("CategoryDefinition").
		SetWorkOrderTypeID(workOrderType.ID).
		SaveX(ctx)

	cliDef := c.CheckListItemDefinition.Create().
		SetTitle("ItemDefinition").
		SetType("simple").
		SetCheckListCategoryDefinitionID(clcDef.ID).
		SaveX(ctx)

	createItem := func(ctx context.Context) error {
		_, err := c.CheckListItemDefinition.Create().
			SetTitle("ItemDefinition").
			SetType("simple").
			SetCheckListCategoryDefinitionID(clcDef.ID).
			Save(ctx)
		return err
	}
	updateItem := func(ctx context.Context) error {
		return c.CheckListItemDefinition.UpdateOne(cliDef).
			SetTitle("NewTitle").
			Exec(ctx)
	}
	deleteItem := func(ctx context.Context) error {
		return c.CheckListItemDefinition.DeleteOne(cliDef).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Templates.Update.IsAllowed = models2.PermissionValueYes
		},
		create: createItem,
		update: updateItem,
		delete: deleteItem,
	})
}

func TestCheckListItemWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)

	workOrder := c.WorkOrder.Create().
		SetName("WorkOrder").
		SetTypeID(workOrderType.ID).
		SetCreationDate(time.Now()).
		SetOwner(viewer.FromContext(ctx).(*viewer.UserViewer).User()).
		SaveX(ctx)

	clc := c.CheckListCategory.Create().
		SetTitle("Category").
		SetWorkOrderID(workOrder.ID).
		SaveX(ctx)

	clItem := c.CheckListItem.Create().
		SetTitle("Item").
		SetCheckListCategoryID(clc.ID).
		SetType("simple").
		SaveX(ctx)

	createItem := func(ctx context.Context) error {
		_, err := c.CheckListItem.Create().
			SetTitle("Item").
			SetCheckListCategoryID(clc.ID).
			SetType("simple").
			Save(ctx)
		return err
	}
	updateItem := func(ctx context.Context) error {
		return c.CheckListItem.UpdateOne(clItem).
			SetTitle("NewTitle").
			Exec(ctx)
	}
	deleteItem := func(ctx context.Context) error {
		return c.CheckListItem.DeleteOne(clItem).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueByCondition
			p.WorkforcePolicy.Data.Update.WorkOrderTypeIds = []int{workOrderType.ID}
		},
		create: createItem,
		update: updateItem,
		delete: deleteItem,
	})
}

func TestCheckListCategoryReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	_, wo2 := prepareWorkOrderData(ctx, c)
	c.CheckListCategory.Create().
		SetTitle("Category1").
		SetWorkOrder(wo1).
		SaveX(ctx)
	c.CheckListCategory.Create().
		SetTitle("Category2").
		SetWorkOrder(wo2).
		SaveX(ctx)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.CheckListCategory.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Zero(t, count)
	})
	t.Run("PartialPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueByCondition
		permissions.WorkforcePolicy.Read.WorkOrderTypeIds = []int{woType1.ID}
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.CheckListCategory.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
	t.Run("FullPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.CheckListCategory.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

func TestCheckListItemReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	_, wo2 := prepareWorkOrderData(ctx, c)
	caetgory1 := c.CheckListCategory.Create().
		SetTitle("Category1").
		SetWorkOrder(wo1).
		SaveX(ctx)
	c.CheckListItem.Create().
		SetTitle("Item1").
		SetCheckListCategory(caetgory1).
		SetType("simple").
		SaveX(ctx)
	caetgory2 := c.CheckListCategory.Create().
		SetTitle("Category2").
		SetWorkOrder(wo2).
		SaveX(ctx)
	c.CheckListItem.Create().
		SetTitle("Item1").
		SetCheckListCategory(caetgory2).
		SetType("simple").
		SaveX(ctx)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.CheckListItem.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Zero(t, count)
	})
	t.Run("PartialPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueByCondition
		permissions.WorkforcePolicy.Read.WorkOrderTypeIds = []int{woType1.ID}
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.CheckListItem.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
	t.Run("FullPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.CheckListItem.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}
