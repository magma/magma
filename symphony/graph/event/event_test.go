// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package event

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/enttest"
	"github.com/facebookincubator/symphony/pkg/ent/migrate"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/pubsub"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/suite"
)

type eventTestSuite struct {
	suite.Suite
	ctx        context.Context
	logger     log.Logger
	client     *ent.Client
	user       *ent.User
	subscriber pubsub.Subscriber
}

func (s *eventTestSuite) SetupSuite(opts ...viewertest.Option) {
	db, name, err := testdb.Open()
	s.Require().NoError(err)
	db.SetMaxOpenConns(1)

	s.client = enttest.NewClient(s.T(),
		enttest.WithOptions(ent.Driver(sql.OpenDB(name, db))),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	s.ctx = viewertest.NewContext(context.Background(), s.client, opts...)
	s.user = viewer.FromContext(s.ctx).(*viewer.UserViewer).User()
	s.logger = logtest.NewTestLogger(s.T())

	eventer := Eventer{Logger: s.logger}
	eventer.Emitter, s.subscriber = pubsub.Pipe()
	eventer.HookTo(s.client)
}
