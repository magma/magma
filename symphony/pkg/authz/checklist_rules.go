// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent/workorder"

	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"

	"github.com/facebookincubator/symphony/pkg/ent/checklistcategory"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

// CheckListCategoryDefinitionWritePolicyRule grants write permission to checklist category definition based on policy.
func CheckListCategoryDefinitionWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).WorkforcePolicy.Templates.Update)
	})
}

// CheckListItemDefinitionWritePolicyRule grants write permission to checklist item definition based on policy.
func CheckListItemDefinitionWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).WorkforcePolicy.Templates.Update)
	})
}

// CheckListItemCreatePolicyRule grants create permission to checklist category based on policy.
// nolint: dupl
func CheckListCategoryCreatePolicyRule() privacy.MutationRule {
	return privacy.CheckListCategoryMutationRuleFunc(func(ctx context.Context, m *ent.CheckListCategoryMutation) error {
		if !m.Op().Is(ent.OpCreate) {
			return privacy.Skip
		}
		if woID, exists := m.WorkOrderID(); exists {
			workOrder, err := m.Client().WorkOrder.Get(ctx, woID)
			if err != nil {
				return privacy.Denyf("failed to fetch work order: %w", err)
			}
			return allowOrSkipWorkOrder(ctx, FromContext(ctx), workOrder)
		}
		return privacy.Skip
	})
}

// CheckListCategoryWritePolicyRule grants write permission to work order based on policy.
func CheckListCategoryWritePolicyRule() privacy.MutationRule {
	return privacy.CheckListCategoryMutationRuleFunc(func(ctx context.Context, m *ent.CheckListCategoryMutation) error {
		categoryID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		workOrder, err := m.Client().CheckListCategory.Query().
			Where(checklistcategory.ID(categoryID)).
			QueryWorkOrder().
			Only(ctx)
		if err != nil {
			return privacy.Denyf("failed to fetch work order: %w", err)
		}
		return allowOrSkipWorkOrder(ctx, FromContext(ctx), workOrder)
	})
}

// CheckListItemCreatePolicyRule grants create permission to checklist item based on policy.
// nolint: dupl
func CheckListItemCreatePolicyRule() privacy.MutationRule {
	return privacy.CheckListItemMutationRuleFunc(func(ctx context.Context, m *ent.CheckListItemMutation) error {
		if !m.Op().Is(ent.OpCreate) {
			return privacy.Skip
		}
		if categoryID, exists := m.CheckListCategoryID(); exists {
			workOrder, err := m.Client().WorkOrder.Query().
				Where(workorder.HasCheckListCategoriesWith(checklistcategory.ID(categoryID))).
				Only(ctx)
			if err != nil {
				return privacy.Denyf("failed to fetch work order: %w", err)
			}
			return allowOrSkipWorkOrder(ctx, FromContext(ctx), workOrder)
		}
		return privacy.Skip
	})
}

// CheckListItemWritePolicyRule grants write permission to checklist item based on policy.
func CheckListItemWritePolicyRule() privacy.MutationRule {
	return privacy.CheckListItemMutationRuleFunc(func(ctx context.Context, m *ent.CheckListItemMutation) error {
		itemID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		workOrder, err := m.Client().CheckListItem.Query().
			Where(checklistitem.ID(itemID)).
			QueryCheckListCategory().
			QueryWorkOrder().
			Only(ctx)
		if err != nil {
			return privacy.Denyf("failed to fetch work order: %w", err)
		}
		return allowOrSkipWorkOrder(ctx, FromContext(ctx), workOrder)
	})
}

// CheckListCategoryReadPolicyRule grants read permission to checklist category based on policy.
func CheckListCategoryReadPolicyRule() privacy.QueryRule {
	return privacy.CheckListCategoryQueryRuleFunc(func(ctx context.Context, q *ent.CheckListCategoryQuery) error {
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			q.Where(checklistcategory.HasWorkOrderWith(woPredicate))
		}
		return privacy.Skip
	})
}

// CheckListItemReadPolicyRule grants read permission to checklist item based on policy.
func CheckListItemReadPolicyRule() privacy.QueryRule {
	return privacy.CheckListItemQueryRuleFunc(func(ctx context.Context, q *ent.CheckListItemQuery) error {
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			q.Where(
				checklistitem.HasCheckListCategoryWith(checklistcategory.HasWorkOrderWith(woPredicate)))
		}
		return privacy.Skip
	})
}
