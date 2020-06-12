// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"strconv"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/activity"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/suite"
)

type workOrderActivitiesTestSuite struct {
	eventTestSuite
	typ *ent.WorkOrderType
}

func TestWorkOrderActivitiesEvents(t *testing.T) {
	suite.Run(t, &workOrderActivitiesTestSuite{})
}

func (s *workOrderActivitiesTestSuite) SetupSuite() {
	s.eventTestSuite.SetupSuite(viewertest.WithFeatures(viewer.FeatureWorkOrderActivitiesHook))
	s.typ = s.client.WorkOrderType.Create().
		SetName("Chore").
		SaveX(s.ctx)
}

func (s *workOrderActivitiesTestSuite) TestAddWorkOrderActivities() {
	t := time.Now()
	wo := s.client.WorkOrder.Create().
		SetName("wo1").
		SetType(s.typ).
		SetCreationDate(t).
		SetAssignee(s.user).
		SetOwner(s.user).
		SaveX(s.ctx)
	s.Require().Equal(wo.Name, "wo1")
	activities := wo.QueryActivities().AllX(s.ctx)
	s.Require().Len(activities, 5)
	for _, a := range activities {
		s.Require().Equal(a.QueryAuthor().OnlyX(s.ctx).AuthID, s.user.AuthID)
		s.Require().Equal(a.QueryWorkOrder().OnlyX(s.ctx).ID, wo.ID)
		switch a.ChangedField {
		case activity.ChangedFieldCREATIONDATE:
			s.Require().Empty(a.OldValue)
			s.Require().Equal(a.NewValue, strconv.FormatInt(t.Unix(), 10))
			s.Require().True(a.IsCreate)
		case activity.ChangedFieldOWNER:
			fallthrough
		case activity.ChangedFieldASSIGNEE:
			s.Require().Empty(a.OldValue)
			s.Require().Equal(a.NewValue, strconv.Itoa(s.user.ID))
			s.Require().True(a.IsCreate)
		case activity.ChangedFieldSTATUS:
			s.Require().Empty(a.OldValue)
			s.Require().Equal(a.NewValue, models.WorkOrderStatusPlanned.String())
			s.Require().True(a.IsCreate)
		case activity.ChangedFieldPRIORITY:
			s.Require().Empty(a.OldValue)
			s.Require().Equal(a.NewValue, models.WorkOrderPriorityNone.String())
			s.Require().True(a.IsCreate)
		default:
			s.Require().Fail("unsupported changed field")
		}
	}
}

func (s *workOrderActivitiesTestSuite) TestEditWorkOrderActivities() {
	t := time.Now()

	wo := s.client.WorkOrder.Create().
		SetName("wo2").
		SetType(s.typ).
		SetCreationDate(t).
		SetAssignee(s.user).
		SetOwner(s.user).
		SaveX(s.ctx)
	s.Require().Equal(wo.Name, "wo2")
	activities := wo.QueryActivities().AllX(s.ctx)
	s.Require().Len(activities, 5)
	u2 := s.client.User.Create().
		SetAuthID("123").
		SetRole(user.RoleUSER).SaveX(s.ctx)
	s.client.WorkOrder.UpdateOne(wo).
		SetAssignee(u2).
		SetStatus(models.WorkOrderStatusPending.String()).ExecX(s.ctx)

	activities = wo.QueryActivities().AllX(s.ctx)
	s.Require().Len(activities, 7)
	newCount := 0
	for _, a := range activities {
		s.Require().Equal(a.QueryAuthor().OnlyX(s.ctx).AuthID, s.user.AuthID)
		s.Require().Equal(a.QueryWorkOrder().OnlyX(s.ctx).ID, wo.ID)
		if a.OldValue == "" {
			continue
		}
		newCount++
		switch a.ChangedField {
		case activity.ChangedFieldASSIGNEE:
			s.Require().Equal(a.NewValue, strconv.Itoa(u2.ID))
			s.Require().False(a.IsCreate)
			s.Require().Equal(a.OldValue, strconv.Itoa(s.user.ID))
		case activity.ChangedFieldSTATUS:
			s.Require().Equal(a.NewValue, models.WorkOrderStatusPending.String())
			s.Require().False(a.IsCreate)
			s.Require().Equal(a.OldValue, models.WorkOrderStatusPlanned.String())
		default:
			s.Require().Fail("unsupported changed field")
		}
	}
	s.Require().Equal(2, newCount)
}
