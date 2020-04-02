// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package hyperlink

import (
	"time"
)

const (
	// Label holds the string label denoting the hyperlink type in the database.
	Label = "hyperlink"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldURL holds the string denoting the url vertex property in the database.
	FieldURL        = "url"         // FieldName holds the string denoting the name vertex property in the database.
	FieldName       = "name"        // FieldCategory holds the string denoting the category vertex property in the database.
	FieldCategory   = "category"

	// Table holds the table name of the hyperlink in the database.
	Table = "hyperlinks"
)

// Columns holds all SQL columns for hyperlink fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldURL,
	FieldName,
	FieldCategory,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Hyperlink type.
var ForeignKeys = []string{
	"equipment_hyperlinks",
	"location_hyperlinks",
	"work_order_hyperlinks",
}

var (
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime func() time.Time
)
