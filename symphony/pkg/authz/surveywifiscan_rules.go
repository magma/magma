// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent/surveywifiscan"

	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"

	"github.com/facebookincubator/symphony/pkg/ent/checklistcategory"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

// SurveyWiFiScanCreatePolicyRule grants create permission to SurveyWiFiScan based on policy.
func SurveyWiFiScanCreatePolicyRule() privacy.MutationRule {
	return privacy.SurveyWiFiScanMutationRuleFunc(func(ctx context.Context, m *ent.SurveyWiFiScanMutation) error {
		if !m.Op().Is(ent.OpCreate) {
			return privacy.Skip
		}
		if checkListItemID, exists := m.ChecklistItemID(); exists {
			return allowOrSkipCheckListItemCreate(ctx, m.Client(), checkListItemID)
		}
		return privacy.Skip
	})
}

// SurveyWiFiScanReadPolicyRule grants read permission to SurveyWiFiScan based on policy.
func SurveyWiFiScanReadPolicyRule() privacy.QueryRule {
	return privacy.SurveyWiFiScanQueryRuleFunc(func(ctx context.Context, q *ent.SurveyWiFiScanQuery) error {
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			q.Where(
				surveywifiscan.HasChecklistItemWith(checklistitem.HasCheckListCategoryWith(checklistcategory.HasWorkOrderWith(woPredicate))))
		}
		return privacy.Skip
	})
}

// SurveyWiFiScanWritePolicyRule grants write permission to SurveyWiFiScan based on policy.
func SurveyWiFiScanWritePolicyRule() privacy.MutationRule {
	return privacy.SurveyWiFiScanMutationRuleFunc(func(ctx context.Context, m *ent.SurveyWiFiScanMutation) error {
		itemID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		workOrder, err := m.Client().SurveyWiFiScan.Query().
			Where(surveywifiscan.ID(itemID)).
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
