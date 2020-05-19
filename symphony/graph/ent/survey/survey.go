// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package survey

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the survey type in the database.
	Label = "survey"
	// FieldID holds the string denoting the id field in the database.
	FieldID                  = "id"                 // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime          = "create_time"        // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime          = "update_time"        // FieldName holds the string denoting the name vertex property in the database.
	FieldName                = "name"               // FieldOwnerName holds the string denoting the owner_name vertex property in the database.
	FieldOwnerName           = "owner_name"         // FieldCreationTimestamp holds the string denoting the creation_timestamp vertex property in the database.
	FieldCreationTimestamp   = "creation_timestamp" // FieldCompletionTimestamp holds the string denoting the completion_timestamp vertex property in the database.
	FieldCompletionTimestamp = "completion_timestamp"

	// EdgeLocation holds the string denoting the location edge name in mutations.
	EdgeLocation = "location"
	// EdgeSourceFile holds the string denoting the source_file edge name in mutations.
	EdgeSourceFile = "source_file"
	// EdgeQuestions holds the string denoting the questions edge name in mutations.
	EdgeQuestions = "questions"

	// Table holds the table name of the survey in the database.
	Table = "surveys"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "surveys"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "survey_location"
	// SourceFileTable is the table the holds the source_file relation/edge.
	SourceFileTable = "files"
	// SourceFileInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	SourceFileInverseTable = "files"
	// SourceFileColumn is the table column denoting the source_file relation/edge.
	SourceFileColumn = "survey_source_file"
	// QuestionsTable is the table the holds the questions relation/edge.
	QuestionsTable = "survey_questions"
	// QuestionsInverseTable is the table name for the SurveyQuestion entity.
	// It exists in this package in order to avoid circular dependency with the "surveyquestion" package.
	QuestionsInverseTable = "survey_questions"
	// QuestionsColumn is the table column denoting the questions relation/edge.
	QuestionsColumn = "survey_question_survey"
)

// Columns holds all SQL columns for survey fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldOwnerName,
	FieldCreationTimestamp,
	FieldCompletionTimestamp,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Survey type.
var ForeignKeys = []string{
	"survey_location",
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
