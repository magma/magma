// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// HyperlinkReadPolicyRule grants read permission to hyperlink based on policy.
func HyperlinkReadPolicyRule() privacy.QueryRule {
	return privacy.HyperlinkQueryRuleFunc(func(ctx context.Context, q *ent.HyperlinkQuery) error {
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			q.Where(
				hyperlink.Or(
					hyperlink.Not(hyperlink.HasWorkOrder()),
					hyperlink.HasWorkOrderWith(woPredicate)))
		}
		return privacy.Skip
	})
}

// HyperlinkWritePolicyRule grants write permission to hyperlink based on policy.
func HyperlinkWritePolicyRule() privacy.MutationRule {
	return privacy.HyperlinkMutationRuleFunc(func(ctx context.Context, m *ent.HyperlinkMutation) error {
		hyperlinkId, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		hyperLink, err := m.Client().Hyperlink.Query().
			Where(hyperlink.ID(hyperlinkId)).
			WithEquipment().
			WithLocation().
			WithWorkOrder().
			Only(ctx)

		if err != nil {
			if !ent.IsNotFound(err) {
				return privacy.Denyf("failed to fetch hyperlink: %w", err)
			}
			return privacy.Skip
		}
		p := FromContext(ctx)
		switch {
		case hyperLink.Edges.Equipment != nil:
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		case hyperLink.Edges.Location != nil:
			locationTypeID, err := hyperLink.Edges.Location.QueryType().OnlyID(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return privacy.Skip
				}
				return privacy.Denyf("failed to fetch location type id: %w", err)
			}
			return allowOrSkipLocations(p.InventoryPolicy.Location.Update, locationTypeID)
		case hyperLink.Edges.WorkOrder != nil:
			return allowOrSkipWorkOrder(ctx, p, hyperLink.Edges.WorkOrder)
		}
		return privacy.Skip
	})
}

// HyperlinkCreatePolicyRule grants create permission to hyperlink based on policy.
func HyperlinkCreatePolicyRule() privacy.MutationRule {
	return privacy.HyperlinkMutationRuleFunc(func(ctx context.Context, m *ent.HyperlinkMutation) error {
		if !m.Op().Is(ent.OpCreate) {
			return privacy.Skip
		}
		p := FromContext(ctx)
		if _, exists := m.EquipmentID(); exists {
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		}
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
		return privacy.Skip
	})
}
