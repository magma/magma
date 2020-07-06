// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (t *TodoQuery) CollectFields(ctx context.Context, satisfies ...string) *TodoQuery {
	if fc := graphql.GetFieldContext(ctx); fc != nil {
		t = t.collectField(graphql.GetOperationContext(ctx), fc.Field, satisfies...)
	}
	return t
}

func (t *TodoQuery) collectField(ctx *graphql.OperationContext, field graphql.CollectedField, satisfies ...string) *TodoQuery {
	for _, field := range graphql.CollectFields(ctx, field.Selections, satisfies) {
		switch field.Name {
		case "children":
			t = t.WithChildren(func(query *TodoQuery) {
				query.collectField(ctx, field)
			})
		case "parent":
			t = t.WithParent(func(query *TodoQuery) {
				query.collectField(ctx, field)
			})
		}
	}
	return t
}
