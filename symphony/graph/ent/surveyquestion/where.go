// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveyquestion

import (
	"strconv"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i], _ = strconv.Atoi(ids[i])
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i], _ = strconv.Atoi(ids[i])
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		id, _ := strconv.Atoi(id)
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreateTime applies equality check predicate on the "create_time" field. It's identical to CreateTimeEQ.
func CreateTime(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// UpdateTime applies equality check predicate on the "update_time" field. It's identical to UpdateTimeEQ.
func UpdateTime(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// FormName applies equality check predicate on the "form_name" field. It's identical to FormNameEQ.
func FormName(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFormName), v))
	})
}

// FormDescription applies equality check predicate on the "form_description" field. It's identical to FormDescriptionEQ.
func FormDescription(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFormDescription), v))
	})
}

// FormIndex applies equality check predicate on the "form_index" field. It's identical to FormIndexEQ.
func FormIndex(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFormIndex), v))
	})
}

// QuestionType applies equality check predicate on the "question_type" field. It's identical to QuestionTypeEQ.
func QuestionType(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionType), v))
	})
}

// QuestionFormat applies equality check predicate on the "question_format" field. It's identical to QuestionFormatEQ.
func QuestionFormat(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionFormat), v))
	})
}

// QuestionText applies equality check predicate on the "question_text" field. It's identical to QuestionTextEQ.
func QuestionText(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionText), v))
	})
}

// QuestionIndex applies equality check predicate on the "question_index" field. It's identical to QuestionIndexEQ.
func QuestionIndex(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionIndex), v))
	})
}

// BoolData applies equality check predicate on the "bool_data" field. It's identical to BoolDataEQ.
func BoolData(v bool) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBoolData), v))
	})
}

// EmailData applies equality check predicate on the "email_data" field. It's identical to EmailDataEQ.
func EmailData(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEmailData), v))
	})
}

// Latitude applies equality check predicate on the "latitude" field. It's identical to LatitudeEQ.
func Latitude(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitude), v))
	})
}

// Longitude applies equality check predicate on the "longitude" field. It's identical to LongitudeEQ.
func Longitude(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitude), v))
	})
}

// LocationAccuracy applies equality check predicate on the "location_accuracy" field. It's identical to LocationAccuracyEQ.
func LocationAccuracy(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLocationAccuracy), v))
	})
}

// Altitude applies equality check predicate on the "altitude" field. It's identical to AltitudeEQ.
func Altitude(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAltitude), v))
	})
}

// PhoneData applies equality check predicate on the "phone_data" field. It's identical to PhoneDataEQ.
func PhoneData(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPhoneData), v))
	})
}

// TextData applies equality check predicate on the "text_data" field. It's identical to TextDataEQ.
func TextData(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTextData), v))
	})
}

// FloatData applies equality check predicate on the "float_data" field. It's identical to FloatDataEQ.
func FloatData(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFloatData), v))
	})
}

// IntData applies equality check predicate on the "int_data" field. It's identical to IntDataEQ.
func IntData(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIntData), v))
	})
}

// DateData applies equality check predicate on the "date_data" field. It's identical to DateDataEQ.
func DateData(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDateData), v))
	})
}

// CreateTimeEQ applies the EQ predicate on the "create_time" field.
func CreateTimeEQ(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeNEQ applies the NEQ predicate on the "create_time" field.
func CreateTimeNEQ(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreateTime), v))
	})
}

// CreateTimeIn applies the In predicate on the "create_time" field.
func CreateTimeIn(vs ...time.Time) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreateTime), v...))
	})
}

// CreateTimeNotIn applies the NotIn predicate on the "create_time" field.
func CreateTimeNotIn(vs ...time.Time) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreateTime), v...))
	})
}

// CreateTimeGT applies the GT predicate on the "create_time" field.
func CreateTimeGT(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeGTE applies the GTE predicate on the "create_time" field.
func CreateTimeGTE(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLT applies the LT predicate on the "create_time" field.
func CreateTimeLT(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreateTime), v))
	})
}

// CreateTimeLTE applies the LTE predicate on the "create_time" field.
func CreateTimeLTE(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreateTime), v))
	})
}

