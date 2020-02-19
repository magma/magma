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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyUpdate is the builder for updating Survey entities.
type SurveyUpdate struct {
	config

	update_time             *time.Time
	name                    *string
	owner_name              *string
	clearowner_name         bool
	creation_timestamp      *time.Time
	clearcreation_timestamp bool
	completion_timestamp    *time.Time
	location                map[string]struct{}
	source_file             map[string]struct{}
	questions               map[string]struct{}
	clearedLocation         bool
	clearedSourceFile       bool
	removedQuestions        map[string]struct{}
	predicates              []predicate.Survey
}

// Where adds a new predicate for the builder.
func (su *SurveyUpdate) Where(ps ...predicate.Survey) *SurveyUpdate {
	su.predicates = append(su.predicates, ps...)
	return su
}

// SetName sets the name field.
func (su *SurveyUpdate) SetName(s string) *SurveyUpdate {
	su.name = &s
	return su
}

// SetOwnerName sets the owner_name field.
func (su *SurveyUpdate) SetOwnerName(s string) *SurveyUpdate {
	su.owner_name = &s
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
	su.owner_name = nil
	su.clearowner_name = true
	return su
}

// SetCreationTimestamp sets the creation_timestamp field.
func (su *SurveyUpdate) SetCreationTimestamp(t time.Time) *SurveyUpdate {
	su.creation_timestamp = &t
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
	su.creation_timestamp = nil
	su.clearcreation_timestamp = true
	return su
}

// SetCompletionTimestamp sets the completion_timestamp field.
func (su *SurveyUpdate) SetCompletionTimestamp(t time.Time) *SurveyUpdate {
	su.completion_timestamp = &t
	return su
}

