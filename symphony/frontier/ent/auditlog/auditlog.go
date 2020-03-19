// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package auditlog

import (
	"time"
)

const (
	// Label holds the string label denoting the auditlog type in the database.
	Label = "audit_log"
	// FieldID holds the string denoting the id field in the database.
	FieldID                = "id"                // FieldCreatedAt holds the string denoting the created_at vertex property in the database.
	FieldCreatedAt         = "createdAt"         // FieldUpdatedAt holds the string denoting the updated_at vertex property in the database.
	FieldUpdatedAt         = "updatedAt"         // FieldActingUserID holds the string denoting the acting_user_id vertex property in the database.
	FieldActingUserID      = "actingUserId"      // FieldOrganization holds the string denoting the organization vertex property in the database.
	FieldOrganization      = "organization"      // FieldMutationType holds the string denoting the mutation_type vertex property in the database.
	FieldMutationType      = "mutationType"      // FieldObjectID holds the string denoting the object_id vertex property in the database.
	FieldObjectID          = "objectId"          // FieldObjectType holds the string denoting the object_type vertex property in the database.
	FieldObjectType        = "objectType"        // FieldObjectDisplayName holds the string denoting the object_display_name vertex property in the database.
	FieldObjectDisplayName = "objectDisplayName" // FieldMutationData holds the string denoting the mutation_data vertex property in the database.
	FieldMutationData      = "mutationData"      // FieldURL holds the string denoting the url vertex property in the database.
	FieldURL               = "url"               // FieldIPAddress holds the string denoting the ip_address vertex property in the database.
	FieldIPAddress         = "ipAddress"         // FieldStatus holds the string denoting the status vertex property in the database.
	FieldStatus            = "status"            // FieldStatusCode holds the string denoting the status_code vertex property in the database.
	FieldStatusCode        = "statusCode"

	// Table holds the table name of the auditlog in the database.
	Table = "AuditLogEntries"
)

// Columns holds all SQL columns for auditlog fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldActingUserID,
	FieldOrganization,
	FieldMutationType,
	FieldObjectID,
	FieldObjectType,
	FieldObjectDisplayName,
	FieldMutationData,
	FieldURL,
	FieldIPAddress,
	FieldStatus,
	FieldStatusCode,
}

var (
	// DefaultCreatedAt holds the default value on creation for the created_at field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the updated_at field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	UpdateDefaultUpdatedAt func() time.Time
)
