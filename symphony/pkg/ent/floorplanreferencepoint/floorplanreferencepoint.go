// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package floorplanreferencepoint

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the floorplanreferencepoint type in the database.
	Label = "floor_plan_reference_point"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldX holds the string denoting the x vertex property in the database.
	FieldX          = "x"           // FieldY holds the string denoting the y vertex property in the database.
	FieldY          = "y"           // FieldLatitude holds the string denoting the latitude vertex property in the database.
	FieldLatitude   = "latitude"    // FieldLongitude holds the string denoting the longitude vertex property in the database.
	FieldLongitude  = "longitude"

	// Table holds the table name of the floorplanreferencepoint in the database.
	Table = "floor_plan_reference_points"
)

// Columns holds all SQL columns for floorplanreferencepoint fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldX,
	FieldY,
	FieldLatitude,
	FieldLongitude,
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
