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
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
)

// SurveyUpdate is the builder for updating Survey entities.
type SurveyUpdate struct {
	config

	update_time          *time.Time
	name                 *string
	owner_name           *string
	clearowner_name      bool
	completion_timestamp *time.Time
	location             map[string]struct{}
	source_file          map[string]struct{}
	questions            map[string]struct{}
	clearedLocation      bool
	clearedSourceFile    bool
	removedQuestions     map[string]struct{}
	predicates           []predicate.Survey
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
	var (
		builder  = sql.Dialect(su.driver.Dialect())
		selector = builder.Select(survey.FieldID).From(builder.Table(survey.Table))
	)
	for _, p := range su.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = su.driver.Query(ctx, query, args, rows); err != nil {
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

	tx, err := su.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(survey.Table)
	)
	updater = updater.Where(sql.InInts(survey.FieldID, ids...))
	if value := su.update_time; value != nil {
		updater.Set(survey.FieldUpdateTime, *value)
	}
	if value := su.name; value != nil {
		updater.Set(survey.FieldName, *value)
	}
	if value := su.owner_name; value != nil {
		updater.Set(survey.FieldOwnerName, *value)
	}
	if su.clearowner_name {
		updater.SetNull(survey.FieldOwnerName)
	}
	if value := su.completion_timestamp; value != nil {
		updater.Set(survey.FieldCompletionTimestamp, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if su.clearedLocation {
		query, args := builder.Update(survey.LocationTable).
			SetNull(survey.LocationColumn).
			Where(sql.InInts(location.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.location) > 0 {
		for eid := range su.location {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(survey.LocationTable).
				Set(survey.LocationColumn, eid).
				Where(sql.InInts(survey.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if su.clearedSourceFile {
		query, args := builder.Update(survey.SourceFileTable).
			SetNull(survey.SourceFileColumn).
			Where(sql.InInts(file.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.source_file) > 0 {
		for eid := range su.source_file {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(survey.SourceFileTable).
				Set(survey.SourceFileColumn, eid).
				Where(sql.InInts(survey.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if len(su.removedQuestions) > 0 {
		eids := make([]int, len(su.removedQuestions))
		for eid := range su.removedQuestions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(survey.QuestionsTable).
			SetNull(survey.QuestionsColumn).
			Where(sql.InInts(survey.QuestionsColumn, ids...)).
			Where(sql.InInts(surveyquestion.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(su.questions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range su.questions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
				}
				p.Or().EQ(surveyquestion.FieldID, eid)
			}
			query, args := builder.Update(survey.QuestionsTable).
				Set(survey.QuestionsColumn, id).
				Where(sql.And(p, sql.IsNull(survey.QuestionsColumn))).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
			affected, err := res.RowsAffected()
			if err != nil {
				return 0, rollback(tx, err)
			}
			if int(affected) < len(su.questions) {
				return 0, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"questions\" %v already connected to a different \"Survey\"", keys(su.questions))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// SurveyUpdateOne is the builder for updating a single Survey entity.
type SurveyUpdateOne struct {
	config
	id string

	update_time          *time.Time
	name                 *string
	owner_name           *string
	clearowner_name      bool
	completion_timestamp *time.Time
	location             map[string]struct{}
	source_file          map[string]struct{}
	questions            map[string]struct{}
	clearedLocation      bool
	clearedSourceFile    bool
	removedQuestions     map[string]struct{}
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
	var (
		builder  = sql.Dialect(suo.driver.Dialect())
		selector = builder.Select(survey.Columns...).From(builder.Table(survey.Table))
	)
	survey.ID(suo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = suo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		s = &Survey{config: suo.config}
		if err := s.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Survey: %v", err)
		}
		id = s.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Survey with id: %v", suo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Survey with the same id: %v", suo.id)
	}

	tx, err := suo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(survey.Table)
	)
	updater = updater.Where(sql.InInts(survey.FieldID, ids...))
	if value := suo.update_time; value != nil {
		updater.Set(survey.FieldUpdateTime, *value)
		s.UpdateTime = *value
	}
	if value := suo.name; value != nil {
		updater.Set(survey.FieldName, *value)
		s.Name = *value
	}
	if value := suo.owner_name; value != nil {
		updater.Set(survey.FieldOwnerName, *value)
		s.OwnerName = *value
	}
	if suo.clearowner_name {
		var value string
		s.OwnerName = value
		updater.SetNull(survey.FieldOwnerName)
	}
	if value := suo.completion_timestamp; value != nil {
		updater.Set(survey.FieldCompletionTimestamp, *value)
		s.CompletionTimestamp = *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if suo.clearedLocation {
		query, args := builder.Update(survey.LocationTable).
			SetNull(survey.LocationColumn).
			Where(sql.InInts(location.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.location) > 0 {
		for eid := range suo.location {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(survey.LocationTable).
				Set(survey.LocationColumn, eid).
				Where(sql.InInts(survey.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if suo.clearedSourceFile {
		query, args := builder.Update(survey.SourceFileTable).
			SetNull(survey.SourceFileColumn).
			Where(sql.InInts(file.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.source_file) > 0 {
		for eid := range suo.source_file {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			query, args := builder.Update(survey.SourceFileTable).
				Set(survey.SourceFileColumn, eid).
				Where(sql.InInts(survey.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(suo.removedQuestions) > 0 {
		eids := make([]int, len(suo.removedQuestions))
		for eid := range suo.removedQuestions {
			eid, serr := strconv.Atoi(eid)
			if serr != nil {
				err = rollback(tx, serr)
				return
			}
			eids = append(eids, eid)
		}
		query, args := builder.Update(survey.QuestionsTable).
			SetNull(survey.QuestionsColumn).
			Where(sql.InInts(survey.QuestionsColumn, ids...)).
			Where(sql.InInts(surveyquestion.FieldID, eids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(suo.questions) > 0 {
		for _, id := range ids {
			p := sql.P()
			for eid := range suo.questions {
				eid, serr := strconv.Atoi(eid)
				if serr != nil {
					err = rollback(tx, serr)
					return
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
			if int(affected) < len(suo.questions) {
				return nil, rollback(tx, &ConstraintError{msg: fmt.Sprintf("one of \"questions\" %v already connected to a different \"Survey\"", keys(suo.questions))})
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return s, nil
}