// UpdateTimeEQ applies the EQ predicate on the "update_time" field.
func UpdateTimeEQ(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeNEQ applies the NEQ predicate on the "update_time" field.
func UpdateTimeNEQ(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeIn applies the In predicate on the "update_time" field.
func UpdateTimeIn(vs ...time.Time) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdateTime), v...))
	})
}

// UpdateTimeNotIn applies the NotIn predicate on the "update_time" field.
func UpdateTimeNotIn(vs ...time.Time) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdateTime), v...))
	})
}

// UpdateTimeGT applies the GT predicate on the "update_time" field.
func UpdateTimeGT(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeGTE applies the GTE predicate on the "update_time" field.
func UpdateTimeGTE(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLT applies the LT predicate on the "update_time" field.
func UpdateTimeLT(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdateTime), v))
	})
}

// UpdateTimeLTE applies the LTE predicate on the "update_time" field.
func UpdateTimeLTE(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdateTime), v))
	})
}

// FormNameEQ applies the EQ predicate on the "form_name" field.
func FormNameEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFormName), v))
	})
}

// FormNameNEQ applies the NEQ predicate on the "form_name" field.
func FormNameNEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFormName), v))
	})
}

// FormNameIn applies the In predicate on the "form_name" field.
func FormNameIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldFormName), v...))
	})
}

// FormNameNotIn applies the NotIn predicate on the "form_name" field.
func FormNameNotIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldFormName), v...))
	})
}

// FormNameGT applies the GT predicate on the "form_name" field.
func FormNameGT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldFormName), v))
	})
}

// FormNameGTE applies the GTE predicate on the "form_name" field.
func FormNameGTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldFormName), v))
	})
}

// FormNameLT applies the LT predicate on the "form_name" field.
func FormNameLT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldFormName), v))
	})
}

// FormNameLTE applies the LTE predicate on the "form_name" field.
func FormNameLTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldFormName), v))
	})
}

// FormNameContains applies the Contains predicate on the "form_name" field.
func FormNameContains(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldFormName), v))
	})
}

// FormNameHasPrefix applies the HasPrefix predicate on the "form_name" field.
func FormNameHasPrefix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldFormName), v))
	})
}

// FormNameHasSuffix applies the HasSuffix predicate on the "form_name" field.
func FormNameHasSuffix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldFormName), v))
	})
}

// FormNameIsNil applies the IsNil predicate on the "form_name" field.
func FormNameIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldFormName)))
	})
}

// FormNameNotNil applies the NotNil predicate on the "form_name" field.
func FormNameNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldFormName)))
	})
}

// FormNameEqualFold applies the EqualFold predicate on the "form_name" field.
func FormNameEqualFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldFormName), v))
	})
}

// FormNameContainsFold applies the ContainsFold predicate on the "form_name" field.
func FormNameContainsFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldFormName), v))
	})
}

// FormDescriptionEQ applies the EQ predicate on the "form_description" field.
func FormDescriptionEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionNEQ applies the NEQ predicate on the "form_description" field.
func FormDescriptionNEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionIn applies the In predicate on the "form_description" field.
func FormDescriptionIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldFormDescription), v...))
	})
}

// FormDescriptionNotIn applies the NotIn predicate on the "form_description" field.
func FormDescriptionNotIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldFormDescription), v...))
	})
}

// FormDescriptionGT applies the GT predicate on the "form_description" field.
func FormDescriptionGT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionGTE applies the GTE predicate on the "form_description" field.
func FormDescriptionGTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionLT applies the LT predicate on the "form_description" field.
func FormDescriptionLT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionLTE applies the LTE predicate on the "form_description" field.
func FormDescriptionLTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionContains applies the Contains predicate on the "form_description" field.
func FormDescriptionContains(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionHasPrefix applies the HasPrefix predicate on the "form_description" field.
func FormDescriptionHasPrefix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionHasSuffix applies the HasSuffix predicate on the "form_description" field.
func FormDescriptionHasSuffix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionIsNil applies the IsNil predicate on the "form_description" field.
func FormDescriptionIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldFormDescription)))
	})
}

// FormDescriptionNotNil applies the NotNil predicate on the "form_description" field.
func FormDescriptionNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldFormDescription)))
	})
}

// FormDescriptionEqualFold applies the EqualFold predicate on the "form_description" field.
func FormDescriptionEqualFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldFormDescription), v))
	})
}

