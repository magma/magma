// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	models2 "github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
)

func getHyperlinkCudOperations(
	ctx context.Context,
	c *ent.Client,
	setParent func(*ent.HyperlinkCreate) *ent.HyperlinkCreate,
) cudOperations {
	hyperlinkQuery := c.Hyperlink.Create().
		SetName("BaseHyperlink").
		SetURL("BaseHyperLinkURL")
	hyperlinkQuery = setParent(hyperlinkQuery)
	hyperlink := hyperlinkQuery.SaveX(ctx)
	createHyperlink := func(ctx context.Context) error {
		hyperlinkQuery := c.Hyperlink.Create().
			SetName("Hyperlink").
			SetURL("HyperLinkURL")
		hyperlinkQuery = setParent(hyperlinkQuery)
		_, err := hyperlinkQuery.Save(ctx)
		return err
	}
	updateHyperlink := func(ctx context.Context) error {
		return c.Hyperlink.UpdateOne(hyperlink).
			SetName("updatedHyperlink").
			Exec(ctx)
	}
	deleteHyperlink := func(ctx context.Context) error {
		return c.Hyperlink.DeleteOne(hyperlink).
			Exec(ctx)
	}
	return cudOperations{
		create: createHyperlink,
		update: updateHyperlink,
		delete: deleteHyperlink,
	}
}

func TestHyperlinkReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	_, wo2 := prepareWorkOrderData(ctx, c)
	c.Hyperlink.Create().
		SetURL("http://url_1").
		SetWorkOrder(wo1).
		SaveX(ctx)
	c.Hyperlink.Create().
		SetURL("http://url_2").
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
		count, err := c.Hyperlink.Query().Count(permissionsContext)
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
		count, err := c.Hyperlink.Query().Count(permissionsContext)
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
		count, err := c.Hyperlink.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

func TestEquipmentHyperlinkPolicyRule(t *testing.T) {
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
	equipment := c.Equipment.Create().
		SetName("Equipment").
		SetLocation(location).
		SetType(equipType).
		SaveX(ctx)

	cudOperations := getHyperlinkCudOperations(ctx, c, func(fc *ent.HyperlinkCreate) *ent.HyperlinkCreate {
		return fc.SetEquipment(equipment)
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

func TestLocationHyperlinkPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	locationType := c.LocationType.Create().
		SetName("LocationType").
		SaveX(ctx)
	location := c.Location.Create().
		SetName("Location").
		SetType(locationType).
		SaveX(ctx)

	cudOperations := getHyperlinkCudOperations(ctx, c, func(fc *ent.HyperlinkCreate) *ent.HyperlinkCreate {
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

func TestWorkOrderHyperlinkPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleOWNER)
	workOrderType := c.WorkOrderType.Create().
		SetName("WorkOrderType").
		SaveX(ctx)
	workOrder := c.WorkOrder.Create().
		SetName("workOrder").
		SetType(workOrderType).
		SetOwner(u).
		SetCreationDate(time.Now()).
		SaveX(ctx)

	cudOperations := getHyperlinkCudOperations(ctx, c, func(ptc *ent.HyperlinkCreate) *ent.HyperlinkCreate {
		return ptc.SetWorkOrder(workOrder)
	})
	runCudPolicyTest(t, cudPolicyTest{
		initialPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		},
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}
