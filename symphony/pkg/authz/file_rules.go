// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent/checklistcategory"

	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"

	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/locationtype"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/file"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

// FileWritePolicyRule grants write permission to file type based on policy.
// nolint: dupl
func FileWritePolicyRule() privacy.MutationRule {
	return privacy.FileMutationRuleFunc(func(ctx context.Context, m *ent.FileMutation) error {
		fileID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		file, err := m.Client().File.Query().
			Where(file.ID(fileID)).
			WithLocation().
			WithEquipment().
			WithWorkOrder().
			WithChecklistItem().
			WithUser().
			Only(ctx)

		if err != nil {
			if !ent.IsNotFound(err) {
				return privacy.Denyf("failed to fetch file: %w", err)
			}
			return privacy.Skip
		}
		p := FromContext(ctx)

		switch {
		case file.Edges.Location != nil:
			locationTypeID, err := file.Edges.Location.QueryType().OnlyID(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return privacy.Skip
				}
				return privacy.Denyf("failed to fetch location type id: %w", err)
			}
			return allowOrSkipLocations(p.InventoryPolicy.Location.Update, locationTypeID)
		case file.Edges.Equipment != nil:
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		case file.Edges.User != nil:
			return allowOrSkip(p.AdminPolicy.Access)
		case file.Edges.WorkOrder != nil:
			return allowOrSkipWorkOrder(ctx, p, file.Edges.WorkOrder)
		case file.Edges.ChecklistItem != nil:
			wo, err := file.Edges.ChecklistItem.QueryCheckListCategory().QueryWorkOrder().Only(ctx)
			if err != nil {
				return privacy.Denyf("failed to fetch work order : %w", err)
			}
			return allowOrSkipWorkOrder(ctx, p, wo)
		}
		return privacy.Skip
	})
}

// FileCreatePolicyRule grants create permission to file type based on policy.
// nolint: dupl
func FileCreatePolicyRule() privacy.MutationRule {
	return privacy.FileMutationRuleFunc(func(ctx context.Context, m *ent.FileMutation) error {
		if !m.Op().Is(ent.OpCreate) {
			return privacy.Skip
		}
		p := FromContext(ctx)
		if locationID, exists := m.LocationID(); exists {
			locationTypeID, err := m.Client().LocationType.Query().
				Where(locationtype.HasLocationsWith(location.ID(locationID))).
				OnlyID(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return privacy.Skip
				}
				return privacy.Denyf("failed to fetch location type id: %w", err)
			}
			return allowOrSkipLocations(p.InventoryPolicy.Location.Update, locationTypeID)
		}
		if _, exists := m.EquipmentID(); exists {
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		}
		if _, exists := m.UserID(); exists {
			return allowOrSkip(p.AdminPolicy.Access)
		}
		if workOrderID, exists := m.WorkOrderID(); exists {
			workOrder, err := m.Client().WorkOrder.Get(ctx, workOrderID)
			if err != nil {
				if !ent.IsNotFound(err) {
					return privacy.Denyf("failed to fetch work order: %w", err)
				}
				return privacy.Skip
			}
			return allowOrSkipWorkOrder(ctx, p, workOrder)
		}
		if clItemID, exists := m.ChecklistItemID(); exists {
			wo, err := m.Client().WorkOrder.Query().
				Where(workorder.HasCheckListCategoriesWith(checklistcategory.HasCheckListItemsWith(checklistitem.ID(clItemID)))).Only(ctx)
			if err != nil {
				return privacy.Denyf("failed to fetch work order : %w", err)
			}
			return allowOrSkipWorkOrder(ctx, p, wo)
		}
		return privacy.Skip
	})
}

// FileReadPolicyRule grants read permission to file based on policy.
func FileReadPolicyRule() privacy.QueryRule {
	return privacy.FileQueryRuleFunc(func(ctx context.Context, q *ent.FileQuery) error {
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			q.Where(
				file.Or(
					file.Not(file.HasWorkOrder()),
					file.HasWorkOrderWith(woPredicate)),
				file.Or(
					file.Not(file.HasChecklistItem()),
					file.HasChecklistItemWith(
						checklistitem.HasCheckListCategoryWith(
							checklistcategory.HasWorkOrderWith(woPredicate)))),
			)
		}
		return privacy.Skip
	})
}