// FormDescriptionContainsFold applies the ContainsFold predicate on the "form_description" field.
func FormDescriptionContainsFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldFormDescription), v))
	})
}

// FormIndexEQ applies the EQ predicate on the "form_index" field.
func FormIndexEQ(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFormIndex), v))
	})
}

// FormIndexNEQ applies the NEQ predicate on the "form_index" field.
func FormIndexNEQ(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFormIndex), v))
	})
}

// FormIndexIn applies the In predicate on the "form_index" field.
func FormIndexIn(vs ...int) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldFormIndex), v...))
	})
}

// FormIndexNotIn applies the NotIn predicate on the "form_index" field.
func FormIndexNotIn(vs ...int) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldFormIndex), v...))
	})
}

// FormIndexGT applies the GT predicate on the "form_index" field.
func FormIndexGT(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldFormIndex), v))
	})
}

// FormIndexGTE applies the GTE predicate on the "form_index" field.
func FormIndexGTE(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldFormIndex), v))
	})
}

// FormIndexLT applies the LT predicate on the "form_index" field.
func FormIndexLT(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldFormIndex), v))
	})
}

// FormIndexLTE applies the LTE predicate on the "form_index" field.
func FormIndexLTE(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldFormIndex), v))
	})
}

// QuestionTypeEQ applies the EQ predicate on the "question_type" field.
func QuestionTypeEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeNEQ applies the NEQ predicate on the "question_type" field.
func QuestionTypeNEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeIn applies the In predicate on the "question_type" field.
func QuestionTypeIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldQuestionType), v...))
	})
}

// QuestionTypeNotIn applies the NotIn predicate on the "question_type" field.
func QuestionTypeNotIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldQuestionType), v...))
	})
}

// QuestionTypeGT applies the GT predicate on the "question_type" field.
func QuestionTypeGT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeGTE applies the GTE predicate on the "question_type" field.
func QuestionTypeGTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeLT applies the LT predicate on the "question_type" field.
func QuestionTypeLT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeLTE applies the LTE predicate on the "question_type" field.
func QuestionTypeLTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeContains applies the Contains predicate on the "question_type" field.
func QuestionTypeContains(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeHasPrefix applies the HasPrefix predicate on the "question_type" field.
func QuestionTypeHasPrefix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeHasSuffix applies the HasSuffix predicate on the "question_type" field.
func QuestionTypeHasSuffix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeIsNil applies the IsNil predicate on the "question_type" field.
func QuestionTypeIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldQuestionType)))
	})
}

// QuestionTypeNotNil applies the NotNil predicate on the "question_type" field.
func QuestionTypeNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldQuestionType)))
	})
}

// QuestionTypeEqualFold applies the EqualFold predicate on the "question_type" field.
func QuestionTypeEqualFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldQuestionType), v))
	})
}

// QuestionTypeContainsFold applies the ContainsFold predicate on the "question_type" field.
func QuestionTypeContainsFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldQuestionType), v))
	})
}

// QuestionFormatEQ applies the EQ predicate on the "question_format" field.
func QuestionFormatEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatNEQ applies the NEQ predicate on the "question_format" field.
func QuestionFormatNEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatIn applies the In predicate on the "question_format" field.
func QuestionFormatIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldQuestionFormat), v...))
	})
}

// QuestionFormatNotIn applies the NotIn predicate on the "question_format" field.
func QuestionFormatNotIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldQuestionFormat), v...))
	})
}

// QuestionFormatGT applies the GT predicate on the "question_format" field.
func QuestionFormatGT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatGTE applies the GTE predicate on the "question_format" field.
func QuestionFormatGTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatLT applies the LT predicate on the "question_format" field.
func QuestionFormatLT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatLTE applies the LTE predicate on the "question_format" field.
func QuestionFormatLTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatContains applies the Contains predicate on the "question_format" field.
func QuestionFormatContains(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatHasPrefix applies the HasPrefix predicate on the "question_format" field.
func QuestionFormatHasPrefix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatHasSuffix applies the HasSuffix predicate on the "question_format" field.
func QuestionFormatHasSuffix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatIsNil applies the IsNil predicate on the "question_format" field.
func QuestionFormatIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldQuestionFormat)))
	})
}

// QuestionFormatNotNil applies the NotNil predicate on the "question_format" field.
func QuestionFormatNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldQuestionFormat)))
	})
}

