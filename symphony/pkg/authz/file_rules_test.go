// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/authz"
	models2 "github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func getFileCudOperations(ctx context.Context, c *ent.Client, setParent func(*ent.FileCreate) *ent.FileCreate) cudOperations {
	fileQuery := c.File.Create().
		SetName("name").
		SetType(models.FileTypeImage.String()).
		SetStoreKey("abc").
		SetContentType("text/html")

	fileQuery = setParent(fileQuery)
	file := fileQuery.SaveX(ctx)

	createFile := func(ctx context.Context) error {
		fileQuery := c.File.Create().
			SetName("name2").
			SetType(models.FileTypeImage.String()).
			SetStoreKey("abcd").
			SetContentType("text/html")
		fileQuery = setParent(fileQuery)
		_, err := fileQuery.Save(ctx)

		return err
	}
	updateFile := func(ctx context.Context) error {
		return c.File.UpdateOne(file).
			SetName("newName").
			Exec(ctx)
	}
	deleteFile := func(ctx context.Context) error {
		return c.File.DeleteOne(file).
			Exec(ctx)
	}
	return cudOperations{
		create: createFile,
		update: updateFile,
		delete: deleteFile,
	}
}

func TestLocationFilePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	location := c.Location.Create().
		SetName("Location").
		SetType(locationType).
		SaveX(ctx)

	cudOperations := getFileCudOperations(ctx, c, func(fc *ent.FileCreate) *ent.FileCreate {
		return fc.SetLocation(location)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models2.PermissionSettings) {
			p.InventoryPolicy.Location.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestEquipmentFilePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	location := c.Location.Create().
		SetName("Location").
		SetType(locationType).
		SaveX(ctx)
	equipType := c.EquipmentType.Create().
		SetName("EquipmentType").
		SaveX(ctx)
	equip := c.Equipment.Create().
		SetName("Equip").
		SetLocation(location).
		SetType(equipType).
		SaveX(ctx)

	cudOperations := getFileCudOperations(ctx, c, func(fc *ent.FileCreate) *ent.FileCreate {
		return fc.SetEquipment(equip)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models2.PermissionSettings) {
			p.InventoryPolicy.Equipment.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestUserFilePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	user := c.User.Create().
		SetAuthID("authID").
		SaveX(ctx)

	file := c.File.Create().
		SetName("name1").
		SetType(models.FileTypeImage.String()).
		SetStoreKey("abc").
		SetUser(user).
		SetContentType("text/html").
		SaveX(ctx)

	user2 := c.User.Create().
		SetAuthID("authID2").
		SaveX(ctx)
	createFile := func(ctx context.Context) error {
		_, err := c.File.Create().
			SetName("name2").
			SetType(models.FileTypeImage.String()).
			SetStoreKey("abcd").
			SetContentType("text/html").
			SetUser(user2).
			Save(ctx)
		return err
	}
	updateFile := func(ctx context.Context) error {
		return c.File.UpdateOne(file).
			SetName("newName").
			Exec(ctx)
	}
	deleteFile := func(ctx context.Context) error {
		return c.File.DeleteOne(file).
			Exec(ctx)
	}

	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models2.PermissionSettings) {
			p.AdminPolicy.Access.IsAllowed = models2.PermissionValueYes
		},
		create: createFile,
		update: updateFile,
		delete: deleteFile,
	})
}

func TestWOFilePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	_, workOrder := prepareWorkOrderData(ctx, c)
	cudOperations := getFileCudOperations(ctx, c, func(fc *ent.FileCreate) *ent.FileCreate {
		return fc.SetWorkOrder(workOrder)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models2.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
		initialPermissions: func(p *models2.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		},
	})
}

func TestChecklistItemFilePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)

	_, workOrder := prepareWorkOrderData(ctx, c)
	clc := c.CheckListCategory.Create().
		SetTitle("Category").
		SetWorkOrderID(workOrder.ID).
		SaveX(ctx)

	clItem := c.CheckListItem.Create().
		SetTitle("Item").
		SetCheckListCategoryID(clc.ID).
		SetType("simple").
		SaveX(ctx)

	cudOperations := getFileCudOperations(ctx, c, func(fc *ent.FileCreate) *ent.FileCreate {
		return fc.SetChecklistItem(clItem)
	})
	runCudPolicyTest(t, cudPolicyTest{
		appendPermissions: func(p *models2.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
		},
		initialPermissions: func(p *models2.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}

func TestFileOfWorkOrderReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	_, wo2 := prepareWorkOrderData(ctx, c)
	c.File.Create().
		SetType("image/png").
		SetName("image1.png").
		SetContentType("image/png").
		SetStoreKey("1111").
		SetWorkOrder(wo1).
		SaveX(ctx)
	c.File.Create().
		SetType("image/png").
		SetName("image2.png").
		SetContentType("image/png").
		SetStoreKey("2222").
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
		count, err := c.File.Query().Count(permissionsContext)
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
		count, err := c.File.Query().Count(permissionsContext)
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
		count, err := c.File.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

func TestFileOfCheckListItemReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	_, wo2 := prepareWorkOrderData(ctx, c)
	caetgory1 := c.CheckListCategory.Create().
		SetTitle("Category1").
		SetWorkOrder(wo1).
		SaveX(ctx)
	item1 := c.CheckListItem.Create().
		SetTitle("Item1").
		SetCheckListCategory(caetgory1).
		SetType("simple").
		SaveX(ctx)
	c.File.Create().
		SetType("image/png").
		SetName("image1.png").
		SetContentType("image/png").
		SetStoreKey("1111").
		SetChecklistItem(item1).
		SaveX(ctx)
	caetgory2 := c.CheckListCategory.Create().
		SetTitle("Category2").
		SetWorkOrder(wo2).
		SaveX(ctx)
	item2 := c.CheckListItem.Create().
		SetTitle("Item1").
		SetCheckListCategory(caetgory2).
		SetType("simple").
		SaveX(ctx)
	c.File.Create().
		SetType("image/png").
		SetName("image2.png").
		SetContentType("image/png").
		SetStoreKey("2222").
		SetChecklistItem(item2).
		SaveX(ctx)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.File.Query().Count(permissionsContext)
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
		count, err := c.File.Query().Count(permissionsContext)
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
		count, err := c.File.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}
