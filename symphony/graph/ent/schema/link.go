// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Link defines the link schema.
type Link struct {
	schema
}

// Fields returns link fields.
func (Link) Fields() []ent.Field {
	return []ent.Field{
		field.String("future_state").
			Optional(),
	}
}

// Edges returns link edges.
func (Link) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("ports", EquipmentPort.Type).
			Ref("link"),
		edge.To("work_order", WorkOrder.Type).
			Unique(),
		edge.To("properties", Property.Type),
		edge.From("service", Service.Type).
			Ref("links"),
	}
}
