// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package hyperlink

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the hyperlink type in the database.
	Label = "hyperlink"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldURL holds the string denoting the url vertex property in the database.
	FieldURL = "url"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldCategory holds the string denoting the category vertex property in the database.
	FieldCategory = "category"

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
	"equipment_hyperlink_id",
	"location_hyperlink_id",
	"work_order_hyperlink_id",
}

var (
	mixin       = schema.Hyperlink{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.Hyperlink{}.Fields()

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
