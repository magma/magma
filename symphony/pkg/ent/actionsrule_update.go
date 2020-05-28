// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/ent/actionsrule"
	"github.com/facebookincubator/symphony/pkg/ent/predicate"
)

// ActionsRuleUpdate is the builder for updating ActionsRule entities.
type ActionsRuleUpdate struct {
	config
	hooks      []Hook
	mutation   *ActionsRuleMutation
	predicates []predicate.ActionsRule
}

// Where adds a new predicate for the builder.
func (aru *ActionsRuleUpdate) Where(ps ...predicate.ActionsRule) *ActionsRuleUpdate {
	aru.predicates = append(aru.predicates, ps...)
	return aru
}

// SetName sets the name field.
func (aru *ActionsRuleUpdate) SetName(s string) *ActionsRuleUpdate {
	aru.mutation.SetName(s)
	return aru
}

// SetTriggerID sets the triggerID field.
func (aru *ActionsRuleUpdate) SetTriggerID(s string) *ActionsRuleUpdate {
	aru.mutation.SetTriggerID(s)
	return aru
}

// SetRuleFilters sets the ruleFilters field.
func (aru *ActionsRuleUpdate) SetRuleFilters(crf []*core.ActionsRuleFilter) *ActionsRuleUpdate {
	aru.mutation.SetRuleFilters(crf)
	return aru
}

// SetRuleActions sets the ruleActions field.
func (aru *ActionsRuleUpdate) SetRuleActions(cra []*core.ActionsRuleAction) *ActionsRuleUpdate {
	aru.mutation.SetRuleActions(cra)
	return aru
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (aru *ActionsRuleUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := aru.mutation.UpdateTime(); !ok {
		v := actionsrule.UpdateDefaultUpdateTime()
		aru.mutation.SetUpdateTime(v)
	}
	var (
		err      error
		affected int
	)
	if len(aru.hooks) == 0 {
		affected, err = aru.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ActionsRuleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			aru.mutation = mutation
			affected, err = aru.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(aru.hooks) - 1; i >= 0; i-- {
			mut = aru.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, aru.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   actionsrule.Table,
			Columns: actionsrule.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: actionsrule.FieldID,
			},
		},
	}
	if ps := aru.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := aru.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: actionsrule.FieldUpdateTime,
		})
	}
	if value, ok := aru.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: actionsrule.FieldName,
		})
	}
	if value, ok := aru.mutation.TriggerID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: actionsrule.FieldTriggerID,
		})
	}
	if value, ok := aru.mutation.RuleFilters(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: actionsrule.FieldRuleFilters,
		})
	}
	if value, ok := aru.mutation.RuleActions(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: actionsrule.FieldRuleActions,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, aru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{actionsrule.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ActionsRuleUpdateOne is the builder for updating a single ActionsRule entity.
type ActionsRuleUpdateOne struct {
	config
	hooks    []Hook
	mutation *ActionsRuleMutation
}

// SetName sets the name field.
func (aruo *ActionsRuleUpdateOne) SetName(s string) *ActionsRuleUpdateOne {
	aruo.mutation.SetName(s)
	return aruo
}

// SetTriggerID sets the triggerID field.
func (aruo *ActionsRuleUpdateOne) SetTriggerID(s string) *ActionsRuleUpdateOne {
	aruo.mutation.SetTriggerID(s)
	return aruo
}

// SetRuleFilters sets the ruleFilters field.
func (aruo *ActionsRuleUpdateOne) SetRuleFilters(crf []*core.ActionsRuleFilter) *ActionsRuleUpdateOne {
	aruo.mutation.SetRuleFilters(crf)
	return aruo
}

// SetRuleActions sets the ruleActions field.
func (aruo *ActionsRuleUpdateOne) SetRuleActions(cra []*core.ActionsRuleAction) *ActionsRuleUpdateOne {
	aruo.mutation.SetRuleActions(cra)
	return aruo
}

// Save executes the query and returns the updated entity.
func (aruo *ActionsRuleUpdateOne) Save(ctx context.Context) (*ActionsRule, error) {
	if _, ok := aruo.mutation.UpdateTime(); !ok {
		v := actionsrule.UpdateDefaultUpdateTime()
		aruo.mutation.SetUpdateTime(v)
	}
	var (
		err  error
		node *ActionsRule
	)
	if len(aruo.hooks) == 0 {
		node, err = aruo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ActionsRuleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			aruo.mutation = mutation
			node, err = aruo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(aruo.hooks) - 1; i >= 0; i-- {
			mut = aruo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, aruo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   actionsrule.Table,
			Columns: actionsrule.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: actionsrule.FieldID,
			},
		},
	}
	id, ok := aruo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing ActionsRule.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := aruo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: actionsrule.FieldUpdateTime,
		})
	}
	if value, ok := aruo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: actionsrule.FieldName,
		})
	}
	if value, ok := aruo.mutation.TriggerID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: actionsrule.FieldTriggerID,
		})
	}
	if value, ok := aruo.mutation.RuleFilters(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: actionsrule.FieldRuleFilters,
		})
	}
	if value, ok := aruo.mutation.RuleActions(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: actionsrule.FieldRuleActions,
		})
	}
	ar = &ActionsRule{config: aruo.config}
	_spec.Assign = ar.assignValues
	_spec.ScanValues = ar.scanValues()
	if err = sqlgraph.UpdateNode(ctx, aruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{actionsrule.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return ar, nil
}
