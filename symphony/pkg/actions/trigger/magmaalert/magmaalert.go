// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package magmaalert

import (
	"github.com/facebookincubator/symphony/pkg/actions/core"
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
	return "an orc8r alert is fired"
}

// SupportedActionIDs returns the ActionsIDs supported by this trigger
func (*trigger) SupportedActionIDs() []core.ActionID {
	return []core.ActionID{
		core.MagmaRebootNodeActionID,
	}
}

func (*trigger) SupportedFilters() []core.Filter {
	return []core.Filter{
		core.NewStringFieldFilter(
			"alertname",
			"the alert's name",
		),
		core.NewStringFieldFilter(
			"networkID",
			"the alert's networkID",
		),
		core.NewStringFieldFilter(
			"gatewayID",
			"the alert's gatewayID",
		),
	}
}
