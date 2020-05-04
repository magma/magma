// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"

	models2 "github.com/facebookincubator/symphony/graph/authz/models"
)

func isUserWOOwner(ctx context.Context, userID int, workOrder *ent.WorkOrder) (bool, error) {
	ownerID, err := workOrder.QueryOwner().OnlyID(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return false, fmt.Errorf("failed to fetch work order owner: %w", err)
		}
		return false, nil
	}
	return ownerID == userID, nil
}

func isUserWOAssignee(ctx context.Context, userID int, workOrder *ent.WorkOrder) (bool, error) {
	assigneeID, err := workOrder.QueryAssignee().OnlyID(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return false, fmt.Errorf("failed to fetch work order assignee: %w", err)
		}
		return false, nil
	}
	return assigneeID == userID, nil
}

func AllowIfWorkOrderOwnerOrAssignee() privacy.MutationRule {
	return privacy.WorkOrderMutationRuleFunc(func(ctx context.Context, m *ent.WorkOrderMutation) error {
		workOrderID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		userViewer, ok := viewer.FromContext(ctx).(*viewer.UserViewer)
		if !ok {
			return privacy.Skip
		}
		client := ent.FromContext(ctx)
		workOrder, err := client.WorkOrder.Get(ctx, workOrderID)
		if err != nil {
			if !ent.IsNotFound(err) {
				return privacy.Denyf("failed to fetch work order: %w", err)
			}
			return privacy.Skip
		}
		isOwner, err := isUserWOOwner(ctx, userViewer.User().ID, workOrder)
		if err != nil {
			return privacy.Denyf(err.Error())
		}
		if isOwner {
			return privacy.Allow
		}
		isAssignee, err := isUserWOAssignee(ctx, userViewer.User().ID, workOrder)
		if err != nil {
			return privacy.Denyf(err.Error())
		}
		_, owned := m.OwnerID()
		if isAssignee && !m.Op().Is(ent.OpDeleteOne) && !owned && !m.OwnerCleared() {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

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
