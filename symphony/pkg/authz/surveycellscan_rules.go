// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent/workorder"

	"github.com/facebookincubator/symphony/pkg/ent/surveycellscan"

	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"

	"github.com/facebookincubator/symphony/pkg/ent/checklistcategory"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

func allowOrSkipCheckListItemCreate(ctx context.Context, client *ent.Client, checkListItemID int) error {
	workOrder, err := client.WorkOrder.Query().
		Where(workorder.HasCheckListCategoriesWith(checklistcategory.HasCheckListItemsWith(checklistitem.ID(checkListItemID)))).
		Only(ctx)
	if err != nil {
		return privacy.Denyf("failed to fetch work order: %w", err)
	}
	return allowOrSkipWorkOrder(ctx, FromContext(ctx), workOrder)
}

// SurveyCellScanCreatePolicyRule grants create permission to SurveyCellScan based on policy.
func SurveyCellScanCreatePolicyRule() privacy.MutationRule {
	return privacy.SurveyCellScanMutationRuleFunc(func(ctx context.Context, m *ent.SurveyCellScanMutation) error {
		if !m.Op().Is(ent.OpCreate) {
			return privacy.Skip
		}
		if checkListItemID, exists := m.ChecklistItemID(); exists {
			return allowOrSkipCheckListItemCreate(ctx, m.Client(), checkListItemID)
		}
		return privacy.Skip
	})
}

// SurveyCellScanReadPolicyRule grants read permission to SurveyCellScan based on policy.
func SurveyCellScanReadPolicyRule() privacy.QueryRule {
	return privacy.SurveyCellScanQueryRuleFunc(func(ctx context.Context, q *ent.SurveyCellScanQuery) error {
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			q.Where(
				surveycellscan.HasChecklistItemWith(checklistitem.HasCheckListCategoryWith(checklistcategory.HasWorkOrderWith(woPredicate))))
		}
		return privacy.Skip
	})
}

// SurveyCellScanWritePolicyRule grants write permission to SurveyCellScan based on policy.
func SurveyCellScanWritePolicyRule() privacy.MutationRule {
	return privacy.SurveyCellScanMutationRuleFunc(func(ctx context.Context, m *ent.SurveyCellScanMutation) error {
		itemID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		workOrder, err := m.Client().SurveyCellScan.Query().
			Where(surveycellscan.ID(itemID)).
			QueryChecklistItem().
			QueryCheckListCategory().
			QueryWorkOrder().
			Only(ctx)
		if err != nil {
			if !ent.IsNotFound(err) {
				return privacy.Denyf("failed to fetch work order: %w", err)
			}
			return privacy.Skip
		}
		return allowOrSkipWorkOrder(ctx, FromContext(ctx), workOrder)
	})
}
