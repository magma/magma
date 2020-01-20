// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// AuditLog defines audit log schema.
type AuditLog struct {
	schema
}

// Config configures audit log schema.
func (AuditLog) Config() ent.Config {
	return ent.Config{
		Table: "AuditLogEntries",
	}
}

// Fields of audit log entity.
func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("acting_user_id").
			StorageKey("actingUserId"),
		field.String("organization"),
		field.String("mutation_type").
			StorageKey("mutationType"),
		field.String("object_id").
			StorageKey("objectId"),
		field.String("object_type").
			StorageKey("objectType"),
		field.String("object_display_name").
			StorageKey("objectDisplayName"),
		field.JSON("mutation_data", map[string]string{}).
			StorageKey("mutationData"),
		field.String("url"),
		field.String("ip_address").
			StorageKey("ipAddress"),
		field.String("status"),
		field.String("status_code").
			StorageKey("statusCode"),
	}
}
