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
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/protos"
	"magma/orc8r/cloud/go/mproto"
	mproto_protos "magma/orc8r/cloud/go/mproto/protos"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

const defaultSubProfile = "default"

// GetDigest returns a deterministic digest of the current configurations of a
// network, which is a concatenation of its subscribers digest and apn resources digest.
func GetDigest(networkID string) (string, error) {
	// HACK: Workaround to decouple apn resources data from subscriber data
	// despite current construction logic.
	//
	// The subscribers and apn resources digests are separately generated so
	// that the digests are per-network, instead of having to account for
	// gateway-specific apn resources configs for subscribers. We can do this
	// because APN resource IDs are unique within each network.
	subscribersDigest, err := getSubscribersDigest(networkID)
	if err != nil {
		return "", err
	}

	apnResourcesDigest, err := getApnResourcesDigest(networkID)
	if err != nil {
		return "", err
	}

	digest := subscribersDigest + apnResourcesDigest
	return digest, nil
}

// getSubscribersDigest returns a deterministic digest of all subscribers in the network.
func getSubscribersDigest(networkID string) (string, error) {
	apnsByName, err := LoadApnsByName(networkID)
	if err != nil {
		return "", err
	}

	digestsByPage := map[string][]byte{}
	token := ""
	curPage := 0
	foundEmptyToken := false

	for !foundEmptyToken {
		subProtosById := map[string]proto.Message{}
		subProtos, nextToken, err := LoadSubProtosPage(0, token, networkID, apnsByName, lte_models.ApnResources{})
		if err != nil {
			return "", err
		}
		for _, subProto := range subProtos {
			subProtosById[subProto.Sid.Id] = subProto
		}

		// Take a digest per page to be combined in the end, to avoid saving all
		// subscriber data in the memory at once
		pageDigest, err := mproto.HashManyDeterministic(subProtosById)
		if err != nil {
			return "", nil
		}
		digestsByPage[string(curPage)] = []byte(pageDigest)

		foundEmptyToken = nextToken == ""
		token = nextToken
		curPage++
	}
	digestProto := &mproto_protos.ProtosByID{BytesById: digestsByPage}
	return mproto.HashDeterministic(digestProto)
}

// getApnResourcesDigest returns a deterministic digest of the apn resources configurations
// in a network.
func getApnResourcesDigest(networkID string) (string, error) {
	apnResourceEnts, _, err := configurator.LoadAllEntitiesOfType(
		networkID, lte.APNResourceEntityType,
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	if err != nil {
		return "", err
	}

	apnResourceInternalProtosByID := map[string]proto.Message{}
	for _, apnResourceEnt := range apnResourceEnts {
		apnResource, ok := apnResourceEnt.Config.(*lte_models.ApnResource)
		if !ok {
			glog.Errorf("Attempt to convert entity %+v of type %T into apn resource failed.", apnResourceEnt.Key, apnResourceEnt)
			continue
		}
		apnResourceProto := apnResource.ToProto()

		// HACK: use the ApnResourceInternal proto to capture ingoing and outgoing
		// associations of the apn_resource
		parentGateways := apnResourceEnt.ParentAssociations.Filter(lte.CellularGatewayEntityType)
		childAPNs := apnResourceEnt.Associations.Filter(lte.APNEntityType)
		apnResourceInternalProto := &protos.ApnResourceInternal{
			AssocApns:     childAPNs.Keys(),
			AssocGateways: parentGateways.Keys(),
			ApnResource:   apnResourceProto,
		}

		apnResourceInternalProtosByID[apnResource.ID] = apnResourceInternalProto
	}

	digest, err := mproto.HashManyDeterministic(apnResourceInternalProtosByID)
	if err != nil {
		return "", err
	}
	return digest, nil
}
