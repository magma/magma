// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// Rule represents metadata for a configured rule.  This is typically
// configured by a user, and will execute when the trigger happens
type Rule struct {
	// ID is a unique id for the rule
	ID string
	// TriggerID is the trigger this rule listens for
	TriggerID TriggerID
	// ActionIDs are a list of actions to execute
	ActionIDs []ActionID
}
