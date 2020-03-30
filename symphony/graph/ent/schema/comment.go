// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Comment defines the location type schema.
type Comment struct {
	schema
}

// Fields returns comment fields.
func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.String("text"),
	}
}

// Edges returns work order edges.
func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("author", User.Type).
			Required().
			Unique(),
	}
}
