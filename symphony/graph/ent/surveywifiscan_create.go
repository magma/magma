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
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
)

// SurveyWiFiScanCreate is the builder for creating a SurveyWiFiScan entity.
type SurveyWiFiScanCreate struct {
	config
	create_time     *time.Time
	update_time     *time.Time
	ssid            *string
	bssid           *string
	timestamp       *time.Time
	frequency       *int
	channel         *int
	band            *string
	channel_width   *int
	capabilities    *string
	strength        *int
	latitude        *float64
	longitude       *float64
	survey_question map[string]struct{}
	location        map[string]struct{}
}

// SetCreateTime sets the create_time field.
func (swfsc *SurveyWiFiScanCreate) SetCreateTime(t time.Time) *SurveyWiFiScanCreate {
	swfsc.create_time = &t
	return swfsc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableCreateTime(t *time.Time) *SurveyWiFiScanCreate {
	if t != nil {
		swfsc.SetCreateTime(*t)
	}
	return swfsc
}

// SetUpdateTime sets the update_time field.
func (swfsc *SurveyWiFiScanCreate) SetUpdateTime(t time.Time) *SurveyWiFiScanCreate {
	swfsc.update_time = &t
	return swfsc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableUpdateTime(t *time.Time) *SurveyWiFiScanCreate {
	if t != nil {
		swfsc.SetUpdateTime(*t)
	}
	return swfsc
}

// SetSsid sets the ssid field.
func (swfsc *SurveyWiFiScanCreate) SetSsid(s string) *SurveyWiFiScanCreate {
	swfsc.ssid = &s
	return swfsc
}

// SetNillableSsid sets the ssid field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableSsid(s *string) *SurveyWiFiScanCreate {
	if s != nil {
		swfsc.SetSsid(*s)
	}
	return swfsc
}

// SetBssid sets the bssid field.
func (swfsc *SurveyWiFiScanCreate) SetBssid(s string) *SurveyWiFiScanCreate {
	swfsc.bssid = &s
	return swfsc
}

// SetTimestamp sets the timestamp field.
func (swfsc *SurveyWiFiScanCreate) SetTimestamp(t time.Time) *SurveyWiFiScanCreate {
	swfsc.timestamp = &t
	return swfsc
}

// SetFrequency sets the frequency field.
func (swfsc *SurveyWiFiScanCreate) SetFrequency(i int) *SurveyWiFiScanCreate {
	swfsc.frequency = &i
	return swfsc
}

// SetChannel sets the channel field.
func (swfsc *SurveyWiFiScanCreate) SetChannel(i int) *SurveyWiFiScanCreate {
	swfsc.channel = &i
	return swfsc
}

// SetBand sets the band field.
func (swfsc *SurveyWiFiScanCreate) SetBand(s string) *SurveyWiFiScanCreate {
	swfsc.band = &s
	return swfsc
}

// SetNillableBand sets the band field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableBand(s *string) *SurveyWiFiScanCreate {
	if s != nil {
		swfsc.SetBand(*s)
	}
	return swfsc
}

// SetChannelWidth sets the channel_width field.
func (swfsc *SurveyWiFiScanCreate) SetChannelWidth(i int) *SurveyWiFiScanCreate {
	swfsc.channel_width = &i
	return swfsc
}

// SetNillableChannelWidth sets the channel_width field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableChannelWidth(i *int) *SurveyWiFiScanCreate {
	if i != nil {
		swfsc.SetChannelWidth(*i)
	}
	return swfsc
}

// SetCapabilities sets the capabilities field.
func (swfsc *SurveyWiFiScanCreate) SetCapabilities(s string) *SurveyWiFiScanCreate {
	swfsc.capabilities = &s
	return swfsc
}

// SetNillableCapabilities sets the capabilities field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableCapabilities(s *string) *SurveyWiFiScanCreate {
	if s != nil {
		swfsc.SetCapabilities(*s)
	}
	return swfsc
}

// SetStrength sets the strength field.
func (swfsc *SurveyWiFiScanCreate) SetStrength(i int) *SurveyWiFiScanCreate {
	swfsc.strength = &i
	return swfsc
}

// SetLatitude sets the latitude field.
func (swfsc *SurveyWiFiScanCreate) SetLatitude(f float64) *SurveyWiFiScanCreate {
	swfsc.latitude = &f
	return swfsc
}

// SetNillableLatitude sets the latitude field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableLatitude(f *float64) *SurveyWiFiScanCreate {
	if f != nil {
		swfsc.SetLatitude(*f)
	}
	return swfsc
}

// SetLongitude sets the longitude field.
func (swfsc *SurveyWiFiScanCreate) SetLongitude(f float64) *SurveyWiFiScanCreate {
	swfsc.longitude = &f
	return swfsc
}

// SetNillableLongitude sets the longitude field if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableLongitude(f *float64) *SurveyWiFiScanCreate {
	if f != nil {
		swfsc.SetLongitude(*f)
	}
	return swfsc
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (swfsc *SurveyWiFiScanCreate) SetSurveyQuestionID(id string) *SurveyWiFiScanCreate {
	if swfsc.survey_question == nil {
		swfsc.survey_question = make(map[string]struct{})
	}
	swfsc.survey_question[id] = struct{}{}
	return swfsc
}

// SetNillableSurveyQuestionID sets the survey_question edge to SurveyQuestion by id if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableSurveyQuestionID(id *string) *SurveyWiFiScanCreate {
	if id != nil {
		swfsc = swfsc.SetSurveyQuestionID(*id)
	}
	return swfsc
}

