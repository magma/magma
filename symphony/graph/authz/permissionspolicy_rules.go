package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// PermissionsPolicyWritePolicyRule grants write permission to permissions policy based on policy.
func PermissionsPolicyWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, _ ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).AdminPolicy.Access)
	})
}
