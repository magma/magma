// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"context"

	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/pkg/errors"
)

// Executor will execute all Actions defined in Rules for a Trigger
type Executor struct {
	Context context.Context
	Registry
	DataLoader
	onError func(error)
}

// Execute runs all workflows for the specified object/trigger
func (exc Executor) Execute(ctx context.Context, objectID string, triggerToPayload map[core.TriggerID]map[string]interface{}) {

	// Note that we should keep this interface serializable, so if we need to eventually
	// offload this to workers, we can

	for triggerID, inputPayload := range triggerToPayload {
		trigger, err := exc.Registry.TriggerForID(triggerID)
		if err != nil {
			// TODO: Should we bail here, or just log an error and continue
			exc.onError(errors.Errorf("could not find trigger: %s", triggerID))
			continue
		}

		for _, rule := range exc.DataLoader.QueryRules(triggerID) {
			shouldExecute, err := trigger.Evaluate(rule)
			if err != nil {
				exc.onError(errors.Errorf("evaluating rule %s: %v", rule.ID, err))
				continue
			}
			if !shouldExecute {
				continue
			}
			for _, actionID := range rule.ActionIDs {
				err := exc.executeAction(rule, actionID, inputPayload)
				if err != nil {
					exc.onError(errors.Errorf("executing action %s: %v", actionID, err))
				}
			}
		}
	}
}

func (exc Executor) executeAction(rule core.Rule, actionID core.ActionID, inputPayload map[string]interface{}) error {
	action, err := exc.Registry.ActionForID(actionID)
	if err != nil {
		return errors.Errorf("could not find action %v, skipping: %v", actionID, err)
	}
	actionContext := core.ActionContext{
		TriggerPayload: inputPayload,
		Rule:           rule,
	}
	err = action.Execute(actionContext)
	if err != nil {
		return errors.Errorf("executing %v: %v", actionID, err)
	}
	return nil
}
