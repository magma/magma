// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package user

import (
	"fmt"
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldAuthID holds the string denoting the auth_id vertex property in the database.
	FieldAuthID = "auth_id"
	// FieldFirstName holds the string denoting the first_name vertex property in the database.
	FieldFirstName = "first_name"
	// FieldLastName holds the string denoting the last_name vertex property in the database.
	FieldLastName = "last_name"
	// FieldEmail holds the string denoting the email vertex property in the database.
	FieldEmail = "email"
	// FieldStatus holds the string denoting the status vertex property in the database.
	FieldStatus = "status"
	// FieldRole holds the string denoting the role vertex property in the database.
	FieldRole = "role"

	// Table holds the table name of the user in the database.
	Table = "users"
	// ProfilePhotoTable is the table the holds the profile_photo relation/edge.
	ProfilePhotoTable = "users"
	// ProfilePhotoInverseTable is the table name for the File entity.
	// It exists in this package in order to avoid circular dependency with the "file" package.
	ProfilePhotoInverseTable = "files"
	// ProfilePhotoColumn is the table column denoting the profile_photo relation/edge.
	ProfilePhotoColumn = "user_profile_photo"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldAuthID,
	FieldFirstName,
	FieldLastName,
	FieldEmail,
	FieldStatus,
	FieldRole,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the User type.
var ForeignKeys = []string{
	"user_profile_photo",
}

var (
	mixin       = schema.User{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.User{}.Fields()

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

	// descAuthID is the schema descriptor for auth_id field.
	descAuthID = fields[0].Descriptor()
	// AuthIDValidator is a validator for the "auth_id" field. It is called by the builders before save.
	AuthIDValidator = descAuthID.Validators[0].(func(string) error)

	// descFirstName is the schema descriptor for first_name field.
	descFirstName = fields[1].Descriptor()
	// FirstNameValidator is a validator for the "first_name" field. It is called by the builders before save.
	FirstNameValidator = descFirstName.Validators[0].(func(string) error)

	// descLastName is the schema descriptor for last_name field.
	descLastName = fields[2].Descriptor()
	// LastNameValidator is a validator for the "last_name" field. It is called by the builders before save.
	LastNameValidator = descLastName.Validators[0].(func(string) error)

	// descEmail is the schema descriptor for email field.
	descEmail = fields[3].Descriptor()
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator = descEmail.Validators[0].(func(string) error)
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
		return fmt.Errorf("user: invalid enum value for status field: %q", s)
	}
}

// Role defines the type for the role enum field.
type Role string

// RoleUSER is the default Role.
const DefaultRole = RoleUSER

// Role values.
const (
	RoleUSER  Role = "USER"
	RoleADMIN Role = "ADMIN"
	RoleOWNER Role = "OWNER"
)

func (s Role) String() string {
	return string(s)
}

// RoleValidator is a validator for the "r" field enum values. It is called by the builders before save.
func RoleValidator(r Role) error {
	switch r {
	case RoleUSER, RoleADMIN, RoleOWNER:
		return nil
	default:
		return fmt.Errorf("user: invalid enum value for role field: %q", r)
	}
}
