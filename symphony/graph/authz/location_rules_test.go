// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestLocationWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	locationType2 := c.LocationType.Create().
		SetName("LocationType2").
		SaveX(ctx)
	location := c.Location.Create().
		SetName("Location").
		SetType(locationType).
		SaveX(ctx)
	location2 := c.Location.Create().
		SetName("Location2").
		SetType(locationType2).
		SaveX(ctx)
	createLocation := func(ctx context.Context) error {
		_, err := c.Location.Create().
			SetName("NewLocation").
			SetType(locationType).
			Save(ctx)
		return err
	}
	createLocation2 := func(ctx context.Context) error {
		_, err := c.Location.Create().
			SetName("NewLocation").
			SetType(locationType2).
			Save(ctx)
		return err
	}
	updateLocation := func(ctx context.Context) error {
		return c.Location.UpdateOne(location).
			SetName("NewName").
			Exec(ctx)
	}
	updateLocation2 := func(ctx context.Context) error {
		return c.Location.UpdateOne(location2).
			SetName("NewName").
			Exec(ctx)
	}
	deleteLocation := func(ctx context.Context) error {
		return c.Location.DeleteOne(location).
			Exec(ctx)
	}
	deleteLocation2 := func(ctx context.Context) error {
		return c.Location.DeleteOne(location2).
			Exec(ctx)
	}
	tests := []policyTest{
		{
			operationName: "Create",
			appendPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Create.IsAllowed = models2.PermissionValueYes
			},
			operation: createLocation,
		},
		{
			operationName: "CreateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Create.IsAllowed = models2.PermissionValueByCondition
				p.InventoryPolicy.Location.Create.LocationTypeIds = []int{locationType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Create.LocationTypeIds = []int{locationType.ID, locationType2.ID}
			},
			operation: createLocation2,
		},
		{
			operationName: "Update",
			appendPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Update.IsAllowed = models2.PermissionValueYes
			},
			operation: updateLocation,
		},
		{
			operationName: "UpdateWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Update.IsAllowed = models2.PermissionValueByCondition
				p.InventoryPolicy.Location.Update.LocationTypeIds = []int{locationType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Update.LocationTypeIds = []int{locationType.ID, locationType2.ID}
			},
			operation: updateLocation2,
		},
		{
			operationName: "Delete",
			appendPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Delete.IsAllowed = models2.PermissionValueYes
			},
			operation: deleteLocation,
		},
		{
			operationName: "DeleteWithType",
			initialPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Delete.IsAllowed = models2.PermissionValueByCondition
				p.InventoryPolicy.Location.Delete.LocationTypeIds = []int{locationType.ID}
			},
			appendPermissions: func(p *models.PermissionSettings) {
				p.InventoryPolicy.Location.Delete.LocationTypeIds = []int{locationType.ID, locationType2.ID}
			},
			operation: deleteLocation2,
		},
	}
	runPolicyTest(t, tests)
}

func TestLocationTypeWritePolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	createLocationType := func(ctx context.Context) error {
		_, err := c.LocationType.Create().
			SetName("NewLocationType").
			Save(ctx)
		return err
	}
	updateLocationType := func(ctx context.Context) error {
		return c.LocationType.UpdateOne(locationType).
			SetName("NewName").
			Exec(ctx)
	}
	deleteLocationType := func(ctx context.Context) error {
		return c.LocationType.DeleteOne(locationType).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		getCud: func(p *models.PermissionSettings) *models.Cud {
			return p.InventoryPolicy.LocationType
		},
		create: createLocationType,
		update: updateLocationType,
		delete: deleteLocationType,
	})
}
