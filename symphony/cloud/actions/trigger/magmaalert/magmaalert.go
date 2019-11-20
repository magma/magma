// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package magmaalert

import (
	"github.com/facebookincubator/symphony/cloud/actions/core"
)

// trigger is the magmalert trigger
type trigger struct{}

// New returns a new trigger
func New() core.Trigger {
	return &trigger{}
}

// ID returns the string identifier for this trigger
func (*trigger) ID() core.TriggerID {
	return core.MagmaAlertTriggerID
}

// Description returns the description
func (*trigger) Description() string {
	return "alert fired relating to magma"
}

// SupportedActionIDs returns the ActionsIDs supported by this trigger
func (*trigger) SupportedActionIDs() []core.ActionID {
	return []core.ActionID{
		core.MagmaRebootNodeActionID,
	}
}

// Evaluate evaluates the user-supplied rule for if this rule
// should be executed or not
func (*trigger) Evaluate(rule core.Rule) (bool, error) {
	return true, nil
}