// QuestionFormatEqualFold applies the EqualFold predicate on the "question_format" field.
func QuestionFormatEqualFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldQuestionFormat), v))
	})
}

// QuestionFormatContainsFold applies the ContainsFold predicate on the "question_format" field.
func QuestionFormatContainsFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldQuestionFormat), v))
	})
}

// QuestionTextEQ applies the EQ predicate on the "question_text" field.
func QuestionTextEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionText), v))
	})
}

// QuestionTextNEQ applies the NEQ predicate on the "question_text" field.
func QuestionTextNEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldQuestionText), v))
	})
}

// QuestionTextIn applies the In predicate on the "question_text" field.
func QuestionTextIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldQuestionText), v...))
	})
}

// QuestionTextNotIn applies the NotIn predicate on the "question_text" field.
func QuestionTextNotIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldQuestionText), v...))
	})
}

// QuestionTextGT applies the GT predicate on the "question_text" field.
func QuestionTextGT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldQuestionText), v))
	})
}

// QuestionTextGTE applies the GTE predicate on the "question_text" field.
func QuestionTextGTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldQuestionText), v))
	})
}

// QuestionTextLT applies the LT predicate on the "question_text" field.
func QuestionTextLT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldQuestionText), v))
	})
}

// QuestionTextLTE applies the LTE predicate on the "question_text" field.
func QuestionTextLTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldQuestionText), v))
	})
}

// QuestionTextContains applies the Contains predicate on the "question_text" field.
func QuestionTextContains(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldQuestionText), v))
	})
}

// QuestionTextHasPrefix applies the HasPrefix predicate on the "question_text" field.
func QuestionTextHasPrefix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldQuestionText), v))
	})
}

// QuestionTextHasSuffix applies the HasSuffix predicate on the "question_text" field.
func QuestionTextHasSuffix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldQuestionText), v))
	})
}

// QuestionTextIsNil applies the IsNil predicate on the "question_text" field.
func QuestionTextIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldQuestionText)))
	})
}

// QuestionTextNotNil applies the NotNil predicate on the "question_text" field.
func QuestionTextNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldQuestionText)))
	})
}

// QuestionTextEqualFold applies the EqualFold predicate on the "question_text" field.
func QuestionTextEqualFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldQuestionText), v))
	})
}

// QuestionTextContainsFold applies the ContainsFold predicate on the "question_text" field.
func QuestionTextContainsFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldQuestionText), v))
	})
}

// QuestionIndexEQ applies the EQ predicate on the "question_index" field.
func QuestionIndexEQ(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldQuestionIndex), v))
	})
}

// QuestionIndexNEQ applies the NEQ predicate on the "question_index" field.
func QuestionIndexNEQ(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldQuestionIndex), v))
	})
}

// QuestionIndexIn applies the In predicate on the "question_index" field.
func QuestionIndexIn(vs ...int) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldQuestionIndex), v...))
	})
}

// QuestionIndexNotIn applies the NotIn predicate on the "question_index" field.
func QuestionIndexNotIn(vs ...int) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldQuestionIndex), v...))
	})
}

// QuestionIndexGT applies the GT predicate on the "question_index" field.
func QuestionIndexGT(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldQuestionIndex), v))
	})
}

// QuestionIndexGTE applies the GTE predicate on the "question_index" field.
func QuestionIndexGTE(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldQuestionIndex), v))
	})
}

// QuestionIndexLT applies the LT predicate on the "question_index" field.
func QuestionIndexLT(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldQuestionIndex), v))
	})
}

// QuestionIndexLTE applies the LTE predicate on the "question_index" field.
func QuestionIndexLTE(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldQuestionIndex), v))
	})
}

// BoolDataEQ applies the EQ predicate on the "bool_data" field.
func BoolDataEQ(v bool) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldBoolData), v))
	})
}

// BoolDataNEQ applies the NEQ predicate on the "bool_data" field.
func BoolDataNEQ(v bool) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldBoolData), v))
	})
}

// BoolDataIsNil applies the IsNil predicate on the "bool_data" field.
func BoolDataIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldBoolData)))
	})
}

// BoolDataNotNil applies the NotNil predicate on the "bool_data" field.
func BoolDataNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldBoolData)))
	})
}

