// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// ActionID is an action identifier -- things that should be executed
type ActionID string

// TriggerID is an identifier for the trigger being executed
type TriggerID string

const (
	// Actions

	// MagmaRebootNodeActionID is the id for magmarebootnode
	MagmaRebootNodeActionID ActionID = "magma_reboot_node"

	// Triggers

	// MagmaAlertTriggerID is the id for magmaalert
	MagmaAlertTriggerID TriggerID = "magma_alert"
)

var (
	// AllTriggerIDs contains all core triggers
	AllTriggerIDs = []TriggerID{
		MagmaAlertTriggerID,
	}

	// AllActionIDs contains all core actions
	AllActionIDs = []ActionID{
		MagmaRebootNodeActionID,
	}
)
