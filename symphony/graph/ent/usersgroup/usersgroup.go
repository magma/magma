// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package usersgroup

import (
	"fmt"
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the usersgroup type in the database.
	Label = "users_group"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"          // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time" // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time" // FieldName holds the string denoting the name vertex property in the database.
	FieldName        = "name"        // FieldDescription holds the string denoting the description vertex property in the database.
	FieldDescription = "description" // FieldStatus holds the string denoting the status vertex property in the database.
	FieldStatus      = "status"

	// EdgeMembers holds the string denoting the members edge name in mutations.
	EdgeMembers = "members"
	// EdgePolicies holds the string denoting the policies edge name in mutations.
	EdgePolicies = "policies"

	// Table holds the table name of the usersgroup in the database.
	Table = "users_groups"
	// MembersTable is the table the holds the members relation/edge. The primary key declared below.
	MembersTable = "users_group_members"
	// MembersInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	MembersInverseTable = "users"
	// PoliciesTable is the table the holds the policies relation/edge. The primary key declared below.
	PoliciesTable = "users_group_policies"
	// PoliciesInverseTable is the table name for the PermissionsPolicy entity.
	// It exists in this package in order to avoid circular dependency with the "permissionspolicy" package.
	PoliciesInverseTable = "permissions_policies"
)

// Columns holds all SQL columns for usersgroup fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldDescription,
	FieldStatus,
}

var (
	// MembersPrimaryKey and MembersColumn2 are the table columns denoting the
	// primary key for the members relation (M2M).
	MembersPrimaryKey = []string{"users_group_id", "user_id"}
	// PoliciesPrimaryKey and PoliciesColumn2 are the table columns denoting the
	// primary key for the policies relation (M2M).
	PoliciesPrimaryKey = []string{"users_group_id", "permissions_policy_id"}
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
)

// Status defines the type for the status enum field.
type Status string

// StatusACTIVE is the default Status.
const DefaultStatus = StatusACTIVE

// Status values.
const (
	StatusACTIVE      Status = "ACTIVE"
	StatusDEACTIVATED Status = "DEACTIVATED"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "s" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusACTIVE, StatusDEACTIVATED:
		return nil
	default:
		return fmt.Errorf("usersgroup: invalid enum value for status field: %q", s)
	}
}
