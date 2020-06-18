// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/pkg/errors"
)

type (
	actionsTriggerResolver struct{}
)

func (actionsTriggerResolver) SupportedActions(
	ctx context.Context, obj *models.ActionsTrigger,
) ([]*models.ActionsAction, error) {
	ac := actions.FromContext(ctx)
	trigger, err := ac.TriggerForID(obj.TriggerID)
	if err != nil {
		return nil, errors.Wrap(err, "supported actions")
	}

	actionIDs := trigger.SupportedActionIDs()
	actions := make([]*models.ActionsAction, 0, len(actionIDs))
	for _, actionID := range actionIDs {
		action, err := ac.ActionForID(actionID)
		if err != nil {
			return nil, errors.Wrap(err, "supported actions")
		}
		actions = append(actions, &models.ActionsAction{
			ActionID:    actionID,
			Description: action.Description(),
			DataType:    action.DataType(),
		})
	}
	return actions, nil
}

func (actionsTriggerResolver) SupportedFilters(ctx context.Context, obj *models.ActionsTrigger) ([]*models.ActionsFilter, error) {
	ac := actions.FromContext(ctx)
	trigger, err := ac.TriggerForID(obj.TriggerID)
	if err != nil {
		return nil, errors.Wrap(err, "supported actions")
	}

	filters := trigger.SupportedFilters()
	ret := make([]*models.ActionsFilter, 0, len(filters))
	for _, filter := range filters {
		ret = append(ret, newActionsFilter(filter))
	}
	return ret, nil
}

func newActionsFilter(filter core.Filter) *models.ActionsFilter {
	operators := make([]*models.ActionsOperator, 0, len(filter.SupportedOperators()))
	for _, operator := range filter.SupportedOperators() {
		operators = append(operators, newSupportedOperator(operator))
	}
	return &models.ActionsFilter{
		FilterID:           filter.FilterID(),
		Description:        filter.Description(),
		SupportedOperators: operators,
	}
}

func newSupportedOperator(operator core.Operator) *models.ActionsOperator {
	return &models.ActionsOperator{
		OperatorID:  operator.OperatorID(),
		Description: operator.Description(),
		DataType:    operator.DataType(),
	}
}
