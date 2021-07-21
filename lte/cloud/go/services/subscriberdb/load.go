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

package subscriberdb

import (
	"sort"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// LoadSubProtosPage streams one page (of configurable size) of subscriber data and
// converts them into proto messages.
func LoadSubProtosPage(
	pageSize uint32, pageToken string, networkID string,
	apnsByName map[string]*lte_models.ApnConfiguration,
	apnResourcesByAPN lte_models.ApnResources,
) ([]*lte_protos.SubscriberData, string, error) {
	lc := configurator.EntityLoadCriteria{
		PageSize:           pageSize,
		PageToken:          pageToken,
		LoadConfig:         true,
		LoadAssocsToThis:   true,
		LoadAssocsFromThis: true,
	}

	subEnts, nextToken, err := configurator.LoadAllEntitiesOfType(networkID, lte.SubscriberEntityType, lc, serdes.Entity)
	if err != nil {
		return nil, "", errors.Wrapf(err, "load subscribers in network of gateway %s", networkID)
	}

	subProtos := make([]*lte_protos.SubscriberData, 0, len(subEnts))
	for _, sub := range subEnts {
		subProto, err := ConvertSubEntsToProtos(sub, apnsByName, apnResourcesByAPN)
		if err != nil {
			return nil, "", err
		}
		subProto.NetworkId = &protos.NetworkID{Id: networkID}
		subProtos = append(subProtos, subProto)
	}

	return subProtos, nextToken, nil
}

func LoadSubProtosByID(
	sids []string, networkID string,
	apnsByName map[string]*lte_models.ApnConfiguration,
	apnResourcesByAPN lte_models.ApnResources,
) ([]*lte_protos.SubscriberData, error) {
	lc := configurator.EntityLoadCriteria{
		LoadConfig:         true,
		LoadAssocsToThis:   true,
		LoadAssocsFromThis: true,
	}

	subEnts, _, err := configurator.LoadEntities(networkID,
		swag.String(lte.SubscriberEntityType), nil, nil,
		storage.MakeTKs(lte.SubscriberEntityType, sids),
		lc, serdes.Entity,
	)
	if err != nil {
		return nil, errors.Wrap(err, "load added/modified subscriber entities")
	}

	subProtos := []*lte_protos.SubscriberData{}
	for _, subEnt := range subEnts {
		subProto, err := ConvertSubEntsToProtos(subEnt, apnsByName, apnResourcesByAPN)
		if err != nil {
			return nil, errors.Wrap(err, "convert subscriber entity into proto object")
		}
		subProto.NetworkId = &protos.NetworkID{Id: networkID}
		subProtos = append(subProtos, subProto)
	}
	return subProtos, nil
}

func LoadApnsByName(networkID string) (map[string]*lte_models.ApnConfiguration, error) {
	apns, _, err := configurator.LoadAllEntitiesOfType(
		networkID, lte.APNEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		return nil, err
	}
	apnsByName := map[string]*lte_models.ApnConfiguration{}
	for _, ent := range apns {
		apn, ok := ent.Config.(*lte_models.ApnConfiguration)
		if !ok {
			glog.Errorf("Attempt to convert entity %+v of type %T into apn config failed.", ent.Key, ent)
			continue
		}
		apnsByName[ent.Key] = apn
	}
	return apnsByName, err
}

func ConvertSubEntsToProtos(ent configurator.NetworkEntity, apnConfigs map[string]*lte_models.ApnConfiguration, apnResources lte_models.ApnResources) (*lte_protos.SubscriberData, error) {
	subData := &lte_protos.SubscriberData{}
	t, err := lte_protos.SidProto(ent.Key)
	if err != nil {
		return nil, err
	}

	subData.Sid = t
	if ent.Config == nil {
		return subData, nil
	}

	cfg := ent.Config.(*models.SubscriberConfig)
	subData.Lte = &lte_protos.LTESubscription{
		State:    lte_protos.LTESubscription_LTESubscriptionState(lte_protos.LTESubscription_LTESubscriptionState_value[cfg.Lte.State]),
		AuthAlgo: lte_protos.LTESubscription_LTEAuthAlgo(lte_protos.LTESubscription_LTEAuthAlgo_value[cfg.Lte.AuthAlgo]),
		AuthKey:  cfg.Lte.AuthKey,
		AuthOpc:  cfg.Lte.AuthOpc,
	}

	if cfg.Lte.SubProfile != "" {
		subData.SubProfile = string(cfg.Lte.SubProfile)
	} else {
		subData.SubProfile = defaultSubProfile
	}

	for _, assoc := range ent.ParentAssociations {
		if assoc.Type == lte.BaseNameEntityType {
			subData.Lte.AssignedBaseNames = append(subData.Lte.AssignedBaseNames, assoc.Key)
		} else if assoc.Type == lte.PolicyRuleEntityType {
			subData.Lte.AssignedPolicies = append(subData.Lte.AssignedPolicies, assoc.Key)
		}
	}

	// Construct the non-3gpp profile
	non3gpp := &lte_protos.Non3GPPUserProfile{}
	apns := ent.Associations.Filter(lte.APNEntityType)
	for _, assoc := range apns {
		apnConfig, apnFound := apnConfigs[assoc.Key]
		if !apnFound {
			continue
		}
		var apnResource *lte_protos.APNConfiguration_APNResource
		if apnResourceModel, ok := apnResources[assoc.Key]; ok {
			apnResource = apnResourceModel.ToProto()
		}
		apnProto := &lte_protos.APNConfiguration{
			ServiceSelection: assoc.Key,
			Ambr: &lte_protos.AggregatedMaximumBitrate{
				MaxBandwidthUl: *(apnConfig.Ambr.MaxBandwidthUl),
				MaxBandwidthDl: *(apnConfig.Ambr.MaxBandwidthDl),
			},
			QosProfile: &lte_protos.APNConfiguration_QoSProfile{
				ClassId:                 swag.Int32Value(apnConfig.QosProfile.ClassID),
				PriorityLevel:           swag.Uint32Value(apnConfig.QosProfile.PriorityLevel),
				PreemptionCapability:    swag.BoolValue(apnConfig.QosProfile.PreemptionCapability),
				PreemptionVulnerability: swag.BoolValue(apnConfig.QosProfile.PreemptionVulnerability),
			},
			Resource: apnResource,
		}
		if staticIP, found := cfg.StaticIps[assoc.Key]; found {
			apnProto.AssignedStaticIp = string(staticIP)
		}
		non3gpp.ApnConfig = append(non3gpp.ApnConfig, apnProto)
	}
	sort.Slice(non3gpp.ApnConfig, func(i, j int) bool {
		return non3gpp.ApnConfig[i].ServiceSelection < non3gpp.ApnConfig[j].ServiceSelection
	})
	subData.Non_3Gpp = non3gpp

	return subData, nil
}
