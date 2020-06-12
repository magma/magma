// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/locationtype"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/suite"
)

type logTestSuite struct {
	eventTestSuite
	toUpdate *ent.LocationType
	toDelete *ent.LocationType
}

func TestLogEvents(t *testing.T) {
	suite.Run(t, &logTestSuite{})
}

func (s *logTestSuite) SetupSuite() {
	s.eventTestSuite.SetupSuite(viewertest.WithFeatures(viewer.FeatureGraphEventLogging))
	s.toUpdate = s.client.LocationType.Create().
		SetName("LocationTypeToUpdate").
		SaveX(s.ctx)
	s.toDelete = s.client.LocationType.Create().
		SetName("LocationTypeToDelete").
		SaveX(s.ctx)
}

func (s *logTestSuite) subscribeForOneEvent(wg *sync.WaitGroup, expect func(entry pubsub.LogEntry)) {
	wg.Add(1)
	ctx, cancel := context.WithCancel(s.ctx)
	events := []string{pubsub.EntMutation}
	listener, err := pubsub.NewListener(s.ctx, pubsub.ListenerConfig{
		Subscriber: s.subscriber,
		Logger:     s.logger.Background(),
		Events:     events,
		Handler: pubsub.HandlerFunc(func(_ context.Context, tenant, name string, body []byte) error {
			s.Assert().NotEmpty(body)
			s.Assert().Equal(viewertest.DefaultTenant, tenant)
			s.Assert().Equal(pubsub.EntMutation, name)
			var entry pubsub.LogEntry
			err := pubsub.Unmarshal(body, &entry)
			s.NoError(err)
			expect(entry)
			cancel()
			return nil
		}),
	})
	s.Assert().NoError(err)
	go func() {
		defer wg.Done()
		defer listener.Shutdown(ctx)
		err := listener.Listen(ctx)
		s.Require().True(errors.Is(err, context.Canceled))
	}()
}

func (s *logTestSuite) TestCreateEnt() {
	var wg sync.WaitGroup
	s.subscribeForOneEvent(&wg, func(entry pubsub.LogEntry) {
		s.Assert().Equal(s.user.AuthID, entry.UserName)
		s.Assert().Equal(s.user.ID, *entry.UserID)
		s.Assert().Equal(ent.OpCreate, entry.Operation)
		s.Assert().Nil(entry.PrevState)
		s.Assert().NotNil(entry.CurrState)
		s.Assert().Equal("LocationType", entry.CurrState.Type)
		found := 0
		for _, field := range entry.CurrState.Fields {
			switch field.Name {
			case locationtype.FieldName:
				s.Assert().Equal("SomeName", field.MustGetString())
				s.Assert().Equal("string", field.Type)
				found++
			case locationtype.FieldIndex:
				s.Assert().Equal(3, field.MustGetInt())
				s.Assert().Equal("int", field.Type)
				found++
			}
		}
		s.Assert().Equal(2, found)
	})
	s.client.LocationType.Create().
		SetName("SomeName").
		SetIndex(3).
		SaveX(s.ctx)
	wg.Wait()
}

func (s *logTestSuite) TestUpdateEnt() {
	var wg sync.WaitGroup
	s.subscribeForOneEvent(&wg, func(entry pubsub.LogEntry) {
		s.Assert().Equal(s.user.AuthID, entry.UserName)
		s.Assert().Equal(s.user.ID, *entry.UserID)
		s.Assert().Equal(ent.OpUpdateOne, entry.Operation)
		s.Assert().NotNil(entry.PrevState)
		found := 0
		for _, field := range entry.PrevState.Fields {
			if field.Name == locationtype.FieldName {
				s.Assert().Equal("LocationTypeToUpdate", field.MustGetString())
				s.Assert().Equal("string", field.Type)
				found++
			}
		}
		s.Assert().NotNil(entry.CurrState)
		for _, field := range entry.CurrState.Fields {
			if field.Name == locationtype.FieldName {
				s.Assert().Equal("NewName", field.MustGetString())
				s.Assert().Equal("string", field.Type)
				found++
			}
		}
		s.Assert().Equal(2, found)
	})
	s.client.LocationType.UpdateOne(s.toUpdate).
		SetName("NewName").
		SaveX(s.ctx)
	wg.Wait()
}

func (s *logTestSuite) TestDeleteEnt() {
	var wg sync.WaitGroup
	s.subscribeForOneEvent(&wg, func(entry pubsub.LogEntry) {
		s.Assert().Equal(s.user.AuthID, entry.UserName)
		s.Assert().Equal(s.user.ID, *entry.UserID)
		s.Assert().Equal(ent.OpDeleteOne, entry.Operation)
		s.Assert().NotNil(entry.PrevState)
		found := 0
		for _, field := range entry.PrevState.Fields {
			if field.Name == locationtype.FieldName {
				s.Assert().Equal("LocationTypeToDelete", field.MustGetString())
				s.Assert().Equal("string", field.Type)
				found++
			}
		}
		s.Assert().Equal(1, found)
		s.Assert().Nil(entry.CurrState)
	})
	s.client.LocationType.DeleteOne(s.toDelete).
		ExecX(s.ctx)
	wg.Wait()
}
