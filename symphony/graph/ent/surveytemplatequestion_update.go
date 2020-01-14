// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateQuestionUpdate is the builder for updating SurveyTemplateQuestion entities.
type SurveyTemplateQuestionUpdate struct {
	config

	update_time          *time.Time
	question_title       *string
	question_description *string
	question_type        *string
	index                *int
	addindex             *int
	category             map[string]struct{}
	clearedCategory      bool
	predicates           []predicate.SurveyTemplateQuestion
}

// Where adds a new predicate for the builder.
func (stqu *SurveyTemplateQuestionUpdate) Where(ps ...predicate.SurveyTemplateQuestion) *SurveyTemplateQuestionUpdate {
	stqu.predicates = append(stqu.predicates, ps...)
	return stqu
}

// SetQuestionTitle sets the question_title field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionTitle(s string) *SurveyTemplateQuestionUpdate {
	stqu.question_title = &s
	return stqu
}

// SetQuestionDescription sets the question_description field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionDescription(s string) *SurveyTemplateQuestionUpdate {
	stqu.question_description = &s
	return stqu
}

// SetQuestionType sets the question_type field.
func (stqu *SurveyTemplateQuestionUpdate) SetQuestionType(s string) *SurveyTemplateQuestionUpdate {
	stqu.question_type = &s
	return stqu
}

// SetIndex sets the index field.
func (stqu *SurveyTemplateQuestionUpdate) SetIndex(i int) *SurveyTemplateQuestionUpdate {
	stqu.index = &i
	stqu.addindex = nil
	return stqu
}

// AddIndex adds i to index.
func (stqu *SurveyTemplateQuestionUpdate) AddIndex(i int) *SurveyTemplateQuestionUpdate {
	if stqu.addindex == nil {
		stqu.addindex = &i
	} else {
		*stqu.addindex += i
	}
	return stqu
}

// SetCategoryID sets the category edge to SurveyTemplateCategory by id.
func (stqu *SurveyTemplateQuestionUpdate) SetCategoryID(id string) *SurveyTemplateQuestionUpdate {
	if stqu.category == nil {
		stqu.category = make(map[string]struct{})
	}
	stqu.category[id] = struct{}{}
	return stqu
}

// SetNillableCategoryID sets the category edge to SurveyTemplateCategory by id if the given value is not nil.
func (stqu *SurveyTemplateQuestionUpdate) SetNillableCategoryID(id *string) *SurveyTemplateQuestionUpdate {
	if id != nil {
		stqu = stqu.SetCategoryID(*id)
	}
	return stqu
}

// SetCategory sets the category edge to SurveyTemplateCategory.
func (stqu *SurveyTemplateQuestionUpdate) SetCategory(s *SurveyTemplateCategory) *SurveyTemplateQuestionUpdate {
	return stqu.SetCategoryID(s.ID)
}

// ClearCategory clears the category edge to SurveyTemplateCategory.
func (stqu *SurveyTemplateQuestionUpdate) ClearCategory() *SurveyTemplateQuestionUpdate {
	stqu.clearedCategory = true
	return stqu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (stqu *SurveyTemplateQuestionUpdate) Save(ctx context.Context) (int, error) {
	if stqu.update_time == nil {
		v := surveytemplatequestion.UpdateDefaultUpdateTime()
		stqu.update_time = &v
	}
	if len(stqu.category) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return stqu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stqu *SurveyTemplateQuestionUpdate) SaveX(ctx context.Context) int {
	affected, err := stqu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (stqu *SurveyTemplateQuestionUpdate) Exec(ctx context.Context) error {
	_, err := stqu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stqu *SurveyTemplateQuestionUpdate) ExecX(ctx context.Context) {
	if err := stqu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stqu *SurveyTemplateQuestionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatequestion.Table,
			Columns: surveytemplatequestion.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: surveytemplatequestion.FieldID,
			},
		},
	}
	if ps := stqu.predicates; len(ps) > 0 {
		spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value := stqu.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveytemplatequestion.FieldUpdateTime,
		})
	}
	if value := stqu.question_title; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatequestion.FieldQuestionTitle,
		})
	}
	if value := stqu.question_description; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatequestion.FieldQuestionDescription,
		})
	}
	if value := stqu.question_type; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatequestion.FieldQuestionType,
		})
	}
	if value := stqu.index; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveytemplatequestion.FieldIndex,
		})
	}
	if value := stqu.addindex; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveytemplatequestion.FieldIndex,
		})
	}
	if stqu.clearedCategory {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatequestion.CategoryTable,
			Columns: []string{surveytemplatequestion.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := stqu.category; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatequestion.CategoryTable,
			Columns: []string{surveytemplatequestion.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return 0, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, stqu.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyTemplateQuestionUpdateOne is the builder for updating a single SurveyTemplateQuestion entity.
type SurveyTemplateQuestionUpdateOne struct {
	config
	id string

	update_time          *time.Time
	question_title       *string
	question_description *string
	question_type        *string
	index                *int
	addindex             *int
	category             map[string]struct{}
	clearedCategory      bool
}

// SetQuestionTitle sets the question_title field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionTitle(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.question_title = &s
	return stquo
}

// SetQuestionDescription sets the question_description field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionDescription(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.question_description = &s
	return stquo
}

// SetQuestionType sets the question_type field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetQuestionType(s string) *SurveyTemplateQuestionUpdateOne {
	stquo.question_type = &s
	return stquo
}

