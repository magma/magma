// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/badoux/checkmail"
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Technician defines the technician schema.
type Technician struct {
	schema
}

// Fields returns technician fields.
func (Technician) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("email").
			Unique().
			Validate(checkmail.ValidateFormat),
	}
}

// Edges returns technician edges.
func (Technician) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("work_orders", WorkOrder.Type).
			Ref("technician"),
	}
}
