// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package reportfilter

import (
	"fmt"
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the reportfilter type in the database.
	Label = "report_filter"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name vertex property in the database.
	FieldName = "name"
	// FieldEntity holds the string denoting the entity vertex property in the database.
	FieldEntity = "entity"
	// FieldFilters holds the string denoting the filters vertex property in the database.
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

var (
	mixin       = schema.ReportFilter{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.ReportFilter{}.Fields()

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

	// descName is the schema descriptor for name field.
	descName = fields[0].Descriptor()
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator = descName.Validators[0].(func(string) error)

	// descFilters is the schema descriptor for filters field.
	descFilters = fields[2].Descriptor()
	// DefaultFilters holds the default value on creation for the filters field.
	DefaultFilters = descFilters.Default.(string)
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
