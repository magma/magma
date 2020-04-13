// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveywifiscan

import (
	"time"
)

const (
	// Label holds the string label denoting the surveywifiscan type in the database.
	Label = "survey_wi_fi_scan"
	// FieldID holds the string denoting the id field in the database.
	FieldID           = "id"            // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime   = "create_time"   // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime   = "update_time"   // FieldSsid holds the string denoting the ssid vertex property in the database.
	FieldSsid         = "ssid"          // FieldBssid holds the string denoting the bssid vertex property in the database.
	FieldBssid        = "bssid"         // FieldTimestamp holds the string denoting the timestamp vertex property in the database.
	FieldTimestamp    = "timestamp"     // FieldFrequency holds the string denoting the frequency vertex property in the database.
	FieldFrequency    = "frequency"     // FieldChannel holds the string denoting the channel vertex property in the database.
	FieldChannel      = "channel"       // FieldBand holds the string denoting the band vertex property in the database.
	FieldBand         = "band"          // FieldChannelWidth holds the string denoting the channel_width vertex property in the database.
	FieldChannelWidth = "channel_width" // FieldCapabilities holds the string denoting the capabilities vertex property in the database.
	FieldCapabilities = "capabilities"  // FieldStrength holds the string denoting the strength vertex property in the database.
	FieldStrength     = "strength"      // FieldLatitude holds the string denoting the latitude vertex property in the database.
	FieldLatitude     = "latitude"      // FieldLongitude holds the string denoting the longitude vertex property in the database.
	FieldLongitude    = "longitude"

	// EdgeChecklistItem holds the string denoting the checklist_item edge name in mutations.
	EdgeChecklistItem = "checklist_item"
	// EdgeSurveyQuestion holds the string denoting the survey_question edge name in mutations.
	EdgeSurveyQuestion = "survey_question"
	// EdgeLocation holds the string denoting the location edge name in mutations.
	EdgeLocation = "location"

	// Table holds the table name of the surveywifiscan in the database.
	Table = "survey_wi_fi_scans"
	// ChecklistItemTable is the table the holds the checklist_item relation/edge.
	ChecklistItemTable = "survey_wi_fi_scans"
	// ChecklistItemInverseTable is the table name for the CheckListItem entity.
	// It exists in this package in order to avoid circular dependency with the "checklistitem" package.
	ChecklistItemInverseTable = "check_list_items"
	// ChecklistItemColumn is the table column denoting the checklist_item relation/edge.
	ChecklistItemColumn = "survey_wi_fi_scan_checklist_item"
	// SurveyQuestionTable is the table the holds the survey_question relation/edge.
	SurveyQuestionTable = "survey_wi_fi_scans"
	// SurveyQuestionInverseTable is the table name for the SurveyQuestion entity.
	// It exists in this package in order to avoid circular dependency with the "surveyquestion" package.
	SurveyQuestionInverseTable = "survey_questions"
	// SurveyQuestionColumn is the table column denoting the survey_question relation/edge.
	SurveyQuestionColumn = "survey_wi_fi_scan_survey_question"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "survey_wi_fi_scans"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "survey_wi_fi_scan_location"
)

// Columns holds all SQL columns for surveywifiscan fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldSsid,
	FieldBssid,
	FieldTimestamp,
	FieldFrequency,
	FieldChannel,
	FieldBand,
	FieldChannelWidth,
	FieldCapabilities,
	FieldStrength,
	FieldLatitude,
	FieldLongitude,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the SurveyWiFiScan type.
var ForeignKeys = []string{
	"survey_wi_fi_scan_checklist_item",
	"survey_wi_fi_scan_survey_question",
	"survey_wi_fi_scan_location",
}

var (
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime func() time.Time
)