// SetSurveyQuestion sets the survey_question edge to SurveyQuestion.
func (swfsc *SurveyWiFiScanCreate) SetSurveyQuestion(s *SurveyQuestion) *SurveyWiFiScanCreate {
	return swfsc.SetSurveyQuestionID(s.ID)
}

// SetLocationID sets the location edge to Location by id.
func (swfsc *SurveyWiFiScanCreate) SetLocationID(id string) *SurveyWiFiScanCreate {
	if swfsc.location == nil {
		swfsc.location = make(map[string]struct{})
	}
	swfsc.location[id] = struct{}{}
	return swfsc
}

// SetNillableLocationID sets the location edge to Location by id if the given value is not nil.
func (swfsc *SurveyWiFiScanCreate) SetNillableLocationID(id *string) *SurveyWiFiScanCreate {
	if id != nil {
		swfsc = swfsc.SetLocationID(*id)
	}
	return swfsc
}

// SetLocation sets the location edge to Location.
func (swfsc *SurveyWiFiScanCreate) SetLocation(l *Location) *SurveyWiFiScanCreate {
	return swfsc.SetLocationID(l.ID)
}

// Save creates the SurveyWiFiScan in the database.
func (swfsc *SurveyWiFiScanCreate) Save(ctx context.Context) (*SurveyWiFiScan, error) {
	if swfsc.create_time == nil {
		v := surveywifiscan.DefaultCreateTime()
		swfsc.create_time = &v
	}
	if swfsc.update_time == nil {
		v := surveywifiscan.DefaultUpdateTime()
		swfsc.update_time = &v
	}
	if swfsc.bssid == nil {
		return nil, errors.New("ent: missing required field \"bssid\"")
	}
	if swfsc.timestamp == nil {
		return nil, errors.New("ent: missing required field \"timestamp\"")
	}
	if swfsc.frequency == nil {
		return nil, errors.New("ent: missing required field \"frequency\"")
	}
	if swfsc.channel == nil {
		return nil, errors.New("ent: missing required field \"channel\"")
	}
	if swfsc.strength == nil {
		return nil, errors.New("ent: missing required field \"strength\"")
	}
	if len(swfsc.survey_question) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"survey_question\"")
	}
	if len(swfsc.location) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"location\"")
	}
	return swfsc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (swfsc *SurveyWiFiScanCreate) SaveX(ctx context.Context) *SurveyWiFiScan {
	v, err := swfsc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (swfsc *SurveyWiFiScanCreate) sqlSave(ctx context.Context) (*SurveyWiFiScan, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(swfsc.driver.Dialect())
		swfs    = &SurveyWiFiScan{config: swfsc.config}
	)
	tx, err := swfsc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(surveywifiscan.Table).Default()
	if value := swfsc.create_time; value != nil {
		insert.Set(surveywifiscan.FieldCreateTime, *value)
		swfs.CreateTime = *value
	}
	if value := swfsc.update_time; value != nil {
		insert.Set(surveywifiscan.FieldUpdateTime, *value)
		swfs.UpdateTime = *value
	}
	if value := swfsc.ssid; value != nil {
		insert.Set(surveywifiscan.FieldSsid, *value)
		swfs.Ssid = *value
	}
	if value := swfsc.bssid; value != nil {
		insert.Set(surveywifiscan.FieldBssid, *value)
		swfs.Bssid = *value
	}
	if value := swfsc.timestamp; value != nil {
		insert.Set(surveywifiscan.FieldTimestamp, *value)
		swfs.Timestamp = *value
	}
	if value := swfsc.frequency; value != nil {
		insert.Set(surveywifiscan.FieldFrequency, *value)
		swfs.Frequency = *value
	}
	if value := swfsc.channel; value != nil {
		insert.Set(surveywifiscan.FieldChannel, *value)
		swfs.Channel = *value
	}
	if value := swfsc.band; value != nil {
		insert.Set(surveywifiscan.FieldBand, *value)
		swfs.Band = *value
	}
	if value := swfsc.channel_width; value != nil {
		insert.Set(surveywifiscan.FieldChannelWidth, *value)
		swfs.ChannelWidth = *value
	}
	if value := swfsc.capabilities; value != nil {
		insert.Set(surveywifiscan.FieldCapabilities, *value)
		swfs.Capabilities = *value
	}
	if value := swfsc.strength; value != nil {
		insert.Set(surveywifiscan.FieldStrength, *value)
		swfs.Strength = *value
	}
	if value := swfsc.latitude; value != nil {
		insert.Set(surveywifiscan.FieldLatitude, *value)
		swfs.Latitude = *value
	}
	if value := swfsc.longitude; value != nil {
		insert.Set(surveywifiscan.FieldLongitude, *value)
		swfs.Longitude = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(surveywifiscan.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	swfs.ID = strconv.FormatInt(id, 10)
	if len(swfsc.survey_question) > 0 {
		for eid := range swfsc.survey_question {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(surveywifiscan.SurveyQuestionTable).
				Set(surveywifiscan.SurveyQuestionColumn, eid).
				Where(sql.EQ(surveywifiscan.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if len(swfsc.location) > 0 {
		for eid := range swfsc.location {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			query, args := builder.Update(surveywifiscan.LocationTable).
				Set(surveywifiscan.LocationColumn, eid).
				Where(sql.EQ(surveywifiscan.FieldID, id)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return swfs, nil
}
