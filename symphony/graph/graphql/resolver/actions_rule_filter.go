// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/actions/core"
)

type actionsRuleFilterResolver struct{}

func (actionsRuleFilterResolver) Operator(ctx context.Context, ar *core.ActionsRuleFilter) (*models.ActionsOperator, error) {
	operator, ok := core.AllOperators[ar.OperatorID]
	if !ok {
		return nil, fmt.Errorf("operator %s does not exist", ar.OperatorID)
	}
	return &models.ActionsOperator{
		OperatorID:  operator.OperatorID(),
		Description: operator.Description(),
		DataType:    operator.DataType(),
	}, nil
}
