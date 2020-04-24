// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package permissionspolicy

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the permissionspolicy type in the database.
	Label = "permissions_policy"
	// FieldID holds the string denoting the id field in the database.
	FieldID              = "id"               // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime      = "create_time"      // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime      = "update_time"      // FieldName holds the string denoting the name vertex property in the database.
	FieldName            = "name"             // FieldDescription holds the string denoting the description vertex property in the database.
	FieldDescription     = "description"      // FieldIsGlobal holds the string denoting the is_global vertex property in the database.
	FieldIsGlobal        = "is_global"        // FieldInventoryPolicy holds the string denoting the inventory_policy vertex property in the database.
	FieldInventoryPolicy = "inventory_policy" // FieldWorkforcePolicy holds the string denoting the workforce_policy vertex property in the database.
	FieldWorkforcePolicy = "workforce_policy"

	// EdgeGroups holds the string denoting the groups edge name in mutations.
	EdgeGroups = "groups"

	// Table holds the table name of the permissionspolicy in the database.
	Table = "permissions_policies"
	// GroupsTable is the table the holds the groups relation/edge. The primary key declared below.
	GroupsTable = "users_group_policies"
	// GroupsInverseTable is the table name for the UsersGroup entity.
	// It exists in this package in order to avoid circular dependency with the "usersgroup" package.
	GroupsInverseTable = "users_groups"
)

// Columns holds all SQL columns for permissionspolicy fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldDescription,
	FieldIsGlobal,
	FieldInventoryPolicy,
	FieldWorkforcePolicy,
}

var (
	// GroupsPrimaryKey and GroupsColumn2 are the table columns denoting the
	// primary key for the groups relation (M2M).
	GroupsPrimaryKey = []string{"users_group_id", "permissions_policy_id"}
)

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
	// DefaultIsGlobal holds the default value on creation for the is_global field.
	DefaultIsGlobal bool
)
