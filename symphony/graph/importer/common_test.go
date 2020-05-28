// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"testing"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/enttest"
	"github.com/facebookincubator/symphony/pkg/ent/migrate"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/testdb"

	"github.com/stretchr/testify/require"
)

const (
	svcName  = "serviceName"
	svc2Name = "serviceName2"
	svc3Name = "serviceName3"
)

type TestImporterResolver struct {
	drv      dialect.Driver
	client   *ent.Client
	importer importer
}

func newImporterTestResolver(t *testing.T) *TestImporterResolver {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	return newResolver(t, sql.OpenDB(name, db))
}

func newResolver(t *testing.T, drv dialect.Driver) *TestImporterResolver {
	client := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	r := resolver.New(resolver.Config{
		Logger:     logtest.NewTestLogger(t),
		Subscriber: event.NewNopSubscriber(),
	})
	return &TestImporterResolver{
		drv:    drv,
		client: client,
		importer: importer{
			logger: logtest.NewTestLogger(t),
			r:      r,
		},
	}
}

func prepareSvcData(ctx context.Context, t *testing.T, r TestImporterResolver) {
	mr := r.importer.r.Mutation()
	serviceType, _ := mr.AddServiceType(ctx, models.ServiceTypeCreateData{Name: "L2 Service", HasCustomer: false})
	_, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          svcName,
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          svc2Name,
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          svc3Name,
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
}