// EmailDataEQ applies the EQ predicate on the "email_data" field.
func EmailDataEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEmailData), v))
	})
}

// EmailDataNEQ applies the NEQ predicate on the "email_data" field.
func EmailDataNEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEmailData), v))
	})
}

// EmailDataIn applies the In predicate on the "email_data" field.
func EmailDataIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldEmailData), v...))
	})
}

// EmailDataNotIn applies the NotIn predicate on the "email_data" field.
func EmailDataNotIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldEmailData), v...))
	})
}

// EmailDataGT applies the GT predicate on the "email_data" field.
func EmailDataGT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldEmailData), v))
	})
}

// EmailDataGTE applies the GTE predicate on the "email_data" field.
func EmailDataGTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldEmailData), v))
	})
}

// EmailDataLT applies the LT predicate on the "email_data" field.
func EmailDataLT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldEmailData), v))
	})
}

// EmailDataLTE applies the LTE predicate on the "email_data" field.
func EmailDataLTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldEmailData), v))
	})
}

// EmailDataContains applies the Contains predicate on the "email_data" field.
func EmailDataContains(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldEmailData), v))
	})
}

// EmailDataHasPrefix applies the HasPrefix predicate on the "email_data" field.
func EmailDataHasPrefix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldEmailData), v))
	})
}

// EmailDataHasSuffix applies the HasSuffix predicate on the "email_data" field.
func EmailDataHasSuffix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldEmailData), v))
	})
}

// EmailDataIsNil applies the IsNil predicate on the "email_data" field.
func EmailDataIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldEmailData)))
	})
}

// EmailDataNotNil applies the NotNil predicate on the "email_data" field.
func EmailDataNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldEmailData)))
	})
}

// EmailDataEqualFold applies the EqualFold predicate on the "email_data" field.
func EmailDataEqualFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldEmailData), v))
	})
}

// EmailDataContainsFold applies the ContainsFold predicate on the "email_data" field.
func EmailDataContainsFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldEmailData), v))
	})
}

// LatitudeEQ applies the EQ predicate on the "latitude" field.
func LatitudeEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLatitude), v))
	})
}

// LatitudeNEQ applies the NEQ predicate on the "latitude" field.
func LatitudeNEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLatitude), v))
	})
}

// LatitudeIn applies the In predicate on the "latitude" field.
func LatitudeIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLatitude), v...))
	})
}

// LatitudeNotIn applies the NotIn predicate on the "latitude" field.
func LatitudeNotIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLatitude), v...))
	})
}

// LatitudeGT applies the GT predicate on the "latitude" field.
func LatitudeGT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLatitude), v))
	})
}

// LatitudeGTE applies the GTE predicate on the "latitude" field.
func LatitudeGTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLatitude), v))
	})
}

// LatitudeLT applies the LT predicate on the "latitude" field.
func LatitudeLT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLatitude), v))
	})
}

// LatitudeLTE applies the LTE predicate on the "latitude" field.
func LatitudeLTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLatitude), v))
	})
}

// LatitudeIsNil applies the IsNil predicate on the "latitude" field.
func LatitudeIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLatitude)))
	})
}

// LatitudeNotNil applies the NotNil predicate on the "latitude" field.
func LatitudeNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLatitude)))
	})
}

// LongitudeEQ applies the EQ predicate on the "longitude" field.
func LongitudeEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLongitude), v))
	})
}

// LongitudeNEQ applies the NEQ predicate on the "longitude" field.
func LongitudeNEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLongitude), v))
	})
}

// LongitudeIn applies the In predicate on the "longitude" field.
func LongitudeIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLongitude), v...))
	})
}

// LongitudeNotIn applies the NotIn predicate on the "longitude" field.
func LongitudeNotIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLongitude), v...))
	})
}

// LongitudeGT applies the GT predicate on the "longitude" field.
func LongitudeGT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLongitude), v))
	})
}

// LongitudeGTE applies the GTE predicate on the "longitude" field.
func LongitudeGTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLongitude), v))
	})
}

// LongitudeLT applies the LT predicate on the "longitude" field.
func LongitudeLT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLongitude), v))
	})
}

// LongitudeLTE applies the LTE predicate on the "longitude" field.
func LongitudeLTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLongitude), v))
	})
}

