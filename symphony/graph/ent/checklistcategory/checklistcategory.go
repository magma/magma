// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package checklistcategory

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the checklistcategory type in the database.
	Label = "check_list_category"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldTitle holds the string denoting the title vertex property in the database.
	FieldTitle = "title"
	// FieldDescription holds the string denoting the description vertex property in the database.
	FieldDescription = "description"

	// Table holds the table name of the checklistcategory in the database.
	Table = "check_list_categories"
	// CheckListItemsTable is the table the holds the check_list_items relation/edge.
	CheckListItemsTable = "check_list_items"
	// CheckListItemsInverseTable is the table name for the CheckListItem entity.
	// It exists in this package in order to avoid circular dependency with the "checklistitem" package.
	CheckListItemsInverseTable = "check_list_items"
	// CheckListItemsColumn is the table column denoting the check_list_items relation/edge.
	CheckListItemsColumn = "check_list_category_check_list_items"
)

// Columns holds all SQL columns for checklistcategory fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldTitle,
	FieldDescription,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the CheckListCategory type.
var ForeignKeys = []string{
	"work_order_check_list_categories",
	"work_order_type_check_list_categories",
}

var (
	mixin       = schema.CheckListCategory{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.CheckListCategory{}.Fields()

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
