package jobs

import (
	"testing"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/enttest"
	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/stretchr/testify/require"
)

func newJobsTestResolver(t *testing.T) *TestJobsResolver {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	return newResolver(t, sql.OpenDB(name, db))
}

func newResolver(t *testing.T, drv dialect.Driver) *TestJobsResolver {
	client := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	r := resolver.New(resolver.Config{
		Logger:     logtest.NewTestLogger(t),
		Subscriber: event.NewNopSubscriber(),
	})
	return &TestJobsResolver{
		drv:    drv,
		client: client,
		jobsRunner: jobs{
			logger: logtest.NewTestLogger(t),
			r:      r,
		},
	}
}
