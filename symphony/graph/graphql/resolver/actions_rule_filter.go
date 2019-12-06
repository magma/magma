// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/schema"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type actionsRuleFilterResolver struct{}

func (actionsRuleFilterResolver) Operator(ctx context.Context, ar *schema.ActionsRuleFilter) (*models.ActionsOperator, error) {
	// TODO: stub
	return &models.ActionsOperator{
		OperatorID:  ar.OperatorID,
		Description: "blah",
		DataInput:   "{}",
	}, nil
}
