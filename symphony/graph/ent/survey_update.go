// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyUpdate is the builder for updating Survey entities.
type SurveyUpdate struct {
	config
	hooks      []Hook
	mutation   *SurveyMutation
	predicates []predicate.Survey
}

// Where adds a new predicate for the builder.
func (su *SurveyUpdate) Where(ps ...predicate.Survey) *SurveyUpdate {
	su.predicates = append(su.predicates, ps...)
	return su
}

// SetName sets the name field.
func (su *SurveyUpdate) SetName(s string) *SurveyUpdate {
	su.mutation.SetName(s)
	return su
}

// SetOwnerName sets the owner_name field.
func (su *SurveyUpdate) SetOwnerName(s string) *SurveyUpdate {
	su.mutation.SetOwnerName(s)
	return su
}

// SetNillableOwnerName sets the owner_name field if the given value is not nil.
func (su *SurveyUpdate) SetNillableOwnerName(s *string) *SurveyUpdate {
	if s != nil {
		su.SetOwnerName(*s)
	}
	return su
}

// ClearOwnerName clears the value of owner_name.
func (su *SurveyUpdate) ClearOwnerName() *SurveyUpdate {
	su.mutation.ClearOwnerName()
	return su
}

// SetCreationTimestamp sets the creation_timestamp field.
func (su *SurveyUpdate) SetCreationTimestamp(t time.Time) *SurveyUpdate {
	su.mutation.SetCreationTimestamp(t)
	return su
}

// SetNillableCreationTimestamp sets the creation_timestamp field if the given value is not nil.
func (su *SurveyUpdate) SetNillableCreationTimestamp(t *time.Time) *SurveyUpdate {
	if t != nil {
		su.SetCreationTimestamp(*t)
	}
	return su
}

// ClearCreationTimestamp clears the value of creation_timestamp.
func (su *SurveyUpdate) ClearCreationTimestamp() *SurveyUpdate {
	su.mutation.ClearCreationTimestamp()
	return su
}

// SetCompletionTimestamp sets the completion_timestamp field.
func (su *SurveyUpdate) SetCompletionTimestamp(t time.Time) *SurveyUpdate {
	su.mutation.SetCompletionTimestamp(t)
	return su
}

// SetLocationID sets the location edge to Location by id.
func (su *SurveyUpdate) SetLocationID(id int) *SurveyUpdate {
	su.mutation.SetLocationID(id)
	return su
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (su *SurveyUpdate) SetNillableLocationID(id *int) *SurveyUpdate {
	if id != nil {
		su = su.SetLocationID(*id)
	}
	return su
}

// SetLocation sets the location edge to Location.
func (su *SurveyUpdate) SetLocation(l *Location) *SurveyUpdate {
	return su.SetLocationID(l.ID)
}

// SetSourceFileID sets the source_file edge to File by id.
func (su *SurveyUpdate) SetSourceFileID(id int) *SurveyUpdate {
	su.mutation.SetSourceFileID(id)
	return su
}

// SetNillableSourceFileID sets the source_file edge to File by id if the given value is not nil.
func (su *SurveyUpdate) SetNillableSourceFileID(id *int) *SurveyUpdate {
	if id != nil {
		su = su.SetSourceFileID(*id)
	}
	return su
}

// SetSourceFile sets the source_file edge to File.
func (su *SurveyUpdate) SetSourceFile(f *File) *SurveyUpdate {
	return su.SetSourceFileID(f.ID)
}

// AddQuestionIDs adds the questions edge to SurveyQuestion by ids.
func (su *SurveyUpdate) AddQuestionIDs(ids ...int) *SurveyUpdate {
	su.mutation.AddQuestionIDs(ids...)
	return su
}

// AddQuestions adds the questions edges to SurveyQuestion.
func (su *SurveyUpdate) AddQuestions(s ...*SurveyQuestion) *SurveyUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddQuestionIDs(ids...)
}

// ClearLocation clears the location edge to Location.
func (su *SurveyUpdate) ClearLocation() *SurveyUpdate {
	su.mutation.ClearLocation()
	return su
}

// ClearSourceFile clears the source_file edge to File.
func (su *SurveyUpdate) ClearSourceFile() *SurveyUpdate {
	su.mutation.ClearSourceFile()
	return su
}

