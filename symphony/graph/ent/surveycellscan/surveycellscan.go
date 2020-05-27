// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveycellscan

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the surveycellscan type in the database.
	Label = "survey_cell_scan"
	// FieldID holds the string denoting the id field in the database.
	FieldID                    = "id"                      // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime            = "create_time"             // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime            = "update_time"             // FieldNetworkType holds the string denoting the network_type vertex property in the database.
	FieldNetworkType           = "network_type"            // FieldSignalStrength holds the string denoting the signal_strength vertex property in the database.
	FieldSignalStrength        = "signal_strength"         // FieldTimestamp holds the string denoting the timestamp vertex property in the database.
	FieldTimestamp             = "timestamp"               // FieldBaseStationID holds the string denoting the base_station_id vertex property in the database.
	FieldBaseStationID         = "base_station_id"         // FieldNetworkID holds the string denoting the network_id vertex property in the database.
	FieldNetworkID             = "network_id"              // FieldSystemID holds the string denoting the system_id vertex property in the database.
	FieldSystemID              = "system_id"               // FieldCellID holds the string denoting the cell_id vertex property in the database.
	FieldCellID                = "cell_id"                 // FieldLocationAreaCode holds the string denoting the location_area_code vertex property in the database.
	FieldLocationAreaCode      = "location_area_code"      // FieldMobileCountryCode holds the string denoting the mobile_country_code vertex property in the database.
	FieldMobileCountryCode     = "mobile_country_code"     // FieldMobileNetworkCode holds the string denoting the mobile_network_code vertex property in the database.
	FieldMobileNetworkCode     = "mobile_network_code"     // FieldPrimaryScramblingCode holds the string denoting the primary_scrambling_code vertex property in the database.
	FieldPrimaryScramblingCode = "primary_scrambling_code" // FieldOperator holds the string denoting the operator vertex property in the database.
	FieldOperator              = "operator"                // FieldArfcn holds the string denoting the arfcn vertex property in the database.
	FieldArfcn                 = "arfcn"                   // FieldPhysicalCellID holds the string denoting the physical_cell_id vertex property in the database.
	FieldPhysicalCellID        = "physical_cell_id"        // FieldTrackingAreaCode holds the string denoting the tracking_area_code vertex property in the database.
	FieldTrackingAreaCode      = "tracking_area_code"      // FieldTimingAdvance holds the string denoting the timing_advance vertex property in the database.
	FieldTimingAdvance         = "timing_advance"          // FieldEarfcn holds the string denoting the earfcn vertex property in the database.
	FieldEarfcn                = "earfcn"                  // FieldUarfcn holds the string denoting the uarfcn vertex property in the database.
	FieldUarfcn                = "uarfcn"                  // FieldLatitude holds the string denoting the latitude vertex property in the database.
	FieldLatitude              = "latitude"                // FieldLongitude holds the string denoting the longitude vertex property in the database.
	FieldLongitude             = "longitude"

	// EdgeChecklistItem holds the string denoting the checklist_item edge name in mutations.
	EdgeChecklistItem = "checklist_item"
	// EdgeSurveyQuestion holds the string denoting the survey_question edge name in mutations.
	EdgeSurveyQuestion = "survey_question"
	// EdgeLocation holds the string denoting the location edge name in mutations.
	EdgeLocation = "location"

	// Table holds the table name of the surveycellscan in the database.
	Table = "survey_cell_scans"
	// ChecklistItemTable is the table the holds the checklist_item relation/edge.
	ChecklistItemTable = "survey_cell_scans"
	// ChecklistItemInverseTable is the table name for the CheckListItem entity.
	// It exists in this package in order to avoid circular dependency with the "checklistitem" package.
	ChecklistItemInverseTable = "check_list_items"
	// ChecklistItemColumn is the table column denoting the checklist_item relation/edge.
	ChecklistItemColumn = "survey_cell_scan_checklist_item"
	// SurveyQuestionTable is the table the holds the survey_question relation/edge.
	SurveyQuestionTable = "survey_cell_scans"
	// SurveyQuestionInverseTable is the table name for the SurveyQuestion entity.
	// It exists in this package in order to avoid circular dependency with the "surveyquestion" package.
	SurveyQuestionInverseTable = "survey_questions"
	// SurveyQuestionColumn is the table column denoting the survey_question relation/edge.
	SurveyQuestionColumn = "survey_cell_scan_survey_question"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "survey_cell_scans"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "survey_cell_scan_location"
)

// Columns holds all SQL columns for surveycellscan fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldNetworkType,
	FieldSignalStrength,
	FieldTimestamp,
	FieldBaseStationID,
	FieldNetworkID,
	FieldSystemID,
	FieldCellID,
	FieldLocationAreaCode,
	FieldMobileCountryCode,
	FieldMobileNetworkCode,
	FieldPrimaryScramblingCode,
	FieldOperator,
	FieldArfcn,
	FieldPhysicalCellID,
	FieldTrackingAreaCode,
	FieldTimingAdvance,
	FieldEarfcn,
	FieldUarfcn,
	FieldLatitude,
	FieldLongitude,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the SurveyCellScan type.
var ForeignKeys = []string{
	"survey_cell_scan_checklist_item",
	"survey_cell_scan_survey_question",
	"survey_cell_scan_location",
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
