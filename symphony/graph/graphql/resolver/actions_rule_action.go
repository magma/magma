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

type actionsRuleActionResolver struct{}

func (actionsRuleActionResolver) Action(ctx context.Context, ar *core.ActionsRuleAction) (*models.ActionsAction, error) {
	ac := actions.FromContext(ctx)

	action, err := ac.ActionForID(ar.ActionID)
	if err != nil {
		return nil, errors.Errorf("actionID %s not a registered action", ar.ActionID)
	}

	return &models.ActionsAction{
		Description: action.Description(),
		DataType:    models.ActionsDataType(action.DataType()),
	}, nil
}

func (actionsRuleActionResolver) ActionID(ctx context.Context, action *core.ActionsRuleAction) (models.ActionID, error) {
	value := models.ActionID(action.ActionID)
	if !value.IsValid() {
		return "", errors.Errorf("not a valid actionID: %s", action.ActionID)
	}
	return value, nil
}
