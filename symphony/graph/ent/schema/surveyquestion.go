// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// SurveyQuestion holds the schema definition for the SurveyQuestion entity.
type SurveyQuestion struct {
	schema
}

// Fields of the SurveyQuestion.
func (SurveyQuestion) Fields() []ent.Field {
	return []ent.Field{
		field.String("form_name").Optional(),
		field.String("form_description").Optional(),
		field.Int("form_index"),
		field.String("question_type").
			Comment("Type of data collected by this question").
			Optional(),
		field.String("question_format").
			Comment("Format of data collected by this question").
			Optional(),
		field.String("question_text").
			Comment("The actual question that was asked").
			Optional(),
		field.Int("question_index"),

		field.Bool("bool_data").
			Comment("Yes/No or True/False data collected for this question").
			Optional(),
		field.String("email_data").
			Comment("Email address collected for this question").
			Optional(),
		field.Float("latitude").Optional(),
		field.Float("longitude").Optional(),
		field.Float("location_accuracy").
			Comment("Accuracy of the GPS location in meters").
			Optional(),
		field.Float("altitude").
			Comment("Altitude in meters above the WGS 84 reference ellipsoid").
			Optional(),
		field.String("phone_data").
			Comment("Phone number collected for this question").
			Optional(),
		field.String("text_data").
			Comment("Text data collected for this question").
			Optional(),
		field.Float("float_data").
			Comment("Float data collected for this question").
			Optional(),
		field.Int("int_data").
			Comment("Int data collected for this question").
			Optional(),
		field.Time("date_data").
			Comment("Date data collected for this question").
			Optional(),
	}
}

// Edges of the SurveyQuestion.
func (SurveyQuestion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("survey", Survey.Type).
			Unique().
			Required(),
		edge.From("wifi_scan", SurveyWiFiScan.Type).
			Ref("survey_question"),
		edge.From("cell_scan", SurveyCellScan.Type).
			Ref("survey_question"),
		edge.To("photo_data", File.Type),
		edge.To("images", File.Type),
	}
}
