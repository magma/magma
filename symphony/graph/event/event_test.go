// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"
	"testing"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/enttest"
	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type eventTestSuite struct {
	suite.Suite
	ctx        context.Context
	logger     log.Logger
	client     *ent.Client
	user       *ent.User
	subscriber Subscriber
}

func (s *eventTestSuite) SetupSuite() {
	db, name, err := testdb.Open()
	s.Require().NoError(err)
	db.SetMaxOpenConns(1)

	s.client = enttest.NewClient(s.T(),
		enttest.WithOptions(ent.Driver(sql.OpenDB(name, db))),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	s.ctx = viewertest.NewContext(context.Background(), s.client)
	s.user = viewer.FromContext(s.ctx).(*viewer.UserViewer).User()
	s.logger = logtest.NewTestLogger(s.T())

	eventer := Eventer{Logger: s.logger}
	eventer.Emitter, s.subscriber = Pipe()
	eventer.HookTo(s.client)
}

func TestPipe(t *testing.T) {
	emitter, subscriber := Pipe()
	require.NotNil(t, emitter)
	require.NotNil(t, subscriber)
	subscription, err := subscriber.Subscribe(context.Background())
	require.NoError(t, err)

	err = emitter.Emit(context.Background(), t.Name(), t.Name(), nil)
	require.NoError(t, err)
	msg, err := subscription.Receive(context.Background())
	require.NoError(t, err)
	require.Equal(t, t.Name(), msg.Metadata[TenantHeader])
	require.Equal(t, t.Name(), msg.Metadata[NameHeader])
	require.Empty(t, msg.Body)
}
