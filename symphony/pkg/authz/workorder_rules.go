// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
	"github.com/facebookincubator/symphony/pkg/viewer"
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

func isViewerWorkOrderOwnerOrAssignee(ctx context.Context, workOrder *ent.WorkOrder) (bool, error) {
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

func getWorkOrderType(ctx context.Context, m *ent.WorkOrderMutation) (*int, error) {
	id, exists := m.ID()
	if !exists {
		return nil, nil
	}
	workOrderTypeID, err := m.Client().WorkOrderType.Query().
		Where(workordertype.HasWorkOrdersWith(workorder.ID(id))).
		OnlyID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch work order type id: %w", err)
	}
	return &workOrderTypeID, nil
}

func workOrderCudBasedCheck(ctx context.Context, cud *models.WorkforceCud, m *ent.WorkOrderMutation) (bool, error) {
	if m.Op().Is(ent.OpCreate) {
		typeID, exists := m.TypeID()
		if !exists {
			return false, errors.New("creating work order with no type")
		}
		return checkWorkforce(cud.Create, &typeID, nil), nil
	}
	workOrderTypeID, err := getWorkOrderType(ctx, m)
	if err != nil {
		return false, err
	}
	if m.Op().Is(ent.OpUpdateOne) {
		return checkWorkforce(cud.Update, workOrderTypeID, nil), nil
	}
	return checkWorkforce(cud.Delete, workOrderTypeID, nil), nil
}

func workOrderReadPredicate(ctx context.Context) predicate.WorkOrder {
	var predicates []predicate.WorkOrder
	rule := FromContext(ctx).WorkforcePolicy.Read
	switch rule.IsAllowed {
	case models.PermissionValueYes:
		return nil
	case models.PermissionValueByCondition:
		predicates = append(predicates,
			workorder.HasTypeWith(workordertype.IDIn(rule.WorkOrderTypeIds...)))
	}
	if v, exists := viewer.FromContext(ctx).(*viewer.UserViewer); exists {
		predicates = append(predicates,
			workorder.HasOwnerWith(user.ID(v.User().ID)),
			workorder.HasAssigneeWith(user.ID(v.User().ID)),
		)
	}
	return workorder.Or(predicates...)
}

func isAssigneeChanged(ctx context.Context, m *ent.WorkOrderMutation) (bool, error) {
	var currAssigneeID *int
	assigneeIDToSet, assigned := m.AssigneeID()
	assigneeCleared := m.AssigneeCleared()
	if !assigned && !assigneeCleared {
		return false, nil
	}
	workOrderID, exists := m.ID()
	if !exists {
		return assigned, nil
	}
	assigneeID, err := m.Client().User.Query().
		Where(user.HasAssignedWorkOrdersWith(workorder.ID(workOrderID))).
		OnlyID(ctx)
	if err == nil {
		currAssigneeID = &assigneeID
	}
	if err != nil && !ent.IsNotFound(err) {
		return false, privacy.Denyf("failed to fetch assignee: %w", err)
	}
	switch {
	case currAssigneeID == nil && assigned:
		return true, nil
	case currAssigneeID != nil && assigned && *currAssigneeID != assigneeIDToSet:
		return true, nil
	case currAssigneeID != nil && assigneeCleared:
		return true, nil
	}
	return false, nil
}

func isOwnerChanged(ctx context.Context, m *ent.WorkOrderMutation) (bool, error) {
	var currOwnerID *int
	ownerIDToSet, owned := m.OwnerID()
	ownerCleared := m.OwnerCleared()
	if !owned && !ownerCleared {
		return false, nil
	}
	workOrderID, exists := m.ID()
	if !exists {
		return owned, nil
	}
	ownerID, err := m.Client().User.Query().
		Where(user.HasOwnedWorkOrdersWith(workorder.ID(workOrderID))).
		OnlyID(ctx)
	if err == nil {
		currOwnerID = &ownerID
	}
	if err != nil && !ent.IsNotFound(err) {
		return false, privacy.Denyf("failed to fetch owner: %w", err)
	}
	switch {
	case currOwnerID == nil && owned:
		return true, nil
	case currOwnerID != nil && owned && *currOwnerID != ownerIDToSet:
		return true, nil
	case currOwnerID != nil && ownerCleared:
		return true, nil
	}
	return false, nil
}

// AllowWorkOrderOwnerOrAssigneeWrite grants write permission if user is owner or assignee of workorder
func AllowWorkOrderOwnerOrAssigneeWrite() privacy.MutationRule {
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
		ownerChanged, err := isOwnerChanged(ctx, m)
		if err != nil {
			return privacy.Denyf(err.Error())
		}
		if isAssignee && !m.Op().Is(ent.OpDeleteOne) && !ownerChanged {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// WorkOrderWritePolicyRule grants write permission to work order based on policy.
func WorkOrderWritePolicyRule() privacy.MutationRule {
	return privacy.WorkOrderMutationRuleFunc(func(ctx context.Context, m *ent.WorkOrderMutation) error {
		cud := FromContext(ctx).WorkforcePolicy.Data
		allowed, err := workOrderCudBasedCheck(ctx, cud, m)
		if err != nil {
			return privacy.Denyf(err.Error())
		}
		if !m.Op().Is(ent.OpCreate) {
			workOrderTypeID, err := getWorkOrderType(ctx, m)
			if err != nil {
				return err
			}
			assigneeChanged, err := isAssigneeChanged(ctx, m)
			if err != nil {
				return privacy.Denyf(err.Error())
			}
			if assigneeChanged {
				allowed = allowed && checkWorkforce(cud.Assign, workOrderTypeID, nil)
			}
			ownerChanged, err := isOwnerChanged(ctx, m)
			if err != nil {
				return privacy.Denyf(err.Error())
			}
			if ownerChanged {
				allowed = allowed && checkWorkforce(cud.TransferOwnership, workOrderTypeID, nil)
			}
		}
		if allowed {
			return privacy.Allow
		}
		return privacy.Skip
	})
}

// WorkOrderReadPolicyRule grants read permission to work order based on policy.
func WorkOrderReadPolicyRule() privacy.QueryRule {
	return privacy.WorkOrderQueryRuleFunc(func(ctx context.Context, q *ent.WorkOrderQuery) error {
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			q.Where(woPredicate)
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
