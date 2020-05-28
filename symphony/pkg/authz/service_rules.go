// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	models2 "github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

// ServiceTypeWritePolicyRule grants write permission to service type based on policy.
func ServiceTypeWritePolicyRule() privacy.MutationRule {
	return privacy.ServiceTypeMutationRuleFunc(func(ctx context.Context, m *ent.ServiceTypeMutation) error {
		deleted, deletedExists := m.IsDeleted()
		allowed := true
		switch {
		case m.Op().Is(ent.OpCreate):
			allowed = allowed && (FromContext(ctx).InventoryPolicy.ServiceType.Create.IsAllowed == models2.PermissionValueYes)
		case m.Op().Is(ent.OpUpdateOne | ent.OpUpdate):
			allowed = allowed && (FromContext(ctx).InventoryPolicy.ServiceType.Update.IsAllowed == models2.PermissionValueYes)
		case m.Op().Is(ent.OpDeleteOne | ent.OpDelete):
			allowed = allowed && (FromContext(ctx).InventoryPolicy.ServiceType.Delete.IsAllowed == models2.PermissionValueYes)
		}
		if deletedExists && deleted {
			allowed = allowed && (FromContext(ctx).InventoryPolicy.ServiceType.Delete.IsAllowed == models2.PermissionValueYes)
		}
		return privacyDecision(allowed)
	})
}

// ServiceWritePolicyRule grants write permission to service based on policy.
func ServiceWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return cudBasedRule(FromContext(ctx).InventoryPolicy.Equipment, m)
	})
}

// ServiceEndpointWritePolicyRule grants write permission to service endpoint based on policy.
func ServiceEndpointWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).InventoryPolicy.Equipment.Update)
	})
}

// ServiceEndpointDefinitionWritePolicyRule grants write permission to service endpoint definition based on policy.
func ServiceEndpointDefinitionWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).InventoryPolicy.ServiceType.Update)
	})
}