// RemoveQuestionIDs removes the questions edge to SurveyQuestion by ids.
func (su *SurveyUpdate) RemoveQuestionIDs(ids ...int) *SurveyUpdate {
	su.mutation.RemoveQuestionIDs(ids...)
	return su
}

// RemoveQuestions removes questions edges to SurveyQuestion.
func (su *SurveyUpdate) RemoveQuestions(s ...*SurveyQuestion) *SurveyUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveQuestionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (su *SurveyUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := su.mutation.UpdateTime(); !ok {
		v := survey.UpdateDefaultUpdateTime()
		su.mutation.SetUpdateTime(v)
	}

	var (
		err      error
		affected int
	)
	if len(su.hooks) == 0 {
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			su.mutation = mutation
			affected, err = su.sqlSave(ctx)
			return affected, err
		})
		for i := len(su.hooks) - 1; i >= 0; i-- {
			mut = su.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, su.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (su *SurveyUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SurveyUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SurveyUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

func (su *SurveyUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   survey.Table,
			Columns: survey.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: survey.FieldID,
			},
		},
	}
	if ps := su.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldUpdateTime,
		})
	}
	if value, ok := su.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: survey.FieldName,
		})
	}
	if value, ok := su.mutation.OwnerName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: survey.FieldOwnerName,
		})
	}
	if su.mutation.OwnerNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: survey.FieldOwnerName,
		})
	}
	if value, ok := su.mutation.CreationTimestamp(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldCreationTimestamp,
		})
	}
	if su.mutation.CreationTimestampCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: survey.FieldCreationTimestamp,
		})
	}
	if value, ok := su.mutation.CompletionTimestamp(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldCompletionTimestamp,
		})
	}
	if su.mutation.LocationCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.SourceFileCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.SourceFileIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.mutation.RemovedQuestionsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.QuestionsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{survey.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// SurveyUpdateOne is the builder for updating a single Survey entity.
type SurveyUpdateOne struct {
	config
	hooks    []Hook
	mutation *SurveyMutation
}

// SetName sets the name field.
func (suo *SurveyUpdateOne) SetName(s string) *SurveyUpdateOne {
	suo.mutation.SetName(s)
	return suo
}

// SetOwnerName sets the owner_name field.
func (suo *SurveyUpdateOne) SetOwnerName(s string) *SurveyUpdateOne {
	suo.mutation.SetOwnerName(s)
	return suo
}

// SetNillableOwnerName sets the owner_name field if the given value is not nil.
func (suo *SurveyUpdateOne) SetNillableOwnerName(s *string) *SurveyUpdateOne {
	if s != nil {
		suo.SetOwnerName(*s)
	}
	return suo
}

// ClearOwnerName clears the value of owner_name.
func (suo *SurveyUpdateOne) ClearOwnerName() *SurveyUpdateOne {
	suo.mutation.ClearOwnerName()
	return suo
}

// SetCreationTimestamp sets the creation_timestamp field.
func (suo *SurveyUpdateOne) SetCreationTimestamp(t time.Time) *SurveyUpdateOne {
	suo.mutation.SetCreationTimestamp(t)
	return suo
}

// SetNillableCreationTimestamp sets the creation_timestamp field if the given value is not nil.
func (suo *SurveyUpdateOne) SetNillableCreationTimestamp(t *time.Time) *SurveyUpdateOne {
	if t != nil {
		suo.SetCreationTimestamp(*t)
	}
	return suo
}

// ClearCreationTimestamp clears the value of creation_timestamp.
func (suo *SurveyUpdateOne) ClearCreationTimestamp() *SurveyUpdateOne {
	suo.mutation.ClearCreationTimestamp()
	return suo
}

// SetCompletionTimestamp sets the completion_timestamp field.
func (suo *SurveyUpdateOne) SetCompletionTimestamp(t time.Time) *SurveyUpdateOne {
	suo.mutation.SetCompletionTimestamp(t)
	return suo
}

// SetLocationID sets the location edge to Location by id.
func (suo *SurveyUpdateOne) SetLocationID(id int) *SurveyUpdateOne {
	suo.mutation.SetLocationID(id)
	return suo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (suo *SurveyUpdateOne) SetNillableLocationID(id *int) *SurveyUpdateOne {
	if id != nil {
		suo = suo.SetLocationID(*id)
	}
	return suo
}

// SetLocation sets the location edge to Location.
func (suo *SurveyUpdateOne) SetLocation(l *Location) *SurveyUpdateOne {
	return suo.SetLocationID(l.ID)
}

// SetSourceFileID sets the source_file edge to File by id.
func (suo *SurveyUpdateOne) SetSourceFileID(id int) *SurveyUpdateOne {
	suo.mutation.SetSourceFileID(id)
	return suo
}

// SetNillableSourceFileID sets the source_file edge to File by id if the given value is not nil.
func (suo *SurveyUpdateOne) SetNillableSourceFileID(id *int) *SurveyUpdateOne {
	if id != nil {
		suo = suo.SetSourceFileID(*id)
	}
	return suo
}

// SetSourceFile sets the source_file edge to File.
func (suo *SurveyUpdateOne) SetSourceFile(f *File) *SurveyUpdateOne {
	return suo.SetSourceFileID(f.ID)
}

// AddQuestionIDs adds the questions edge to SurveyQuestion by ids.
func (suo *SurveyUpdateOne) AddQuestionIDs(ids ...int) *SurveyUpdateOne {
	suo.mutation.AddQuestionIDs(ids...)
	return suo
}

// AddQuestions adds the questions edges to SurveyQuestion.
func (suo *SurveyUpdateOne) AddQuestions(s ...*SurveyQuestion) *SurveyUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddQuestionIDs(ids...)
}

// ClearLocation clears the location edge to Location.
func (suo *SurveyUpdateOne) ClearLocation() *SurveyUpdateOne {
	suo.mutation.ClearLocation()
	return suo
}

// ClearSourceFile clears the source_file edge to File.
func (suo *SurveyUpdateOne) ClearSourceFile() *SurveyUpdateOne {
	suo.mutation.ClearSourceFile()
	return suo
}

// RemoveQuestionIDs removes the questions edge to SurveyQuestion by ids.
func (suo *SurveyUpdateOne) RemoveQuestionIDs(ids ...int) *SurveyUpdateOne {
	suo.mutation.RemoveQuestionIDs(ids...)
	return suo
}

// RemoveQuestions removes questions edges to SurveyQuestion.
func (suo *SurveyUpdateOne) RemoveQuestions(s ...*SurveyQuestion) *SurveyUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveQuestionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (suo *SurveyUpdateOne) Save(ctx context.Context) (*Survey, error) {
	if _, ok := suo.mutation.UpdateTime(); !ok {
		v := survey.UpdateDefaultUpdateTime()
		suo.mutation.SetUpdateTime(v)
	}

	var (
		err  error
		node *Survey
	)
	if len(suo.hooks) == 0 {
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SurveyMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			suo.mutation = mutation
			node, err = suo.sqlSave(ctx)
			return node, err
		})
		for i := len(suo.hooks) - 1; i >= 0; i-- {
			mut = suo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, suo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SurveyUpdateOne) SaveX(ctx context.Context) *Survey {
	s, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return s
}

// Exec executes the query on the entity.
func (suo *SurveyUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SurveyUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (suo *SurveyUpdateOne) sqlSave(ctx context.Context) (s *Survey, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   survey.Table,
			Columns: survey.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: survey.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing Survey.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := suo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldUpdateTime,
		})
	}
	if value, ok := suo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: survey.FieldName,
		})
	}
	if value, ok := suo.mutation.OwnerName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: survey.FieldOwnerName,
		})
	}
	if suo.mutation.OwnerNameCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: survey.FieldOwnerName,
		})
	}
	if value, ok := suo.mutation.CreationTimestamp(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldCreationTimestamp,
		})
	}
	if suo.mutation.CreationTimestampCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: survey.FieldCreationTimestamp,
		})
	}
	if value, ok := suo.mutation.CompletionTimestamp(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: survey.FieldCompletionTimestamp,
		})
	}
	if suo.mutation.LocationCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.LocationIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.SourceFileCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.SourceFileIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.mutation.RemovedQuestionsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.QuestionsIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	s = &Survey{config: suo.config}
	_spec.Assign = s.assignValues
	_spec.ScanValues = s.scanValues()
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{survey.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return s, nil
}
