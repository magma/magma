// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package migrate

import (
	"github.com/facebookincubator/symphony/frontier/ent/tenant"
	"github.com/facebookincubator/symphony/frontier/ent/user"

	"github.com/facebookincubator/ent/dialect/sql/schema"
	"github.com/facebookincubator/ent/schema/field"
)

var (
	// AuditLogEntriesColumns holds the columns for the "AuditLogEntries" table.
	AuditLogEntriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "createdAt", Type: field.TypeTime},
		{Name: "updatedAt", Type: field.TypeTime},
		{Name: "actingUserId", Type: field.TypeInt},
		{Name: "organization", Type: field.TypeString},
		{Name: "mutationType", Type: field.TypeString},
		{Name: "objectId", Type: field.TypeString},
		{Name: "objectType", Type: field.TypeString},
		{Name: "objectDisplayName", Type: field.TypeString},
		{Name: "mutationData", Type: field.TypeJSON},
		{Name: "url", Type: field.TypeString},
		{Name: "ipAddress", Type: field.TypeString},
		{Name: "status", Type: field.TypeString},
		{Name: "statusCode", Type: field.TypeString},
	}
	// AuditLogEntriesTable holds the schema information for the "AuditLogEntries" table.
	AuditLogEntriesTable = &schema.Table{
		Name:        "AuditLogEntries",
		Columns:     AuditLogEntriesColumns,
		PrimaryKey:  []*schema.Column{AuditLogEntriesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// OrganizationsColumns holds the columns for the "Organizations" table.
	OrganizationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "createdAt", Type: field.TypeTime},
		{Name: "updatedAt", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "customDomains", Type: field.TypeJSON},
		{Name: "networkIDs", Type: field.TypeJSON},
		{Name: "tabs", Type: field.TypeJSON, Nullable: true},
		{Name: "ssoCert", Type: field.TypeString, Size: 2147483647, Default: tenant.DefaultSSOCert},
		{Name: "ssoEntrypoint", Type: field.TypeString, Default: tenant.DefaultSSOEntryPoint},
		{Name: "ssoIssuer", Type: field.TypeString, Default: tenant.DefaultSSOIssuer},
	}
	// OrganizationsTable holds the schema information for the "Organizations" table.
	OrganizationsTable = &schema.Table{
		Name:        "Organizations",
		Columns:     OrganizationsColumns,
		PrimaryKey:  []*schema.Column{OrganizationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// TokensColumns holds the columns for the "tokens" table.
	TokensColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "createdAt", Type: field.TypeTime},
		{Name: "updatedAt", Type: field.TypeTime},
		{Name: "value", Type: field.TypeString},
		{Name: "user_tokens", Type: field.TypeInt, Nullable: true},
	}
	// TokensTable holds the schema information for the "tokens" table.
	TokensTable = &schema.Table{
		Name:       "tokens",
		Columns:    TokensColumns,
		PrimaryKey: []*schema.Column{TokensColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "tokens_Users_tokens",
				Columns: []*schema.Column{TokensColumns[4]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "token_value_user_tokens",
				Unique:  true,
				Columns: []*schema.Column{TokensColumns[3], TokensColumns[4]},
			},
		},
	}
	// UsersColumns holds the columns for the "Users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "createdAt", Type: field.TypeTime},
		{Name: "updatedAt", Type: field.TypeTime},
		{Name: "email", Type: field.TypeString},
		{Name: "password", Type: field.TypeString},
		{Name: "role", Type: field.TypeInt, Default: user.DefaultRole},
		{Name: "organization", Type: field.TypeString, Default: user.DefaultTenant},
		{Name: "networkIDs", Type: field.TypeJSON},
		{Name: "tabs", Type: field.TypeJSON, Nullable: true},
	}
	// UsersTable holds the schema information for the "Users" table.
	UsersTable = &schema.Table{
		Name:        "Users",
		Columns:     UsersColumns,
		PrimaryKey:  []*schema.Column{UsersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
		Indexes: []*schema.Index{
			{
				Name:    "user_email_organization",
				Unique:  true,
				Columns: []*schema.Column{UsersColumns[3], UsersColumns[6]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AuditLogEntriesTable,
		OrganizationsTable,
		TokensTable,
		UsersTable,
	}
)

func init() {
	TokensTable.ForeignKeys[0].RefTable = UsersTable
}
