// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package file

import (
	"time"
)

const (
	// Label holds the string label denoting the file type in the database.
	Label = "file"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"           // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time"  // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time"  // FieldType holds the string denoting the type vertex property in the database.
	FieldType        = "type"         // FieldName holds the string denoting the name vertex property in the database.
	FieldName        = "name"         // FieldSize holds the string denoting the size vertex property in the database.
	FieldSize        = "size"         // FieldModifiedAt holds the string denoting the modified_at vertex property in the database.
	FieldModifiedAt  = "modified_at"  // FieldUploadedAt holds the string denoting the uploaded_at vertex property in the database.
	FieldUploadedAt  = "uploaded_at"  // FieldContentType holds the string denoting the content_type vertex property in the database.
	FieldContentType = "content_type" // FieldStoreKey holds the string denoting the store_key vertex property in the database.
	FieldStoreKey    = "store_key"    // FieldCategory holds the string denoting the category vertex property in the database.
	FieldCategory    = "category"

	// Table holds the table name of the file in the database.
	Table = "files"
)

// Columns holds all SQL columns for file fields.
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

// ForeignKeys holds the SQL foreign-keys that are owned by the File type.
var ForeignKeys = []string{
	"check_list_item_files",
	"equipment_files",
	"location_files",
	"survey_question_photo_data",
	"survey_question_images",
	"work_order_files",
}

var (
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime func() time.Time
	// SizeValidator is a validator for the "size" field. It is called by the builders before save.
	SizeValidator func(int) error
)
