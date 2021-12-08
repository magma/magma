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
	"context"
	"fmt"
	"sort"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
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

	subEnts, nextToken, err := configurator.LoadAllEntitiesOfType(context.Background(), networkID, lte.SubscriberEntityType, lc, serdes.Entity)
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

func SerializeSubscribers(subProtos []*lte_protos.SubscriberData) (map[string][]byte, error) {
	subsSerialized := map[string][]byte{}
	for _, subProto := range subProtos {
		sid := lte_protos.SidString(subProto.Sid)
		serialized, err := proto.Marshal(subProto)
		if err != nil {
			return nil, errors.Wrap(err, "serialize subscriber proto")
		}
		subsSerialized[sid] = serialized
	}
	return subsSerialized, nil
}

// DeserializeSubscribers deserializes the given list of serialized representations of subscribers.
func DeserializeSubscribers(subProtosSerialized [][]byte) ([]*lte_protos.SubscriberData, error) {
	subs := []*lte_protos.SubscriberData{}
	for _, serialized := range subProtosSerialized {
		subProto := &lte_protos.SubscriberData{}
		err := proto.Unmarshal(serialized, subProto)
		if err != nil {
			return nil, errors.Wrap(err, "deserialize subscriber proto")
		}
		subs = append(subs, subProto)
	}
	return subs, nil
}

func LoadSubProtosByID(
	ctx context.Context,
	sids []string, networkID string,
	apnsByName map[string]*lte_models.ApnConfiguration,
	apnResourcesByAPN lte_models.ApnResources,
) ([]*lte_protos.SubscriberData, error) {
	lc := configurator.EntityLoadCriteria{
		LoadConfig:         true,
		LoadAssocsToThis:   true,
		LoadAssocsFromThis: true,
	}

	subEnts, _, err := configurator.LoadEntities(
		ctx,
		networkID,
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
		context.Background(),
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

	const coreNwTypePrefix = "NT_"
	subNetwork := &lte_protos.CoreNetworkType{}
	subNetwork.ForbiddenNetworkTypes = make([]lte_protos.CoreNetworkType_CoreNetworkTypes, len(cfg.ForbiddenNetworkTypes))
	for i, nwType := range cfg.ForbiddenNetworkTypes {
		subNetwork.ForbiddenNetworkTypes[i] = lte_protos.CoreNetworkType_CoreNetworkTypes(lte_protos.CoreNetworkType_CoreNetworkTypes_value[fmt.Sprintf("%v%v", coreNwTypePrefix, nwType)])
	}
	subData.SubNetwork = subNetwork

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
			Pdn:      lte_protos.APNConfiguration_PDNType(apnConfig.PdnType),
			Resource: apnResource,
		}
		if staticIP, found := cfg.StaticIps[assoc.Key]; found {
			apnProto.AssignedStaticIp = staticIP
		}
		non3gpp.ApnConfig = append(non3gpp.ApnConfig, apnProto)
	}
	sort.Slice(non3gpp.ApnConfig, func(i, j int) bool {
		return non3gpp.ApnConfig[i].ServiceSelection < non3gpp.ApnConfig[j].ServiceSelection
	})
	subData.Non_3Gpp = non3gpp

	return subData, nil
}

func LoadSuciProtos(ctx context.Context, networkID string) ([]*lte_protos.SuciProfile, error) {
	network, err := configurator.LoadNetwork(ctx, networkID, true, true, serdes.Network)
	if err != nil {
		return nil, errors.Wrapf(err, "network loading failed")
	}

	ngcModel := &lte_models.NetworkNgcConfigs{}
	ngcConfig := ngcModel.GetFromNetwork(network)
	if ngcConfig == nil {
		return nil, errors.Wrapf(err, "ngcConfig is nil")
	}

	suciProfiles := ngcConfig.(*lte_models.NetworkNgcConfigs).SuciProfiles
	suciProtos := []*lte_protos.SuciProfile{}
	for _, suciProfile := range suciProfiles {
		suciProtos = append(suciProtos, ngcModel.ConvertSuciEntsToProtos(suciProfile))
	}

	return suciProtos, nil
}
