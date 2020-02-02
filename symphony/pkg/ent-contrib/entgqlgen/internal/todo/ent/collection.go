// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// WithFieldCollection tells the query-builder to eagerly load connected nodes by resolver context.
func (t *TodoQuery) WithFieldCollection(ctx context.Context, satisfies ...string) *TodoQuery {
	if resctx := graphql.GetResolverContext(ctx); resctx != nil {
		t = t.withField(graphql.GetRequestContext(ctx), resctx.Field, satisfies...)
	}
	return t
}

func (t *TodoQuery) withField(reqctx *graphql.RequestContext, field graphql.CollectedField, satisfies ...string) *TodoQuery {
	for _, field := range graphql.CollectFields(reqctx, field.Selections, satisfies) {
		switch field.Name {
		case "children":
			t = t.WithChildren(func(query *TodoQuery) {
				query.withField(reqctx, field)
			})
		case "parent":
			t = t.WithParent(func(query *TodoQuery) {
				query.withField(reqctx, field)
			})
		}
	}
	return t
}
