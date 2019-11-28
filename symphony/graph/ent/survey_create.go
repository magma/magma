// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
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
		res     sql.Result
		builder = sql.Dialect(sc.driver.Dialect())
		s       = &Survey{config: sc.config}
	)
	tx, err := sc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(survey.Table).Default()
	if value := sc.create_time; value != nil {
		insert.Set(survey.FieldCreateTime, *value)
		s.CreateTime = *value
	}
	if value := sc.update_time; value != nil {
		insert.Set(survey.FieldUpdateTime, *value)
		s.UpdateTime = *value
	}
	if value := sc.name; value != nil {
		insert.Set(survey.FieldName, *value)
		s.Name = *value
	}
	if value := sc.owner_name; value != nil {
		insert.Set(survey.FieldOwnerName, *value)
		s.OwnerName = *value
	}
	if value := sc.completion_timestamp; value != nil {
		insert.Set(survey.FieldCompletionTimestamp, *value)
		s.CompletionTimestamp = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(survey.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	s.ID = strconv.FormatInt(id, 10)
	if len(sc.location) > 0 {
		for eid := range sc.location {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(survey.LocationTable).
				Set(survey.LocationColumn, eid).
				Where(sql.EQ(survey.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(sc.source_file) > 0 {
		for eid := range sc.source_file {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(survey.SourceFileTable).
				Set(survey.SourceFileColumn, eid).
				Where(sql.EQ(survey.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(sc.questions) > 0 {
		p := sql.P()
		for eid := range sc.questions {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(surveyquestion.FieldID, eid)
		}
		query, args := builder.Update(survey.QuestionsTable).
			Set(survey.QuestionsColumn, id).
			Where(sql.And(p, sql.IsNull(survey.QuestionsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(sc.questions) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"questions\" %v already connected to a different \"Survey\"", keys(sc.questions))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s, nil
}
