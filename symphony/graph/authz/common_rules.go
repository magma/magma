// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
)

func mutationWithPermissionRule(rule func(context.Context, ent.Mutation, *models.PermissionSettings) error) privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		p := FromContext(ctx)
		if p == nil {
			return privacy.Skip
		}
		return rule(ctx, m, p)
	})
}

func workOrderMutationWithPermissionRule(rule func(context.Context, *ent.WorkOrderMutation, *models.PermissionSettings) error) privacy.MutationRule {
	return privacy.WorkOrderMutationRuleFunc(func(ctx context.Context, m *ent.WorkOrderMutation) error {
		p := FromContext(ctx)
		if p == nil {
			return privacy.Skip
		}
		return rule(ctx, m, p)
	})
}

func projectMutationWithPermissionRule(rule func(context.Context, *ent.ProjectMutation, *models.PermissionSettings) error) privacy.MutationRule {
	return privacy.ProjectMutationRuleFunc(func(ctx context.Context, m *ent.ProjectMutation) error {
		p := FromContext(ctx)
		if p == nil {
			return privacy.Skip
		}
		return rule(ctx, m, p)
	})
}

func cudBasedCheck(cud *models.Cud, m ent.Mutation) bool {
	if m.Op().Is(ent.OpDeleteOne) || m.Op().Is(ent.OpDelete) {
		return cud.Delete.IsAllowed == models2.PermissionValueYes
	}
	if m.Op().Is(ent.OpUpdateOne) || m.Op().Is(ent.OpUpdate) {
		return cud.Update.IsAllowed == models2.PermissionValueYes
	}
	if m.Op().Is(ent.OpCreate) {
		return cud.Create.IsAllowed == models2.PermissionValueYes
	}
	return false
}

func cudBasedRule(cud *models.Cud, m ent.Mutation) error {
	if cudBasedCheck(cud, m) {
		return privacy.Allow
	}
	return privacy.Skip
}

// AllowWritePermissionsRule grants write permission.
func AllowWritePermissionsRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, _ ent.Mutation, p *models.PermissionSettings) error {
		if p.CanWrite {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// AlwaysAllowIfNoPermissionRule grants access if no permissions is present on context.
func AlwaysAllowIfNoPermissionRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, _ ent.Mutation) error {
		if FromContext(ctx) == nil {
			return privacy.Allow
		}
		return privacy.Skip
	})
}
