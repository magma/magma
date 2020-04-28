// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveytemplatecategory

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the surveytemplatecategory type in the database.
	Label = "survey_template_category"
	// FieldID holds the string denoting the id field in the database.
	FieldID                  = "id"             // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime          = "create_time"    // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime          = "update_time"    // FieldCategoryTitle holds the string denoting the category_title vertex property in the database.
	FieldCategoryTitle       = "category_title" // FieldCategoryDescription holds the string denoting the category_description vertex property in the database.
	FieldCategoryDescription = "category_description"

	// EdgeSurveyTemplateQuestions holds the string denoting the survey_template_questions edge name in mutations.
	EdgeSurveyTemplateQuestions = "survey_template_questions"

	// Table holds the table name of the surveytemplatecategory in the database.
	Table = "survey_template_categories"
	// SurveyTemplateQuestionsTable is the table the holds the survey_template_questions relation/edge.
	SurveyTemplateQuestionsTable = "survey_template_questions"
	// SurveyTemplateQuestionsInverseTable is the table name for the SurveyTemplateQuestion entity.
	// It exists in this package in order to avoid circular dependency with the "surveytemplatequestion" package.
	SurveyTemplateQuestionsInverseTable = "survey_template_questions"
	// SurveyTemplateQuestionsColumn is the table column denoting the survey_template_questions relation/edge.
	SurveyTemplateQuestionsColumn = "survey_template_category_survey_template_questions"
)

// Columns holds all SQL columns for surveytemplatecategory fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldCategoryTitle,
	FieldCategoryDescription,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the SurveyTemplateCategory type.
var ForeignKeys = []string{
	"location_type_survey_template_categories",
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/facebookincubator/symphony/graph/ent/runtime"
//
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime func() time.Time
)
