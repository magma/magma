// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package reportfilter

import (
	"fmt"
	"io"
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the reportfilter type in the database.
	Label = "report_filter"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldEntity holds the string denoting the entity field in the database.
	FieldEntity = "entity"
	// FieldFilters holds the string denoting the filters field in the database.
	FieldFilters = "filters"

	// Table holds the table name of the reportfilter in the database.
	Table = "report_filters"
)

// Columns holds all SQL columns for reportfilter fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldEntity,
	FieldFilters,
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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultFilters holds the default value on creation for the filters field.
	DefaultFilters string
)

// Entity defines the type for the entity enum field.
type Entity string

// Entity values.
const (
	EntityWORKORDER Entity = "WORK_ORDER"
	EntityPORT      Entity = "PORT"
	EntityEQUIPMENT Entity = "EQUIPMENT"
	EntityLINK      Entity = "LINK"
	EntityLOCATION  Entity = "LOCATION"
	EntitySERVICE   Entity = "SERVICE"
)

func (e Entity) String() string {
	return string(e)
}

// EntityValidator is a validator for the "e" field enum values. It is called by the builders before save.
func EntityValidator(e Entity) error {
	switch e {
	case EntityWORKORDER, EntityPORT, EntityEQUIPMENT, EntityLINK, EntityLOCATION, EntitySERVICE:
		return nil
	default:
		return fmt.Errorf("reportfilter: invalid enum value for entity field: %q", e)
	}
}

// MarshalGQL implements graphql.Marshaler interface.
func (e Entity) MarshalGQL(w io.Writer) {
	writeQuotedStringer(w, e)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (e *Entity) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", v)
	}
	*e = Entity(str)
	if err := EntityValidator(*e); err != nil {
		return fmt.Errorf("%s is not a valid Entity", str)
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
