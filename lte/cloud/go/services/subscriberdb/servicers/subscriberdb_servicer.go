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
	"time"

	"magma/lte/cloud/go/services/subscriberdb"
	"magma/orc8r/cloud/go/orc8r/math"
	"magma/orc8r/cloud/go/syncstore"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type subscriberdbServicer struct {
	subscriberdb.Config
	store syncstore.SyncStoreReader
}

func NewSubscriberdbServicer(config subscriberdb.Config, store syncstore.SyncStoreReader) lte_protos.SubscriberDBCloudServer {
	return &subscriberdbServicer{store: store, Config: config}
}

func (s *subscriberdbServicer) CheckInSync(
	ctx context.Context,
	req *lte_protos.CheckInSyncRequest,
) (*lte_protos.CheckInSyncResponse, error) {
	if !s.DigestsEnabled {
		return &lte_protos.CheckInSyncResponse{InSync: false}, nil
	}

	gateway := protos.GetClientGateway(ctx)
	if gateway == nil {
		return nil, status.Errorf(codes.PermissionDenied, "missing gateway identity")
	}
	if !gateway.Registered() {
		return nil, status.Errorf(codes.PermissionDenied, "gateway is not registered")
	}

	networkID := gateway.NetworkId
	gatewayID := gateway.LogicalId
	if s.shouldResync(networkID, gatewayID) {
		return &lte_protos.CheckInSyncResponse{InSync: false}, nil
	}

	_, inSync := s.getDigestInfo(req.RootDigest, networkID)
	res := &lte_protos.CheckInSyncResponse{InSync: inSync}
	return res, nil
}

func (s *subscriberdbServicer) Sync(
	ctx context.Context,
	req *lte_protos.SyncRequest,
) (*lte_protos.SyncResponse, error) {
	if !s.DigestsEnabled {
		return &lte_protos.SyncResponse{Resync: true}, nil
	}

	gateway := protos.GetClientGateway(ctx)
	if gateway == nil {
		return nil, status.Errorf(codes.PermissionDenied, "missing gateway identity")
	}
	if !gateway.Registered() {
		return nil, status.Errorf(codes.PermissionDenied, "gateway is not registered")
	}

	networkID := gateway.NetworkId
	gatewayID := gateway.LogicalId
	if s.shouldResync(networkID, gatewayID) {
		return &lte_protos.SyncResponse{Resync: true}, nil
	}

	digestTree, err := syncstore.GetDigestTree(s.store, networkID)
	if err != nil {
		return nil, err
	}
	resync, renewed, deleted, err := s.getSubscribersChangeset(networkID, req.LeafDigests, digestTree.LeafDigests)
	if err != nil {
		return nil, err
	}
	if resync {
		return &lte_protos.SyncResponse{Resync: true}, nil
	}

	// Since the cached protos don't contain gateway-specific information, inject
	// the apn resource configs related to the gateway
	renewed, err = injectAPNResources(renewed, gateway)
	if err != nil {
		return nil, err
	}
	renewedMarshaled := []*any.Any{}
	for _, subProto := range renewed {
		anyVal, err := ptypes.MarshalAny(subProto)
		if err != nil {
			return nil, errors.Wrapf(err, "marshal subscriber protos for network %+v", networkID)
		}
		renewedMarshaled = append(renewedMarshaled, anyVal)
	}

	res := &lte_protos.SyncResponse{
		Digests: digestTree,
		Changeset: &protos.Changeset{
			ToRenew: renewedMarshaled,
			Deleted: deleted,
		},
		Resync: false,
	}
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

	apnsByName, apnResourcesByAPN, err := loadAPNs(gateway)
	if err != nil {
		return nil, err
	}

	var subProtos []*lte_protos.SubscriberData
	var nextToken string
	if s.DigestsEnabled {
		subProtos, nextToken, err = s.loadSubscribersPageFromCache(networkID, req, gateway)
		if err != nil {
			return nil, err
		}
	} else {
		subProtos, nextToken, err = subscriberdb.LoadSubProtosPage(req.PageSize, req.PageToken, networkID, apnsByName, apnResourcesByAPN)
		if err != nil {
			return nil, err
		}
	}

	rootDigest := &protos.Digest{Md5Base64Digest: ""}
	leafDigests := []*protos.LeafDigest{}
	// The digests are sent back during the request for the first page of subscriber data
	if req.PageToken == "" && s.DigestsEnabled {
		digestTree, _ := s.getDigestInfo(&protos.Digest{Md5Base64Digest: ""}, networkID)
		rootDigest = digestTree.RootDigest
		leafDigests = digestTree.LeafDigests
	}

	// At the AGW request for the last page, update the lastResyncTime of the gateway to the current time.
	// NOTE: Since the resync is Orc8r-directed, and Orc8r doesn't track the request status on the AGW side,
	// Orc8r takes the AGW request for the last page as an approximate indication of the completion of a resync.
	if nextToken == "" {
		err := s.store.RecordResync(networkID, gatewayID, time.Now().Unix())
		if err != nil {
			glog.Errorf("Failed to set last resync time for gateway %+v of network %+v: %+v", gatewayID, networkID, err)
		}
	}

	listRes := &lte_protos.ListSubscribersResponse{
		Subscribers:   subProtos,
		NextPageToken: nextToken,
		Digests: &protos.DigestTree{
			RootDigest:  rootDigest,
			LeafDigests: leafDigests,
		},
	}
	return listRes, nil
}

