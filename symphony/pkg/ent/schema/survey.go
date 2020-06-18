// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Survey holds the schema definition for the Survey entity.
type Survey struct {
	schema
}

// Fields of the Survey.
func (Survey) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("owner_name").Optional(),
		field.Time("creation_timestamp").
			StructTag(`gqlgen:"creationTimestamp"`).Optional(),
		field.Time("completion_timestamp").
			StructTag(`gqlgen:"completionTimestamp"`),
	}
}

// Edges of the Survey.
func (Survey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("location", Location.Type).Unique(),
		edge.To("source_file", File.Type).Unique(),
		edge.From("questions", SurveyQuestion.Type).
			Ref("survey"),
	}
}
