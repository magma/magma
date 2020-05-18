package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// UsersGroupWritePolicyRule grants write permission to users group based on policy.
func UsersGroupWritePolicyRule() privacy.MutationRule {
	return privacy.MutationRuleFunc(func(ctx context.Context, _ ent.Mutation) error {
		return allowOrSkip(FromContext(ctx).AdminPolicy.Access)
	})
}
