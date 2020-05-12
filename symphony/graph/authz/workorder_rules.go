// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
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

func workOrderIsEditable(ctx context.Context, workOrder *ent.WorkOrder) (bool, error) {
	userViewer, ok := viewer.FromContext(ctx).(*viewer.UserViewer)
	if !ok {
		return false, nil
	}
	isOwner, err := isUserWOOwner(ctx, userViewer.User().ID, workOrder)
	if err != nil {
		return false, err
	}
	if isOwner {
		return true, nil
	}
	isAssignee, err := isUserWOAssignee(ctx, userViewer.User().ID, workOrder)
	if err != nil {
		return false, err
	}
	if isAssignee {
		return true, nil
	}
	return false, nil
}

func workOrderCudBasedCheck(ctx context.Context, cud *models.WorkforceCud, m *ent.WorkOrderMutation) (bool, error) {
	if m.Op().Is(ent.OpCreate) {
		typeID, exists := m.TypeID()
		if !exists {
			return false, errors.New("creating work order with no type")
		}
		return checkWorkforce(cud.Create, &typeID, nil), nil
	}
	id, exists := m.ID()
	if !exists {
		return false, nil
	}
	workOrderTypeID, err := m.Client().WorkOrderType.Query().
		Where(workordertype.HasWorkOrdersWith(workorder.ID(id))).
		OnlyID(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to fetch work order type id: %w", err)
	}
	if m.Op().Is(ent.OpUpdateOne) {
		return checkWorkforce(cud.Update, &workOrderTypeID, nil), nil
	}
	return checkWorkforce(cud.Delete, &workOrderTypeID, nil), nil
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
		workOrder, err := m.Client().WorkOrder.Get(ctx, workOrderID)
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
	return privacy.WorkOrderMutationRuleFunc(func(ctx context.Context, m *ent.WorkOrderMutation) error {
		cud := FromContext(ctx).WorkforcePolicy.Data
		allowed, err := workOrderCudBasedCheck(ctx, cud, m)
		if err != nil {
			return privacy.Denyf(err.Error())
		}
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
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return cudBasedRule(FromContext(ctx).WorkforcePolicy.Templates, m)
	})
}
