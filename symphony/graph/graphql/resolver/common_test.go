// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"flag"
	"net/http"
	"os"
	"testing"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/enttest"
	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/directive"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/testdb"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/handler"
	"github.com/stretchr/testify/require"
)

var debug = flag.Bool("debug", false, "run database driver on debug mode")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

type TestResolver struct {
	generated.ResolverRoot
	drv    dialect.Driver
	client *ent.Client
}

func newTestResolver(t *testing.T, opts ...Option) *TestResolver {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	return newResolver(t, sql.OpenDB(name, db), opts...)
}

func newResolver(t *testing.T, drv dialect.Driver, opts ...Option) *TestResolver {
	if *debug {
		drv = dialect.Debug(drv)
	}
	c := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)

	emitter, subscriber := event.Pipe()
	logger := logtest.NewTestLogger(t)
	eventer := event.Eventer{Logger: logger, Emitter: emitter}
	eventer.HookTo(c)

	r := New(Config{Logger: logger, Subscriber: subscriber}, opts...)
	return &TestResolver{r, drv, c}
}

func newGraphClient(t *testing.T, resolver *TestResolver) *client.Client {
	gql := handler.GraphQL(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers:  resolver,
				Directives: directive.New(logtest.NewTestLogger(t)),
			},
		),
	)
	return client.New(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := viewertest.NewContext(context.Background(), resolver.client)
		gql.ServeHTTP(w, r.WithContext(ctx))
	}))
}

func resolverctx(t *testing.T) (generated.ResolverRoot, context.Context) {
	r := newTestResolver(t)
	return r, viewertest.NewContext(context.Background(), r.client)
}

func mutationctx(t *testing.T) (generated.MutationResolver, context.Context) {
	r, ctx := resolverctx(t)
	return r.Mutation(), ctx
}
