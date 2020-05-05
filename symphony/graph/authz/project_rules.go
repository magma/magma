// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

// ProjectWritePolicyRule grants write permission to project based on policy.
func ProjectWritePolicyRule() privacy.MutationRule {
	return projectMutationWithPermissionRule(func(ctx context.Context, m *ent.ProjectMutation, p *models.PermissionSettings) error {
		cud := p.WorkforcePolicy.Data
		allowed := cudBasedCheck(&models.Cud{
			Create: cud.Create,
			Update: cud.Update,
			Delete: cud.Delete,
		}, m)
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
	return mutationWithPermissionRule(func(ctx context.Context, m ent.Mutation, p *models.PermissionSettings) error {
		return cudBasedRule(p.WorkforcePolicy.Templates, m)
	})
}
