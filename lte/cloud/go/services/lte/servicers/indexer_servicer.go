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

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_api "magma/lte/cloud/go/services/lte"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
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
	indexerTypes = []string{lte.EnodebStateType}
)

type indexerServicer struct{}

// NewIndexerServicer returns the state indexer for the lte service.
//
// Enodeb state is reported as ENB SN -> EnodebState json model for a given
// networkID. Since multiple gateways can report this state, index the state to
// add gatewayID as an additional primary key. This allows for differentiation
// between state reported from different gateways for the same ENB SN.
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
	return setEnodebState(networkID, states)
}

// setEnodebState stores EnodebState with reporterID as an additional PK
func setEnodebState(networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	stateErrors := state_types.StateErrors{}
	for id, st := range states {
		// Set time reported before storing
		enbState, ok := st.ReportedState.(*lte_models.EnodebState)
		if !ok {
			stateErrors[id] = fmt.Errorf("error converting state for deviceID %s to EnodebModel", id.DeviceID)
			continue
		}
		enbState.TimeReported = st.TimeMs
		serializedState, err := serde.Serialize(st.ReportedState, lte.EnodebStateType, serdes.State)
		if err != nil {
			stateErrors[id] = fmt.Errorf("error serializing EnodebState for deviceID %s", id.DeviceID)
			continue
		}
		gwEnt, err := configurator.LoadEntityForPhysicalID(st.ReporterID, configurator.EntityLoadCriteria{}, serdes.Entity)
		if err != nil {
			stateErrors[id] = errors.Wrap(err, "error loading gatewayID")
			continue
		}
		err = lte_api.SetEnodebState(networkID, gwEnt.Key, id.DeviceID, serializedState)
		if err != nil {
			stateErrors[id] = errors.Wrap(err, "error setting enodeb state")
			continue
		}
		glog.V(2).Infof("successfully stored ENB state for eNB SN: %s, gatewayID: %s:w", id.DeviceID, gwEnt.Key)
	}
	return stateErrors, nil
}