// getSubscribersChangeset compares the cloud and AGW digests and returns
// 1. Whether a resync is required for this AGW.
// 2. If no resync, the list of subscriber configs to be renewed.
// 3. If no resync, the list of subscriber IDs to be deleted.
// 4. Any error that occurred.
func (s *subscriberdbServicer) getSubscribersChangeset(networkID string, clientDigests []*protos.LeafDigest, cloudDigests []*protos.LeafDigest) (bool, []*lte_protos.SubscriberData, []string, error) {
	toRenew, deleted := syncstore.GetLeafDigestsDiff(clientDigests, cloudDigests)
	if len(toRenew) > s.ChangesetSizeThreshold || len(toRenew) > int(s.MaxProtosLoadSize) {
		return true, nil, nil, nil
	}

	sids := funk.Keys(toRenew).([]string)
	renewedSerialized, err := s.store.GetCachedByID(networkID, sids)
	if err != nil {
		return true, nil, nil, err
	}
	renewed, err := subscriberdb.DeserializeSubscribers(renewedSerialized)
	if err != nil {
		return true, nil, nil, err
	}
	return false, renewed, deleted, nil
}

func (s *subscriberdbServicer) loadSubscribersPageFromCache(networkID string, req *lte_protos.ListSubscribersRequest, gateway *protos.Identity_Gateway) ([]*lte_protos.SubscriberData, string, error) {
	// If request page size is 0, return max entity load size
	pageSize := uint64(req.PageSize)
	if req.PageSize == 0 {
		pageSize = s.MaxProtosLoadSize
	}
	subProtosSerialized, nextToken, err := s.store.GetCachedByPage(networkID, req.PageToken, pageSize)
	if err != nil {
		return nil, "", err
	}
	subProtos, err := subscriberdb.DeserializeSubscribers(subProtosSerialized)
	if err != nil {
		return nil, "", err
	}
	subProtos, err = injectAPNResources(subProtos, gateway)
	if err != nil {
		return nil, "", err
	}

	return subProtos, nextToken, nil
}

// getDigestInfo returns the correctly formatted Digest and NoUpdates values
// according to the client digest.
func (s *subscriberdbServicer) getDigestInfo(clientDigest *protos.Digest, networkID string) (*protos.DigestTree, bool) {
	digestTree, err := syncstore.GetDigestTree(s.store, networkID)
	// If digest generation fails, the error is swallowed to not affect the main functionality
	if err != nil {
		glog.Errorf("Load digest for network %s failed: %+v", networkID, err)
		return &protos.DigestTree{RootDigest: &protos.Digest{Md5Base64Digest: ""}}, false
	}
	rootDigest := digestTree.RootDigest.GetMd5Base64Digest()
	noUpdates := rootDigest != "" && rootDigest == clientDigest.GetMd5Base64Digest()

	return digestTree, noUpdates
}

// shouldResync returns whether a gateway requires an Orc8r-directed resync by checking its
// last resync time.
func (l *subscriberdbServicer) shouldResync(network string, gateway string) bool {
	lastResyncTime, err := l.store.GetLastResync(network, gateway)
	// If check last resync time in store fails, swallow the error and stick to the original callpath
	if err != nil {
		glog.Errorf("check last resync time of gateway %+v of network %+v: %+v", gateway, network, err)
		return false
	}
	// Jitter the AGW sync interval by a fraction in the range of [0, 0.5] to ameliorate the thundering herd effect
	shouldResync := time.Now().Unix()-lastResyncTime > math.JitterInt64(l.ResyncIntervalSecs, gateway, 0.5)
	return shouldResync
}

func loadAPNs(gateway *protos.Identity_Gateway) (map[string]*lte_models.ApnConfiguration, lte_models.ApnResources, error) {
	networkID := gateway.NetworkId
	gatewayID := gateway.LogicalId
	lteGateway, err := configurator.LoadEntity(
		networkID, lte.CellularGatewayEntityType, gatewayID,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
		serdes.Entity,
	)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "load cellular gateway for gateway %s", gatewayID)
	}

	apnsByName, err := subscriberdb.LoadApnsByName(networkID)
	if err != nil {
		return nil, nil, err
	}
	apnResources, err := lte_models.LoadAPNResources(networkID, lteGateway.Associations.Filter(lte.APNResourceEntityType).Keys())
	if err != nil {
		return nil, nil, err
	}

	return apnsByName, apnResources, nil
}

// injectAPNResources adds the gateway-specific apn resources data to subscriber
// protos before returning to AGWs.
func injectAPNResources(subProtos []*lte_protos.SubscriberData, gateway *protos.Identity_Gateway) ([]*lte_protos.SubscriberData, error) {
	_, apnResources, err := loadAPNs(gateway)
	if err != nil {
		return nil, err
	}

	for _, subProto := range subProtos {
		if subProto.GetNon_3Gpp().GetApnConfig() == nil {
			continue
		}
		for _, apnConfig := range subProto.Non_3Gpp.ApnConfig {
			if apnResourceModel, ok := apnResources[apnConfig.ServiceSelection]; ok {
				apnConfig.Resource = apnResourceModel.ToProto()
			}
		}
	}
	return subProtos, nil
}
