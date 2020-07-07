/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"time"

	"magma/lte/cloud/go/lte"
	policymodels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/thoas/go-funk"
)

func (m *Subscriber) FromBackendModels(ent configurator.NetworkEntity, statesByType map[string]state_types.State) *Subscriber {
	m.ID = policymodels.SubscriberID(ent.Key)
	m.Lte = ent.Config.(*LteSubscription)
	// If no profile in backend, return "default"
	if m.Lte.SubProfile == "" {
		m.Lte.SubProfile = "default"
	}
	for _, tk := range ent.Associations {
		if tk.Type == lte.ApnEntityType {
			m.ActiveApns = append(m.ActiveApns, tk.Key)
		}
	}

	if !funk.IsEmpty(statesByType) {
		m.Monitoring = &SubscriberStatus{}
		m.State = &SubscriberState{}
	}

	for stateType, stateVal := range statesByType {
		switch stateType {
		case lte.ICMPStateType:
			reportedState := stateVal.ReportedState.(*IcmpStatus)
			// reported time is unix timestamp in seconds, so divide ms by 1k
			reportedState.LastReportedTime = int64(stateVal.TimeMs / uint64(time.Second/time.Millisecond))
			m.Monitoring.Icmp = reportedState
		case lte.SPGWStateType:
			reportedState := stateVal.ReportedState.(*state.ArbitaryJSON)
			m.State.Spgw = reportedState
		case lte.MMEStateType:
			reportedState := stateVal.ReportedState.(*state.ArbitaryJSON)
			m.State.Mme = reportedState
		case lte.S1APStateType:
			reportedState := stateVal.ReportedState.(*state.ArbitaryJSON)
			m.State.S1ap = reportedState
		default:
			glog.Errorf("Loaded unrecognized subscriber state type %s", stateType)
		}
	}
	return m
}

func (m *SubProfile) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m ApnList) ToAssocs() []storage.TypeAndKey {
	return funk.Map(
		m,
		func(rn string) storage.TypeAndKey {
			return storage.TypeAndKey{Type: lte.ApnEntityType, Key: rn}
		},
	).([]storage.TypeAndKey)
}
