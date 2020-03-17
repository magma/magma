// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyCreate is the builder for creating a Survey entity.
type SurveyCreate struct {
	config
	mutation *SurveyMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (sc *SurveyCreate) SetCreateTime(t time.Time) *SurveyCreate {
	sc.mutation.SetCreateTime(t)
	return sc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (sc *SurveyCreate) SetNillableCreateTime(t *time.Time) *SurveyCreate {
	if t != nil {
		sc.SetCreateTime(*t)
	}
	return sc
}

// SetUpdateTime sets the update_time field.
func (sc *SurveyCreate) SetUpdateTime(t time.Time) *SurveyCreate {
	sc.mutation.SetUpdateTime(t)
	return sc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (sc *SurveyCreate) SetNillableUpdateTime(t *time.Time) *SurveyCreate {
	if t != nil {
		sc.SetUpdateTime(*t)
	}
	return sc
}

// SetName sets the name field.
func (sc *SurveyCreate) SetName(s string) *SurveyCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetOwnerName sets the owner_name field.
func (sc *SurveyCreate) SetOwnerName(s string) *SurveyCreate {
	sc.mutation.SetOwnerName(s)
	return sc
}

// SetNillableOwnerName sets the owner_name field if the given value is not nil.
func (sc *SurveyCreate) SetNillableOwnerName(s *string) *SurveyCreate {
	if s != nil {
		sc.SetOwnerName(*s)
	}
	return sc
}

// SetCreationTimestamp sets the creation_timestamp field.
func (sc *SurveyCreate) SetCreationTimestamp(t time.Time) *SurveyCreate {
	sc.mutation.SetCreationTimestamp(t)
	return sc
}

// SetNillableCreationTimestamp sets the creation_timestamp field if the given value is not nil.
func (sc *SurveyCreate) SetNillableCreationTimestamp(t *time.Time) *SurveyCreate {
	if t != nil {
		sc.SetCreationTimestamp(*t)
	}
	return sc
}

// SetCompletionTimestamp sets the completion_timestamp field.
func (sc *SurveyCreate) SetCompletionTimestamp(t time.Time) *SurveyCreate {
	sc.mutation.SetCompletionTimestamp(t)
	return sc
}

// SetLocationID sets the location edge to Location by id.
func (sc *SurveyCreate) SetLocationID(id int) *SurveyCreate {
	sc.mutation.SetLocationID(id)
	return sc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (sc *SurveyCreate) SetNillableLocationID(id *int) *SurveyCreate {
	if id != nil {
		sc = sc.SetLocationID(*id)
	}
	return sc
}

// SetLocation sets the location edge to Location.
func (sc *SurveyCreate) SetLocation(l *Location) *SurveyCreate {
	return sc.SetLocationID(l.ID)
}

// SetSourceFileID sets the source_file edge to File by id.
func (sc *SurveyCreate) SetSourceFileID(id int) *SurveyCreate {
	sc.mutation.SetSourceFileID(id)
	return sc
}

// SetNillableSourceFileID sets the source_file edge to File by id if the given value is not nil.
func (sc *SurveyCreate) SetNillableSourceFileID(id *int) *SurveyCreate {
	if id != nil {
		sc = sc.SetSourceFileID(*id)
	}
	return sc
}

// SetSourceFile sets the source_file edge to File.
func (sc *SurveyCreate) SetSourceFile(f *File) *SurveyCreate {
	return sc.SetSourceFileID(f.ID)
}

// AddQuestionIDs adds the questions edge to SurveyQuestion by ids.
func (sc *SurveyCreate) AddQuestionIDs(ids ...int) *SurveyCreate {
	sc.mutation.AddQuestionIDs(ids...)
	return sc
}

// AddQuestions adds the questions edges to SurveyQuestion.
func (sc *SurveyCreate) AddQuestions(s ...*SurveyQuestion) *SurveyCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddQuestionIDs(ids...)
}

// Save creates the Survey in the database.
func (sc *SurveyCreate) Save(ctx context.Context) (*Survey, error) {
	if _, ok := sc.mutation.CreateTime(); !ok {
		v := survey.DefaultCreateTime()
		sc.mutation.SetCreateTime(v)
	}
	if _, ok := sc.mutation.UpdateTime(); !ok {
		v := survey.DefaultUpdateTime()
		sc.mutation.SetUpdateTime(v)
	}
	if _, ok := sc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if _, ok := sc.mutation.CompletionTimestamp(); !ok {
		return nil, errors.New("ent: missing required field \"completion_timestamp\"")
	}
	var (
		err  error
		node *Survey
	)
	if len(sc.hooks) == 0 {
		node, err = sc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			sc.mutation = mutation
			node, err = sc.sqlSave(ctx)
			return node, err
		})
		for i := len(sc.hooks); i > 0; i-- {
			mut = sc.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, sc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SurveyCreate) SaveX(ctx context.Context) *Survey {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (sc *SurveyCreate) sqlSave(ctx context.Context) (*Survey, error) {
	var (
		s     = &Survey{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: survey.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: survey.FieldID,
			},
		}
	)
	if value, ok := sc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldCreateTime,
		})
		s.CreateTime = value
	}
	if value, ok := sc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldUpdateTime,
		})
		s.UpdateTime = value
	}
	if value, ok := sc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: survey.FieldName,
		})
		s.Name = value
	}
	if value, ok := sc.mutation.OwnerName(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: survey.FieldOwnerName,
		})
		s.OwnerName = value
	}
	if value, ok := sc.mutation.CreationTimestamp(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldCreationTimestamp,
		})
		s.CreationTimestamp = value
	}
	if value, ok := sc.mutation.CompletionTimestamp(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldCompletionTimestamp,
		})
		s.CompletionTimestamp = value
	}
	if nodes := sc.mutation.LocationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.LocationTable,
			Columns: []string{survey.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: location.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.SourceFileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.SourceFileTable,
			Columns: []string{survey.SourceFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: file.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := sc.mutation.QuestionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   survey.QuestionsTable,
			Columns: []string{survey.QuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: surveyquestion.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	s.ID = int(id)
	return s, nil
}
