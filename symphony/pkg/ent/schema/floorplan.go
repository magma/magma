// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// A floor plan is an image file that is mapped to a location on a map by
// maintaining reference point and scale. The reference point tells us which
// (x, y) pair in the image maps to a (lat, lon) pair on a map. The provides two
// (x, y) pairs from the image and the distance between them in meters

type FloorPlan struct {
	schema
}

// Fields of the FloorPlan
func (FloorPlan) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

// Edges of the FloorPlan
func (FloorPlan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("location", Location.Type).Unique(),
		edge.To("reference_point", FloorPlanReferencePoint.Type).Unique(),
		edge.To("scale", FloorPlanScale.Type).Unique(),
		edge.To("image", File.Type).Unique(),
	}
}

type FloorPlanReferencePoint struct {
	schema
}

// Fields of the FloorPlanReferencePoint
func (FloorPlanReferencePoint) Fields() []ent.Field {
	return []ent.Field{
		field.Int("x"),
		field.Int("y"),
		field.Float("latitude"),
		field.Float("longitude"),
	}
}

// Edges of the FloorPlanReferencePoint
func (FloorPlanReferencePoint) Edges() []ent.Edge {
	return []ent.Edge{}
}

type FloorPlanScale struct {
	schema
}

// Fields of the FloorPlanScale
func (FloorPlanScale) Fields() []ent.Field {
	return []ent.Field{
		field.Int("reference_point1_x"),
		field.Int("reference_point1_y"),
		field.Int("reference_point2_x"),
		field.Int("reference_point2_y"),
		field.Float("scale_in_meters"),
	}
}

// Edges of the FloorPlanScale
func (FloorPlanScale) Edges() []ent.Edge {
	return []ent.Edge{}
}
