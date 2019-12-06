// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

// ActionsRuleCreate is the builder for creating a ActionsRule entity.
type ActionsRuleCreate struct {
	config
	create_time *time.Time
	update_time *time.Time
	name        *string
	triggerID   *string
	ruleFilters *[]*schema.ActionsRuleFilter
	ruleActions *[]*schema.ActionsRuleAction
}

// SetCreateTime sets the create_time field.
func (arc *ActionsRuleCreate) SetCreateTime(t time.Time) *ActionsRuleCreate {
	arc.create_time = &t
	return arc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (arc *ActionsRuleCreate) SetNillableCreateTime(t *time.Time) *ActionsRuleCreate {
	if t != nil {
		arc.SetCreateTime(*t)
	}
	return arc
}

// SetUpdateTime sets the update_time field.
func (arc *ActionsRuleCreate) SetUpdateTime(t time.Time) *ActionsRuleCreate {
	arc.update_time = &t
	return arc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (arc *ActionsRuleCreate) SetNillableUpdateTime(t *time.Time) *ActionsRuleCreate {
	if t != nil {
		arc.SetUpdateTime(*t)
	}
	return arc
}

// SetName sets the name field.
func (arc *ActionsRuleCreate) SetName(s string) *ActionsRuleCreate {
	arc.name = &s
	return arc
}

// SetTriggerID sets the triggerID field.
func (arc *ActionsRuleCreate) SetTriggerID(s string) *ActionsRuleCreate {
	arc.triggerID = &s
	return arc
}

// SetRuleFilters sets the ruleFilters field.
func (arc *ActionsRuleCreate) SetRuleFilters(srf []*schema.ActionsRuleFilter) *ActionsRuleCreate {
	arc.ruleFilters = &srf
	return arc
}

// SetRuleActions sets the ruleActions field.
func (arc *ActionsRuleCreate) SetRuleActions(sra []*schema.ActionsRuleAction) *ActionsRuleCreate {
	arc.ruleActions = &sra
	return arc
}

// Save creates the ActionsRule in the database.
func (arc *ActionsRuleCreate) Save(ctx context.Context) (*ActionsRule, error) {
	if arc.create_time == nil {
		v := actionsrule.DefaultCreateTime()
		arc.create_time = &v
	}
	if arc.update_time == nil {
		v := actionsrule.DefaultUpdateTime()
		arc.update_time = &v
	}
	if arc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if arc.triggerID == nil {
		return nil, errors.New("ent: missing required field \"triggerID\"")
	}
	if arc.ruleFilters == nil {
		return nil, errors.New("ent: missing required field \"ruleFilters\"")
	}
	if arc.ruleActions == nil {
		return nil, errors.New("ent: missing required field \"ruleActions\"")
	}
	return arc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (arc *ActionsRuleCreate) SaveX(ctx context.Context) *ActionsRule {
	v, err := arc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (arc *ActionsRuleCreate) sqlSave(ctx context.Context) (*ActionsRule, error) {
	var (
		builder = sql.Dialect(arc.driver.Dialect())
		ar      = &ActionsRule{config: arc.config}
	)
	tx, err := arc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(actionsrule.Table).Default()
	if value := arc.create_time; value != nil {
		insert.Set(actionsrule.FieldCreateTime, *value)
		ar.CreateTime = *value
	}
	if value := arc.update_time; value != nil {
		insert.Set(actionsrule.FieldUpdateTime, *value)
		ar.UpdateTime = *value
	}
	if value := arc.name; value != nil {
		insert.Set(actionsrule.FieldName, *value)
		ar.Name = *value
	}
	if value := arc.triggerID; value != nil {
		insert.Set(actionsrule.FieldTriggerID, *value)
		ar.TriggerID = *value
	}
	if value := arc.ruleFilters; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		insert.Set(actionsrule.FieldRuleFilters, buf)
		ar.RuleFilters = *value
	}
	if value := arc.ruleActions; value != nil {
		buf, err := json.Marshal(*value)
		if err != nil {
			return nil, err
		}
		insert.Set(actionsrule.FieldRuleActions, buf)
		ar.RuleActions = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(actionsrule.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	ar.ID = strconv.FormatInt(id, 10)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return ar, nil
}
