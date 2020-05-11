// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/stretchr/testify/suite"
)

type workOrderTestSuite struct {
	eventTestSuite
	typ *ent.WorkOrderType
}

func TestWorkOrderEvents(t *testing.T) {
	suite.Run(t, &workOrderTestSuite{})
}

func (s *workOrderTestSuite) SetupSuite() {
	s.eventTestSuite.SetupSuite()
	s.typ = s.client.WorkOrderType.Create().
		SetName("Chore").
		SaveX(s.ctx)
}

func (s *workOrderTestSuite) TestWorkOrderCreate() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithCancel(s.ctx)
		events := []string{WorkOrderAdded, WorkOrderDone}
		emitted := make(map[string]struct{}, len(events))
		for i := range events {
			emitted[events[i]] = struct{}{}
		}
		err := SubscribeAndListen(ctx, ListenerConfig{
			Subscriber: s.subscriber,
			Logger:     s.logger.Background(),
			Tenant:     viewertest.DefaultTenant,
			Events:     events,
			Handler: HandlerFunc(func(_ context.Context, name string, body []byte) error {
				s.Assert().NotEmpty(body)
				_, ok := emitted[name]
				s.Assert().True(ok)
				delete(emitted, name)
				if len(emitted) == 0 {
					cancel()
				}
				return nil
			}),
		})
		s.Require().True(errors.Is(err, context.Canceled))
	}()
	woType := s.client.WorkOrderType.Create().
		SetName("CleanType").
		SaveX(s.ctx)
	s.client.WorkOrder.Create().
		SetName("Clean").
		SetType(woType).
		SetCreationDate(time.Now()).
		SetOwner(s.user).
		SetStatus(models.WorkOrderStatusDone.String()).
		SaveX(s.ctx)
	wg.Wait()
}

func (s *workOrderTestSuite) TestWorkOrderUpdate() {
	err := s.client.WorkOrder.Update().
		SetStatus(models.WorkOrderStatusDone.String()).
		Exec(s.ctx)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "work order status update to done by predicate not allowed")
}

func (s *workOrderTestSuite) TestWorkOrderUpdateOne() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithCancel(s.ctx)
		err := SubscribeAndListen(ctx, ListenerConfig{
			Subscriber: s.subscriber,
			Logger:     s.logger.Background(),
			Tenant:     viewertest.DefaultTenant,
			Events:     []string{WorkOrderDone},
			Handler: HandlerFunc(func(_ context.Context, name string, body []byte) error {
				s.Assert().Equal(WorkOrderDone, name)
				s.Assert().NotEmpty(body)
				cancel()
				return nil
			}),
		})
		s.Require().True(errors.Is(err, context.Canceled))
	}()

	woType := s.client.WorkOrderType.Create().
		SetName("VacuumType").
		SaveX(s.ctx)
	wo := s.client.WorkOrder.Create().
		SetName("Vacuum").
		SetType(woType).
		SetCreationDate(time.Now()).
		SetOwner(s.user).
		SaveX(s.ctx)
	tx, err := s.client.Tx(s.ctx)
	s.Require().NoError(err)
	ctx := ent.NewTxContext(s.ctx, tx)
	tx.WorkOrder.UpdateOne(wo).
		SetStatus(models.WorkOrderStatusDone.String()).
		ExecX(ctx)
	err = tx.Commit()
	s.Require().NoError(err)
	wg.Wait()
}
