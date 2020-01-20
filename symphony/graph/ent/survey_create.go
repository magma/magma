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
	create_time          *time.Time
	update_time          *time.Time
	name                 *string
	owner_name           *string
	creation_timestamp   *time.Time
	completion_timestamp *time.Time
	location             map[string]struct{}
	source_file          map[string]struct{}
	questions            map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (sc *SurveyCreate) SetCreateTime(t time.Time) *SurveyCreate {
	sc.create_time = &t
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
	sc.update_time = &t
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
	sc.name = &s
	return sc
}

// SetOwnerName sets the owner_name field.
func (sc *SurveyCreate) SetOwnerName(s string) *SurveyCreate {
	sc.owner_name = &s
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
	sc.creation_timestamp = &t
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
	sc.completion_timestamp = &t
	return sc
}

// SetLocationID sets the location edge to Location by id.
func (sc *SurveyCreate) SetLocationID(id string) *SurveyCreate {
	if sc.location == nil {
		sc.location = make(map[string]struct{})
	}
	sc.location[id] = struct{}{}
	return sc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (sc *SurveyCreate) SetNillableLocationID(id *string) *SurveyCreate {
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
func (sc *SurveyCreate) SetSourceFileID(id string) *SurveyCreate {
	if sc.source_file == nil {
		sc.source_file = make(map[string]struct{})
	}
	sc.source_file[id] = struct{}{}
	return sc
}

// SetNillableSourceFileID sets the source_file edge to File by id if the given value is not nil.
func (sc *SurveyCreate) SetNillableSourceFileID(id *string) *SurveyCreate {
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
func (sc *SurveyCreate) AddQuestionIDs(ids ...string) *SurveyCreate {
	if sc.questions == nil {
		sc.questions = make(map[string]struct{})
	}
	for i := range ids {
		sc.questions[ids[i]] = struct{}{}
	}
	return sc
}

// AddQuestions adds the questions edges to SurveyQuestion.
func (sc *SurveyCreate) AddQuestions(s ...*SurveyQuestion) *SurveyCreate {
	ids := make([]string, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return sc.AddQuestionIDs(ids...)
}

// Save creates the Survey in the database.
func (sc *SurveyCreate) Save(ctx context.Context) (*Survey, error) {
	if sc.create_time == nil {
		v := survey.DefaultCreateTime()
		sc.create_time = &v
	}
	if sc.update_time == nil {
		v := survey.DefaultUpdateTime()
		sc.update_time = &v
	}
	if sc.name == nil {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if sc.completion_timestamp == nil {
		return nil, errors.New("ent: missing required field \"completion_timestamp\"")
	}
	if len(sc.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	if len(sc.source_file) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"source_file\"")
	}
	return sc.sqlSave(ctx)
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
		s    = &Survey{config: sc.config}
		spec = &sqlgraph.CreateSpec{
			Table: survey.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: survey.FieldID,
			},
		}
	)
	if value := sc.create_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldCreateTime,
		})
		s.CreateTime = *value
	}
	if value := sc.update_time; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldUpdateTime,
		})
		s.UpdateTime = *value
	}
	if value := sc.name; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: survey.FieldName,
		})
		s.Name = *value
	}
	if value := sc.owner_name; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: survey.FieldOwnerName,
		})
		s.OwnerName = *value
	}
	if value := sc.creation_timestamp; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldCreationTimestamp,
		})
		s.CreationTimestamp = *value
	}
	if value := sc.completion_timestamp; value != nil {
		spec.Fields = append(spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: survey.FieldCompletionTimestamp,
		})
		s.CompletionTimestamp = *value
	}
	if nodes := sc.location; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := sc.source_file; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if nodes := sc.questions; len(nodes) > 0 {
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
		spec.Edges = append(spec.Edges, edge)
	}
	if err := sqlgraph.CreateNode(ctx, sc.driver, spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := spec.ID.Value.(int64)
	s.ID = strconv.FormatInt(id, 10)
	return s, nil
}
