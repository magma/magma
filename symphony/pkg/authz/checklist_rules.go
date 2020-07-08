// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategory"
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
)

func allowOrSkipChecklistCategoryDefinition(ctx context.Context, client *ent.Client, checklistCategoryDefinitionID int) error {
	checklistCategoryDefinition, err := client.CheckListCategoryDefinition.Query().
		Where(checklistcategorydefinition.ID(checklistCategoryDefinitionID)).
		WithWorkOrderTemplate().
		WithWorkOrderType().
		Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return privacy.Denyf("failed to query check list category definition: %w", err)
		}
		return privacy.Skip
	}
	switch {
	case checklistCategoryDefinition.Edges.WorkOrderTemplate != nil:
		return privacy.Allow
	case checklistCategoryDefinition.Edges.WorkOrderType != nil:
		return allowOrSkip(FromContext(ctx).WorkforcePolicy.Templates.Update)
	}
	return privacy.Skip
}

// CheckListCategoryDefinitionWritePolicyRule grants write permission to checklist category definition based on policy.
func CheckListCategoryDefinitionWritePolicyRule() privacy.MutationRule {
	return privacy.CheckListCategoryDefinitionMutationRuleFunc(func(ctx context.Context, m *ent.CheckListCategoryDefinitionMutation) error {
		if m.Op().Is(ent.OpCreate) {
			if _, exists := m.WorkOrderTemplateID(); exists {
				return privacy.Allow
			}
			if _, exists := m.WorkOrderTypeID(); exists {
				return allowOrSkip(FromContext(ctx).WorkforcePolicy.Templates.Update)
			}
		} else {
			checklistCategoryDefinitionID, exists := m.ID()
			if !exists {
				return privacy.Skip
			}
			return allowOrSkipChecklistCategoryDefinition(ctx, m.Client(), checklistCategoryDefinitionID)
		}
		return privacy.Skip
	})
}

// CheckListItemDefinitionWritePolicyRule grants write permission to checklist item definition based on policy.
func CheckListItemDefinitionWritePolicyRule() privacy.MutationRule {
	return privacy.CheckListItemDefinitionMutationRuleFunc(func(ctx context.Context, m *ent.CheckListItemDefinitionMutation) error {
		var (
			checklistCategoryDefinitionID int
			exists                        bool
			err                           error
		)
		if m.Op().Is(ent.OpCreate) {
			checklistCategoryDefinitionID, exists = m.CheckListCategoryDefinitionID()
			if !exists {
				return privacy.Skip
			}
		} else {
			checklistItemDefinitionID, exists := m.ID()
			if !exists {
				return privacy.Skip
			}
			checklistCategoryDefinitionID, err = m.Client().CheckListCategoryDefinition.Query().
				Where(checklistcategorydefinition.HasCheckListItemDefinitionsWith(checklistitemdefinition.ID(checklistItemDefinitionID))).
				OnlyID(ctx)
			if err != nil {
				if !ent.IsNotFound(err) {
					return privacy.Denyf("failed to query check list category definition: %w", err)
				}
				return privacy.Skip
			}
		}
		return allowOrSkipChecklistCategoryDefinition(ctx, m.Client(), checklistCategoryDefinitionID)
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
