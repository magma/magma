/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package streamer

import (
	"sort"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

// SubscribersProvider provides the implementation for subscriber streaming.
type SubscribersProvider struct{}

func (p *SubscribersProvider) GetStreamName() string {
	return lte.SubscriberStreamName
}

func (p *SubscribersProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	ent, err := configurator.LoadEntityForPhysicalID(gatewayId, configurator.EntityLoadCriteria{})
	if err != nil {
		return nil, err
	}
	// Collect all subscribers in one RPC call
	subEnts, err := configurator.LoadAllEntitiesInNetwork(ent.NetworkID, lte.SubscriberEntityType, configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsToThis: true, LoadAssocsFromThis: true})
	if err != nil {
		return nil, err
	}
	// Collect all APNs in one RPC call
	apnEnts, err := configurator.LoadAllEntitiesInNetwork(ent.NetworkID, lte.APNEntityType, configurator.EntityLoadCriteria{LoadConfig: true})
	// Create a map to avoid for loops in function calls to populate subscriber data from subscriber associations
	apnConfigMap := make(map[string]*lte_models.ApnConfiguration, len(apnEnts))
	for _, apnEnt := range apnEnts {
		apnConfigMap[apnEnt.Key] = apnEnt.Config.(*lte_models.ApnConfiguration)
	}

	subProtos := make([]*lte_protos.SubscriberData, 0, len(subEnts))
	for _, sub := range subEnts {
		subProto := &lte_protos.SubscriberData{}
		subProto, err = subscriberToMconfig(sub, apnConfigMap)
		if err != nil {
			return nil, err
		}
		subProto.NetworkId = &protos.NetworkID{Id: ent.NetworkID}
		subProtos = append(subProtos, subProto)
	}
	return subscribersToUpdates(subProtos)
}

func subscribersToUpdates(subs []*lte_protos.SubscriberData) ([]*protos.DataUpdate, error) {
	ret := make([]*protos.DataUpdate, 0, len(subs))
	for _, sub := range subs {
		marshaledProto, err := proto.Marshal(sub)
		if err != nil {
			return nil, err
		}
		update := &protos.DataUpdate{Key: lte_protos.SidString(sub.Sid), Value: marshaledProto}
		ret = append(ret, update)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Key < ret[j].Key })
	return ret, nil
}

func subscriberToMconfig(ent configurator.NetworkEntity, apnConfigs map[string]*lte_models.ApnConfiguration) (*lte_protos.SubscriberData, error) {
	sub := &lte_protos.SubscriberData{}
	t, err := lte_protos.SidProto(ent.Key)
	if err != nil {
		return nil, err
	}

	sub.Sid = t
	if ent.Config == nil {
		return sub, nil
	}

	cfg := ent.Config.(*models.SubscriberConfig)
	sub.Lte = &lte_protos.LTESubscription{
		State:    lte_protos.LTESubscription_LTESubscriptionState(lte_protos.LTESubscription_LTESubscriptionState_value[cfg.Lte.State]),
		AuthAlgo: lte_protos.LTESubscription_LTEAuthAlgo(lte_protos.LTESubscription_LTEAuthAlgo_value[cfg.Lte.AuthAlgo]),
		AuthKey:  cfg.Lte.AuthKey,
		AuthOpc:  cfg.Lte.AuthOpc,
	}

	if cfg.Lte.SubProfile != "" {
		sub.SubProfile = string(cfg.Lte.SubProfile)
	} else {
		sub.SubProfile = "default"
	}

	for _, assoc := range ent.ParentAssociations {
		if assoc.Type == lte.BaseNameEntityType {
			sub.Lte.AssignedBaseNames = append(sub.Lte.AssignedBaseNames, assoc.Key)
		} else if assoc.Type == lte.PolicyRuleEntityType {
			sub.Lte.AssignedPolicies = append(sub.Lte.AssignedPolicies, assoc.Key)
		}
	}

	// Construct the non-3gpp profile
	non3gpp := &lte_protos.Non3GPPUserProfile{
		ApnConfig: make([]*lte_protos.APNConfiguration, 0, len(ent.Associations)),
	}
	for _, assoc := range ent.Associations {
		apnConfig, apnFound := apnConfigs[assoc.Key]
		if !apnFound {
			continue
		}

		apnProto := &lte_protos.APNConfiguration{
			ServiceSelection: assoc.Key,
			Ambr: &lte_protos.AggregatedMaximumBitrate{
				MaxBandwidthUl: *(apnConfig.Ambr.MaxBandwidthUl),
				MaxBandwidthDl: *(apnConfig.Ambr.MaxBandwidthDl),
			},
			QosProfile: &lte_protos.APNConfiguration_QoSProfile{
				ClassId:                 *(apnConfig.QosProfile.ClassID),
				PriorityLevel:           *(apnConfig.QosProfile.PriorityLevel),
				PreemptionCapability:    *(apnConfig.QosProfile.PreemptionCapability),
				PreemptionVulnerability: *(apnConfig.QosProfile.PreemptionVulnerability),
			},
		}
		if staticIP, found := cfg.StaticIps[assoc.Key]; found {
			apnProto.AssignedStaticIp = string(staticIP)
		}
		non3gpp.ApnConfig = append(non3gpp.ApnConfig, apnProto)
	}
	sort.Slice(non3gpp.ApnConfig, func(i, j int) bool {
		return non3gpp.ApnConfig[i].ServiceSelection < non3gpp.ApnConfig[j].ServiceSelection
	})
	sub.Non_3Gpp = non3gpp

	return sub, nil
}
