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
)

// SurveyTemplateQuestion is the model entity for the SurveyTemplateQuestion schema.
type SurveyTemplateQuestion struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
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
}

// FromRows scans the sql response data into SurveyTemplateQuestion.
func (stq *SurveyTemplateQuestion) FromRows(rows *sql.Rows) error {
	var scanstq struct {
		ID                  int
		CreateTime          sql.NullTime
		UpdateTime          sql.NullTime
		QuestionTitle       sql.NullString
		QuestionDescription sql.NullString
		QuestionType        sql.NullString
		Index               sql.NullInt64
	}
	// the order here should be the same as in the `surveytemplatequestion.Columns`.
	if err := rows.Scan(
		&scanstq.ID,
		&scanstq.CreateTime,
		&scanstq.UpdateTime,
		&scanstq.QuestionTitle,
		&scanstq.QuestionDescription,
		&scanstq.QuestionType,
		&scanstq.Index,
	); err != nil {
		return err
	}
	stq.ID = strconv.Itoa(scanstq.ID)
	stq.CreateTime = scanstq.CreateTime.Time
	stq.UpdateTime = scanstq.UpdateTime.Time
	stq.QuestionTitle = scanstq.QuestionTitle.String
	stq.QuestionDescription = scanstq.QuestionDescription.String
	stq.QuestionType = scanstq.QuestionType.String
	stq.Index = int(scanstq.Index.Int64)
	return nil
}

// QueryCategory queries the category edge of the SurveyTemplateQuestion.
func (stq *SurveyTemplateQuestion) QueryCategory() *SurveyTemplateCategoryQuery {
	return (&SurveyTemplateQuestionClient{stq.config}).QueryCategory(stq)
}

// Update returns a builder for updating this SurveyTemplateQuestion.
// Note that, you need to call SurveyTemplateQuestion.Unwrap() before calling this method, if this SurveyTemplateQuestion
// was returned from a transaction, and the transaction was committed or rolled back.
func (stq *SurveyTemplateQuestion) Update() *SurveyTemplateQuestionUpdateOne {
	return (&SurveyTemplateQuestionClient{stq.config}).UpdateOne(stq)
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

// id returns the int representation of the ID field.
func (stq *SurveyTemplateQuestion) id() int {
	id, _ := strconv.Atoi(stq.ID)
	return id
}

// SurveyTemplateQuestions is a parsable slice of SurveyTemplateQuestion.
type SurveyTemplateQuestions []*SurveyTemplateQuestion

// FromRows scans the sql response data into SurveyTemplateQuestions.
func (stq *SurveyTemplateQuestions) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanstq := &SurveyTemplateQuestion{}
		if err := scanstq.FromRows(rows); err != nil {
			return err
		}
		*stq = append(*stq, scanstq)
	}
	return nil
}

func (stq SurveyTemplateQuestions) config(cfg config) {
	for _i := range stq {
		stq[_i].config = cfg
	}
}
