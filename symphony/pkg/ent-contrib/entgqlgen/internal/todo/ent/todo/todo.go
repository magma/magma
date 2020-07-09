// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package todo

import (
	"fmt"
	"io"
)

const (
	// Label holds the string label denoting the todo type in the database.
	Label = "todo"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldText holds the string denoting the text field in the database.
	FieldText = "text"

	// EdgeParent holds the string denoting the parent edge name in mutations.
	EdgeParent = "parent"
	// EdgeChildren holds the string denoting the children edge name in mutations.
	EdgeChildren = "children"

	// Table holds the table name of the todo in the database.
	Table = "todos"
	// ParentTable is the table the holds the parent relation/edge.
	ParentTable = "todos"
	// ParentColumn is the table column denoting the parent relation/edge.
	ParentColumn = "todo_children"
	// ChildrenTable is the table the holds the children relation/edge.
	ChildrenTable = "todos"
	// ChildrenColumn is the table column denoting the children relation/edge.
	ChildrenColumn = "todo_children"
)

// Columns holds all SQL columns for todo fields.
var Columns = []string{
	FieldID,
	FieldStatus,
	FieldText,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the Todo type.
var ForeignKeys = []string{
	"todo_children",
}

var (
	// TextValidator is a validator for the "text" field. It is called by the builders before save.
	TextValidator func(string) error
)

// Status defines the type for the status enum field.
type Status string

// Status values.
const (
	StatusINPROGRESS Status = "IN_PROGRESS"
	StatusCOMPLETED  Status = "COMPLETED"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "s" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusINPROGRESS, StatusCOMPLETED:
		return nil
	default:
		return fmt.Errorf("todo: invalid enum value for status field: %q", s)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (s Status) MarshalGQL(w io.Writer) {
	writeQuotedStringer(w, s)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (s *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", v)
	}
	*s = Status(str)
	if err := StatusValidator(*s); err != nil {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func writeQuotedStringer(w io.Writer, s fmt.Stringer) {
	const quote = '"'
	switch w := w.(type) {
	case io.ByteWriter:
		w.WriteByte(quote)
		defer w.WriteByte(quote)
	default:
		w.Write([]byte{quote})
		defer w.Write([]byte{quote})
	}
	io.WriteString(w, s.String())
}
