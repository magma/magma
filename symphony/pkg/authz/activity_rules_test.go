// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

//TODO(T67933416): Return these rules
/*
func getActivityCudOperations(
	ctx context.Context,
	c *ent.Client,
	setParent func(*ent.ActivityCreate) *ent.ActivityCreate,
) cudOperations {
	author := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleOWNER)
	activityQuery := c.Activity.Create().
		SetAuthor(author).
		SetChangedField(activity.ChangedFieldASSIGNEE).
		SetNewValue("a").
		SetOldValue("b")
	activityQuery = setParent(activityQuery)
	activityEntity := activityQuery.SaveX(ctx)
	createActivity := func(ctx context.Context) error {
		activityQuery := c.Activity.Create().
			SetChangedField(activity.ChangedFieldASSIGNEE).
			SetNewValue("a").
			SetOldValue("b").
			SetAuthor(author)
		activityQuery = setParent(activityQuery)
		_, err := activityQuery.Save(ctx)
		return err
	}
	updateActivity := func(ctx context.Context) error {
		return c.Activity.UpdateOne(activityEntity).
			SetNewValue("c").
			Exec(ctx)
	}
	deleteActivity := func(ctx context.Context) error {
		return c.Activity.DeleteOne(activityEntity).
			Exec(ctx)
	}
	return cudOperations{
		create: createActivity,
		update: updateActivity,
		delete: deleteActivity,
	}
}

func TestActivityOfWorkOrderReadPolicyRule(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	u := viewer.MustGetOrCreateUser(ctx, "AuthID", user.RoleUSER)
	woType1, wo1 := prepareWorkOrderData(ctx, c)
	_, wo2 := prepareWorkOrderData(ctx, c)
	c.Activity.Create().
		SetAuthor(u).
		SetWorkOrder(wo1).
		SetChangedField(activity.ChangedFieldASSIGNEE).
		SetNewValue("a").
		SetOldValue("b").
		SaveX(ctx)
	c.Activity.Create().
		SetAuthor(u).
		SetWorkOrder(wo2).
		SetChangedField(activity.ChangedFieldASSIGNEE).
		SetNewValue("a").
		SetOldValue("b").
		SaveX(ctx)
	t.Run("EmptyPermissions", func(t *testing.T) {
		permissions := authz.EmptyPermissions()
		permissionsContext := viewertest.NewContext(
			context.Background(),
			c,
			viewertest.WithUser("user"),
			viewertest.WithRole(user.RoleUSER),
			viewertest.WithPermissions(permissions))
		count, err := c.Activity.Query().Count(permissionsContext)
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
		count, err := c.Activity.Query().Count(permissionsContext)
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
		count, err := c.Activity.Query().Count(permissionsContext)
		require.NoError(t, err)
		require.Equal(t, 2, count)
	})
}

func TestWorkOrderActivityPolicyRule(t *testing.T) {
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

	cudOperations := getActivityCudOperations(ctx, c, func(ptc *ent.ActivityCreate) *ent.ActivityCreate {
		return ptc.SetWorkOrder(workOrder)
	})
	runCudPolicyTest(t, cudPolicyTest{
		initialPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Read.IsAllowed = models.PermissionValueYes
		},
		appendPermissions: func(p *models.PermissionSettings) {
			p.WorkforcePolicy.Data.Update.IsAllowed = models.PermissionValueYes
		},
		create: cudOperations.create,
		update: cudOperations.update,
		delete: cudOperations.delete,
	})
}
*/
