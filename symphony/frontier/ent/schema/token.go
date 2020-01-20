// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
)

// Token defines token schema.
type Token struct {
	schema
}

// Fields of token entity.
func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("value").
			Sensitive().
			Immutable().
			NotEmpty(),
	}
}

// Edges of token entity.
func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("tokens").
			Required().
			Unique(),
	}
}

// Indexes of token entity.
func (Token) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("value").
			Edges("user").
			Unique(),
	}
}
