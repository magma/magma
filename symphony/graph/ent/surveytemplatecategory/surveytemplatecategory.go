// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveytemplatecategory

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the surveytemplatecategory type in the database.
	Label = "survey_template_category"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldCategoryTitle holds the string denoting the category_title vertex property in the database.
	FieldCategoryTitle = "category_title"
	// FieldCategoryDescription holds the string denoting the category_description vertex property in the database.
	FieldCategoryDescription = "category_description"

	// Table holds the table name of the surveytemplatecategory in the database.
	Table = "survey_template_categories"
	// SurveyTemplateQuestionsTable is the table the holds the survey_template_questions relation/edge.
	SurveyTemplateQuestionsTable = "survey_template_questions"
	// SurveyTemplateQuestionsInverseTable is the table name for the SurveyTemplateQuestion entity.
	// It exists in this package in order to avoid circular dependency with the "surveytemplatequestion" package.
	SurveyTemplateQuestionsInverseTable = "survey_template_questions"
	// SurveyTemplateQuestionsColumn is the table column denoting the survey_template_questions relation/edge.
	SurveyTemplateQuestionsColumn = "category_id"
)

// Columns holds all SQL columns are surveytemplatecategory fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldCategoryTitle,
	FieldCategoryDescription,
}

var (
	mixin       = schema.SurveyTemplateCategory{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.SurveyTemplateCategory{}.Fields()

	// descCreateTime is the schema descriptor for create_time field.
	descCreateTime = mixinFields[0][0].Descriptor()
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime = descCreateTime.Default.(func() time.Time)

	// descUpdateTime is the schema descriptor for update_time field.
	descUpdateTime = mixinFields[0][1].Descriptor()
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime = descUpdateTime.Default.(func() time.Time)
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime = descUpdateTime.UpdateDefault.(func() time.Time)
)
