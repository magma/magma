// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// ActionsRuleAction are actions that will be executed for the rule
type ActionsRuleAction struct {
	ActionID ActionID `json:"actionID"`
	Data     string   `json:"data"`
}

// ActionsRuleFilter are filters that are applied to a rule
type ActionsRuleFilter struct {
	FilterID   string `json:"filterID"`
	OperatorID string `json:"operatorID"`
	Data       string `json:"data"`
}

// Rule represents metadata for a configured rule.  This is typically
// configured by a user, and will execute when the trigger happens
type Rule struct {
	// ID is a unique id for the rule
	ID string
	// Name is the name given to this particular rule
	Name string
	// TriggerID is the trigger this rule listens for
	TriggerID TriggerID
	// RuleFilters are the filters that are specified
	RuleFilters []*ActionsRuleFilter
	// RuleActions are the actions and metadata that should be executed
	RuleActions []*ActionsRuleAction
}
