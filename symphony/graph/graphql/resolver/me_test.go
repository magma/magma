// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/handler"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/stretchr/testify/assert"
)

func TestQueryMe(t *testing.T) {
	resolver := newTestResolver(t)
	defer resolver.drv.Close()

	v := &viewer.Viewer{Tenant: "testing", User: "tester@example.com"}
	h := handler.GraphQL(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
		handler.RequestMiddleware(
			func(ctx context.Context, next func(context.Context) []byte) []byte {
				return next(viewer.NewContext(ctx, v))
			},
		),
	)

	var rsp struct {
		Me struct {
			Tenant string
			Email  string
		}
	}
	client.New(h).MustPost("query { me { tenant, email } }", &rsp)
	assert.Equal(t, v.Tenant, rsp.Me.Tenant)
	assert.Equal(t, v.User, rsp.Me.Email)
}
