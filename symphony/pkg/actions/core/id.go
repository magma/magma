// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"io"
	"strconv"
)

// ActionID is an action identifier -- things that should be executed
type ActionID string

func (e ActionID) IsValid() bool {
	return e == MagmaRebootNodeActionID
}

func (e ActionID) String() string {
	return string(e)
}

func (e *ActionID) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ActionID(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ActionID", str)
	}
	return nil
}

func (e ActionID) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// TriggerID is an identifier for the trigger being executed
type TriggerID string

func (e TriggerID) IsValid() bool {
	return e == MagmaAlertTriggerID
}

func (e TriggerID) String() string {
	return string(e)
}

func (e *TriggerID) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TriggerID(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TriggerID", str)
	}
	return nil
}

func (e TriggerID) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

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
