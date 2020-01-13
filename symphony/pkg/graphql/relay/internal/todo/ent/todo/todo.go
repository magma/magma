// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package todo

import (
	"github.com/facebookincubator/symphony/pkg/graphql/relay/internal/todo/ent/schema"
)

const (
	// Label holds the string label denoting the todo type in the database.
	Label = "todo"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldText holds the string denoting the text vertex property in the database.
	FieldText = "text"

	// Table holds the table name of the todo in the database.
	Table = "todos"
)

// Columns holds all SQL columns are todo fields.
var Columns = []string{
	FieldID,
	FieldText,
}

var (
	fields = schema.Todo{}.Fields()

	// descText is the schema descriptor for text field.
	descText = fields[0].Descriptor()
	// TextValidator is a validator for the "text" field. It is called by the builders before save.
	TextValidator = descText.Validators[0].(func(string) error)
)
