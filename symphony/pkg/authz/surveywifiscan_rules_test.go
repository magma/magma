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

	models2 "github.com/facebookincubator/symphony/pkg/authz/models"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
)

func TestSurveyWiFiScanWritePolicyRule(t *testing.T) {
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

	surveyWiFiScan := c.SurveyWiFiScan.Create().
		SetChecklistItem(checkListItem).
		SetSsid("WiFi").
		SetBssid("bssid").
		SetFrequency(1).
		SetChannel(2).
		SetTimestamp(time.Now()).
		SetStrength(10).
		SaveX(ctx)

	createSurveyWiFiScan := func(ctx context.Context) error {
		_, err := c.SurveyWiFiScan.Create().
			SetChecklistItem(checkListItem).
			SetSsid("WiFi").
			SetBssid("bssid").
			SetFrequency(1).
			SetChannel(2).
			SetTimestamp(time.Now()).
			SetStrength(10).
			Save(ctx)
		return err
	}
	updateSurveyWiFiScan := func(ctx context.Context) error {
		return c.SurveyWiFiScan.UpdateOne(surveyWiFiScan).
			SetStrength(5).
			Exec(ctx)
	}
	deleteSurveyWiFiScan := func(ctx context.Context) error {
		return c.SurveyWiFiScan.DeleteOne(surveyWiFiScan).
			Exec(ctx)
	}
	runCudPolicyTest(t, cudPolicyTest{
		initialPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueYes
		},
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models2.PermissionValueByCondition
			p.WorkforcePolicy.Data.Update.WorkOrderTypeIds = []int{workOrderType.ID}
		},
		create: createSurveyWiFiScan,
		update: updateSurveyWiFiScan,
		delete: deleteSurveyWiFiScan,
	})
}

func TestSurveyWiFiScanReadPolicyRule(t *testing.T) {
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
	c.SurveyWiFiScan.Create().
		SetChecklistItem(checkListItem1).
		SetSsid("WiFi").
		SetBssid("bssid").
		SetFrequency(1).
		SetChannel(2).
		SetTimestamp(time.Now()).
		SetStrength(10).
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
	c.SurveyWiFiScan.Create().
		SetChecklistItem(checkListItem2).
		SetSsid("WiFi").
		SetBssid("bssid").
		SetFrequency(1).
		SetChannel(2).
		SetTimestamp(time.Now()).
		SetStrength(10).
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
		permissions.WorkforcePolicy.Read.IsAllowed = models2.PermissionValueByCondition
		permissions.WorkforcePolicy.Read.WorkOrderTypeIds = []int{woType1.ID}
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.SurveyWiFiScan.Query().Count(permissionsContext)
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
		count, err := c.SurveyWiFiScan.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}
