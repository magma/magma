// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// SurveyCellScan holds the schema definition for the SurveyCellScan entity.
type SurveyCellScan struct {
	schema
}

// Fields of the SurveyCellScan.
func (SurveyCellScan) Fields() []ent.Field {
	return []ent.Field{
		field.String("network_type").
			Comment("The type of the cellular network"),
		field.Int("signal_strength").
			Comment("The strength of the cellular network in dBm"),

		field.Time("timestamp").
			Comment("Time at which cellular network was scanned").
			Optional(),

		field.String("base_station_id").
			Comment("Base Station Identity Code").
			Optional(),
		field.String("network_id").
			Comment("CDMA Network ID").
			Optional(),
		field.String("system_id").
			Comment("CDMA System ID").
			Optional(),
		field.String("cell_id").
			Comment("The Cell Identity (cid) of the tower as described in TS 27.007").
			Optional(),

		field.String("location_area_code").
			Comment("GSM 16-bit Location Area Code (lac)").
			Optional(),
		field.String("mobile_country_code").
			Comment("3-digit Mobile Country Code (mcc)").
			Optional(),
		field.String("mobile_network_code").
			Comment("2 or 3-digit Mobile Network Code (mnc)").
			Optional(),
		field.String("primary_scrambling_code").
			Comment("UMTS Primary Scrambling Code described in TS 25.331").
			Optional(),
		field.String("operator").
			Comment("Operator name of the cellular network").
			Optional(),
		field.Int("arfcn").
			Comment("GSM Absolute RF Channel Number (arfcn)").
			Optional(),

		field.String("physical_cell_id").
			Comment("LTE Physical Cell Id (pci)").
			Optional(),
		field.String("tracking_area_code").
			Comment("LTE 16-bit Tracking Area Code (tac)").
			Optional(),
		field.Int("timing_advance").
			Comment("LTE timing advance described in 3GPP 36.213 Sec 4.2.3").
			Optional(),

		field.Int("earfcn").
			Comment("LTE Absolute RF Channel Number (earfcn)").
			Optional(),
		field.Int("uarfcn").
			Comment("UMTS Absolute RF Channel Number described in TS 25.101 sec. 5.4.4 (uarfcn)s").
			Optional(),

		field.Float("latitude").
			Comment("Latitude of where cell data was collected").
			Optional(),
		field.Float("longitude").
			Comment("Longitude of where cell data was collected").
			Optional(),
	}
}

// Edges of the SurveyCellScan.
func (SurveyCellScan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("survey_question", SurveyQuestion.Type).
			Unique(),
		edge.To("location", Location.Type).
			Unique(),
	}
}
