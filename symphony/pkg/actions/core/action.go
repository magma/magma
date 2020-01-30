// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// Action is an interface for implementing an action
type Action interface {
	ID() ActionID
	Description() string
	Execute(ActionContext) error
	DataType() DataType
}

// ActionContext is a context passed to the an actions Execute method
type ActionContext struct {
	// TriggerPayload is the initial data used for the trigger
	TriggerPayload map[string]interface{}
	// Rule is the user-configured rule this action is executing for
	Rule Rule
	// RuleAction is the particular action + data this action is executing for
	RuleAction *ActionsRuleAction
}
