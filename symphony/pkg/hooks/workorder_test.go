// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hooks_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type workOrderTestSuite struct {
	suite.Suite
	ctx    context.Context
	client *ent.WorkOrderClient
	user   *ent.User
	typ    *ent.WorkOrderType
}

func (s *workOrderTestSuite) SetupSuite() {
	client := viewertest.NewTestClient(s.T())
	s.ctx = viewertest.NewContext(
		context.Background(),
		client,
	)
	s.client = client.WorkOrder
	u, ok := viewer.FromContext(s.ctx).(*viewer.UserViewer)
	s.Require().True(ok)
	s.user = u.User()
	var err error
	s.typ, err = client.WorkOrderType.
		Create().
		SetName("deploy").
		Save(s.ctx)
	s.Require().NoError(err)
}

func (s *workOrderTestSuite) CreateWorkOrder() *ent.WorkOrderCreate {
	return s.client.Create().
		SetCreationDate(time.Now()).
		SetOwner(s.user).
		SetType(s.typ)
}

func (s *workOrderTestSuite) TestWorkOrderCloseDate() {
	s.Run("CreateDoneAndReopen", func() {
		order, err := s.CreateWorkOrder().
			SetName("antenna").
			SetStatus(workorder.StatusDONE).
			Save(s.ctx)
		s.Require().NoError(err)
		s.Assert().False(order.CloseDate.IsZero())

		order, err = order.Update().
			SetStatus(workorder.StatusPENDING).
			Save(s.ctx)
		s.Require().NoError(err)
		s.Assert().True(order.CloseDate.IsZero())
	})
	s.Run("CreateDoneWithCloseDate", func() {
		now := time.Now()
		order, err := s.CreateWorkOrder().
			SetName("pole").
			SetStatus(workorder.StatusDONE).
			SetCloseDate(now).
			Save(s.ctx)
		s.Require().NoError(err)
		s.Assert().True(order.CloseDate.Equal(now))

		order, err = order.Update().
			SetStatus(workorder.StatusDONE).
			Save(s.ctx)
		s.Require().NoError(err)
		s.Assert().True(
			now.Equal(order.CloseDate),
			"close date modified on status reapply",
		)
	})
	s.Run("CreatePlannedSetCloseDate", func() {
		order, err := s.CreateWorkOrder().
			SetName("tower").
			SetStatus(workorder.StatusPENDING).
			Save(s.ctx)
		s.Require().NoError(err)
		s.Assert().True(order.CloseDate.IsZero())

		now := time.Now()
		order, err = order.Update().
			SetCloseDate(now).
			Save(s.ctx)
		s.Require().NoError(err)
		s.Assert().True(order.CloseDate.Equal(now))
	})
	s.Run("UpdateDoneMany", func() {
		var ids []int
		for _, name := range []string{"foo", "bar", "baz"} {
			order, err := s.CreateWorkOrder().
				SetName(name).
				Save(s.ctx)
			s.Require().NoError(err)
			s.Assert().True(order.CloseDate.IsZero())
			ids = append(ids, order.ID)
		}
		n, err := s.client.Update().
			SetStatus(workorder.StatusDONE).
			Where(workorder.IDIn(ids...)).
			Save(privacy.DecisionContext(
				s.ctx, privacy.Allow,
			))
		s.Require().NoError(err)
		s.Assert().Equal(len(ids), n)
		count, err := s.client.
			Query().
			Where(
				workorder.IDIn(ids...),
				workorder.StatusEQ(workorder.StatusDONE),
				workorder.CloseDateNotNil(),
			).
			Count(s.ctx)
		s.Require().NoError(err)
		s.Assert().Equal(len(ids), count)
	})
}

func TestWorkOrderHooks(t *testing.T) {
	suite.Run(t, &workOrderTestSuite{})
}
