// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Todo defines the todo type schema.
type Todo struct {
	ent.Schema
}

// Fields returns todo fields.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.Text("text").
			NotEmpty(),
	}
}

// Edges returns todo edges.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Todo.Type).
			StructTag(`gqlgen:"children"`).
			From("parent").
			StructTag(`gqlgen:"parent"`).
			Unique(),
	}
}
