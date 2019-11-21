// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package locationtype

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the locationtype type in the database.
	Label = "location_type"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldSite holds the string denoting the site vertex property in the database.
	FieldSite = "site"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldMapType holds the string denoting the map_type vertex property in the database.
	FieldMapType = "map_type"
	// FieldMapZoomLevel holds the string denoting the map_zoom_level vertex property in the database.
	FieldMapZoomLevel = "map_zoom_level"
	// FieldIndex holds the string denoting the index vertex property in the database.
	FieldIndex = "index"

	// Table holds the table name of the locationtype in the database.
	Table = "location_types"
	// LocationsTable is the table the holds the locations relation/edge.
	LocationsTable = "locations"
	// LocationsInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationsInverseTable = "locations"
	// LocationsColumn is the table column denoting the locations relation/edge.
	LocationsColumn = "type_id"
	// PropertyTypesTable is the table the holds the property_types relation/edge.
	PropertyTypesTable = "property_types"
	// PropertyTypesInverseTable is the table name for the PropertyType entity.
	// It exists in this package in order to avoid circular dependency with the "propertytype" package.
	PropertyTypesInverseTable = "property_types"
	// PropertyTypesColumn is the table column denoting the property_types relation/edge.
	PropertyTypesColumn = "location_type_id"
	// SurveyTemplateCategoriesTable is the table the holds the survey_template_categories relation/edge.
	SurveyTemplateCategoriesTable = "survey_template_categories"
	// SurveyTemplateCategoriesInverseTable is the table name for the SurveyTemplateCategory entity.
	// It exists in this package in order to avoid circular dependency with the "surveytemplatecategory" package.
	SurveyTemplateCategoriesInverseTable = "survey_template_categories"
	// SurveyTemplateCategoriesColumn is the table column denoting the survey_template_categories relation/edge.
	SurveyTemplateCategoriesColumn = "location_type_survey_template_category_id"
)

// Columns holds all SQL columns are locationtype fields.
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

var (
	mixin       = schema.LocationType{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.LocationType{}.Fields()

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

	// descSite is the schema descriptor for site field.
	descSite = fields[0].Descriptor()
	// DefaultSite holds the default value on creation for the site field.
	DefaultSite = descSite.Default.(bool)

	// descMapZoomLevel is the schema descriptor for map_zoom_level field.
	descMapZoomLevel = fields[3].Descriptor()
	// DefaultMapZoomLevel holds the default value on creation for the map_zoom_level field.
	DefaultMapZoomLevel = descMapZoomLevel.Default.(int)

	// descIndex is the schema descriptor for index field.
	descIndex = fields[4].Descriptor()
	// DefaultIndex holds the default value on creation for the index field.
	DefaultIndex = descIndex.Default.(int)
)
