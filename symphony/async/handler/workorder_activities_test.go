// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent/activity"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/event"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestAddWorkOrderActivities(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	c.Use(event.LogHook(event.HandleActivityLog, log.NewNopLogger()))
	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()

	now := time.Now()
	typ := c.WorkOrderType.Create().
		SetName("Chore").
		SaveX(ctx)
	wo := c.WorkOrder.Create().
		SetName("wo1").
		SetType(typ).
		SetCreationDate(now).
		SetAssignee(u).
		SetOwner(u).
		SaveX(ctx)
	require.Equal(t, wo.Name, "wo1")
	activities := wo.QueryActivities().AllX(ctx)
	require.Len(t, activities, 5)
	for _, a := range activities {
		require.Equal(t, a.QueryAuthor().OnlyX(ctx).AuthID, u.AuthID)
		require.Equal(t, a.QueryWorkOrder().OnlyX(ctx).ID, wo.ID)
		switch a.ChangedField {
		case activity.ChangedFieldCREATIONDATE:
			require.Empty(t, a.OldValue)
			require.Equal(t, a.NewValue, strconv.FormatInt(now.Unix(), 10))
			require.True(t, a.IsCreate)
		case activity.ChangedFieldOWNER, activity.ChangedFieldASSIGNEE:
			require.Empty(t, a.OldValue)
			require.Equal(t, a.NewValue, strconv.Itoa(u.ID))
			require.True(t, a.IsCreate)
		case activity.ChangedFieldSTATUS:
			require.Empty(t, a.OldValue)
			require.Equal(t, a.NewValue, models.WorkOrderStatusPlanned.String())
			require.True(t, a.IsCreate)
		case activity.ChangedFieldPRIORITY:
			require.Empty(t, a.OldValue)
			require.Equal(t, a.NewValue, models.WorkOrderPriorityNone.String())
			require.True(t, a.IsCreate)
		default:
			require.Fail(t, "unsupported changed field")
		}
	}
}

func TestEditWorkOrderActivities(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c)
	c.Use(event.LogHook(event.HandleActivityLog, log.NewNopLogger()))
	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()

	now := time.Now()
	typ := c.WorkOrderType.Create().
		SetName("Chore").
		SaveX(ctx)
	wo := c.WorkOrder.Create().
		SetName("wo2").
		SetType(typ).
		SetCreationDate(now).
		SetAssignee(u).
		SetOwner(u).
		SaveX(ctx)
	require.Equal(t, wo.Name, "wo2")
	activities := wo.QueryActivities().AllX(ctx)
	require.Len(t, activities, 5)
	u2 := c.User.Create().
		SetAuthID("123").
		SetRole(user.RoleUSER).SaveX(ctx)
	c.WorkOrder.UpdateOne(wo).
		SetAssignee(u2).
		SetStatus(models.WorkOrderStatusPending.String()).ExecX(ctx)

	activities = wo.QueryActivities().AllX(ctx)
	require.Len(t, activities, 7)
	newCount := 0
	for _, a := range activities {
		require.Equal(t, a.QueryAuthor().OnlyX(ctx).AuthID, u.AuthID)
		require.Equal(t, a.QueryWorkOrder().OnlyX(ctx).ID, wo.ID)
		if a.OldValue == "" {
			continue
		}
		newCount++
		switch a.ChangedField {
		case activity.ChangedFieldASSIGNEE:
			require.Equal(t, a.NewValue, strconv.Itoa(u2.ID))
			require.False(t, a.IsCreate)
			require.Equal(t, a.OldValue, strconv.Itoa(u.ID))
		case activity.ChangedFieldSTATUS:
			require.Equal(t, a.NewValue, models.WorkOrderStatusPending.String())
			require.False(t, a.IsCreate)
			require.Equal(t, a.OldValue, models.WorkOrderStatusPlanned.String())
		default:
			require.Fail(t, "unsupported changed field")
		}
	}
	require.Equal(t, 2, newCount)
}
