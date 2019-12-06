// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

// ActionsRuleUpdate is the builder for updating ActionsRule entities.
type ActionsRuleUpdate struct {
	config

	update_time *time.Time
	name        *string
	triggerID   *string
	ruleFilters *[]*schema.ActionsRuleFilter
	ruleActions *[]*schema.ActionsRuleAction
	predicates  []predicate.ActionsRule
}

// Where adds a new predicate for the builder.
func (aru *ActionsRuleUpdate) Where(ps ...predicate.ActionsRule) *ActionsRuleUpdate {
	aru.predicates = append(aru.predicates, ps...)
	return aru
}

// SetName sets the name field.
func (aru *ActionsRuleUpdate) SetName(s string) *ActionsRuleUpdate {
	aru.name = &s
	return aru
}

// SetTriggerID sets the triggerID field.
func (aru *ActionsRuleUpdate) SetTriggerID(s string) *ActionsRuleUpdate {
	aru.triggerID = &s
	return aru
}

// SetRuleFilters sets the ruleFilters field.
func (aru *ActionsRuleUpdate) SetRuleFilters(srf []*schema.ActionsRuleFilter) *ActionsRuleUpdate {
	aru.ruleFilters = &srf
	return aru
}

// SetRuleActions sets the ruleActions field.
func (aru *ActionsRuleUpdate) SetRuleActions(sra []*schema.ActionsRuleAction) *ActionsRuleUpdate {
	aru.ruleActions = &sra
	return aru
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (aru *ActionsRuleUpdate) Save(ctx context.Context) (int, error) {
	if aru.update_time == nil {
		v := actionsrule.UpdateDefaultUpdateTime()
		aru.update_time = &v
	}
	return aru.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (aru *ActionsRuleUpdate) SaveX(ctx context.Context) int {
	affected, err := aru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (aru *ActionsRuleUpdate) Exec(ctx context.Context) error {
	_, err := aru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aru *ActionsRuleUpdate) ExecX(ctx context.Context) {
	if err := aru.Exec(ctx); err != nil {
		panic(err)
	}
}

func (aru *ActionsRuleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(aru.driver.Dialect())
		selector = builder.Select(actionsrule.FieldID).From(builder.Table(actionsrule.Table))
	)
	for _, p := range aru.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = aru.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := aru.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(actionsrule.Table)
	)
	updater = updater.Where(sql.InInts(actionsrule.FieldID, ids...))
	if value := aru.update_time; value != nil {
		updater.Set(actionsrule.FieldUpdateTime, *value)
	}
	if value := aru.name; value != nil {
		updater.Set(actionsrule.FieldName, *value)
	}
	if value := aru.triggerID; value != nil {
		updater.Set(actionsrule.FieldTriggerID, *value)
	}
	if value := aru.ruleFilters; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return 0, err
		}
		updater.Set(actionsrule.FieldRuleFilters, buf)
	}
	if value := aru.ruleActions; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return 0, err
		}
		updater.Set(actionsrule.FieldRuleActions, buf)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// ActionsRuleUpdateOne is the builder for updating a single ActionsRule entity.
type ActionsRuleUpdateOne struct {
	config
	id string

	update_time *time.Time
	name        *string
	triggerID   *string
	ruleFilters *[]*schema.ActionsRuleFilter
	ruleActions *[]*schema.ActionsRuleAction
}

// SetName sets the name field.
func (aruo *ActionsRuleUpdateOne) SetName(s string) *ActionsRuleUpdateOne {
	aruo.name = &s
	return aruo
}

// SetTriggerID sets the triggerID field.
func (aruo *ActionsRuleUpdateOne) SetTriggerID(s string) *ActionsRuleUpdateOne {
	aruo.triggerID = &s
	return aruo
}

// SetRuleFilters sets the ruleFilters field.
func (aruo *ActionsRuleUpdateOne) SetRuleFilters(srf []*schema.ActionsRuleFilter) *ActionsRuleUpdateOne {
	aruo.ruleFilters = &srf
	return aruo
}

// SetRuleActions sets the ruleActions field.
func (aruo *ActionsRuleUpdateOne) SetRuleActions(sra []*schema.ActionsRuleAction) *ActionsRuleUpdateOne {
	aruo.ruleActions = &sra
	return aruo
}

// Save executes the query and returns the updated entity.
func (aruo *ActionsRuleUpdateOne) Save(ctx context.Context) (*ActionsRule, error) {
	if aruo.update_time == nil {
		v := actionsrule.UpdateDefaultUpdateTime()
		aruo.update_time = &v
	}
	return aruo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (aruo *ActionsRuleUpdateOne) SaveX(ctx context.Context) *ActionsRule {
	ar, err := aruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return ar
}

// Exec executes the query on the entity.
func (aruo *ActionsRuleUpdateOne) Exec(ctx context.Context) error {
	_, err := aruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aruo *ActionsRuleUpdateOne) ExecX(ctx context.Context) {
	if err := aruo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (aruo *ActionsRuleUpdateOne) sqlSave(ctx context.Context) (ar *ActionsRule, err error) {
	var (
		builder  = sql.Dialect(aruo.driver.Dialect())
		selector = builder.Select(actionsrule.Columns...).From(builder.Table(actionsrule.Table))
	)
	actionsrule.ID(aruo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = aruo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		ar = &ActionsRule{config: aruo.config}
		if err := ar.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into ActionsRule: %v", err)
		}
		id = ar.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("ActionsRule with id: %v", aruo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one ActionsRule with the same id: %v", aruo.id)
	}

	tx, err := aruo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(actionsrule.Table)
	)
	updater = updater.Where(sql.InInts(actionsrule.FieldID, ids...))
	if value := aruo.update_time; value != nil {
		updater.Set(actionsrule.FieldUpdateTime, *value)
		ar.UpdateTime = *value
	}
	if value := aruo.name; value != nil {
		updater.Set(actionsrule.FieldName, *value)
		ar.Name = *value
	}
	if value := aruo.triggerID; value != nil {
		updater.Set(actionsrule.FieldTriggerID, *value)
		ar.TriggerID = *value
	}
	if value := aruo.ruleFilters; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		updater.Set(actionsrule.FieldRuleFilters, buf)
		ar.RuleFilters = *value
	}
	if value := aruo.ruleActions; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		updater.Set(actionsrule.FieldRuleActions, buf)
		ar.RuleActions = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return ar, nil
}
