// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testRule struct {
	mock.Mock
}

func (m *testRule) EvalQuery(ctx context.Context, query ent.Query) error {
	return m.Called(ctx, query).Error(0)
}

func (m *testRule) EvalMutation(ctx context.Context, mutation ent.Mutation) error {
	return m.Called(ctx, mutation).Error(0)
}

func TestPolicy(t *testing.T) {
	var preRule, midRule, postRule testRule
	for _, rule := range []*testRule{&preRule, &midRule, &postRule} {
		rule.On("EvalQuery", mock.Anything, mock.Anything).
			Return(privacy.Skip).
			Once()
		rule.On("EvalMutation", mock.Anything, mock.Anything).
			Return(privacy.Skip).
			Once()
	}
	defer func() {
		preRule.AssertExpectations(t)
		midRule.AssertExpectations(t)
		postRule.AssertExpectations(t)
	}()

	policy := authz.NewPolicy(
		authz.WithQueryRules(&midRule),
		authz.WithMutationRules(&midRule),
		authz.WithPrePolicy(privacy.Policy{
			Query:    privacy.QueryPolicy{&preRule},
			Mutation: privacy.MutationPolicy{&preRule},
		}),
		authz.WithPostPolicy(privacy.Policy{
			Query:    privacy.QueryPolicy{&postRule},
			Mutation: privacy.MutationPolicy{&postRule},
		}),
	)
	t.Run("Query", func(t *testing.T) {
		err := policy.EvalQuery(context.Background(), nil)
		assert.NoError(t, err)
	})
	t.Run("Mutation", func(t *testing.T) {
		err := policy.EvalMutation(context.Background(), nil)
		assert.NoError(t, err)
	})
}
