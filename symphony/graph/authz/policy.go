// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

var (
	// prePolicy is executed before privacy policy.
	prePolicy = privacy.Policy{
		Query: privacy.QueryPolicy{
			AlwaysDenyQueryIfNoPermissionRule(),
		},
		Mutation: privacy.MutationPolicy{
			AlwaysDenyMutationIfNoPermissionRule(),
		},
	}

	// postPolicy is executed after privacy policy.
	postPolicy = privacy.Policy{
		Query: privacy.QueryPolicy{
			privacy.AlwaysAllowRule(),
		},
		Mutation: privacy.MutationPolicy{
			AllowWritePermissionsRule(),
			privacy.AlwaysDenyRule(),
		},
	}
)

type (
	// PolicyOption configures policy creation.
	PolicyOption func(*policies)

	// policies aggregate policy options.
	policies struct {
		query     privacy.QueryPolicy
		mutation  privacy.MutationPolicy
		pre, post privacy.Policy
	}
)

// WithQueryRules adds query rules to policy.
func WithQueryRules(rules ...privacy.QueryRule) PolicyOption {
	return func(policies *policies) {
		policies.query = append(policies.query, rules...)
	}
}

// WithMutationRules adds mutation rules to policy.
func WithMutationRules(rules ...privacy.MutationRule) PolicyOption {
	return func(policies *policies) {
		policies.mutation = append(policies.mutation, rules...)
	}
}

// WithPrePolicy overrides the pre-policy to be executed.
func WithPrePolicy(policy privacy.Policy) PolicyOption {
	return func(policies *policies) {
		policies.pre = policy
	}
}

// WithPostPolicy overrides the post-policy to be executed.
func WithPostPolicy(policy privacy.Policy) PolicyOption {
	return func(policies *policies) {
		policies.post = policy
	}
}

// NewPolicy creates a privacy policy.
func NewPolicy(opts ...PolicyOption) ent.Policy {
	policies := policies{
		pre:  prePolicy,
		post: postPolicy,
	}
	for _, opt := range opts {
		opt(&policies)
	}
	return privacy.Policy{
		Query:    policies.queryPolicy(),
		Mutation: policies.mutationPolicy(),
	}
}

func (policies policies) queryPolicy() privacy.QueryPolicy {
	policy := append(privacy.QueryPolicy(nil), policies.pre.Query...)
	policy = append(policy, policies.query...)
	policy = append(policy, policies.post.Query...)
	return policy
}

func (policies policies) mutationPolicy() privacy.MutationPolicy {
	policy := append(privacy.MutationPolicy(nil), policies.pre.Mutation...)
	policy = append(policy, policies.mutation...)
	policy = append(policy, policies.post.Mutation...)
	return policy
}
