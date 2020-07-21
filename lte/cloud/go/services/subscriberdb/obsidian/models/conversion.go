/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"encoding/base64"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"magma/lte/cloud/go/lte"
	policymodels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func (m *Subscriber) FromBackendModels(ent configurator.NetworkEntity, statesByID state_types.StatesByID) *Subscriber {
	m.ID = policymodels.SubscriberID(ent.Key)
	m.Name = ent.Name
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

	if !funk.IsEmpty(statesByID) {
		m.Monitoring = &SubscriberStatus{}
		m.State = &SubscriberState{}
	}

	for stateID, stateVal := range statesByID {
		switch stateID.Type {
		case lte.ICMPStateType:
			reportedState := stateVal.ReportedState.(*IcmpStatus)
			// reported time is unix timestamp in seconds, so divide ms by 1k
			reportedState.LastReportedTime = int64(stateVal.TimeMs / uint64(time.Second/time.Millisecond))
			m.Monitoring.Icmp = reportedState
		case lte.SPGWStateType:
			reportedState := stateVal.ReportedState.(*state.ArbitraryJSON)
			m.State.Spgw = reportedState
		case lte.MMEStateType:
			reportedState := stateVal.ReportedState.(*state.ArbitraryJSON)
			m.State.Mme = reportedState
		case lte.S1APStateType:
			reportedState := stateVal.ReportedState.(*state.ArbitraryJSON)
			m.State.S1ap = reportedState
		case lte.MobilitydStateType:
			reportedState := stateVal.ReportedState.(*state.ArbitraryJSON)
			if reportedState == nil {
				break
			}
			// We swallow and log errors because we don't want to block an API
			// request if some AGW is sending buggy/malformed mobilityd state
			reportedIP, err := getAssignedIPAddress(*reportedState)
			if err != nil {
				glog.Errorf("failed to retrieve allocated IP for state key %s: %s", stateID.DeviceID, err)
			}
			// The state ID is the IMSI with the APN appended after a dot
			ipAPN := strings.TrimPrefix(stateID.DeviceID, fmt.Sprintf("%s.", ent.Key))
			m.State.Mobility = append(m.State.Mobility, &SubscriberIPAllocation{Apn: ipAPN, IP: strfmt.IPv4(reportedIP)})
		default:
			glog.Errorf("Loaded unrecognized subscriber state type %s", stateID.Type)
		}
	}
	// Sort mobility state by APN for determinism
	if m.State != nil && !funk.IsEmpty(m.State.Mobility) {
		sort.Slice(m.State.Mobility, func(i, j int) bool {
			return m.State.Mobility[i].Apn < m.State.Mobility[j].Apn
		})
	}
	return m
}

// We expect something along the lines of:
// {
//   "state": "ALLOCATED",
//   "sid": {"id": "IMSI001010000000001.magma.ipv4"},
//   "ipBlock": {"netAddress": "wKiAAA==", "prefixLen": 24},
//   "ip": {"address": "wKiArg=="}
//  }
// The IP addresses are base64 encoded versions of the packed bytes
func getAssignedIPAddress(mobilitydState state.ArbitraryJSON) (string, error) {
	ipField, ipExists := mobilitydState["ip"]
	if !ipExists {
		return "", errors.New("no ip field found in mobilityd state")
	}
	ipFieldAsMap, castOK := ipField.(map[string]interface{})
	if !castOK {
		return "", errors.New("could not cast ip field of mobilityd state to arbitrary JSON map type")
	}
	ipAddress, addrExists := ipFieldAsMap["address"]
	if !addrExists {
		return "", errors.New("no IP address found in mobilityd state")
	}
	ipAddressAsString, castOK := ipAddress.(string)
	if !castOK {
		return "", errors.New("encoded IP address is not a string as expected")
	}

	return base64DecodeIPAddress(ipAddressAsString)
}

func base64DecodeIPAddress(encodedIP string) (string, error) {
	ipBytes, err := base64.StdEncoding.DecodeString(encodedIP)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode mobilityd IP address")
	}
	if len(ipBytes) != 4 {
		return "", errors.Errorf("expected IP address to decode to 4 bytes, got %d", len(ipBytes))
	}
	return net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3]).String(), nil
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
