// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveyquestion

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the surveyquestion type in the database.
	Label = "survey_question"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldFormName holds the string denoting the form_name vertex property in the database.
	FieldFormName = "form_name"
	// FieldFormDescription holds the string denoting the form_description vertex property in the database.
	FieldFormDescription = "form_description"
	// FieldFormIndex holds the string denoting the form_index vertex property in the database.
	FieldFormIndex = "form_index"
	// FieldQuestionType holds the string denoting the question_type vertex property in the database.
	FieldQuestionType = "question_type"
	// FieldQuestionFormat holds the string denoting the question_format vertex property in the database.
	FieldQuestionFormat = "question_format"
	// FieldQuestionText holds the string denoting the question_text vertex property in the database.
	FieldQuestionText = "question_text"
	// FieldQuestionIndex holds the string denoting the question_index vertex property in the database.
	FieldQuestionIndex = "question_index"
	// FieldBoolData holds the string denoting the bool_data vertex property in the database.
	FieldBoolData = "bool_data"
	// FieldEmailData holds the string denoting the email_data vertex property in the database.
	FieldEmailData = "email_data"
	// FieldLatitude holds the string denoting the latitude vertex property in the database.
	FieldLatitude = "latitude"
	// FieldLongitude holds the string denoting the longitude vertex property in the database.
	FieldLongitude = "longitude"
	// FieldLocationAccuracy holds the string denoting the location_accuracy vertex property in the database.
	FieldLocationAccuracy = "location_accuracy"
	// FieldAltitude holds the string denoting the altitude vertex property in the database.
	FieldAltitude = "altitude"
	// FieldPhoneData holds the string denoting the phone_data vertex property in the database.
	FieldPhoneData = "phone_data"
	// FieldTextData holds the string denoting the text_data vertex property in the database.
	FieldTextData = "text_data"
	// FieldFloatData holds the string denoting the float_data vertex property in the database.
	FieldFloatData = "float_data"
	// FieldIntData holds the string denoting the int_data vertex property in the database.
	FieldIntData = "int_data"
	// FieldDateData holds the string denoting the date_data vertex property in the database.
	FieldDateData = "date_data"

	// Table holds the table name of the surveyquestion in the database.
	Table = "survey_questions"
	// SurveyTable is the table the holds the survey relation/edge.
	SurveyTable = "survey_questions"
	// SurveyInverseTable is the table name for the Survey entity.
	// It exists in this package in order to avoid circular dependency with the "survey" package.
	SurveyInverseTable = "surveys"
	// SurveyColumn is the table column denoting the survey relation/edge.
	SurveyColumn = "survey_id"
	// WifiScanTable is the table the holds the wifi_scan relation/edge.
	WifiScanTable = "survey_wi_fi_scans"
	// WifiScanInverseTable is the table name for the SurveyWiFiScan entity.
	// It exists in this package in order to avoid circular dependency with the "surveywifiscan" package.
	WifiScanInverseTable = "survey_wi_fi_scans"
	// WifiScanColumn is the table column denoting the wifi_scan relation/edge.
	WifiScanColumn = "survey_question_id"
	// CellScanTable is the table the holds the cell_scan relation/edge.
	CellScanTable = "survey_cell_scans"
	// CellScanInverseTable is the table name for the SurveyCellScan entity.
	// It exists in this package in order to avoid circular dependency with the "surveycellscan" package.
	CellScanInverseTable = "survey_cell_scans"
	// CellScanColumn is the table column denoting the cell_scan relation/edge.
	CellScanColumn = "survey_question_id"
	// PhotoDataTable is the table the holds the photo_data relation/edge.
	PhotoDataTable = "files"
	// PhotoDataInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	PhotoDataInverseTable = "files"
	// PhotoDataColumn is the table column denoting the photo_data relation/edge.
	PhotoDataColumn = "survey_question_photo_datum_id"
)

// Columns holds all SQL columns are surveyquestion fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldFormName,
	FieldFormDescription,
	FieldFormIndex,
	FieldQuestionType,
	FieldQuestionFormat,
	FieldQuestionText,
	FieldQuestionIndex,
	FieldBoolData,
	FieldEmailData,
	FieldLatitude,
	FieldLongitude,
	FieldLocationAccuracy,
	FieldAltitude,
	FieldPhoneData,
	FieldTextData,
	FieldFloatData,
	FieldIntData,
	FieldDateData,
}

var (
	mixin       = schema.SurveyQuestion{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.SurveyQuestion{}.Fields()

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
