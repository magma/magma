// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package locationtype

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the locationtype type in the database.
	Label = "location_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID           = "id"             // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime   = "create_time"    // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime   = "update_time"    // FieldSite holds the string denoting the site vertex property in the database.
	FieldSite         = "site"           // FieldName holds the string denoting the name vertex property in the database.
	FieldName         = "name"           // FieldMapType holds the string denoting the map_type vertex property in the database.
	FieldMapType      = "map_type"       // FieldMapZoomLevel holds the string denoting the map_zoom_level vertex property in the database.
	FieldMapZoomLevel = "map_zoom_level" // FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex        = "index"

	// EdgeLocations holds the string denoting the locations edge name in mutations.
	EdgeLocations = "locations"
	// EdgePropertyTypes holds the string denoting the property_types edge name in mutations.
	EdgePropertyTypes = "property_types"
	// EdgeSurveyTemplateCategories holds the string denoting the survey_template_categories edge name in mutations.
	EdgeSurveyTemplateCategories = "survey_template_categories"

	// Table holds the table name of the locationtype in the database.
	Table = "location_types"
	// LocationsTable is the table the holds the locations relation/edge.
	LocationsTable = "locations"
	// LocationsInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationsInverseTable = "locations"
	// LocationsColumn is the table column denoting the locations relation/edge.
	LocationsColumn = "location_type"
	// PropertyTypesTable is the table the holds the property_types relation/edge.
	PropertyTypesTable = "property_types"
	// PropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	PropertyTypesInverseTable = "property_types"
	// PropertyTypesColumn is the table column denoting the property_types relation/edge.
	PropertyTypesColumn = "location_type_property_types"
	// SurveyTemplateCategoriesTable is the table the holds the survey_template_categories relation/edge.
	SurveyTemplateCategoriesTable = "survey_template_categories"
	// SurveyTemplateCategoriesInverseTable is the table name for the SurveyTemplateCategory entity.
	// It exists in this package in order to avoid circular dependency with the "surveytemplatecategory" package.
	SurveyTemplateCategoriesInverseTable = "survey_template_categories"
	// SurveyTemplateCategoriesColumn is the table column denoting the survey_template_categories relation/edge.
	SurveyTemplateCategoriesColumn = "location_type_survey_template_categories"
)

// Columns holds all SQL columns for locationtype fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldSite,
	FieldName,
	FieldMapType,
	FieldMapZoomLevel,
	FieldIndex,
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/facebookincubator/symphony/pkg/ent/runtime"
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
	// DefaultSite holds the default value on creation for the site field.
	DefaultSite bool
	// DefaultMapZoomLevel holds the default value on creation for the map_zoom_level field.
	DefaultMapZoomLevel int
	// DefaultIndex holds the default value on creation for the index field.
	DefaultIndex int
)
