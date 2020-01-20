// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// SurveyTemplateCategory holds the schema definition for the SurveyTemplateCategory entity.
type SurveyTemplateCategory struct {
	schema
}

// Fields of the SurveyTemplateCategory.
func (SurveyTemplateCategory) Fields() []ent.Field {
	return []ent.Field{
		field.String("category_title"),
		field.String("category_description"),
	}
}

// Edges of the SurveyTemplateCategory.
func (SurveyTemplateCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("survey_template_questions", SurveyTemplateQuestion.Type),
	}
}
