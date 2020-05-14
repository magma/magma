// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/enttest"
	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/directive"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	_ "github.com/mattn/go-sqlite3"
)

type TestResolver struct {
	generated.ResolverRoot
	logger     log.Logger
	client     *ent.Client
	emitter    *event.PipeEmitter
	subscriber *event.PipeSubscriber
}

func newTestResolver(t *testing.T, opts ...Option) *TestResolver {
	c := enttest.Open(t, dialect.SQLite,
		fmt.Sprintf("file:%s-%d?mode=memory&cache=shared&_fk=1",
			t.Name(), time.Now().UnixNano(),
		),
		enttest.WithMigrateOptions(
			migrate.WithGlobalUniqueID(true),
		),
	)

	emitter, subscriber := event.Pipe()
	logger := logtest.NewTestLogger(t)
	eventer := event.Eventer{Logger: logger, Emitter: emitter}
	eventer.HookTo(c)

	r := New(Config{
		Logger:     logger,
		Subscriber: subscriber,
	}, opts...)

	return &TestResolver{
		ResolverRoot: r,
		logger:       logger,
		client:       c,
		emitter:      emitter,
		subscriber:   subscriber,
	}
}

func (tr *TestResolver) Close() error {
	var (
		shutdowners = []func(context.Context) error{
			func(context.Context) error { return tr.client.Close() },
			tr.emitter.Shutdown,
			tr.subscriber.Shutdown,
		}
		ctx = context.Background()
		wg  sync.WaitGroup
	)
	wg.Add(len(shutdowners))
	for _, shutdowner := range shutdowners {
		go func(shutdowner func(context.Context) error) {
			defer wg.Done()
			_ = shutdowner(ctx)
		}(shutdowner)
	}
	wg.Wait()
	return nil
}

func (tr *TestResolver) GraphClient() *client.Client {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers:  tr.ResolverRoot,
				Directives: directive.New(tr.logger),
			},
		),
	)
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		ctx = viewertest.NewContext(ctx, tr.client)
		return next(ctx)
	})
	return client.New(srv)
}

func resolverctx(t *testing.T) (generated.ResolverRoot, context.Context) {
	r := newTestResolver(t)
	return r, viewertest.NewContext(context.Background(), r.client)
}

func mutationctx(t *testing.T) (generated.MutationResolver, context.Context) {
	r, ctx := resolverctx(t)
	return r.Mutation(), ctx
}