// LongitudeIsNil applies the IsNil predicate on the "longitude" field.
func LongitudeIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLongitude)))
	})
}

// LongitudeNotNil applies the NotNil predicate on the "longitude" field.
func LongitudeNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLongitude)))
	})
}

// LocationAccuracyEQ applies the EQ predicate on the "location_accuracy" field.
func LocationAccuracyEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLocationAccuracy), v))
	})
}

// LocationAccuracyNEQ applies the NEQ predicate on the "location_accuracy" field.
func LocationAccuracyNEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLocationAccuracy), v))
	})
}

// LocationAccuracyIn applies the In predicate on the "location_accuracy" field.
func LocationAccuracyIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLocationAccuracy), v...))
	})
}

// LocationAccuracyNotIn applies the NotIn predicate on the "location_accuracy" field.
func LocationAccuracyNotIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLocationAccuracy), v...))
	})
}

// LocationAccuracyGT applies the GT predicate on the "location_accuracy" field.
func LocationAccuracyGT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLocationAccuracy), v))
	})
}

// LocationAccuracyGTE applies the GTE predicate on the "location_accuracy" field.
func LocationAccuracyGTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLocationAccuracy), v))
	})
}

// LocationAccuracyLT applies the LT predicate on the "location_accuracy" field.
func LocationAccuracyLT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLocationAccuracy), v))
	})
}

// LocationAccuracyLTE applies the LTE predicate on the "location_accuracy" field.
func LocationAccuracyLTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLocationAccuracy), v))
	})
}

// LocationAccuracyIsNil applies the IsNil predicate on the "location_accuracy" field.
func LocationAccuracyIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldLocationAccuracy)))
	})
}

// LocationAccuracyNotNil applies the NotNil predicate on the "location_accuracy" field.
func LocationAccuracyNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldLocationAccuracy)))
	})
}

// AltitudeEQ applies the EQ predicate on the "altitude" field.
func AltitudeEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAltitude), v))
	})
}

// AltitudeNEQ applies the NEQ predicate on the "altitude" field.
func AltitudeNEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldAltitude), v))
	})
}

// AltitudeIn applies the In predicate on the "altitude" field.
func AltitudeIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldAltitude), v...))
	})
}

// AltitudeNotIn applies the NotIn predicate on the "altitude" field.
func AltitudeNotIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldAltitude), v...))
	})
}

// AltitudeGT applies the GT predicate on the "altitude" field.
func AltitudeGT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldAltitude), v))
	})
}

// AltitudeGTE applies the GTE predicate on the "altitude" field.
func AltitudeGTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldAltitude), v))
	})
}

// AltitudeLT applies the LT predicate on the "altitude" field.
func AltitudeLT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldAltitude), v))
	})
}

// AltitudeLTE applies the LTE predicate on the "altitude" field.
func AltitudeLTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldAltitude), v))
	})
}

// AltitudeIsNil applies the IsNil predicate on the "altitude" field.
func AltitudeIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldAltitude)))
	})
}

// AltitudeNotNil applies the NotNil predicate on the "altitude" field.
func AltitudeNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldAltitude)))
	})
}

// PhoneDataEQ applies the EQ predicate on the "phone_data" field.
func PhoneDataEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPhoneData), v))
	})
}

// PhoneDataNEQ applies the NEQ predicate on the "phone_data" field.
func PhoneDataNEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPhoneData), v))
	})
}

// PhoneDataIn applies the In predicate on the "phone_data" field.
func PhoneDataIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldPhoneData), v...))
	})
}

// PhoneDataNotIn applies the NotIn predicate on the "phone_data" field.
func PhoneDataNotIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldPhoneData), v...))
	})
}

// PhoneDataGT applies the GT predicate on the "phone_data" field.
func PhoneDataGT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPhoneData), v))
	})
}

// PhoneDataGTE applies the GTE predicate on the "phone_data" field.
func PhoneDataGTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPhoneData), v))
	})
}

// PhoneDataLT applies the LT predicate on the "phone_data" field.
func PhoneDataLT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPhoneData), v))
	})
}

// PhoneDataLTE applies the LTE predicate on the "phone_data" field.
func PhoneDataLTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPhoneData), v))
	})
}

// PhoneDataContains applies the Contains predicate on the "phone_data" field.
func PhoneDataContains(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPhoneData), v))
	})
}

