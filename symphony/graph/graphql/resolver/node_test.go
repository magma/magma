// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryNode(t *testing.T) {
	resolver := newTestResolver(t)
	defer resolver.drv.Close()
	c := newGraphClient(t, resolver)

	var lt struct{ AddLocationType struct{ ID string } }
	c.MustPost(
		`mutation($input: AddLocationTypeInput!) { addLocationType(input: $input) { id } }`,
		&lt,
		client.Var("input", models.AddLocationTypeInput{Name: "city"}),
	)

	var l struct{ AddLocation struct{ ID string } }
	c.MustPost(
		`mutation($input: AddLocationInput!) { addLocation(input: $input) { id } }`,
		&l,
		client.Var("input", models.AddLocationInput{Name: "tlv", Type: lt.AddLocationType.ID}),
	)

	t.Run("LocationType", func(t *testing.T) {
		var rsp struct{ Node struct{ Name string } }
		c.MustPost(
			`query($id: ID!) { node(id: $id) { ... on LocationType { name } } }`,
			&rsp,
			client.Var("id", lt.AddLocationType.ID),
		)
		assert.Equal(t, "city", rsp.Node.Name)
	})
	t.Run("Location", func(t *testing.T) {
		var rsp struct{ Node struct{ Name string } }
		c.MustPost(
			`query($id: ID!) { node(id: $id) { ... on Location { name } } }`,
			&rsp,
			client.Var("id", l.AddLocation.ID),
		)
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
		rsp, err := c.RawPost(`query { node(id: "-1") { id } }`)
		require.NoError(t, err)
		assert.Empty(t, rsp.Errors)
		v, ok := rsp.Data.(map[string]interface{})["node"]
		assert.True(t, ok)
		assert.Nil(t, v)
	})
}
