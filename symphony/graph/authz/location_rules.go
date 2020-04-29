package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

// LocationWritePolicyRule grants write permission to location based on policy.
func LocationWritePolicyRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, m ent.Mutation, p *models.PermissionSettings) error {
		return cudBasedRule(p.InventoryPolicy.Location, m)
	})
}

// LocationTypeWritePolicyRule grants write permission to location type based on policy.
func LocationTypeWritePolicyRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, m ent.Mutation, p *models.PermissionSettings) error {
		return cudBasedRule(p.InventoryPolicy.LocationType, m)
	})
}
