// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package surveywifiscan

import (
	"time"

	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/symphony/graph/ent/schema"
)

const (
	// Label holds the string label denoting the surveywifiscan type in the database.
	Label = "survey_wi_fi_scan"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time vertex property in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time vertex property in the database.
	FieldUpdateTime = "update_time"
	// FieldSsid holds the string denoting the ssid vertex property in the database.
	FieldSsid = "ssid"
	// FieldBssid holds the string denoting the bssid vertex property in the database.
	FieldBssid = "bssid"
	// FieldTimestamp holds the string denoting the timestamp vertex property in the database.
	FieldTimestamp = "timestamp"
	// FieldFrequency holds the string denoting the frequency vertex property in the database.
	FieldFrequency = "frequency"
	// FieldChannel holds the string denoting the channel vertex property in the database.
	FieldChannel = "channel"
	// FieldBand holds the string denoting the band vertex property in the database.
	FieldBand = "band"
	// FieldChannelWidth holds the string denoting the channel_width vertex property in the database.
	FieldChannelWidth = "channel_width"
	// FieldCapabilities holds the string denoting the capabilities vertex property in the database.
	FieldCapabilities = "capabilities"
	// FieldStrength holds the string denoting the strength vertex property in the database.
	FieldStrength = "strength"
	// FieldLatitude holds the string denoting the latitude vertex property in the database.
	FieldLatitude = "latitude"
	// FieldLongitude holds the string denoting the longitude vertex property in the database.
	FieldLongitude = "longitude"

	// Table holds the table name of the surveywifiscan in the database.
	Table = "survey_wi_fi_scans"
	// SurveyQuestionTable is the table the holds the survey_question relation/edge.
	SurveyQuestionTable = "survey_wi_fi_scans"
	// SurveyQuestionInverseTable is the table name for the SurveyQuestion entity.
	// It exists in this package in order to avoid circular dependency with the "surveyquestion" package.
	SurveyQuestionInverseTable = "survey_questions"
	// SurveyQuestionColumn is the table column denoting the survey_question relation/edge.
	SurveyQuestionColumn = "survey_question_id"
	// LocationTable is the table the holds the location relation/edge.
	LocationTable = "survey_wi_fi_scans"
	// LocationInverseTable is the table name for the Location entity.
	// It exists in this package in order to avoid circular dependency with the "location" package.
	LocationInverseTable = "locations"
	// LocationColumn is the table column denoting the location relation/edge.
	LocationColumn = "location_id"
)

// Columns holds all SQL columns for surveywifiscan fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldSsid,
	FieldBssid,
	FieldTimestamp,
	FieldFrequency,
	FieldChannel,
	FieldBand,
	FieldChannelWidth,
	FieldCapabilities,
	FieldStrength,
	FieldLatitude,
	FieldLongitude,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the SurveyWiFiScan type.
var ForeignKeys = []string{
	"survey_question_id",
	"location_id",
}

var (
	mixin       = schema.SurveyWiFiScan{}.Mixin()
	mixinFields = [...][]ent.Field{
		mixin[0].Fields(),
	}
	fields = schema.SurveyWiFiScan{}.Fields()

	// descCreateTime is the schema descriptor for create_time field.
	descCreateTime = mixinFields[0][0].Descriptor()
	// DefaultCreateTime holds the default value on creation for the create_time field.
	DefaultCreateTime = descCreateTime.Default.(func() time.Time)

	// descUpdateTime is the schema descriptor for update_time field.
	descUpdateTime = mixinFields[0][1].Descriptor()
	// DefaultUpdateTime holds the default value on creation for the update_time field.
	DefaultUpdateTime = descUpdateTime.Default.(func() time.Time)
	// UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	UpdateDefaultUpdateTime = descUpdateTime.UpdateDefault.(func() time.Time)
)
