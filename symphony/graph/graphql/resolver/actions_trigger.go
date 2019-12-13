// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/cloud/actions"
	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/pkg/errors"
)

type (
	actionsTriggerResolver struct{}
)

func (actionsTriggerResolver) SupportedActions(
	ctx context.Context, obj *models.ActionsTrigger,
) ([]*models.ActionsAction, error) {
	ac := actions.FromContext(ctx)
	triggerID := core.TriggerID(obj.TriggerID)
	trigger, err := ac.TriggerForID(triggerID)
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
		modelActionID := models.ActionID(actionID)
		if !modelActionID.IsValid() {
			return nil, errors.Errorf("action %s is not a valid models.ActionID", action.ID())
		}
		actions = append(actions, &models.ActionsAction{
			ActionID:    modelActionID,
			Description: action.Description(),
			DataType:    "string", // TODO
		})
	}
	return actions, nil
}

func (actionsTriggerResolver) SupportedFilters(ctx context.Context, obj *models.ActionsTrigger) ([]*models.ActionsFilter, error) {
	ac := actions.FromContext(ctx)
	triggerID := core.TriggerID(obj.TriggerID)
	trigger, err := ac.TriggerForID(triggerID)
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
		DataType:    models.ActionsDataType(operator.DataType()),
	}
}
