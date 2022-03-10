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
	// GetSubscriberData returns subscriber data of given subscriber ID
	GetSubscriberData(*lteprotos.SubscriberID, *protos.NetworkID) (*lteprotos.SubscriberData, error)
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

// GetSubscriberData returns subscriber data of given subscriber ID
func (s *subscriberDBStorageImpl) GetSubscriberData(
	sid *lteprotos.SubscriberID, nid *protos.NetworkID) (*lteprotos.SubscriberData, error) {

	lc := configurator.EntityLoadCriteria{
		PageSize:           0,
		PageToken:          "",
		LoadConfig:         true,
		LoadAssocsToThis:   false,
		LoadAssocsFromThis: false,
	}

	ent, err := configurator.LoadEntity(
		context.Background(), nid.GetId(), lte.SubscriberEntityType, lteprotos.SidString(sid), lc, serdes.Entity)

	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			"error loading subscriber entity for NID: %s, SID: %s: %v", nid.GetId(), sid.GetId(), err)
	}
	subData := &lteprotos.SubscriberData{
		Sid:       sid,
		NetworkId: nid,
	}
	if ent.Config == nil {
		return nil, status.Errorf(
			codes.NotFound, "missing subscriber configuration for NID: %s, SID: %s", nid.GetId(), sid.GetId())
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
