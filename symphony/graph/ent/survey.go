// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/survey"
)

// Survey is the model entity for the Survey schema.
type Survey struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// OwnerName holds the value of the "owner_name" field.
	OwnerName string `json:"owner_name,omitempty"`
	// CreationTimestamp holds the value of the "creation_timestamp" field.
	CreationTimestamp time.Time `json:"creation_timestamp,omitempty" gqlgen:"creationTimestamp"`
	// CompletionTimestamp holds the value of the "completion_timestamp" field.
	CompletionTimestamp time.Time `json:"completion_timestamp,omitempty" gqlgen:"completionTimestamp"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SurveyQuery when eager-loading is set.
	Edges struct {
		// Location holds the value of the location edge.
		Location *Location
		// SourceFile holds the value of the source_file edge.
		SourceFile *File
		// Questions holds the value of the questions edge.
		Questions []*SurveyQuestion
	} `json:"edges"`
	location_id           *string
	survey_source_file_id *string
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Survey) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // name
		&sql.NullString{}, // owner_name
		&sql.NullTime{},   // creation_timestamp
		&sql.NullTime{},   // completion_timestamp
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*Survey) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // location_id
		&sql.NullInt64{}, // survey_source_file_id
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Survey fields.
func (s *Survey) assignValues(values ...interface{}) error {
	if m, n := len(values), len(survey.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	s.ID = strconv.FormatInt(value.Int64, 10)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		s.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		s.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		s.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field owner_name", values[3])
	} else if value.Valid {
		s.OwnerName = value.String
	}
	if value, ok := values[4].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field creation_timestamp", values[4])
	} else if value.Valid {
		s.CreationTimestamp = value.Time
	}
	if value, ok := values[5].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field completion_timestamp", values[5])
	} else if value.Valid {
		s.CompletionTimestamp = value.Time
	}
	values = values[6:]
	if len(values) == len(survey.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field location_id", value)
		} else if value.Valid {
			s.location_id = new(string)
			*s.location_id = strconv.FormatInt(value.Int64, 10)
		}
		if value, ok := values[1].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field survey_source_file_id", value)
		} else if value.Valid {
			s.survey_source_file_id = new(string)
			*s.survey_source_file_id = strconv.FormatInt(value.Int64, 10)
		}
	}
	return nil
}

// QueryLocation queries the location edge of the Survey.
func (s *Survey) QueryLocation() *LocationQuery {
	return (&SurveyClient{s.config}).QueryLocation(s)
}

// QuerySourceFile queries the source_file edge of the Survey.
func (s *Survey) QuerySourceFile() *FileQuery {
	return (&SurveyClient{s.config}).QuerySourceFile(s)
}

// QueryQuestions queries the questions edge of the Survey.
func (s *Survey) QueryQuestions() *SurveyQuestionQuery {
	return (&SurveyClient{s.config}).QueryQuestions(s)
}

// Update returns a builder for updating this Survey.
// Note that, you need to call Survey.Unwrap() before calling this method, if this Survey
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Survey) Update() *SurveyUpdateOne {
	return (&SurveyClient{s.config}).UpdateOne(s)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (s *Survey) Unwrap() *Survey {
	tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Survey is not a transactional entity")
	}
	s.config.driver = tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Survey) String() string {
	var builder strings.Builder
	builder.WriteString("Survey(")
	builder.WriteString(fmt.Sprintf("id=%v", s.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(s.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(s.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(s.Name)
	builder.WriteString(", owner_name=")
	builder.WriteString(s.OwnerName)
	builder.WriteString(", creation_timestamp=")
	builder.WriteString(s.CreationTimestamp.Format(time.ANSIC))
	builder.WriteString(", completion_timestamp=")
	builder.WriteString(s.CompletionTimestamp.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (s *Survey) id() int {
	id, _ := strconv.Atoi(s.ID)
	return id
}

// Surveys is a parsable slice of Survey.
type Surveys []*Survey

func (s Surveys) config(cfg config) {
	for _i := range s {
		s[_i].config = cfg
	}
}
