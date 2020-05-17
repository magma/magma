// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// PropertyTypeWritePolicyRule grants write permission to property type based on policy.
// nolint: dupl
func PropertyTypeWritePolicyRule() privacy.MutationRule {
	return privacy.PropertyTypeMutationRuleFunc(func(ctx context.Context, m *ent.PropertyTypeMutation) error {
		propertyTypeID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		propType, err := m.Client().PropertyType.Query().
			Where(propertytype.ID(propertyTypeID)).
			WithLocationType().
			WithEquipmentType().
			WithEquipmentPortType().
			WithLinkEquipmentPortType().
			WithServiceType().
			WithWorkOrderType().
			WithProjectType().
			Only(ctx)

		if err != nil {
			if !ent.IsNotFound(err) {
				return privacy.Denyf("failed to fetch property type: %w", err)
			}
			return privacy.Skip
		}
		p := FromContext(ctx)
		switch {
		case propType.Edges.LocationType != nil:
			return allowOrSkip(p.InventoryPolicy.LocationType.Update)
		case propType.Edges.EquipmentType != nil:
			return allowOrSkip(p.InventoryPolicy.EquipmentType.Update)
		case propType.Edges.EquipmentPortType != nil:
			return allowOrSkip(p.InventoryPolicy.PortType.Update)
		case propType.Edges.LinkEquipmentPortType != nil:
			return allowOrSkip(p.InventoryPolicy.PortType.Update)
		case propType.Edges.ServiceType != nil:
			return allowOrSkip(p.InventoryPolicy.ServiceType.Update)
		case propType.Edges.WorkOrderType != nil:
			return allowOrSkip(p.WorkforcePolicy.Templates.Update)
		case propType.Edges.ProjectType != nil:
			return allowOrSkip(p.WorkforcePolicy.Templates.Update)
		}
		return privacy.Skip
	})
}

// PropertyTypeCreatePolicyRule grants create permission to property type based on policy.
// nolint: dupl
func PropertyTypeCreatePolicyRule() privacy.MutationRule {
	return privacy.PropertyTypeMutationRuleFunc(func(ctx context.Context, m *ent.PropertyTypeMutation) error {
		if !m.Op().Is(ent.OpCreate) {
			return privacy.Skip
		}
		p := FromContext(ctx)
		if _, exists := m.LocationTypeID(); exists {
			return allowOrSkip(p.InventoryPolicy.LocationType.Update)
		}
		if _, exists := m.EquipmentTypeID(); exists {
			return allowOrSkip(p.InventoryPolicy.EquipmentType.Update)
		}
		if _, exists := m.EquipmentPortTypeID(); exists {
			return allowOrSkip(p.InventoryPolicy.PortType.Update)
		}
		if _, exists := m.LinkEquipmentPortTypeID(); exists {
			return allowOrSkip(p.InventoryPolicy.LocationType.Update)
		}
		if _, exists := m.ServiceTypeID(); exists {
			return allowOrSkip(p.InventoryPolicy.ServiceType.Update)
		}
		if _, exists := m.WorkOrderTypeID(); exists {
			return allowOrSkip(p.WorkforcePolicy.Templates.Update)
		}
		if _, exists := m.ProjectTypeID(); exists {
			return allowOrSkip(p.WorkforcePolicy.Templates.Update)
		}
		return privacy.Skip
	})
}

// PropertyWritePolicyRule grants write permission to property based on policy.
// nolint: dupl
func PropertyWritePolicyRule() privacy.MutationRule {
	return privacy.PropertyMutationRuleFunc(func(ctx context.Context, m *ent.PropertyMutation) error {
		propertyID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		prop, err := m.Client().Property.Query().
			Where(property.ID(propertyID)).
			WithLocation().
			WithEquipment().
			WithEquipmentPort().
			WithLink().
			WithService().
			WithWorkOrder().
			WithProject().
			Only(ctx)

		if err != nil {
			if !ent.IsNotFound(err) {
				return privacy.Denyf("failed to fetch property: %w", err)
			}
			return privacy.Skip
		}
		p := FromContext(ctx)
		switch {
		case prop.Edges.Location != nil:
			locationTypeID, err := prop.Edges.Location.QueryType().OnlyID(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return privacy.Skip
				}
				return privacy.Denyf("failed to fetch location type id: %w", err)
			}
			return allowOrSkipLocations(p.InventoryPolicy.Location.Update, locationTypeID)
		case prop.Edges.Equipment != nil:
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		case prop.Edges.EquipmentPort != nil:
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		case prop.Edges.Link != nil:
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		case prop.Edges.Service != nil:
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		case prop.Edges.WorkOrder != nil:
			return allowOrSkipWorkOrder(ctx, p, prop.Edges.WorkOrder)
		case prop.Edges.Project != nil:
			projectTypeID, err := prop.Edges.Project.QueryType().OnlyID(ctx)
			if err != nil {
				return privacy.Denyf("failed to fetch project type id: %w", err)
			}
			return privacyDecision(checkWorkforce(p.WorkforcePolicy.Data.Update, nil, &projectTypeID))
		}
		return privacy.Skip
	})
}

// PropertyCreatePolicyRule grants create permission to property based on policy.
// nolint: dupl
func PropertyCreatePolicyRule() privacy.MutationRule {
	return privacy.PropertyMutationRuleFunc(func(ctx context.Context, m *ent.PropertyMutation) error {
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
		if _, exists := m.EquipmentPortID(); exists {
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		}
		if _, exists := m.LinkID(); exists {
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
		}
		if _, exists := m.ServiceID(); exists {
			return allowOrSkip(p.InventoryPolicy.Equipment.Update)
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
		if projectID, exists := m.ProjectID(); exists {
			projectTypeID, err := m.Client().ProjectType.Query().
				Where(projecttype.HasProjectsWith(project.ID(projectID))).
				OnlyID(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return privacy.Skip
				}
				return privacy.Denyf("failed to fetch project type id: %w", err)
			}
			return privacyDecision(checkWorkforce(p.WorkforcePolicy.Data.Update, nil, &projectTypeID))
		}
		return privacy.Skip
	})
}
