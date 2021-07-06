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

package servicers

import (
	"context"

	"magma/lte/cloud/go/services/subscriberdb"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"
)

type subscriberdbServicer struct {
	flatDigestEnabled bool
	digestStore       storage.DigestLookup
}

func NewSubscriberdbServicer(config subscriberdb.Config, digestStore storage.DigestLookup) lte_protos.SubscriberDBCloudServer {
	servicer := &subscriberdbServicer{
		flatDigestEnabled: config.FlatDigestEnabled,
		digestStore:       digestStore,
	}
	return servicer
}

func (s *subscriberdbServicer) CheckSubscribersInSync(
	ctx context.Context,
	req *lte_protos.CheckSubscribersInSyncRequest,
) (*lte_protos.CheckSubscribersInSyncResponse, error) {
	gateway := protos.GetClientGateway(ctx)
	if gateway == nil {
		return nil, status.Errorf(codes.PermissionDenied, "missing gateway identity")
	}
	if !gateway.Registered() {
		return nil, status.Errorf(codes.PermissionDenied, "gateway is not registered")
	}
	networkID := gateway.NetworkId
	_, inSync := s.getDigestInfo(req.FlatDigest, networkID)

	res := &lte_protos.CheckSubscribersInSyncResponse{InSync: inSync}
	return res, nil
}

func (s *subscriberdbServicer) SyncSubscribers(
	ctx context.Context,
	req *lte_protos.SyncSubscribersRequest,
) (*lte_protos.SyncSubscribersResponse, error) {
	// TODO(wangyyt1013): implement logic to return subscriber changeset.
	// The current scaffolding code is for backward-compatibility purposes.
	res := &lte_protos.SyncSubscribersResponse{Resync: true}
	return res, nil
}

// ListSubscribers returns a page of subscribers and a token to be used on
// subsequent requests. The page token specified in the request is used to
// determine the first subscriber to include in the page. The page size
// specified in the request determines the maximum number of entities to
// return. If no page size is specified, the maximum size configured in the
// configurator service will be returned.
func (s *subscriberdbServicer) ListSubscribers(ctx context.Context, req *lte_protos.ListSubscribersRequest) (*lte_protos.ListSubscribersResponse, error) {
	gateway := protos.GetClientGateway(ctx)
	if gateway == nil {
		return nil, status.Errorf(codes.PermissionDenied, "missing gateway identity")
	}
	if !gateway.Registered() {
		return nil, status.Errorf(codes.PermissionDenied, "gateway is not registered")
	}
	networkID := gateway.NetworkId
	gatewayID := gateway.LogicalId

	lteGateway, err := configurator.LoadEntity(
		networkID, lte.CellularGatewayEntityType, gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
		serdes.Entity,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "load cellular gateway for gateway %s", gatewayID)
	}
	apnsByName, apnResourcesByAPN, err := loadAPNs(lteGateway)
	if err != nil {
		return nil, err
	}
	subProtos, nextToken, err := subscriberdb.LoadSubProtosPage(req.PageSize, req.PageToken, networkID, apnsByName, apnResourcesByAPN)
	if err != nil {
		return nil, err
	}

	digest, _ := s.getDigestInfo(&lte_protos.Digest{Md5Base64Digest: ""}, networkID)
	listRes := &lte_protos.ListSubscribersResponse{
		Subscribers:   subProtos,
		NextPageToken: nextToken,
		FlatDigest:    digest,
	}
	return listRes, nil
}

// getDigestInfo returns the correctly formatted Digest and NoUpdates values
// according to the client digest.
func (s *subscriberdbServicer) getDigestInfo(clientDigest *lte_protos.Digest, networkID string) (*lte_protos.Digest, bool) {
	// The flat digest functionality is currently placed behind a feature flag
	if !s.flatDigestEnabled {
		return &lte_protos.Digest{Md5Base64Digest: ""}, false
	}

	digest, _, err := storage.GetDigest(s.digestStore, networkID)
	// If digest generation fails, the error is swallowed to not affect the main functionality
	if err != nil {
		glog.Errorf("Generating digest for network %s failed: %+v", networkID, err)
		return &lte_protos.Digest{Md5Base64Digest: ""}, false
	}

	noUpdates := digest != "" && digest == clientDigest.GetMd5Base64Digest()
	digestProto := &lte_protos.Digest{Md5Base64Digest: digest}
	return digestProto, noUpdates
}

func loadAPNs(gateway configurator.NetworkEntity) (map[string]*lte_models.ApnConfiguration, lte_models.ApnResources, error) {
	apnsByName, err := subscriberdb.LoadApnsByName(gateway.NetworkID)
	if err != nil {
		return nil, nil, err
	}
	apnResources, err := lte_models.LoadAPNResources(gateway.NetworkID, gateway.Associations.Filter(lte.APNResourceEntityType).Keys())
	if err != nil {
		return nil, nil, err
	}

	return apnsByName, apnResources, nil
}
