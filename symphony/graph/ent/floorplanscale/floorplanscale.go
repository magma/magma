// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package floorplanscale

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the floorplanscale type in the database.
	Label = "floor_plan_scale"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldReferencePoint1X holds the string denoting the reference_point1_x vertex property in the database.
	FieldReferencePoint1X = "reference_point1_x"
	// FieldReferencePoint1Y holds the string denoting the reference_point1_y vertex property in the database.
	FieldReferencePoint1Y = "reference_point1_y"
	// FieldReferencePoint2X holds the string denoting the reference_point2_x vertex property in the database.
	FieldReferencePoint2X = "reference_point2_x"
	// FieldReferencePoint2Y holds the string denoting the reference_point2_y vertex property in the database.
	FieldReferencePoint2Y = "reference_point2_y"
	// FieldScaleInMeters holds the string denoting the scale_in_meters vertex property in the database.
	FieldScaleInMeters = "scale_in_meters"

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

var (
	mixin       = schema.FloorPlanScale{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.FloorPlanScale{}.Fields()

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