// PhoneDataHasPrefix applies the HasPrefix predicate on the "phone_data" field.
func PhoneDataHasPrefix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPhoneData), v))
	})
}

// PhoneDataHasSuffix applies the HasSuffix predicate on the "phone_data" field.
func PhoneDataHasSuffix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPhoneData), v))
	})
}

// PhoneDataIsNil applies the IsNil predicate on the "phone_data" field.
func PhoneDataIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldPhoneData)))
	})
}

// PhoneDataNotNil applies the NotNil predicate on the "phone_data" field.
func PhoneDataNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldPhoneData)))
	})
}

// PhoneDataEqualFold applies the EqualFold predicate on the "phone_data" field.
func PhoneDataEqualFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPhoneData), v))
	})
}

// PhoneDataContainsFold applies the ContainsFold predicate on the "phone_data" field.
func PhoneDataContainsFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPhoneData), v))
	})
}

// TextDataEQ applies the EQ predicate on the "text_data" field.
func TextDataEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTextData), v))
	})
}

// TextDataNEQ applies the NEQ predicate on the "text_data" field.
func TextDataNEQ(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTextData), v))
	})
}

// TextDataIn applies the In predicate on the "text_data" field.
func TextDataIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTextData), v...))
	})
}

// TextDataNotIn applies the NotIn predicate on the "text_data" field.
func TextDataNotIn(vs ...string) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTextData), v...))
	})
}

// TextDataGT applies the GT predicate on the "text_data" field.
func TextDataGT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTextData), v))
	})
}

// TextDataGTE applies the GTE predicate on the "text_data" field.
func TextDataGTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTextData), v))
	})
}

// TextDataLT applies the LT predicate on the "text_data" field.
func TextDataLT(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTextData), v))
	})
}

// TextDataLTE applies the LTE predicate on the "text_data" field.
func TextDataLTE(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTextData), v))
	})
}

// TextDataContains applies the Contains predicate on the "text_data" field.
func TextDataContains(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTextData), v))
	})
}

// TextDataHasPrefix applies the HasPrefix predicate on the "text_data" field.
func TextDataHasPrefix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTextData), v))
	})
}

// TextDataHasSuffix applies the HasSuffix predicate on the "text_data" field.
func TextDataHasSuffix(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTextData), v))
	})
}

// TextDataIsNil applies the IsNil predicate on the "text_data" field.
func TextDataIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldTextData)))
	})
}

// TextDataNotNil applies the NotNil predicate on the "text_data" field.
func TextDataNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldTextData)))
	})
}

// TextDataEqualFold applies the EqualFold predicate on the "text_data" field.
func TextDataEqualFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTextData), v))
	})
}

// TextDataContainsFold applies the ContainsFold predicate on the "text_data" field.
func TextDataContainsFold(v string) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTextData), v))
	})
}

// FloatDataEQ applies the EQ predicate on the "float_data" field.
func FloatDataEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFloatData), v))
	})
}

// FloatDataNEQ applies the NEQ predicate on the "float_data" field.
func FloatDataNEQ(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFloatData), v))
	})
}

// FloatDataIn applies the In predicate on the "float_data" field.
func FloatDataIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldFloatData), v...))
	})
}

// FloatDataNotIn applies the NotIn predicate on the "float_data" field.
func FloatDataNotIn(vs ...float64) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldFloatData), v...))
	})
}

// FloatDataGT applies the GT predicate on the "float_data" field.
func FloatDataGT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldFloatData), v))
	})
}

// FloatDataGTE applies the GTE predicate on the "float_data" field.
func FloatDataGTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldFloatData), v))
	})
}

// FloatDataLT applies the LT predicate on the "float_data" field.
func FloatDataLT(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldFloatData), v))
	})
}

// FloatDataLTE applies the LTE predicate on the "float_data" field.
func FloatDataLTE(v float64) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldFloatData), v))
	})
}

// FloatDataIsNil applies the IsNil predicate on the "float_data" field.
func FloatDataIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldFloatData)))
	})
}

// FloatDataNotNil applies the NotNil predicate on the "float_data" field.
func FloatDataNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldFloatData)))
	})
}

// IntDataEQ applies the EQ predicate on the "int_data" field.
func IntDataEQ(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldIntData), v))
	})
}

