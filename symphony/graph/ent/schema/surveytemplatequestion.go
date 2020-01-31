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

// SurveyTemplateQuestion holds the schema definition for the SurveyTemplateQuestion entity.
type SurveyTemplateQuestion struct {
	schema
}

// Fields of the SurveyTemplateQuestion.
func (SurveyTemplateQuestion) Fields() []ent.Field {
	return []ent.Field{
		field.String("question_title"),
		field.String("question_description"),
		field.String("question_type"),
		field.Int("index"),
	}
}

// Edges of the SurveyTemplateQuestion.
func (SurveyTemplateQuestion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", SurveyTemplateCategory.Type).
			Ref("survey_template_questions").
			Unique(),
	}
}

// Indexes of the SurveyTemplateQuestion.
func (SurveyTemplateQuestion) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Edges("category").
			Unique(),
	}
}
