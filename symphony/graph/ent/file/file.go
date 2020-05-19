// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package file

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the file type in the database.
	Label = "file"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"           // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time"  // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time"  // FieldType holds the string denoting the type vertex property in the database.
	FieldType        = "type"         // FieldName holds the string denoting the name vertex property in the database.
	FieldName        = "name"         // FieldSize holds the string denoting the size vertex property in the database.
	FieldSize        = "size"         // FieldModifiedAt holds the string denoting the modified_at vertex property in the database.
	FieldModifiedAt  = "modified_at"  // FieldUploadedAt holds the string denoting the uploaded_at vertex property in the database.
	FieldUploadedAt  = "uploaded_at"  // FieldContentType holds the string denoting the content_type vertex property in the database.
	FieldContentType = "content_type" // FieldStoreKey holds the string denoting the store_key vertex property in the database.
	FieldStoreKey    = "store_key"    // FieldCategory holds the string denoting the category vertex property in the database.
	FieldCategory    = "category"

	// EdgeLocation holds the string denoting the location edge name in mutations.
	EdgeLocation = "location"
	// EdgeEquipment holds the string denoting the equipment edge name in mutations.
	EdgeEquipment = "equipment"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeWorkOrder holds the string denoting the work_order edge name in mutations.
	EdgeWorkOrder = "work_order"
	// EdgeChecklistItem holds the string denoting the checklist_item edge name in mutations.
	EdgeChecklistItem = "checklist_item"
	// EdgeSurvey holds the string denoting the survey edge name in mutations.
	EdgeSurvey = "survey"
	// EdgeFloorPlan holds the string denoting the floor_plan edge name in mutations.
	EdgeFloorPlan = "floor_plan"
	// EdgePhotoSurveyQuestion holds the string denoting the photo_survey_question edge name in mutations.
	EdgePhotoSurveyQuestion = "photo_survey_question"
	// EdgeSurveyQuestion holds the string denoting the survey_question edge name in mutations.
	EdgeSurveyQuestion = "survey_question"

	// Table holds the table name of the file in the database.
	Table = "files"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "files"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "location_files"
	// EquipmentTable is the table the holds the equipment relation/edge.
	EquipmentTable = "files"
	// EquipmentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentInverseTable = "equipment"
	// EquipmentColumn is the table column denoting the equipment relation/edge.
	EquipmentColumn = "equipment_files"
	// UserTable is the table the holds the user relation/edge.
	UserTable = "files"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_profile_photo"
	// WorkOrderTable is the table the holds the work_order relation/edge.
	WorkOrderTable = "files"
	// WorkOrderInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrderInverseTable = "work_orders"
	// WorkOrderColumn is the table column denoting the work_order relation/edge.
	WorkOrderColumn = "work_order_files"
	// ChecklistItemTable is the table the holds the checklist_item relation/edge.
	ChecklistItemTable = "files"
	// ChecklistItemInverseTable is the table name for the CheckListItem entity.
	// It exists in this package in order to avoid circular dependency with the "checklistitem" package.
	ChecklistItemInverseTable = "check_list_items"
	// ChecklistItemColumn is the table column denoting the checklist_item relation/edge.
	ChecklistItemColumn = "check_list_item_files"
	// SurveyTable is the table the holds the survey relation/edge.
	SurveyTable = "files"
	// SurveyInverseTable is the table name for the Survey entity.
	// It exists in this package in order to avoid circular dependency with the "survey" package.
	SurveyInverseTable = "surveys"
	// SurveyColumn is the table column denoting the survey relation/edge.
	SurveyColumn = "survey_source_file"
	// FloorPlanTable is the table the holds the floor_plan relation/edge.
	FloorPlanTable = "files"
	// FloorPlanInverseTable is the table name for the FloorPlan entity.
	// It exists in this package in order to avoid circular dependency with the "floorplan" package.
	FloorPlanInverseTable = "floor_plans"
	// FloorPlanColumn is the table column denoting the floor_plan relation/edge.
	FloorPlanColumn = "floor_plan_image"
	// PhotoSurveyQuestionTable is the table the holds the photo_survey_question relation/edge.
	PhotoSurveyQuestionTable = "files"
	// PhotoSurveyQuestionInverseTable is the table name for the SurveyQuestion entity.
	// It exists in this package in order to avoid circular dependency with the "surveyquestion" package.
	PhotoSurveyQuestionInverseTable = "survey_questions"
	// PhotoSurveyQuestionColumn is the table column denoting the photo_survey_question relation/edge.
	PhotoSurveyQuestionColumn = "survey_question_photo_data"
	// SurveyQuestionTable is the table the holds the survey_question relation/edge.
	SurveyQuestionTable = "files"
	// SurveyQuestionInverseTable is the table name for the SurveyQuestion entity.
	// It exists in this package in order to avoid circular dependency with the "surveyquestion" package.
	SurveyQuestionInverseTable = "survey_questions"
	// SurveyQuestionColumn is the table column denoting the survey_question relation/edge.
	SurveyQuestionColumn = "survey_question_images"
)

// Columns holds all SQL columns for file fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldType,
	FieldName,
	FieldSize,
	FieldModifiedAt,
	FieldUploadedAt,
	FieldContentType,
	FieldStoreKey,
	FieldCategory,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the File type.
var ForeignKeys = []string{
	"check_list_item_files",
	"equipment_files",
	"floor_plan_image",
	"location_files",
	"survey_source_file",
	"survey_question_photo_data",
	"survey_question_images",
	"user_profile_photo",
	"work_order_files",
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
	// SizeValidator is a validator for the "size" field. It is called by the builders before save.
	SizeValidator func(int) error
)
