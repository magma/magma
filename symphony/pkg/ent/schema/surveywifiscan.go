// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/pkg/authz"
)

// SurveyWiFiScan holds the schema definition for the SurveyWifiScan entity.
type SurveyWiFiScan struct {
	schema
}

// Fields of the SurveyWifiScan.
func (SurveyWiFiScan) Fields() []ent.Field {
	return []ent.Field{
		field.String("ssid").
			Comment("The SSID network name of the Wi-Fi AP").
			Optional(),
		field.String("bssid").
			Comment("The broadcast SSID of the AP"),

		field.Time("timestamp").
			Comment("The time at which the Wi-Fi network was scanned"),

		field.Int("frequency").Comment("Frequency of the Wi-Fi channel in Mhz"),
		field.Int("channel").Comment("Channel of the Wi-FI AP"),
		field.String("band").
			Comment("Frequency band of the Wi-FI AP").
			Optional(),
		field.Int("channel_width").
			Comment("Width of the channel in MHz").
			Optional(),
		field.String("capabilities").
			Comment("Additional reported capabilities of the AP").
			Optional(),
		field.Int("strength").
			Comment("The signal strength normalized between 0 and 5"),

		field.Float("latitude").
			Comment("Latitude of where wifi scan data was collected").
			Optional(),
		field.Float("longitude").
			Comment("Longitude of where wifi scan data was collected").
			Optional(),
	}
}

// Edges of the SurveyWiFiScan.
func (SurveyWiFiScan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("checklist_item", CheckListItem.Type).
			Unique(),
		edge.To("survey_question", SurveyQuestion.Type).
			Unique(),
		edge.To("location", Location.Type).
			Unique(),
	}
}

// Policy returns SurveyWiFiScan policy.
func (SurveyWiFiScan) Policy() ent.Policy {
	return authz.NewPolicy(
		authz.WithQueryRules(
			authz.SurveyWiFiScanReadPolicyRule(),
		),
		authz.WithMutationRules(
			authz.SurveyWiFiScanWritePolicyRule(),
			authz.SurveyWiFiScanCreatePolicyRule(),
		),
	)
}
