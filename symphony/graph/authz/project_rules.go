package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

// ProjectTypeWritePolicyRule grants write permission to project type based on policy.
func ProjectTypeWritePolicyRule() privacy.MutationRule {
	return mutationWithPermissionRule(func(ctx context.Context, m ent.Mutation, p *models.PermissionSettings) error {
		return cudBasedRule(p.WorkforcePolicy.Templates, m)
	})
}
