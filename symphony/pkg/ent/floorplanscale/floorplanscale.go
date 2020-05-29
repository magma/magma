// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package floorplanscale

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the floorplanscale type in the database.
	Label = "floor_plan_scale"
	// FieldID holds the string denoting the id field in the database.
	FieldID               = "id"                 // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime       = "create_time"        // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime       = "update_time"        // FieldReferencePoint1X holds the string denoting the reference_point1_x vertex property in the database.
	FieldReferencePoint1X = "reference_point1_x" // FieldReferencePoint1Y holds the string denoting the reference_point1_y vertex property in the database.
	FieldReferencePoint1Y = "reference_point1_y" // FieldReferencePoint2X holds the string denoting the reference_point2_x vertex property in the database.
	FieldReferencePoint2X = "reference_point2_x" // FieldReferencePoint2Y holds the string denoting the reference_point2_y vertex property in the database.
	FieldReferencePoint2Y = "reference_point2_y" // FieldScaleInMeters holds the string denoting the scale_in_meters vertex property in the database.
	FieldScaleInMeters    = "scale_in_meters"

	// Table holds the table name of the floorplanscale in the database.
	Table = "floor_plan_scales"
)

// Columns holds all SQL columns for floorplanscale fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldReferencePoint1X,
	FieldReferencePoint1Y,
	FieldReferencePoint2X,
	FieldReferencePoint2Y,
	FieldScaleInMeters,
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
)
