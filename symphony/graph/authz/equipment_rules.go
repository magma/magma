// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// EquipmentWritePolicyRule grants write permission to equipment based on policy.
func EquipmentWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return cudBasedRule(FromContext(ctx).InventoryPolicy.Equipment, m)
	})
}

// EquipmentTypeWritePolicyRule grants write permission to equipment type based on policy.
func EquipmentTypeWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return cudBasedRule(FromContext(ctx).InventoryPolicy.EquipmentType, m)
	})
}

// EquipmentPortTypeWritePolicyRule grants write permission to equipment port type based on policy.
func EquipmentPortTypeWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return cudBasedRule(FromContext(ctx).InventoryPolicy.PortType, m)
	})
}

// EquipmentPortDefinitionWritePolicyRule grants write permission to equipment port definition based on policy.
func EquipmentPortDefinitionWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).InventoryPolicy.EquipmentType.Update)
	})
}

// EquipmentCategoryWritePolicyRule grants write permission to equipment category based on policy.
func EquipmentCategoryWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).InventoryPolicy.EquipmentType.Update)
	})
}
