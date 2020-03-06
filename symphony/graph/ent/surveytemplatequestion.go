// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
)

// SurveyTemplateQuestion is the model entity for the SurveyTemplateQuestion schema.
type SurveyTemplateQuestion struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// QuestionTitle holds the value of the "question_title" field.
	QuestionTitle string `json:"question_title,omitempty"`
	// QuestionDescription holds the value of the "question_description" field.
	QuestionDescription string `json:"question_description,omitempty"`
	// QuestionType holds the value of the "question_type" field.
	QuestionType string `json:"question_type,omitempty"`
	// Index holds the value of the "index" field.
	Index int `json:"index,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SurveyTemplateQuestionQuery when eager-loading is set.
	Edges                                              SurveyTemplateQuestionEdges `json:"edges"`
	survey_template_category_survey_template_questions *int
}

// SurveyTemplateQuestionEdges holds the relations/edges for other nodes in the graph.
type SurveyTemplateQuestionEdges struct {
	// Category holds the value of the category edge.
	Category *SurveyTemplateCategory
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// CategoryOrErr returns the Category value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SurveyTemplateQuestionEdges) CategoryOrErr() (*SurveyTemplateCategory, error) {
	if e.loadedTypes[0] {
		if e.Category == nil {
			// The edge category was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: surveytemplatecategory.Label}
		}
		return e.Category, nil
	}
	return nil, &NotLoadedError{edge: "category"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*SurveyTemplateQuestion) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // create_time
		&sql.NullTime{},   // update_time
		&sql.NullString{}, // question_title
		&sql.NullString{}, // question_description
		&sql.NullString{}, // question_type
		&sql.NullInt64{},  // index
	}
}

// fkValues returns the types for scanning foreign-keys values from sql.Rows.
func (*SurveyTemplateQuestion) fkValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{}, // survey_template_category_survey_template_questions
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the SurveyTemplateQuestion fields.
func (stq *SurveyTemplateQuestion) assignValues(values ...interface{}) error {
	if m, n := len(values), len(surveytemplatequestion.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	stq.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field create_time", values[0])
	} else if value.Valid {
		stq.CreateTime = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field update_time", values[1])
	} else if value.Valid {
		stq.UpdateTime = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field question_title", values[2])
	} else if value.Valid {
		stq.QuestionTitle = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field question_description", values[3])
	} else if value.Valid {
		stq.QuestionDescription = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field question_type", values[4])
	} else if value.Valid {
		stq.QuestionType = value.String
	}
	if value, ok := values[5].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field index", values[5])
	} else if value.Valid {
		stq.Index = int(value.Int64)
	}
	values = values[6:]
	if len(values) == len(surveytemplatequestion.ForeignKeys) {
		if value, ok := values[0].(*sql.NullInt64); !ok {
			return fmt.Errorf("unexpected type %T for edge-field survey_template_category_survey_template_questions", value)
		} else if value.Valid {
			stq.survey_template_category_survey_template_questions = new(int)
			*stq.survey_template_category_survey_template_questions = int(value.Int64)
		}
	}
	return nil
}

// QueryCategory queries the category edge of the SurveyTemplateQuestion.
func (stq *SurveyTemplateQuestion) QueryCategory() *SurveyTemplateCategoryQuery {
	return (&SurveyTemplateQuestionClient{config: stq.config}).QueryCategory(stq)
}

// Update returns a builder for updating this SurveyTemplateQuestion.
// Note that, you need to call SurveyTemplateQuestion.Unwrap() before calling this method, if this SurveyTemplateQuestion
// was returned from a transaction, and the transaction was committed or rolled back.
func (stq *SurveyTemplateQuestion) Update() *SurveyTemplateQuestionUpdateOne {
	return (&SurveyTemplateQuestionClient{config: stq.config}).UpdateOne(stq)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (stq *SurveyTemplateQuestion) Unwrap() *SurveyTemplateQuestion {
	tx, ok := stq.config.driver.(*txDriver)
	if !ok {
		panic("ent: SurveyTemplateQuestion is not a transactional entity")
	}
	stq.config.driver = tx.drv
	return stq
}

// String implements the fmt.Stringer.
func (stq *SurveyTemplateQuestion) String() string {
	var builder strings.Builder
	builder.WriteString("SurveyTemplateQuestion(")
	builder.WriteString(fmt.Sprintf("id=%v", stq.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(stq.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(stq.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", question_title=")
	builder.WriteString(stq.QuestionTitle)
	builder.WriteString(", question_description=")
	builder.WriteString(stq.QuestionDescription)
	builder.WriteString(", question_type=")
	builder.WriteString(stq.QuestionType)
	builder.WriteString(", index=")
	builder.WriteString(fmt.Sprintf("%v", stq.Index))
	builder.WriteByte(')')
	return builder.String()
}

// SurveyTemplateQuestions is a parsable slice of SurveyTemplateQuestion.
type SurveyTemplateQuestions []*SurveyTemplateQuestion

func (stq SurveyTemplateQuestions) config(cfg config) {
	for _i := range stq {
		stq[_i].config = cfg
	}
}
