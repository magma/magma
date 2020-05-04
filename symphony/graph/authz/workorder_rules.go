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

// WorkOrderWritePolicyRule grants write permission to workorder based on policy.
func WorkOrderWritePolicyRule() privacy.MutationRule {
	return workOrderMutationWithPermissionRule(func(ctx context.Context, m *ent.WorkOrderMutation, p *models.PermissionSettings) error {
		cud := p.WorkforcePolicy.Data
		allowed := cudBasedCheck(&models.Cud{
			Create: cud.Create,
			Update: cud.Update,
			Delete: cud.Delete,
		}, m)
		_, assigned := m.AssigneeID()
		if assigned || m.AssigneeCleared() {
			allowed = allowed && (cud.Assign.IsAllowed == models2.PermissionValueYes)
		}
		_, owned := m.OwnerID()
		if owned || m.OwnerCleared() {
			allowed = allowed && (cud.TransferOwnership.IsAllowed == models2.PermissionValueYes)
		}
		if allowed {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// WorkOrderTypeWritePolicyRule grants write permission to work order type based on policy.
func WorkOrderTypeWritePolicyRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, m ent.Mutation, p *models.PermissionSettings) error {
		return cudBasedRule(p.WorkforcePolicy.Templates, m)
	})
}
