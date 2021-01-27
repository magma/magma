/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"fmt"
	"sort"
	"strings"

	"magma/lte/cloud/go/lte"
	policymodels "magma/lte/cloud/go/services/policydb/obsidian/models"
	subscriberdb_state "magma/lte/cloud/go/services/subscriberdb/state"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/thoas/go-funk"
)

func (m *MutableSubscriber) ToSubscriber() *Subscriber {
	sub := &Subscriber{
		ActiveApns:          m.ActiveApns,
		ActiveBaseNames:     m.ActiveBaseNames,
		ActivePolicies:      m.ActivePolicies,
		ActivePoliciesByApn: m.ActivePoliciesByApn,
		Config:              nil, // handled below
		ID:                  m.ID,
		Lte:                 m.Lte,
		Monitoring:          nil, // augmented field
		Name:                m.Name,
		State:               nil, // augmented field
	}

	// TODO(v1.3.0+): For backwards compatibility we maintain the Lte field
	// on the sub. We can convert to just storing and exposing the Config
	// field after the next minor version.
	sub.Config = &SubscriberConfig{
		Lte:       m.Lte,
		StaticIps: m.StaticIps,
	}

	return sub
}

func (m *Subscriber) FillAugmentedFields(states state_types.StatesByID) {
	if !funk.IsEmpty(states) {
		m.Monitoring = &SubscriberStatus{}
		m.State = &SubscriberState{}
	}

	for stateID, stateVal := range states {
		switch stateID.Type {
		case orc8r.DirectoryRecordType:
			reportedState := stateVal.ReportedState.(*directoryd_types.DirectoryRecord)
			m.State.Directory = &SubscriberDirectoryRecord{LocationHistory: reportedState.LocationHistory}
		case lte.SubscriberStateType:
			m.State.SubscriberState = stateVal.ReportedState.(*state.ArbitraryJSON)
		case lte.ICMPStateType:
			reportedState := stateVal.ReportedState.(*IcmpStatus)
			reportedState.LastReportedTime = int64(stateVal.TimeMs)
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
			reportedIP, err := subscriberdb_state.GetAssignedIPAddress(*reportedState)
			if err != nil {
				glog.Errorf("failed to retrieve allocated IP for state key %s: %s", stateID.DeviceID, err)
			}
			// State ID is the IMSI with the APN appended after a dot
			apn := strings.TrimPrefix(stateID.DeviceID, fmt.Sprintf("%s.", m.ID))
			m.State.Mobility = append(m.State.Mobility, &SubscriberIPAllocation{Apn: apn, IP: strfmt.IPv4(reportedIP)})
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
}

// IsAssignedIP returns true if the subscriber has mobility state assigning it
// the passed IP address, else false.
func (m *Subscriber) IsAssignedIP(ip string) bool {
	if m.State == nil {
		return false
	}
	for _, st := range m.State.Mobility {
		if string(st.IP) == ip {
			return true
		}
	}
	return false
}

func (m *MutableSubscriber) ToTK() storage.TypeAndKey {
	return storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: string(m.ID)}
}

func (m *MutableSubscriber) FromEnt(ent configurator.NetworkEntity, policyProfileEnts configurator.NetworkEntities) (*MutableSubscriber, error) {
	model := &MutableSubscriber{
		ActivePoliciesByApn: policymodels.PolicyIdsByApn{},
		ID:                  policymodels.SubscriberID(ent.Key),
		Name:                ent.Name,
	}

	if ent.Config != nil {
		config := ent.Config.(*SubscriberConfig)
		model.Lte = config.Lte
		model.StaticIps = config.StaticIps
		// If no profile in backend, return "default"
		// TODO(8/21/20): enforce this at the API layer (and include data migration)
		if model.Lte.SubProfile == "" {
			model.Lte.SubProfile = "default"
		}
	}

	for _, tk := range ent.Associations.Filter(lte.APNEntityType) {
		model.ActiveApns = append(model.ActiveApns, tk.Key)
	}

	for _, tk := range ent.Associations.Filter(lte.PolicyRuleEntityType) {
		model.ActivePolicies = append(model.ActivePolicies, policymodels.PolicyID(tk.Key))
	}

	for _, tk := range ent.Associations.Filter(lte.BaseNameEntityType) {
		model.ActiveBaseNames = append(model.ActiveBaseNames, policymodels.BaseName(tk.Key))
	}

	// Each policy profile has 1 apn and n policy_rule
	// Convert these assocs to a map of apn->policy_rules
	for _, p := range policyProfileEnts {
		apn, err := p.Associations.GetFirst(lte.APNEntityType)
		if err != nil {
			return nil, err
		}
		model.ActivePoliciesByApn[apn.Key] = policymodels.PolicyIds{} // place apn key in map
		for _, policyRuleAssoc := range p.Associations.Filter(lte.PolicyRuleEntityType) {
			model.ActivePoliciesByApn[apn.Key] = append(model.ActivePoliciesByApn[apn.Key], policymodels.PolicyID(policyRuleAssoc.Key))
		}
	}

	return model, nil
}

func (m *MutableSubscriber) GetAssocs() []storage.TypeAndKey {
	var assocs []storage.TypeAndKey
	assocs = append(assocs, m.ActivePoliciesByApn.ToTKs(string(m.ID))...)
	assocs = append(assocs, m.ActiveApns.ToTKs()...)
	assocs = append(assocs, m.ActivePolicies.ToTKs()...)
	assocs = append(assocs, m.ActiveBaseNames.ToTKs()...)
	return assocs
}

func (m *SubProfile) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m ApnList) ToTKs() []storage.TypeAndKey {
	var tks []storage.TypeAndKey
	for _, apnName := range m {
		tks = append(tks, storage.TypeAndKey{Type: lte.APNEntityType, Key: apnName})
	}
	return tks
}
