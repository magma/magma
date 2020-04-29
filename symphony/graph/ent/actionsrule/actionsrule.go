// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package actionsrule

import (
	"time"

	"github.com/facebookincubator/ent"
)

const (
	// Label holds the string label denoting the actionsrule type in the database.
	Label = "actions_rule"
	// FieldID holds the string denoting the id field in the database.
	FieldID          = "id"           // FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime  = "create_time"  // FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime  = "update_time"  // FieldName holds the string denoting the name vertex property in the database.
	FieldName        = "name"         // FieldTriggerID holds the string denoting the triggerid vertex property in the database.
	FieldTriggerID   = "trigger_id"   // FieldRuleFilters holds the string denoting the rulefilters vertex property in the database.
	FieldRuleFilters = "rule_filters" // FieldRuleActions holds the string denoting the ruleactions vertex property in the database.
	FieldRuleActions = "rule_actions"

	// Table holds the table name of the actionsrule in the database.
	Table = "actions_rules"
)

// Columns holds all SQL columns for actionsrule fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldTriggerID,
	FieldRuleFilters,
	FieldRuleActions,
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/facebookincubator/symphony/graph/ent/runtime"
//
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime func() time.Time
)
