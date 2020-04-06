// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package user

import (
	"time"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID        = "id"           // FieldCreatedAt holds the string denoting the created_at vertex property in the database.
	FieldCreatedAt = "createdAt"    // FieldUpdatedAt holds the string denoting the updated_at vertex property in the database.
	FieldUpdatedAt = "updatedAt"    // FieldEmail holds the string denoting the email vertex property in the database.
	FieldEmail     = "email"        // FieldPassword holds the string denoting the password vertex property in the database.
	FieldPassword  = "password"     // FieldRole holds the string denoting the role vertex property in the database.
	FieldRole      = "role"         // FieldTenant holds the string denoting the tenant vertex property in the database.
	FieldTenant    = "organization" // FieldNetworks holds the string denoting the networks vertex property in the database.
	FieldNetworks  = "networkIDs"   // FieldTabs holds the string denoting the tabs vertex property in the database.
	FieldTabs      = "tabs"

	// EdgeTokens holds the string denoting the tokens edge name in mutations.
	EdgeTokens = "tokens"

	// Table holds the table name of the user in the database.
	Table = "Users"
	// TokensTable is the table the holds the tokens relation/edge.
	TokensTable = "tokens"
	// TokensInverseTable is the table name for the Token entity.
	// It exists in this package in order to avoid circular dependency with the "token" package.
	TokensInverseTable = "tokens"
	// TokensColumn is the table column denoting the tokens relation/edge.
	TokensColumn = "user_tokens"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldEmail,
	FieldPassword,
	FieldRole,
	FieldTenant,
	FieldNetworks,
	FieldTabs,
}

var (
	// DefaultCreatedAt holds the default value on creation for the created_at field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the updated_at field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	UpdateDefaultUpdatedAt func() time.Time
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	PasswordValidator func(string) error
	// DefaultRole holds the default value on creation for the role field.
	DefaultRole int
	// RoleValidator is a validator for the "role" field. It is called by the builders before save.
	RoleValidator func(int) error
	// DefaultTenant holds the default value on creation for the tenant field.
	DefaultTenant string
)
