// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/handler"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/directive"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryNode(t *testing.T) {
	resolver, err := newTestResolver(t)
	require.NoError(t, err)
	defer resolver.drv.Close()

	c := client.New(handler.GraphQL(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers:  resolver,
				Directives: directive.New(logtest.NewTestLogger(t)),
			},
		),
		handler.RequestMiddleware(
			func(ctx context.Context, next func(context.Context) []byte) []byte {
				ctx = ent.NewContext(ctx, resolver.client)
				return next(ctx)
			},
		),
	))

	var lt struct{ AddLocationType struct{ ID string } }
	err = c.Post(
		`mutation($input: AddLocationTypeInput!) { addLocationType(input: $input) { id } }`,
		&lt,
		client.Var("input", models.AddLocationTypeInput{Name: "city"}),
	)
	require.NoError(t, err)

	var l struct{ AddLocation struct{ ID string } }
	err = c.Post(
		`mutation($input: AddLocationInput!) { addLocation(input: $input) { id } }`,
		&l,
		client.Var("input", models.AddLocationInput{Name: "tlv", Type: lt.AddLocationType.ID}),
	)
	require.NoError(t, err)

	t.Run("LocationType", func(t *testing.T) {
		var rsp struct{ Node struct{ Name string } }
		err := c.Post(
			`query($id: ID!) { node(id: $id) { ... on LocationType { name } } }`,
			&rsp,
			client.Var("id", lt.AddLocationType.ID),
		)
		require.NoError(t, err)
		assert.Equal(t, "city", rsp.Node.Name)
	})
	t.Run("Location", func(t *testing.T) {
		var rsp struct{ Node struct{ Name string } }
		err := c.Post(
			`query($id: ID!) { node(id: $id) { ... on Location { name } } }`,
			&rsp,
			client.Var("id", l.AddLocation.ID),
		)
		require.NoError(t, err)
		assert.Equal(t, "tlv", rsp.Node.Name)
	})
	t.Run("NonExistent", func(t *testing.T) {
		rsp, err := c.RawPost(
			`query($id: ID!) { node(id: $id) { id } }`,
			client.Var("id", func() string {
				id, err := strconv.Atoi(l.AddLocation.ID)
				require.NoError(t, err)
				return strconv.Itoa(id + 42)
			}()),
		)
		require.NoError(t, err)
		assert.Empty(t, rsp.Errors)
		v, ok := rsp.Data.(map[string]interface{})["node"]
		assert.True(t, ok)
		assert.Nil(t, v)
	})
	t.Run("BadID", func(t *testing.T) {
		rsp, err := c.RawPost(`query { node(id: "_") { id } }`)
		require.NoError(t, err)
		assert.Empty(t, rsp.Errors)
		v, ok := rsp.Data.(map[string]interface{})["node"]
		assert.True(t, ok)
		assert.Nil(t, v)
	})
}
