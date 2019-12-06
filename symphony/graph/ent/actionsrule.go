// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

// ActionsRule is the model entity for the ActionsRule schema.
type ActionsRule struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// TriggerID holds the value of the "triggerID" field.
	TriggerID string `json:"triggerID,omitempty"`
	// RuleFilters holds the value of the "ruleFilters" field.
	RuleFilters []*schema.ActionsRuleFilter `json:"ruleFilters,omitempty"`
	// RuleActions holds the value of the "ruleActions" field.
	RuleActions []*schema.ActionsRuleAction `json:"ruleActions,omitempty"`
}

// FromRows scans the sql response data into ActionsRule.
func (ar *ActionsRule) FromRows(rows *sql.Rows) error {
	var scanar struct {
		ID          int
		CreateTime  sql.NullTime
		UpdateTime  sql.NullTime
		Name        sql.NullString
		TriggerID   sql.NullString
		RuleFilters []byte
		RuleActions []byte
	}
	// the order here should be the same as in the `actionsrule.Columns`.
	if err := rows.Scan(
		&scanar.ID,
		&scanar.CreateTime,
		&scanar.UpdateTime,
		&scanar.Name,
		&scanar.TriggerID,
		&scanar.RuleFilters,
		&scanar.RuleActions,
	); err != nil {
		return err
	}
	ar.ID = strconv.Itoa(scanar.ID)
	ar.CreateTime = scanar.CreateTime.Time
	ar.UpdateTime = scanar.UpdateTime.Time
	ar.Name = scanar.Name.String
	ar.TriggerID = scanar.TriggerID.String
	if value := scanar.RuleFilters; len(value) > 0 {
		if err := json.Unmarshal(value, &ar.RuleFilters); err != nil {
			return fmt.Errorf("unmarshal field ruleFilters: %v", err)
		}
	}
	if value := scanar.RuleActions; len(value) > 0 {
		if err := json.Unmarshal(value, &ar.RuleActions); err != nil {
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

// id returns the int representation of the ID field.
func (ar *ActionsRule) id() int {
	id, _ := strconv.Atoi(ar.ID)
	return id
}

// ActionsRules is a parsable slice of ActionsRule.
type ActionsRules []*ActionsRule

// FromRows scans the sql response data into ActionsRules.
func (ar *ActionsRules) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanar := &ActionsRule{}
		if err := scanar.FromRows(rows); err != nil {
			return err
		}
		*ar = append(*ar, scanar)
	}
	return nil
}

func (ar ActionsRules) config(cfg config) {
	for _i := range ar {
		ar[_i].config = cfg
	}
}
