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
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"magma/lte/cloud/go/lte"
	policymodels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	mobilitydStateExpectedMatchCount = 2
)

var (
	// mobilitydStateKeyRe captures the IMSI portion of mobilityd state keys.
	// Mobilityd states are keyed as <IMSI>.<APN>.
	mobilitydStateKeyRe = regexp.MustCompile(`^(?P<imsi>IMSI\d+)\..+$`)

	subscriberLoadCriteria = configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsFromThis: true}
)

var subscriberStateTypes = []string{
	lte.ICMPStateType,
	lte.S1APStateType,
	lte.MMEStateType,
	lte.SPGWStateType,
	lte.MobilitydStateType,
	orc8r.DirectoryRecordType,
}

func (m *Subscriber) Load(networkID, key string) (*Subscriber, error) {
	mutableSub, err := (&MutableSubscriber{}).Load(networkID, key)
	if err != nil {
		return nil, err
	}
	sub := m.FromMutable(mutableSub)

	states, err := state.SearchStates(networkID, subscriberStateTypes, nil, &key)
	if err != nil {
		return nil, err
	}
	err = sub.FillAugmentedFields(states)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (m *Subscriber) LoadAll(networkID string) (map[string]*Subscriber, error) {
	mutableSubs, err := (&MutableSubscriber{}).LoadAll(networkID)
	if err != nil {
		return nil, err
	}

	states, err := state.SearchStates(networkID, subscriberStateTypes, nil, nil)
	if err != nil {
		return nil, err
	}
	// Each entry in this map contains all the states that the SID cares about.
	// The DeviceID fields of the state IDs in the nested maps do not have to
	// match the SID, as in the case of mobilityd state for example.
	statesBySid := map[string]state_types.StatesByID{}
	for stateID, st := range states {
		sidKey := stateID.DeviceID
		if stateID.Type == lte.MobilitydStateType {
			matches := mobilitydStateKeyRe.FindStringSubmatch(stateID.DeviceID)
			if len(matches) != mobilitydStateExpectedMatchCount {
				glog.Errorf("mobilityd state composite ID %s did not match regex", sidKey)
				continue
			}
			sidKey = matches[1]
		}

		if _, exists := statesBySid[sidKey]; !exists {
			statesBySid[sidKey] = state_types.StatesByID{}
		}
		statesBySid[sidKey][stateID] = st
	}

	subs := map[string]*Subscriber{}
	for _, mutableSub := range mutableSubs {
		sub := m.FromMutable(mutableSub)
		err = sub.FillAugmentedFields(statesBySid[string(sub.ID)])
		if err != nil {
			return nil, err
		}
		subs[string(sub.ID)] = sub
	}

	return subs, nil
}

func (m *Subscriber) FromMutable(mut *MutableSubscriber) *Subscriber {
	sub := &Subscriber{
		ActiveApns:          mut.ActiveApns,
		ActiveBaseNames:     mut.ActiveBaseNames,
		ActivePolicies:      mut.ActivePolicies,
		ActivePoliciesByApn: mut.ActivePoliciesByApn,
		Config:              nil, // handled below
		ID:                  mut.ID,
		Lte:                 mut.Lte,
		Monitoring:          nil, // augmented field
		Name:                mut.Name,
		State:               nil, // augmented field
	}

	// TODO(v1.3.0+): For backwards compatibility we maintain the Lte field
	// on the sub. We can convert to just storing and exposing the Config
	// field after the next minor version.
	sub.Config = &SubscriberConfig{
		Lte:       mut.Lte,
		StaticIps: mut.StaticIps,
	}

	return sub
}

func (m *Subscriber) FillAugmentedFields(states state_types.StatesByID) error {
	if !funk.IsEmpty(states) {
		m.Monitoring = &SubscriberStatus{}
		m.State = &SubscriberState{}
	}

	for stateID, stateVal := range states {
		switch stateID.Type {
		case orc8r.DirectoryRecordType:
			reportedState := stateVal.ReportedState.(*directoryd.DirectoryRecord)
			m.State.Directory = &SubscriberDirectoryRecord{LocationHistory: reportedState.LocationHistory}
		case lte.ICMPStateType:
			reportedState := stateVal.ReportedState.(*IcmpStatus)
			// Reported time is unix timestamp in seconds, so divide ms by 1k
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
			// State ID is the IMSI with the APN appended after a dot
			ipAPN := strings.TrimPrefix(stateID.DeviceID, fmt.Sprintf("%s.", m.ID))
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

	return nil
}

func (m *MutableSubscriber) Load(networkID, key string) (*MutableSubscriber, error) {
	ent, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, key, subscriberLoadCriteria)
	if err != nil {
		return nil, err
	}
	return m.FromEnt(ent)
}

func (m *MutableSubscriber) LoadAll(networkID string) (map[string]*MutableSubscriber, error) {
	ents, err := configurator.LoadAllEntitiesInNetwork(networkID, lte.SubscriberEntityType, subscriberLoadCriteria)
	if err != nil {
		return nil, err
	}
	subs := map[string]*MutableSubscriber{}
	for _, ent := range ents {
		sub, err := m.FromEnt(ent)
		if err != nil {
			return nil, err
		}
		subs[ent.Key] = sub
	}
	return subs, nil
}

func (m *MutableSubscriber) Create(networkID string) error {
	// New ents
	//	- active_policies_by_apn
	//		- Assocs: policy_rule..., apn
	//	- subscriber
	//		- Assocs: active_policies_by_apn

	subEnt := configurator.NetworkEntity{
		Type: lte.SubscriberEntityType,
		Key:  string(m.ID),
		Name: m.Name,
		Config: &SubscriberConfig{
			Lte:       m.Lte,
			StaticIps: m.StaticIps,
		},
		Associations: m.GetAssocs(),
	}

	var ents []configurator.NetworkEntity
	ents = append(ents, m.ActivePoliciesByApn.ToEntities(subEnt.Key)...)
	ents = append(ents, subEnt)

	_, err := configurator.CreateEntities(networkID, ents)
	if err != nil {
		return err
	}

	return nil
}

func (m *MutableSubscriber) Update(networkID string) error {
	var writes []configurator.EntityWriteOperation

	existingSub, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, string(m.ID), subscriberLoadCriteria)
	if err != nil {
		return err
	}

	// For simplicity, delete all of subscriber's existing
	// apn_policy_profile, then add new
	policyMapTKs := existingSub.Associations.Filter(lte.APNPolicyProfileEntityType)
	for _, tk := range policyMapTKs {
		writes = append(writes, configurator.EntityUpdateCriteria{Type: tk.Type, Key: tk.Key, DeleteEntity: true})
	}
	for _, e := range m.ActivePoliciesByApn.ToEntities(string(m.ID)) {
		writes = append(writes, e)
	}

	subUpdate := configurator.EntityUpdateCriteria{
		Key:     string(m.ID),
		Type:    lte.SubscriberEntityType,
		NewName: swag.String(m.Name),
		NewConfig: &SubscriberConfig{
			Lte:       m.Lte,
			StaticIps: m.StaticIps,
		},
		AssociationsToSet: m.GetAssocs(),
	}
	writes = append(writes, subUpdate)

	err = configurator.WriteEntities(networkID, writes...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MutableSubscriber) Delete(networkID, key string) error {
	ent, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, key, configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
	if err != nil {
		return err
	}
	sub, err := m.FromEnt(ent)
	if err != nil {
		return err
	}

	var deletes []storage.TypeAndKey
	deletes = append(deletes, sub.ToTK())
	deletes = append(deletes, sub.ActivePoliciesByApn.ToTKs(string(sub.ID))...)

	err = configurator.DeleteEntities(networkID, deletes)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return nil
}

func (m *MutableSubscriber) ToTK() storage.TypeAndKey {
	return storage.TypeAndKey{Type: lte.SubscriberEntityType, Key: string(m.ID)}
}

func (m *MutableSubscriber) FromEnt(ent configurator.NetworkEntity) (*MutableSubscriber, error) {
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

	policyProfileAssocs := ent.Associations.Filter(lte.APNPolicyProfileEntityType)
	if len(policyProfileAssocs) == 0 {
		return model, nil
	}

	// Need to load the policy profile ents to determine their edges.
	// Configurator doesn't currently support loading a specified subgraph,
	// so we have to load the subscriber and its policy profiles in
	// separate calls.
	policyProfileEnts, _, err := configurator.LoadEntities(
		ent.NetworkID, nil, nil, nil,
		policyProfileAssocs,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
	)
	if err != nil {
		return nil, err
	}
	// Each policy profile has 1 apn and n policy_rule
	// Convert these assocs to a map of apn->policy_rules
	for _, p := range policyProfileEnts {
		apn, err := p.Associations.GetFirst(lte.APNEntityType)
		if err != nil {
			return nil, err
		}
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
	return assocs
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

func (m ApnList) ToTKs() []storage.TypeAndKey {
	var tks []storage.TypeAndKey
	for _, apnName := range m {
		tks = append(tks, storage.TypeAndKey{Type: lte.APNEntityType, Key: apnName})
	}
	return tks
}
