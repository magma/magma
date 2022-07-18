/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/lte/cloud/go/lte"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
)

const (
	EpsAuthStateStore = "eps_auth_state_store"
	EpsAuthType       = "eps_auth_state_type"
)

type SubscriberDBStorage interface {
	// GetSubscriberData returns subscriber data of given subscriber ID,
	// GetSubscriberData only fills in auth related fields of the SubscriberProfile
	GetSubscriberData(*lteprotos.SubscriberID, *protos.NetworkID) (*lteprotos.SubscriberData, error)
	// GetSubscriberDataProfile returns subscriber data & its [non]3GPP profile information (APN, static IPs, ...)
	GetSubscriberDataProfile(*lteprotos.SubscriberID, *protos.NetworkID) (*lteprotos.SubscriberData, map[string]string, []string, error)
	// UpdateSubscriberAuthNextSeq sets subscriber's AuthNextSeq to the provided value
	UpdateSubscriberAuthNextSeq(*lteprotos.SubscriberData) (*protos.Void, error)
	// IncrementSubscriberAuthNextSeq increments subscriber's AuthNextSeq
	IncrementSubscriberAuthNextSeq(*lteprotos.SubscriberData) (*protos.Void, error)
}

// needs to be implemented
type subscriberDBStorageImpl struct {
	blobstore.StoreFactory
}

func NewSubscriberDBStorage(storeFactory blobstore.StoreFactory) SubscriberDBStorage {
	return &subscriberDBStorageImpl{StoreFactory: storeFactory}
}

// GetSubscriberData returns subscriber data of given subscriber ID, the subscribers StaticIps map and associated APNs
func (s *subscriberDBStorageImpl) GetSubscriberData(
	sid *lteprotos.SubscriberID, networkID *protos.NetworkID) (*lteprotos.SubscriberData, error) {

	lc := configurator.EntityLoadCriteria{LoadConfig: true}
	ent, err := configurator.LoadEntity(
		context.Background(), networkID.GetId(), lte.SubscriberEntityType, lteprotos.SidString(sid), lc, serdes.Entity)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"error loading subscriber ent for network ID: %s, SID: %s: %v", networkID.GetId(), sid.GetId(), err)
	}
	subData := &lteprotos.SubscriberData{Sid: sid, NetworkId: networkID}
	if ent.Config == nil {
		return nil, status.Errorf(
			codes.NotFound,
			"empty subscriber configuration for network ID: %s, SID: %s", networkID.GetId(), sid.GetId())
	}
	if cfg := ent.Config.(*models.SubscriberConfig); cfg != nil && cfg.Lte != nil {
		subData.SubProfile = string(*cfg.Lte.SubProfile)
		subData.Lte = &lteprotos.LTESubscription{
			State: lteprotos.LTESubscription_LTESubscriptionState(
				lteprotos.LTESubscription_LTESubscriptionState_value[cfg.Lte.State]),
			AuthAlgo: lteprotos.LTESubscription_LTEAuthAlgo(
				lteprotos.LTESubscription_LTEAuthAlgo_value[cfg.Lte.AuthAlgo]),
			AuthKey: cfg.Lte.AuthKey,
			AuthOpc: cfg.Lte.AuthOpc,
		}
	}
	return subData, nil
}

// GetSubscriberData returns subscriber data of given subscriber ID, the subscribers StaticIps map and associated APNs
func (s *subscriberDBStorageImpl) GetSubscriberDataProfile(
	sid *lteprotos.SubscriberID,
	networkID *protos.NetworkID) (*lteprotos.SubscriberData, map[string]string, []string, error) {

	lc := configurator.EntityLoadCriteria{
		LoadConfig:         true,
		LoadAssocsToThis:   false,
		LoadAssocsFromThis: true,
	}
	ent, err := configurator.LoadEntity(
		context.Background(), networkID.GetId(), lte.SubscriberEntityType, lteprotos.SidString(sid), lc, serdes.Entity)

	if err != nil {
		return nil, nil, nil, status.Errorf(
			codes.NotFound,
			"error loading subscriber ent with assocs for network ID: %s, SID: %s: %v",
			networkID.GetId(), sid.GetId(), err)
	}
	subData := &lteprotos.SubscriberData{
		Sid:       sid,
		NetworkId: networkID,
	}
	if ent.Config == nil {
		return nil, nil, nil, status.Errorf(
			codes.NotFound,
			"missing subscriber configuration for network ID: %s, SID: %s", networkID.GetId(), sid.GetId())
	}
	var staticIps models.SubscriberStaticIps
	if cfg := ent.Config.(*models.SubscriberConfig); cfg != nil && cfg.Lte != nil {
		subData.SubProfile = string(*cfg.Lte.SubProfile)
		subData.Lte = &lteprotos.LTESubscription{
			State: lteprotos.LTESubscription_LTESubscriptionState(
				lteprotos.LTESubscription_LTESubscriptionState_value[cfg.Lte.State]),
			AuthAlgo: lteprotos.LTESubscription_LTEAuthAlgo(
				lteprotos.LTESubscription_LTEAuthAlgo_value[cfg.Lte.AuthAlgo]),
			AuthKey: cfg.Lte.AuthKey,
			AuthOpc: cfg.Lte.AuthOpc,
		}
		staticIps = cfg.StaticIps
	}
	return subData, staticIps, ent.Associations.Filter(lte.APNEntityType).Keys(), nil
}

// UpdateSubscriberAuthNextSeq sets subscriber's AuthNextSeq to the provided value
// Note: the implementation uses the blob version as the next sequence
func (s *subscriberDBStorageImpl) UpdateSubscriberAuthNextSeq(sd *lteprotos.SubscriberData) (*protos.Void, error) {
	ret := &protos.Void{}
	if sd == nil {
		return ret, status.Errorf(codes.InvalidArgument, "nil Subscriber Data")
	}
	store, err := s.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return ret, status.Errorf(codes.Unavailable, "error starting update nex seq transaction: %v", err)
	}
	defer store.Rollback()

	err = store.Write(sd.GetNetworkId().GetId(), []blobstore.Blob{{
		Type:    EpsAuthType,
		Key:     sd.GetSid().GetId(),
		Value:   []byte{},
		Version: sd.GetState().GetLteAuthNextSeq(),
	}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "state write error: %v", err)
	}
	return ret, store.Commit()
}

// IncrementSubscriberAuthNextSeq increments subscriber's AuthNextSeq
// Note: the implementation uses the blob version as the next sequence
func (s *subscriberDBStorageImpl) IncrementSubscriberAuthNextSeq(sd *lteprotos.SubscriberData) (*protos.Void, error) {
	ret := &protos.Void{}
	if sd == nil {
		return ret, status.Errorf(codes.InvalidArgument, "nil Subscriber Data")
	}
	store, err := s.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return ret, status.Errorf(codes.Unavailable, "error starting increment nex seq transaction: %v", err)
	}
	defer store.Rollback()

	err = store.IncrementVersion(sd.GetNetworkId().GetId(), storage.TK{Type: EpsAuthType, Key: sd.GetSid().GetId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "auth next seq increment error: %v", err)
	}
	return ret, store.Commit()
}
