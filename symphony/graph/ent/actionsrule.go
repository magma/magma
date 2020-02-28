// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/pkg/actions/core"
)

// ActionsRule is the model entity for the ActionsRule schema.
type ActionsRule struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// TriggerID holds the value of the "triggerID" field.
	TriggerID string `json:"triggerID,omitempty"`
	// RuleFilters holds the value of the "ruleFilters" field.
	RuleFilters []*core.ActionsRuleFilter `json:"ruleFilters,omitempty"`
	// RuleActions holds the value of the "ruleActions" field.
	RuleActions []*core.ActionsRuleAction `json:"ruleActions,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ActionsRule) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // triggerID
		&[]byte{},         // ruleFilters
		&[]byte{},         // ruleActions
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ActionsRule fields.
func (ar *ActionsRule) assignValues(values ...interface{}) error {
	if m, n := len(values), len(actionsrule.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	ar.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		ar.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		ar.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		ar.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field triggerID", values[3])
	} else if value.Valid {
		ar.TriggerID = value.String
	}

	if value, ok := values[4].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field ruleFilters", values[4])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &ar.RuleFilters); err != nil {
			return fmt.Errorf("unmarshal field ruleFilters: %v", err)
		}
	}

	if value, ok := values[5].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field ruleActions", values[5])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &ar.RuleActions); err != nil {
			return fmt.Errorf("unmarshal field ruleActions: %v", err)
		}
	}
	return nil
}

// Update returns a builder for updating this ActionsRule.
// Note that, you need to call ActionsRule.Unwrap() before calling this method, if this ActionsRule
// was returned from a transaction, and the transaction was committed or rolled back.
func (ar *ActionsRule) Update() *ActionsRuleUpdateOne {
	return (&ActionsRuleClient{ar.config}).UpdateOne(ar)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (ar *ActionsRule) Unwrap() *ActionsRule {
	tx, ok := ar.config.driver.(*txDriver)
	if !ok {
		panic("ent: ActionsRule is not a transactional entity")
	}
	ar.config.driver = tx.drv
	return ar
}

// String implements the fmt.Stringer.
func (ar *ActionsRule) String() string {
	var builder strings.Builder
	builder.WriteString("ActionsRule(")
	builder.WriteString(fmt.Sprintf("id=%v", ar.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(ar.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(ar.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(ar.Name)
	builder.WriteString(", triggerID=")
	builder.WriteString(ar.TriggerID)
	builder.WriteString(", ruleFilters=")
	builder.WriteString(fmt.Sprintf("%v", ar.RuleFilters))
	builder.WriteString(", ruleActions=")
	builder.WriteString(fmt.Sprintf("%v", ar.RuleActions))
	builder.WriteByte(')')
	return builder.String()
}

// ActionsRules is a parsable slice of ActionsRule.
type ActionsRules []*ActionsRule

func (ar ActionsRules) config(cfg config) {
	for _i := range ar {
		ar[_i].config = cfg
	}
}
