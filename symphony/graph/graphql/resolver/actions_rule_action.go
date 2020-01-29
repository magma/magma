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

type actionsRuleActionResolver struct{}

func (actionsRuleActionResolver) Action(ctx context.Context, ar *core.ActionsRuleAction) (*models.ActionsAction, error) {
	ac := actions.FromContext(ctx)

	action, err := ac.ActionForID(ar.ActionID)
	if err != nil {
		return nil, errors.Errorf("actionID %s not a registered action", ar.ActionID)
	}

	return &models.ActionsAction{
		Description: action.Description(),
		DataType:    action.DataType(),
	}, nil
}
