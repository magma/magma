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
	"fmt"

	"github.com/golang/glog"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/subscriberdb"
	subscriberdb_protos "magma/lte/cloud/go/services/subscriberdb/protos"
	subscriberdb_state "magma/lte/cloud/go/services/subscriberdb/state"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/protos"
	state_types "magma/orc8r/cloud/go/services/state/types"
)

const (
	indexerVersion indexer.Version = 1
)

var (
	indexerTypes = []string{lte.MobilitydStateType, lte.GatewaySubscriberStateType}
)

type indexerServicer struct {
	subscriberStore storage.SubscriberStorage
}

// NewIndexerServicer returns the state indexer for subscriberdb.
//
// The directoryd indexer performs the following indexing functions:
//   - ipToIMSI: map IP address to IMSI
//
// ipToIMSI
//
// Mobilityd state is reported as IMSI.APN -> arbitrary JSON. The arbitrary
// JSON contains the assigned IP address for the IMSI under the associated APN,
// which we extract to form the reverse map.
// NOTE: the indexer provides a best-effort generation of the IP -> IMSI map,
// meaning
//   - an {IP -> IMSI} mapping only makes sense under non-NAT deployments and/or
//     for non-private IP addresses
//   - an {IP -> IMSI} mapping may be missing even though the IMSI is assigned
//     that IP
//   - an {IP -> IMSI} mapping may be stale (caller should check for staleness)
//
// Gateway Subscriber State is reported as a map IMSI to arbitrary JSON (reported
// state for that IMSI).
// The indexer updates gateway subscriber state in SubscriberStorage.
// It deletes all entries for the gateway ID and writes the current
// IMSI,state pairs into the SubscriberStorage.
func NewIndexerServicer(ss storage.SubscriberStorage) protos.IndexerServer {
	return &indexerServicer{subscriberStore: ss}
}

func (i *indexerServicer) Index(ctx context.Context, req *protos.IndexRequest) (*protos.IndexResponse, error) {
	states, err := state_types.MakeStatesByID(req.States, serdes.State)
	if err != nil {
		return nil, err
	}
	stErrs, err := i.indexImpl(ctx, req.NetworkId, states)
	if err != nil {
		return nil, err
	}
	res := &protos.IndexResponse{StateErrors: state_types.MakeProtoStateErrors(stErrs)}
	return res, nil
}

func (i *indexerServicer) DeIndex(ctx context.Context, req *protos.DeIndexRequest) (*protos.DeIndexResponse, error) {
	return &protos.DeIndexResponse{}, nil
}

func (i *indexerServicer) PrepareReindex(ctx context.Context, req *protos.PrepareReindexRequest) (*protos.PrepareReindexResponse, error) {
	return &protos.PrepareReindexResponse{}, nil
}

func (i *indexerServicer) CompleteReindex(ctx context.Context, req *protos.CompleteReindexRequest) (*protos.CompleteReindexResponse, error) {
	if req.FromVersion == 0 && req.ToVersion == 1 {
		return &protos.CompleteReindexResponse{}, nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "unsupported from/to for CompleteReindex: %v to %v", req.FromVersion, req.ToVersion)
}

func (i *indexerServicer) indexImpl(ctx context.Context, networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	statesMobilityD := state_types.StatesByID{}
	statesSubscribers := state_types.StatesByID{}
	for id, st := range states {
		switch id.Type {
		case lte.MobilitydStateType:
			statesMobilityD[id] = st
		case lte.GatewaySubscriberStateType:
			statesSubscribers[id] = st
		default:
			glog.Errorf("Unsupported state type %s", id.Type)
		}
	}
	var stErrs state_types.StateErrors
	errs := &multierror.Error{}
	if len(statesMobilityD) > 0 {
		stErrsMobilityD, err := setIPMappings(ctx, networkID, statesMobilityD)
		stErrs = stErrsMobilityD
		errs = multierror.Append(errs, err)
	}
	if len(statesSubscribers) > 0 {
		err := i.setGatewaySubscriberStates(networkID, statesSubscribers)
		errs = multierror.Append(errs, err)
	}
	return stErrs, errs.ErrorOrNil()
}

// setIPMappings maps {IP -> IMSI}.
func setIPMappings(ctx context.Context, networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	var ipMappings []*subscriberdb_protos.IPMapping
	stateErrors := state_types.StateErrors{}
	for id, st := range states {
		reportedState := st.ReportedState.(*state.ArbitraryJSON)
		ip, err := subscriberdb_state.GetAssignedIPAddress(*reportedState)
		if err != nil {
			stateErrors[id] = err
			continue
		}
		if ip == "" {
			glog.V(2).Infof("IP missing from mobilityd state for state key %s", id.DeviceID)
			continue
		}
		imsi, apn, err := subscriberdb_state.GetIMSIAndAPNFromMobilitydStateKey(id.DeviceID)
		if err != nil {
			stateErrors[id] = err
			continue
		}
		ipMappings = append(ipMappings, &subscriberdb_protos.IPMapping{Ip: ip, Imsi: imsi, Apn: apn})
	}

	if len(ipMappings) == 0 {
		return stateErrors, nil
	}

	err := subscriberdb.SetIMSIsForIPs(ctx, networkID, ipMappings)
	if err != nil {
		return stateErrors, fmt.Errorf("update directoryd mapping of session IDs to IMSIs %+v: %w", ipMappings, err)
	}

	return stateErrors, nil
}

func (i *indexerServicer) setGatewaySubscriberStates(networkID string, states state_types.StatesByID) *multierror.Error {
	errs := &multierror.Error{}
	for id, st := range states {
		err := i.subscriberStore.SetAllSubscribersForGateway(networkID, id.DeviceID, st.ReportedState.(*storage.GatewaySubscriberState))
		if err != nil {
			glog.Errorf("Error setting subscriber state for gateway %s: %s", id.DeviceID, err)
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}
