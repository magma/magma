// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/pkg/actions/core"
)

// ActionsRuleUpdate is the builder for updating ActionsRule entities.
type ActionsRuleUpdate struct {
	config

	update_time *time.Time
	name        *string
	triggerID   *string
	ruleFilters *[]*core.ActionsRuleFilter
	ruleActions *[]*core.ActionsRuleAction
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
func (aru *ActionsRuleUpdate) SetRuleFilters(crf []*core.ActionsRuleFilter) *ActionsRuleUpdate {
	aru.ruleFilters = &crf
	return aru
}

// SetRuleActions sets the ruleActions field.
func (aru *ActionsRuleUpdate) SetRuleActions(cra []*core.ActionsRuleAction) *ActionsRuleUpdate {
	aru.ruleActions = &cra
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
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   actionsrule.Table,
			Columns: actionsrule.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: actionsrule.FieldID,
			},
		},
	}
	if ps := aru.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := aru.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: actionsrule.FieldUpdateTime,
		})
	}
	if value := aru.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: actionsrule.FieldName,
		})
	}
	if value := aru.triggerID; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: actionsrule.FieldTriggerID,
		})
	}
	if value := aru.ruleFilters; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: actionsrule.FieldRuleFilters,
		})
	}
	if value := aru.ruleActions; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: actionsrule.FieldRuleActions,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, aru.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ActionsRuleUpdateOne is the builder for updating a single ActionsRule entity.
type ActionsRuleUpdateOne struct {
	config
	id string

	update_time *time.Time
	name        *string
	triggerID   *string
	ruleFilters *[]*core.ActionsRuleFilter
	ruleActions *[]*core.ActionsRuleAction
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
func (aruo *ActionsRuleUpdateOne) SetRuleFilters(crf []*core.ActionsRuleFilter) *ActionsRuleUpdateOne {
	aruo.ruleFilters = &crf
	return aruo
}

// SetRuleActions sets the ruleActions field.
func (aruo *ActionsRuleUpdateOne) SetRuleActions(cra []*core.ActionsRuleAction) *ActionsRuleUpdateOne {
	aruo.ruleActions = &cra
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
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   actionsrule.Table,
			Columns: actionsrule.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  aruo.id,
				Type:   field.TypeString,
				Column: actionsrule.FieldID,
			},
		},
	}
	if value := aruo.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: actionsrule.FieldUpdateTime,
		})
	}
	if value := aruo.name; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: actionsrule.FieldName,
		})
	}
	if value := aruo.triggerID; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: actionsrule.FieldTriggerID,
		})
	}
	if value := aruo.ruleFilters; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: actionsrule.FieldRuleFilters,
		})
	}
	if value := aruo.ruleActions; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  *value,
			Column: actionsrule.FieldRuleActions,
		})
	}
	ar = &ActionsRule{config: aruo.config}
	spec.Assign = ar.assignValues
	spec.ScanValues = ar.scanValues()
	if err = sqlgraph.UpdateNode(ctx, aruo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ar, nil
}
