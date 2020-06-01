// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/authz/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/locationtype"
	"github.com/facebookincubator/symphony/pkg/ent/privacy"
)

func locationCudBasedRule(ctx context.Context, cud *models.LocationCud, m *ent.LocationMutation) error {
	if m.Op().Is(ent.OpCreate) {
		typeID, exists := m.TypeID()
		if !exists {
			return privacy.Denyf("creating location with no type")
		}
		return allowOrSkipLocations(cud.Create, typeID)
	}
	id, exists := m.ID()
	if !exists {
		return privacy.Skip
	}
	locationTypeID, err := m.Client().LocationType.Query().
		Where(locationtype.HasLocationsWith(location.ID(id))).
		OnlyID(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return privacy.Skip
		}
		return privacy.Denyf("failed to fetch location type id: %w", err)
	}
	if m.Op().Is(ent.OpUpdateOne) {
		return allowOrSkipLocations(cud.Update, locationTypeID)
	}
	return allowOrSkipLocations(cud.Delete, locationTypeID)
}

// LocationWritePolicyRule grants write permission to location based on policy.
func LocationWritePolicyRule() privacy.MutationRule {
	return privacy.LocationMutationRuleFunc(func(ctx context.Context, m *ent.LocationMutation) error {
		return locationCudBasedRule(ctx, FromContext(ctx).InventoryPolicy.Location, m)
	})
}

// LocationTypeWritePolicyRule grants write permission to location type based on policy.
func LocationTypeWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		return cudBasedRule(FromContext(ctx).InventoryPolicy.LocationType, m)
	})
}
