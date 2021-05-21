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

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/subscriberdb"
	subscriberdb_protos "magma/lte/cloud/go/services/subscriberdb/protos"
	subscriberdb_state "magma/lte/cloud/go/services/subscriberdb/state"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/protos"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	indexerVersion indexer.Version = 1
)

var (
	indexerTypes = []string{lte.MobilitydStateType}
)

type indexerServicer struct{}

// NewIndexerServicer returns the state indexer for subscriberdb.
//
// The directoryd indexer performs the following indexing functions:
//	- ipToIMSI: map IP address to IMSI
//
// ipToIMSI
//
// Mobilityd state is reported as IMSI.APN -> arbitrary JSON. The arbitrary
// JSON contains the assigned IP address for the IMSI under the associated APN,
// which we extract to form the reverse map.
// NOTE: the indexer provides a best-effort generation of the IP -> IMSI map,
// meaning
//	- an {IP -> IMSI} mapping only makes sense under non-NAT deployments and/or
//	  for non-private IP addresses
//	- an {IP -> IMSI} mapping may be missing even though the IMSI is assigned
//	  that IP
//	- an {IP -> IMSI} mapping may be stale (caller should check for staleness)
func NewIndexerServicer() protos.IndexerServer {
	return &indexerServicer{}
}

func (i *indexerServicer) Index(ctx context.Context, req *protos.IndexRequest) (*protos.IndexResponse, error) {
	states, err := state_types.MakeStatesByID(req.States, serdes.State)
	if err != nil {
		return nil, err
	}
	stErrs, err := indexImpl(req.NetworkId, states)
	if err != nil {
		return nil, err
	}
	res := &protos.IndexResponse{StateErrors: state_types.MakeProtoStateErrors(stErrs)}
	return res, nil
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

func indexImpl(networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	return setIPMappings(networkID, states)
}

// setIPMappings maps {IP -> IMSI}.
func setIPMappings(networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
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

	err := subscriberdb.SetIMSIsForIPs(networkID, ipMappings)
	if err != nil {
		return stateErrors, errors.Wrapf(err, "update directoryd mapping of session IDs to IMSIs %+v", ipMappings)
	}

	return stateErrors, nil
}
