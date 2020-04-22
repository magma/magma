package authz

import (
	"context"
	"errors"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// Policy wraps privacy policy with static pre/post policies.
type Policy privacy.Policy

var (
	// PrePolicy is executed before privacy policy.
	PrePolicy = privacy.Policy{}

	// PostPolicy is executed after privacy policy.
	PostPolicy = privacy.Policy{}
)

// EvalQuery evaluates query policy.
func (p Policy) EvalQuery(ctx context.Context, q ent.Query) error {
	for _, policy := range []ent.Policy{PrePolicy, privacy.Policy(p), PostPolicy} {
		if err := policy.EvalQuery(ctx, q); err != nil && !errors.Is(err, privacy.Skip) {
			return err
		}
	}
	return nil
}

// EvalMutation evaluates mutation policy.
func (p Policy) EvalMutation(ctx context.Context, m ent.Mutation) error {
	for _, policy := range []ent.Policy{PrePolicy, privacy.Policy(p), PostPolicy} {
		if err := policy.EvalMutation(ctx, m); err != nil && !errors.Is(err, privacy.Skip) {
			return err
		}
	}
	return nil
}
