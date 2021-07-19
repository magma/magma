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

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
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
	digestsEnabled         bool
	changesetSizeThreshold int
	maxProtosLoadSize      uint64
	resyncIntervalSecs uint64
	digestStore            storage.DigestStore
	perSubDigestStore      *storage.PerSubDigestStore
	subStore               *storage.SubStore
	lastResyncTimeStore   *storage.LastResyncTimeStore
}

func NewSubscriberdbServicer(
	config subscriberdb.Config,
	digestStore storage.DigestStore,
	perSubDigestStore *storage.PerSubDigestStore,
	subStore *storage.SubStore,
	lastResyncTimeStore *storage.LastResyncTimeStore,
) lte_protos.SubscriberDBCloudServer {
	servicer := &subscriberdbServicer{
		digestsEnabled:         config.DigestsEnabled,
		changesetSizeThreshold: config.ChangesetSizeThreshold,
		maxProtosLoadSize:      config.MaxProtosLoadSize,
		resyncIntervalSecs: config.ResyncIntervalSecs,
		digestStore:            digestStore,
		perSubDigestStore:      perSubDigestStore,
		subStore:               subStore,
		lastResyncTimeStore:   lastResyncTimeStore,
	}
	return servicer
}

func (s *subscriberdbServicer) CheckSubscribersInSync(
	ctx context.Context,
	req *lte_protos.CheckSubscribersInSyncRequest,
) (*lte_protos.CheckSubscribersInSyncResponse, error) {
	if !s.digestsEnabled {
		return &lte_protos.CheckSubscribersInSyncResponse{InSync: false}, nil
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
	if s.shouldCloudDirectedResync(networkID, gatewayID) {
		return &lte_protos.CheckSubscribersInSyncResponse{InSync: false}, nil
	}

	_, inSync := s.getDigestInfo(req.FlatDigest, networkID)
	res := &lte_protos.CheckSubscribersInSyncResponse{InSync: inSync}
	return res, nil
}

func (s *subscriberdbServicer) SyncSubscribers(
	ctx context.Context,
	req *lte_protos.SyncSubscribersRequest,
) (*lte_protos.SyncSubscribersResponse, error) {
	if !s.digestsEnabled {
		return &lte_protos.SyncSubscribersResponse{Resync: true}, nil
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
	if s.shouldCloudDirectedResync(networkID, gatewayID) {
		return &lte_protos.SyncSubscribersResponse{Resync: true}, nil
	}

	flatDigest, err := storage.GetDigest(s.digestStore, networkID)
	if err != nil {
		return nil, err
	}

	clientPerSubDigests := req.PerSubDigests
	cloudPerSubDigests, err := s.perSubDigestStore.GetDigest(networkID)
	if err != nil {
		return nil, err
	}
	resync, renewed, deleted, err := s.getSubscribersChangeset(networkID, clientPerSubDigests, cloudPerSubDigests)
	if err != nil {
		return nil, err
	}
	if resync {
		return &lte_protos.SyncSubscribersResponse{Resync: true}, nil
	}

	// Since the cached protos don't contain gateway-specific information, inject
	// the apn resource configs related to the gateway
	renewed, err = injectAPNResources(renewed, gateway)
	if err != nil {
		return nil, err
	}
	res := &lte_protos.SyncSubscribersResponse{
		FlatDigest:    &lte_protos.Digest{Md5Base64Digest: flatDigest},
		PerSubDigests: cloudPerSubDigests,
		ToRenew:       renewed,
		Deleted:       deleted,
		Resync:        false,
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
	if s.digestsEnabled {
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

	flatDigest := &lte_protos.Digest{Md5Base64Digest: ""}
	perSubDigests := []*lte_protos.SubscriberDigestWithID{}
	// The digests are sent back during the request for the first page of subscriber data
	if req.PageToken == "" && s.digestsEnabled {
		flatDigest, _ = s.getDigestInfo(&lte_protos.Digest{Md5Base64Digest: ""}, networkID)
		perSubDigests, err = s.perSubDigestStore.GetDigest(networkID)
		if err != nil {
			glog.Errorf("Failed to get per-sub digests from store for network %+v: %+v", networkID, err)
		}
	}

	// At the AGW request for the last page, update the lastResyncTime of the gateway to the current time
	// NOTE: Since the resync is orc8r-directed, and orc8r doesn't track the request status on the AGW side,
	// orc8r takes the AGW request for the last page as an approximate indication of the completion of a resync
	if nextToken == "" {
		err := s.lastResyncTimeStore.Set(networkID, gatewayID, uint64(time.Now().Unix()))
		if err != nil {
			glog.Errorf("Failed to set last resync time for gateway %+v of network %+v: %+v", gatewayID, networkID, err)
		}
	}

	listRes := &lte_protos.ListSubscribersResponse{
		Subscribers:   subProtos,
		NextPageToken: nextToken,
		FlatDigest:    flatDigest,
		PerSubDigests: perSubDigests,
	}
	return listRes, nil
}

// getSubscribersChangeset compares the cloud and AGW digests and returns
// 1. Whether a resync is required for this AGW.
// 2. If no resync, the list of subscriber configs to be renewed.
// 3. If no resync, the list of subscriber IDs to be deleted.
// 4. Any error that occurred.
func (s *subscriberdbServicer) getSubscribersChangeset(networkID string, clientDigests []*lte_protos.SubscriberDigestWithID, cloudDigests []*lte_protos.SubscriberDigestWithID) (bool, []*lte_protos.SubscriberData, []string, error) {
	toRenew, deleted := subscriberdb.GetPerSubscriberDigestsDiff(clientDigests, cloudDigests)
	if len(toRenew) > s.changesetSizeThreshold || len(toRenew) > int(s.maxProtosLoadSize) {
		return true, nil, nil, nil
	}

	sids := funk.Keys(toRenew).([]string)
	renewed, err := s.subStore.GetSubscribers(networkID, sids)
	if err != nil {
		return true, nil, nil, err
	}
	return false, renewed, deleted, nil
}

func (s *subscriberdbServicer) loadSubscribersPageFromCache(networkID string, req *lte_protos.ListSubscribersRequest, gateway *protos.Identity_Gateway) ([]*lte_protos.SubscriberData, string, error) {
	// If request page size is 0, return max entity load size
	pageSize := uint64(req.PageSize)
	if req.PageSize == 0 {
		pageSize = s.maxProtosLoadSize
	}
	subProtos, nextToken, err := s.subStore.GetSubscribersPage(networkID, req.PageToken, pageSize)
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
func (s *subscriberdbServicer) getDigestInfo(clientDigest *lte_protos.Digest, networkID string) (*lte_protos.Digest, bool) {
	digest, err := storage.GetDigest(s.digestStore, networkID)
	// If digest generation fails, the error is swallowed to not affect the main functionality
	if err != nil {
		glog.Errorf("Generating digest for network %s failed: %+v", networkID, err)
		return &lte_protos.Digest{Md5Base64Digest: ""}, false
	}

	noUpdates := digest != "" && digest == clientDigest.GetMd5Base64Digest()
	digestProto := &lte_protos.Digest{Md5Base64Digest: digest}
	return digestProto, noUpdates
}

// shouldCloudDirectedResync returns whether a gateway requires a orc8r-directed resync by checking its
// last resync time.
func (l *subscriberdbServicer) shouldCloudDirectedResync(network string, gateway string) bool {
	lastResyncTime, err := l.lastResyncTimeStore.Get(network, gateway)
	// If check last resync time in store fails, swallow the error and stick to the original callpath
	if err != nil {
		glog.Errorf("check last resync time of gateway %+v of network %+v: %+v", gateway, network, err)
		return false
	}
	shouldResync := uint64(time.Now().Unix())-lastResyncTime > l.resyncIntervalSecs
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
