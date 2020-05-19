// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package floorplan

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the floorplan type in the database.
	Label = "floor_plan"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName       = "name"

	// EdgeLocation holds the string denoting the location edge name in mutations.
	EdgeLocation = "location"
	// EdgeReferencePoint holds the string denoting the reference_point edge name in mutations.
	EdgeReferencePoint = "reference_point"
	// EdgeScale holds the string denoting the scale edge name in mutations.
	EdgeScale = "scale"
	// EdgeImage holds the string denoting the image edge name in mutations.
	EdgeImage = "image"

	// Table holds the table name of the floorplan in the database.
	Table = "floor_plans"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "floor_plans"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "floor_plan_location"
	// ReferencePointTable is the table the holds the reference_point relation/edge.
	ReferencePointTable = "floor_plans"
	// ReferencePointInverseTable is the table name for the FloorPlanReferencePoint entity.
	// It exists in this package in order to avoid circular dependency with the "floorplanreferencepoint" package.
	ReferencePointInverseTable = "floor_plan_reference_points"
	// ReferencePointColumn is the table column denoting the reference_point relation/edge.
	ReferencePointColumn = "floor_plan_reference_point"
	// ScaleTable is the table the holds the scale relation/edge.
	ScaleTable = "floor_plans"
	// ScaleInverseTable is the table name for the FloorPlanScale entity.
	// It exists in this package in order to avoid circular dependency with the "floorplanscale" package.
	ScaleInverseTable = "floor_plan_scales"
	// ScaleColumn is the table column denoting the scale relation/edge.
	ScaleColumn = "floor_plan_scale"
	// ImageTable is the table the holds the image relation/edge.
	ImageTable = "files"
	// ImageInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	ImageInverseTable = "files"
	// ImageColumn is the table column denoting the image relation/edge.
	ImageColumn = "floor_plan_image"
)

// Columns holds all SQL columns for floorplan fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the FloorPlan type.
var ForeignKeys = []string{
	"floor_plan_location",
	"floor_plan_reference_point",
	"floor_plan_scale",
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
