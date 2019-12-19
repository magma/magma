// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/testdb"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/schema"

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

func newTestResolver(t *testing.T, opts ...ResolveOption) (*TestResolver, error) {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	return newResolver(t, sql.OpenDB(name, db), opts...)
}

func newResolver(t *testing.T, drv dialect.Driver, opts ...ResolveOption) (*TestResolver, error) {
	if *debug {
		drv = dialect.Debug(drv)
	}
	client := ent.NewClient(ent.Driver(drv))
	require.NoError(t, client.Schema.Create(context.Background(), schema.WithGlobalUniqueID(true)))
	r, err := New(logtest.NewTestLogger(t), opts...)
	if err != nil {
		return nil, err
	}
	return &TestResolver{r, drv, client}, nil
}

func resolverctx(t *testing.T) (generated.ResolverRoot, context.Context) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	return r, viewertest.NewContext(r.client)
}

func mutationctx(t *testing.T) (generated.MutationResolver, context.Context) {
	r, ctx := resolverctx(t)
	return r.Mutation(), ctx
}
