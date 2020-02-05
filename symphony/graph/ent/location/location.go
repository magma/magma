// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package location

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the location type in the database.
	Label = "location"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldExternalID holds the string denoting the external_id vertex property in the database.
	FieldExternalID = "external_id"
	// FieldLatitude holds the string denoting the latitude vertex property in the database.
	FieldLatitude = "latitude"
	// FieldLongitude holds the string denoting the longitude vertex property in the database.
	FieldLongitude = "longitude"
	// FieldSiteSurveyNeeded holds the string denoting the site_survey_needed vertex property in the database.
	FieldSiteSurveyNeeded = "site_survey_needed"

	// Table holds the table name of the location in the database.
	Table = "locations"
	// TypeTable is the table the holds the type relation/edge.
	TypeTable = "locations"
	// TypeInverseTable is the table name for the LocationType entity.
	// It exists in this package in order to avoid circular dependency with the "locationtype" package.
	TypeInverseTable = "location_types"
	// TypeColumn is the table column denoting the type relation/edge.
	TypeColumn = "location_type"
	// ParentTable is the table the holds the parent relation/edge.
	ParentTable = "locations"
	// ParentColumn is the table column denoting the parent relation/edge.
	ParentColumn = "location_children"
	// ChildrenTable is the table the holds the children relation/edge.
	ChildrenTable = "locations"
	// ChildrenColumn is the table column denoting the children relation/edge.
	ChildrenColumn = "location_children"
	// FilesTable is the table the holds the files relation/edge.
	FilesTable = "files"
	// FilesInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	FilesInverseTable = "files"
	// FilesColumn is the table column denoting the files relation/edge.
	FilesColumn = "location_files"
	// HyperlinksTable is the table the holds the hyperlinks relation/edge.
	HyperlinksTable = "hyperlinks"
	// HyperlinksInverseTable is the table name for the Hyperlink entity.
	// It exists in this package in order to avoid circular dependency with the "hyperlink" package.
	HyperlinksInverseTable = "hyperlinks"
	// HyperlinksColumn is the table column denoting the hyperlinks relation/edge.
	HyperlinksColumn = "location_hyperlinks"
	// EquipmentTable is the table the holds the equipment relation/edge.
	EquipmentTable = "equipment"
	// EquipmentInverseTable is the table name for the Equipment entity.
	// It exists in this package in order to avoid circular dependency with the "equipment" package.
	EquipmentInverseTable = "equipment"
	// EquipmentColumn is the table column denoting the equipment relation/edge.
	EquipmentColumn = "location_equipment"
	// PropertiesTable is the table the holds the properties relation/edge.
	PropertiesTable = "properties"
	// PropertiesInverseTable is the table name for the Property entity.
	// It exists in this package in order to avoid circular dependency with the "property" package.
	PropertiesInverseTable = "properties"
	// PropertiesColumn is the table column denoting the properties relation/edge.
	PropertiesColumn = "location_properties"
	// SurveyTable is the table the holds the survey relation/edge.
	SurveyTable = "surveys"
	// SurveyInverseTable is the table name for the Survey entity.
	// It exists in this package in order to avoid circular dependency with the "survey" package.
	SurveyInverseTable = "surveys"
	// SurveyColumn is the table column denoting the survey relation/edge.
	SurveyColumn = "survey_location"
	// WifiScanTable is the table the holds the wifi_scan relation/edge.
	WifiScanTable = "survey_wi_fi_scans"
	// WifiScanInverseTable is the table name for the SurveyWiFiScan entity.
	// It exists in this package in order to avoid circular dependency with the "surveywifiscan" package.
	WifiScanInverseTable = "survey_wi_fi_scans"
	// WifiScanColumn is the table column denoting the wifi_scan relation/edge.
	WifiScanColumn = "survey_wi_fi_scan_location"
	// CellScanTable is the table the holds the cell_scan relation/edge.
	CellScanTable = "survey_cell_scans"
	// CellScanInverseTable is the table name for the SurveyCellScan entity.
	// It exists in this package in order to avoid circular dependency with the "surveycellscan" package.
	CellScanInverseTable = "survey_cell_scans"
	// CellScanColumn is the table column denoting the cell_scan relation/edge.
	CellScanColumn = "survey_cell_scan_location"
	// WorkOrdersTable is the table the holds the work_orders relation/edge.
	WorkOrdersTable = "work_orders"
	// WorkOrdersInverseTable is the table name for the WorkOrder entity.
	// It exists in this package in order to avoid circular dependency with the "workorder" package.
	WorkOrdersInverseTable = "work_orders"
	// WorkOrdersColumn is the table column denoting the work_orders relation/edge.
	WorkOrdersColumn = "work_order_location"
	// FloorPlansTable is the table the holds the floor_plans relation/edge.
	FloorPlansTable = "floor_plans"
	// FloorPlansInverseTable is the table name for the FloorPlan entity.
	// It exists in this package in order to avoid circular dependency with the "floorplan" package.
	FloorPlansInverseTable = "floor_plans"
	// FloorPlansColumn is the table column denoting the floor_plans relation/edge.
	FloorPlansColumn = "floor_plan_location"
)

// Columns holds all SQL columns for location fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldExternalID,
	FieldLatitude,
	FieldLongitude,
	FieldSiteSurveyNeeded,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Location type.
var ForeignKeys = []string{
	"location_type",
	"location_children",
}

var (
	mixin       = schema.Location{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.Location{}.Fields()

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

	// descName is the schema descriptor for name field.
	descName = fields[0].Descriptor()
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator = descName.Validators[0].(func(string) error)

	// descLatitude is the schema descriptor for latitude field.
	descLatitude = fields[2].Descriptor()
	// DefaultLatitude holds the default value on creation for the latitude field.
	DefaultLatitude = descLatitude.Default.(float64)
	// LatitudeValidator is a validator for the "latitude" field. It is called by the builders before save.
	LatitudeValidator = descLatitude.Validators[0].(func(float64) error)

	// descLongitude is the schema descriptor for longitude field.
	descLongitude = fields[3].Descriptor()
	// DefaultLongitude holds the default value on creation for the longitude field.
	DefaultLongitude = descLongitude.Default.(float64)
	// LongitudeValidator is a validator for the "longitude" field. It is called by the builders before save.
	LongitudeValidator = descLongitude.Validators[0].(func(float64) error)

	// descSiteSurveyNeeded is the schema descriptor for site_survey_needed field.
	descSiteSurveyNeeded = fields[4].Descriptor()
	// DefaultSiteSurveyNeeded holds the default value on creation for the site_survey_needed field.
	DefaultSiteSurveyNeeded = descSiteSurveyNeeded.Default.(bool)
)
