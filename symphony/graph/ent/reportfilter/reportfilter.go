// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package reportfilter

import (
	"fmt"
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the reportfilter type in the database.
	Label = "report_filter"
	// FieldID holds the string denoting the id field in the database.
	FieldID         = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName       = "name"        // FieldEntity holds the string denoting the entity vertex property in the database.
	FieldEntity     = "entity"      // FieldFilters holds the string denoting the filters vertex property in the database.
	FieldFilters    = "filters"

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
//	import _ "github.com/facebookincubator/symphony/graph/ent/runtime"
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

func (s Entity) String() string {
	return string(s)
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