// SetIndex sets the index field.
func (stquo *SurveyTemplateQuestionUpdateOne) SetIndex(i int) *SurveyTemplateQuestionUpdateOne {
	stquo.index = &i
	stquo.addindex = nil
	return stquo
}

// AddIndex adds i to index.
func (stquo *SurveyTemplateQuestionUpdateOne) AddIndex(i int) *SurveyTemplateQuestionUpdateOne {
	if stquo.addindex == nil {
		stquo.addindex = &i
	} else {
		*stquo.addindex += i
	}
	return stquo
}

// SetCategoryID sets the category edge to SurveyTemplateCategory by id.
func (stquo *SurveyTemplateQuestionUpdateOne) SetCategoryID(id string) *SurveyTemplateQuestionUpdateOne {
	if stquo.category == nil {
		stquo.category = make(map[string]struct{})
	}
	stquo.category[id] = struct{}{}
	return stquo
}

// SetNillableCategoryID sets the category edge to SurveyTemplateCategory by id if the given value is not nil.
func (stquo *SurveyTemplateQuestionUpdateOne) SetNillableCategoryID(id *string) *SurveyTemplateQuestionUpdateOne {
	if id != nil {
		stquo = stquo.SetCategoryID(*id)
	}
	return stquo
}

// SetCategory sets the category edge to SurveyTemplateCategory.
func (stquo *SurveyTemplateQuestionUpdateOne) SetCategory(s *SurveyTemplateCategory) *SurveyTemplateQuestionUpdateOne {
	return stquo.SetCategoryID(s.ID)
}

// ClearCategory clears the category edge to SurveyTemplateCategory.
func (stquo *SurveyTemplateQuestionUpdateOne) ClearCategory() *SurveyTemplateQuestionUpdateOne {
	stquo.clearedCategory = true
	return stquo
}

// Save executes the query and returns the updated entity.
func (stquo *SurveyTemplateQuestionUpdateOne) Save(ctx context.Context) (*SurveyTemplateQuestion, error) {
	if stquo.update_time == nil {
		v := surveytemplatequestion.UpdateDefaultUpdateTime()
		stquo.update_time = &v
	}
	if len(stquo.category) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"category\"")
	}
	return stquo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (stquo *SurveyTemplateQuestionUpdateOne) SaveX(ctx context.Context) *SurveyTemplateQuestion {
	stq, err := stquo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return stq
}

// Exec executes the query on the entity.
func (stquo *SurveyTemplateQuestionUpdateOne) Exec(ctx context.Context) error {
	_, err := stquo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (stquo *SurveyTemplateQuestionUpdateOne) ExecX(ctx context.Context) {
	if err := stquo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (stquo *SurveyTemplateQuestionUpdateOne) sqlSave(ctx context.Context) (stq *SurveyTemplateQuestion, err error) {
	spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   surveytemplatequestion.Table,
			Columns: surveytemplatequestion.Columns,
			ID: &sqlgraph.FieldSpec{
				Value:  stquo.id,
				Type:   field.TypeString,
				Column: surveytemplatequestion.FieldID,
			},
		},
	}
	if value := stquo.update_time; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: surveytemplatequestion.FieldUpdateTime,
		})
	}
	if value := stquo.question_title; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatequestion.FieldQuestionTitle,
		})
	}
	if value := stquo.question_description; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatequestion.FieldQuestionDescription,
		})
	}
	if value := stquo.question_type; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: surveytemplatequestion.FieldQuestionType,
		})
	}
	if value := stquo.index; value != nil {
		spec.Fields.Set = append(spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveytemplatequestion.FieldIndex,
		})
	}
	if value := stquo.addindex; value != nil {
		spec.Fields.Add = append(spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: surveytemplatequestion.FieldIndex,
		})
	}
	if stquo.clearedCategory {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatequestion.CategoryTable,
			Columns: []string{surveytemplatequestion.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		spec.Edges.Clear = append(spec.Edges.Clear, edge)
	}
	if nodes := stquo.category; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   surveytemplatequestion.CategoryTable,
			Columns: []string{surveytemplatequestion.CategoryColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveytemplatecategory.FieldID,
				},
			},
		}
		for k, _ := range nodes {
			k, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		spec.Edges.Add = append(spec.Edges.Add, edge)
	}
	stq = &SurveyTemplateQuestion{config: stquo.config}
	spec.Assign = stq.assignValues
	spec.ScanValues = stq.scanValues()
	if err = sqlgraph.UpdateNode(ctx, stquo.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return stq, nil
}
