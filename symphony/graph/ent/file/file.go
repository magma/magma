// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package file

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the file type in the database.
	Label = "file"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldType holds the string denoting the type vertex property in the database.
	FieldType = "type"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldSize holds the string denoting the size vertex property in the database.
	FieldSize = "size"
	// FieldModifiedAt holds the string denoting the modified_at vertex property in the database.
	FieldModifiedAt = "modified_at"
	// FieldUploadedAt holds the string denoting the uploaded_at vertex property in the database.
	FieldUploadedAt = "uploaded_at"
	// FieldContentType holds the string denoting the content_type vertex property in the database.
	FieldContentType = "content_type"
	// FieldStoreKey holds the string denoting the store_key vertex property in the database.
	FieldStoreKey = "store_key"
	// FieldCategory holds the string denoting the category vertex property in the database.
	FieldCategory = "category"

	// Table holds the table name of the file in the database.
	Table = "files"
)

// Columns holds all SQL columns are file fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldType,
	FieldName,
	FieldSize,
	FieldModifiedAt,
	FieldUploadedAt,
	FieldContentType,
	FieldStoreKey,
	FieldCategory,
}

var (
	mixin       = schema.File{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.File{}.Fields()

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

	// descSize is the schema descriptor for size field.
	descSize = fields[2].Descriptor()
	// SizeValidator is a validator for the "size" field. It is called by the builders before save.
	SizeValidator = descSize.Validators[0].(func(int) error)
)