// SetLocationID sets the location edge to Location by id.
func (su *SurveyUpdate) SetLocationID(id string) *SurveyUpdate {
	if su.location == nil {
		su.location = make(map[string]struct{})
	}
	su.location[id] = struct{}{}
	return su
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (su *SurveyUpdate) SetNillableLocationID(id *string) *SurveyUpdate {
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
func (su *SurveyUpdate) SetSourceFileID(id string) *SurveyUpdate {
	if su.source_file == nil {
		su.source_file = make(map[string]struct{})
	}
	su.source_file[id] = struct{}{}
	return su
}

// SetNillableSourceFileID sets the source_file edge to File by id if the given value is not nil.
func (su *SurveyUpdate) SetNillableSourceFileID(id *string) *SurveyUpdate {
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
func (su *SurveyUpdate) AddQuestionIDs(ids ...string) *SurveyUpdate {
	if su.questions == nil {
		su.questions = make(map[string]struct{})
	}
	for i := range ids {
		su.questions[ids[i]] = struct{}{}
	}
	return su
}

// AddQuestions adds the questions edges to SurveyQuestion.
func (su *SurveyUpdate) AddQuestions(s ...*SurveyQuestion) *SurveyUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddQuestionIDs(ids...)
}

// ClearLocation clears the location edge to Location.
func (su *SurveyUpdate) ClearLocation() *SurveyUpdate {
	su.clearedLocation = true
	return su
}

// ClearSourceFile clears the source_file edge to File.
func (su *SurveyUpdate) ClearSourceFile() *SurveyUpdate {
	su.clearedSourceFile = true
	return su
}

// RemoveQuestionIDs removes the questions edge to SurveyQuestion by ids.
func (su *SurveyUpdate) RemoveQuestionIDs(ids ...string) *SurveyUpdate {
	if su.removedQuestions == nil {
		su.removedQuestions = make(map[string]struct{})
	}
	for i := range ids {
		su.removedQuestions[ids[i]] = struct{}{}
	}
	return su
}

// RemoveQuestions removes questions edges to SurveyQuestion.
func (su *SurveyUpdate) RemoveQuestions(s ...*SurveyQuestion) *SurveyUpdate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveQuestionIDs(ids...)
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (su *SurveyUpdate) Save(ctx context.Context) (int, error) {
	if su.update_time == nil {
		v := survey.UpdateDefaultUpdateTime()
		su.update_time = &v
	}
	if len(su.location) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(su.source_file) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"source_file\"")
	}
	return su.sqlSave(ctx)
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
				Type:   field.TypeString,
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
	if value := su.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldUpdateTime,
		})
	}
	if value := su.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: survey.FieldName,
		})
	}
	if value := su.owner_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: survey.FieldOwnerName,
		})
	}
	if su.clearowner_name {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: survey.FieldOwnerName,
		})
	}
	if value := su.creation_timestamp; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldCreationTimestamp,
		})
	}
	if su.clearcreation_timestamp {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: survey.FieldCreationTimestamp,
		})
	}
	if value := su.completion_timestamp; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldCompletionTimestamp,
		})
	}
	if su.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.LocationTable,
			Columns: []string{survey.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.LocationTable,
			Columns: []string{survey.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.clearedSourceFile {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.SourceFileTable,
			Columns: []string{survey.SourceFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.source_file; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.SourceFileTable,
			Columns: []string{survey.SourceFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := su.removedQuestions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   survey.QuestionsTable,
			Columns: []string{survey.QuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.questions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   survey.QuestionsTable,
			Columns: []string{survey.QuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
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
	id string

	update_time             *time.Time
	name                    *string
	owner_name              *string
	clearowner_name         bool
	creation_timestamp      *time.Time
	clearcreation_timestamp bool
	completion_timestamp    *time.Time
	location                map[string]struct{}
	source_file             map[string]struct{}
	questions               map[string]struct{}
	clearedLocation         bool
	clearedSourceFile       bool
	removedQuestions        map[string]struct{}
}

// SetName sets the name field.
func (suo *SurveyUpdateOne) SetName(s string) *SurveyUpdateOne {
	suo.name = &s
	return suo
}

// SetOwnerName sets the owner_name field.
func (suo *SurveyUpdateOne) SetOwnerName(s string) *SurveyUpdateOne {
	suo.owner_name = &s
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
	suo.owner_name = nil
	suo.clearowner_name = true
	return suo
}

// SetCreationTimestamp sets the creation_timestamp field.
func (suo *SurveyUpdateOne) SetCreationTimestamp(t time.Time) *SurveyUpdateOne {
	suo.creation_timestamp = &t
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
	suo.creation_timestamp = nil
	suo.clearcreation_timestamp = true
	return suo
}

// SetCompletionTimestamp sets the completion_timestamp field.
func (suo *SurveyUpdateOne) SetCompletionTimestamp(t time.Time) *SurveyUpdateOne {
	suo.completion_timestamp = &t
	return suo
}

// SetLocationID sets the location edge to Location by id.
func (suo *SurveyUpdateOne) SetLocationID(id string) *SurveyUpdateOne {
	if suo.location == nil {
		suo.location = make(map[string]struct{})
	}
	suo.location[id] = struct{}{}
	return suo
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (suo *SurveyUpdateOne) SetNillableLocationID(id *string) *SurveyUpdateOne {
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
func (suo *SurveyUpdateOne) SetSourceFileID(id string) *SurveyUpdateOne {
	if suo.source_file == nil {
		suo.source_file = make(map[string]struct{})
	}
	suo.source_file[id] = struct{}{}
	return suo
}

// SetNillableSourceFileID sets the source_file edge to File by id if the given value is not nil.
func (suo *SurveyUpdateOne) SetNillableSourceFileID(id *string) *SurveyUpdateOne {
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
func (suo *SurveyUpdateOne) AddQuestionIDs(ids ...string) *SurveyUpdateOne {
	if suo.questions == nil {
		suo.questions = make(map[string]struct{})
	}
	for i := range ids {
		suo.questions[ids[i]] = struct{}{}
	}
	return suo
}

// AddQuestions adds the questions edges to SurveyQuestion.
func (suo *SurveyUpdateOne) AddQuestions(s ...*SurveyQuestion) *SurveyUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddQuestionIDs(ids...)
}

// ClearLocation clears the location edge to Location.
func (suo *SurveyUpdateOne) ClearLocation() *SurveyUpdateOne {
	suo.clearedLocation = true
	return suo
}

// ClearSourceFile clears the source_file edge to File.
func (suo *SurveyUpdateOne) ClearSourceFile() *SurveyUpdateOne {
	suo.clearedSourceFile = true
	return suo
}

// RemoveQuestionIDs removes the questions edge to SurveyQuestion by ids.
func (suo *SurveyUpdateOne) RemoveQuestionIDs(ids ...string) *SurveyUpdateOne {
	if suo.removedQuestions == nil {
		suo.removedQuestions = make(map[string]struct{})
	}
	for i := range ids {
		suo.removedQuestions[ids[i]] = struct{}{}
	}
	return suo
}

// RemoveQuestions removes questions edges to SurveyQuestion.
func (suo *SurveyUpdateOne) RemoveQuestions(s ...*SurveyQuestion) *SurveyUpdateOne {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveQuestionIDs(ids...)
}

// Save executes the query and returns the updated entity.
func (suo *SurveyUpdateOne) Save(ctx context.Context) (*Survey, error) {
	if suo.update_time == nil {
		v := survey.UpdateDefaultUpdateTime()
		suo.update_time = &v
	}
	if len(suo.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(suo.source_file) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"source_file\"")
	}
	return suo.sqlSave(ctx)
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
				Value:  suo.id,
				Type:   field.TypeString,
				Column: survey.FieldID,
			},
		},
	}
	if value := suo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldUpdateTime,
		})
	}
	if value := suo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: survey.FieldName,
		})
	}
	if value := suo.owner_name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: survey.FieldOwnerName,
		})
	}
	if suo.clearowner_name {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: survey.FieldOwnerName,
		})
	}
	if value := suo.creation_timestamp; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldCreationTimestamp,
		})
	}
	if suo.clearcreation_timestamp {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: survey.FieldCreationTimestamp,
		})
	}
	if value := suo.completion_timestamp; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldCompletionTimestamp,
		})
	}
	if suo.clearedLocation {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.LocationTable,
			Columns: []string{survey.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.location; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.LocationTable,
			Columns: []string{survey.LocationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: location.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.clearedSourceFile {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.SourceFileTable,
			Columns: []string{survey.SourceFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.source_file; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   survey.SourceFileTable,
			Columns: []string{survey.SourceFileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: file.FieldID,
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if nodes := suo.removedQuestions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   survey.QuestionsTable,
			Columns: []string{survey.QuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.questions; len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   survey.QuestionsTable,
			Columns: []string{survey.QuestionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: surveyquestion.FieldID,
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
