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

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
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
	indexerTypes = []string{orc8r.DirectoryRecordType}
)

type indexerServicer struct{}

// NewIndexerServicer returns the state indexer for directoryd.
//
// The directoryd indexer performs the following indexing functions:
//	- sidToIMSI: map session ID to IMSI
//
// sidToIMSI
//
// Directoryd records are reported as {IMSI -> {records...}}. The sidToIMSI
// function is an online generation of the derived reverse map, producing {session ID -> IMSI}.
// NOTE: the indexer provides a best-effort generation of the session ID -> IMSI mapping, meaning
//	- a {session ID -> IMSI} mapping may be missing even though the IMSI has a session ID record
//	- a {session ID -> IMSI} mapping may be stale
func NewIndexerServicer() protos.IndexerServer {
	return &indexerServicer{}
}

func (i *indexerServicer) GetIndexerInfo(ctx context.Context, req *protos.GetIndexerInfoRequest) (*protos.GetIndexerInfoResponse, error) {
	res := &protos.GetIndexerInfoResponse{
		Version:    uint32(indexerVersion),
		StateTypes: indexerTypes,
	}
	return res, nil
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
	return setSessionID(networkID, states)
}

// setSessionID maps {sessionID -> IMSI}.
func setSessionID(networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	sessionIDToIMSI := map[string]string{}
	stateErrors := state_types.StateErrors{}
	for id, st := range states {
		sessionID, imsi, err := getSessionIDAndIMSI(id, st)
		if err != nil {
			stateErrors[id] = err
			continue
		}
		if sessionID == "" {
			glog.V(2).Infof("Session ID not found for record from %s", imsi)
			continue
		}

		sessionIDToIMSI[sessionID] = imsi
	}

	if len(sessionIDToIMSI) == 0 {
		return stateErrors, nil
	}

	err := directoryd.MapSessionIDsToIMSIs(networkID, sessionIDToIMSI)
	if err != nil {
		return stateErrors, errors.Wrapf(err, "update directoryd mapping of session IDs to IMSIs %+v", sessionIDToIMSI)
	}

	return stateErrors, nil
}

// getSessionIDAndIMSI extracts session ID and IMSI from the state.
// Returns (session ID, IMSI, error).
func getSessionIDAndIMSI(id state_types.ID, st state_types.State) (string, string, error) {
	imsi := id.DeviceID

	record, ok := st.ReportedState.(*directoryd_types.DirectoryRecord)
	if !ok {
		return "", "", fmt.Errorf(
			"convert reported state (id: <%+v>, state: <%+v>) to type %s",
			id, st, orc8r.DirectoryRecordType,
		)
	}
	sessionID, err := record.GetSessionID()
	if err != nil {
		return "", "", errors.Wrap(err, "extract session ID from record")
	}

	return sessionID, imsi, nil
}
