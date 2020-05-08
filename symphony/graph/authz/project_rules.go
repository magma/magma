// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// ProjectWritePolicyRule grants write permission to project based on policy.
func ProjectWritePolicyRule() privacy.MutationRule {
	return privacy.ProjectMutationRuleFunc(func(ctx context.Context, m *ent.ProjectMutation) error {
		cud := FromContext(ctx).WorkforcePolicy.Data
		allowed := workforceCudBasedCheck(cud, m)
		_, owned := m.CreatorID()
		if owned || m.CreatorCleared() {
			allowed = allowed && (cud.TransferOwnership.IsAllowed == models2.PermissionValueYes)
		}
		if allowed {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// ProjectTypeWritePolicyRule grants write permission to project type based on policy.
func ProjectTypeWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return cudBasedRule(FromContext(ctx).WorkforcePolicy.Templates, m)
	})
}
