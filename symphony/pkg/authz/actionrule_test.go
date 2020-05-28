// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/pkg/actions/core"

	"github.com/facebookincubator/symphony/pkg/authz"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestActionsRuleCanAlwayBeWritten(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c, viewertest.WithPermissions(authz.EmptyPermissions()))
	filter := &core.ActionsRuleFilter{
		FilterID:   "ID1",
		OperatorID: "ID2",
		Data:       "data",
	}
	action := &core.ActionsRuleAction{
		ActionID: "ID3",
		Data:     "data2",
	}
	actionRule, err := c.ActionsRule.Create().
		SetName("ActionRule").
		SetTriggerID("Trigger1").
		SetRuleFilters([]*core.ActionsRuleFilter{filter}).
		SetRuleActions([]*core.ActionsRuleAction{action}).
		Save(ctx)
	require.NoError(t, err)
	err = c.ActionsRule.UpdateOne(actionRule).
		SetName("NewActionRule").
		Exec(ctx)
	require.NoError(t, err)
	err = c.ActionsRule.DeleteOne(actionRule).
		Exec(ctx)
	require.NoError(t, err)
}
