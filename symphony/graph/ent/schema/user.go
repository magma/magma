// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// User defines the user schema.
type User struct {
	schema
}

// Fields returns user fields.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("auth_id").
			NotEmpty().
			Immutable().
			Unique(),
		field.String("first_name").
			NotEmpty().
			Optional(),
		field.String("last_name").
			NotEmpty().
			Optional(),
		field.String("email").
			NotEmpty().
			Optional(),
		field.Enum("status").
			Values("ACTIVE", "DEACTIVATED").
			Default("ACTIVE"),
		field.Enum("role").
			Values("USER", "ADMIN", "OWNER").
			Default("USER"),
	}
}

// Edges returns user edges.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("profile_photo", File.Type).
			Unique(),
	}
}
