// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/badoux/checkmail"
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/ent/schema/index"
	"github.com/facebookincubator/symphony/frontier/ent/user/role"
)

// User defines user schema.
type User struct {
	schema
}

// Config configures user schema.
func (User) Config() ent.Config {
	return ent.Config{
		Table: "Users",
	}
}

// Fields of user entity.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			Validate(checkmail.ValidateFormat),
		field.String("password").
			Sensitive().
			NotEmpty(),
		field.Int("role").
			Default(int(role.UserRole)).
			Validate(role.ValidateValue),
		field.String("tenant").
			StorageKey("organization").
			Default("fb-test").
			Immutable(),
		field.Strings("networks").
			StorageKey("networkIDs"),
		field.Strings("tabs").
			Optional(),
	}
}

// Edges of user entity.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tokens", Token.Type),
	}
}

// Indexes of user entity.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email", "tenant").
			Unique(),
	}
}