// IntDataNEQ applies the NEQ predicate on the "int_data" field.
func IntDataNEQ(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldIntData), v))
	})
}

// IntDataIn applies the In predicate on the "int_data" field.
func IntDataIn(vs ...int) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldIntData), v...))
	})
}

// IntDataNotIn applies the NotIn predicate on the "int_data" field.
func IntDataNotIn(vs ...int) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldIntData), v...))
	})
}

// IntDataGT applies the GT predicate on the "int_data" field.
func IntDataGT(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldIntData), v))
	})
}

// IntDataGTE applies the GTE predicate on the "int_data" field.
func IntDataGTE(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldIntData), v))
	})
}

// IntDataLT applies the LT predicate on the "int_data" field.
func IntDataLT(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldIntData), v))
	})
}

// IntDataLTE applies the LTE predicate on the "int_data" field.
func IntDataLTE(v int) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldIntData), v))
	})
}

// IntDataIsNil applies the IsNil predicate on the "int_data" field.
func IntDataIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldIntData)))
	})
}

// IntDataNotNil applies the NotNil predicate on the "int_data" field.
func IntDataNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldIntData)))
	})
}

// DateDataEQ applies the EQ predicate on the "date_data" field.
func DateDataEQ(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldDateData), v))
	})
}

// DateDataNEQ applies the NEQ predicate on the "date_data" field.
func DateDataNEQ(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldDateData), v))
	})
}

// DateDataIn applies the In predicate on the "date_data" field.
func DateDataIn(vs ...time.Time) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldDateData), v...))
	})
}

// DateDataNotIn applies the NotIn predicate on the "date_data" field.
func DateDataNotIn(vs ...time.Time) predicate.SurveyQuestion {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(vs) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldDateData), v...))
	})
}

// DateDataGT applies the GT predicate on the "date_data" field.
func DateDataGT(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldDateData), v))
	})
}

// DateDataGTE applies the GTE predicate on the "date_data" field.
func DateDataGTE(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldDateData), v))
	})
}

// DateDataLT applies the LT predicate on the "date_data" field.
func DateDataLT(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldDateData), v))
	})
}

// DateDataLTE applies the LTE predicate on the "date_data" field.
func DateDataLTE(v time.Time) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldDateData), v))
	})
}

// DateDataIsNil applies the IsNil predicate on the "date_data" field.
func DateDataIsNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldDateData)))
	})
}

// DateDataNotNil applies the NotNil predicate on the "date_data" field.
func DateDataNotNil() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldDateData)))
	})
}

// HasSurvey applies the HasEdge predicate on the "survey" edge.
func HasSurvey() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SurveyTable, SurveyColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSurveyWith applies the HasEdge predicate on the "survey" edge with a given conditions (other predicates).
func HasSurveyWith(preds ...predicate.Survey) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SurveyInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, SurveyTable, SurveyColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasWifiScan applies the HasEdge predicate on the "wifi_scan" edge.
func HasWifiScan() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WifiScanTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, WifiScanTable, WifiScanColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasWifiScanWith applies the HasEdge predicate on the "wifi_scan" edge with a given conditions (other predicates).
func HasWifiScanWith(preds ...predicate.SurveyWiFiScan) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(WifiScanInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, WifiScanTable, WifiScanColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCellScan applies the HasEdge predicate on the "cell_scan" edge.
func HasCellScan() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CellScanTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, CellScanTable, CellScanColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCellScanWith applies the HasEdge predicate on the "cell_scan" edge with a given conditions (other predicates).
func HasCellScanWith(preds ...predicate.SurveyCellScan) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CellScanInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, CellScanTable, CellScanColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPhotoData applies the HasEdge predicate on the "photo_data" edge.
func HasPhotoData() predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PhotoDataTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PhotoDataTable, PhotoDataColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPhotoDataWith applies the HasEdge predicate on the "photo_data" edge with a given conditions (other predicates).
func HasPhotoDataWith(preds ...predicate.File) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PhotoDataInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PhotoDataTable, PhotoDataColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.SurveyQuestion) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.SurveyQuestion) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.SurveyQuestion) predicate.SurveyQuestion {
	return predicate.SurveyQuestion(func(s *sql.Selector) {
		p(s.Not())
	})
}
