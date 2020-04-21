// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/authz/models"
)

// PermissionsPolicy defines the policy schema.
type PermissionsPolicy struct {
	schema
}

// Fields returns policy fields.
func (PermissionsPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("description").
			Optional(),
		field.Bool("is_global").
			Optional().
			Default(false),
		field.JSON("inventory_policy", &models.InventoryPolicyInput{}).
			Optional(),
		field.JSON("workforce_policy", &models.WorkforcePolicyInput{}).
			Optional(),
	}
}
