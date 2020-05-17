// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
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
		appendPermissions: func(p *models.PermissionSettings) {
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
		appendPermissions: func(p *models.PermissionSettings) {
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
		appendPermissions: func(p *models.PermissionSettings) {
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
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
		initialPermissions: func(p *models.PermissionSettings) {
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
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
		},
		initialPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}
