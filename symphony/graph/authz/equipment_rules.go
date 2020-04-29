package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

// EquipmentWritePolicyRule grants write permission to equipment based on policy.
func EquipmentWritePolicyRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, m ent.Mutation, p *models.PermissionSettings) error {
		return cudBasedRule(p.InventoryPolicy.Equipment, m)
	})
}

// EquipmentTypeWritePolicyRule grants write permission to equipment type based on policy.
func EquipmentTypeWritePolicyRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, m ent.Mutation, p *models.PermissionSettings) error {
		return cudBasedRule(p.InventoryPolicy.EquipmentType, m)
	})
}

// EquipmentPortTypeWritePolicyRule grants write permission to equipment port type based on policy.
func EquipmentPortTypeWritePolicyRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, m ent.Mutation, p *models.PermissionSettings) error {
		return cudBasedRule(p.InventoryPolicy.PortType, m)
	})
}
