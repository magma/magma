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
	multierrors "magma/orc8r/lib/go/errors"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	indexerVersion indexer.Version = 1
)

var (
	indexerTypes = []string{orc8r.DirectoryRecordType}
)

type directorydRecordParameters struct {
	imsi      string
	sessionId string
	teids     []string
	hwid      string
}

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
	return setSecondaryStates(networkID, states)
}

// setSecondaryState maps {sessionID -> IMSI} and {teid -> HwId}
// Will attempt to update all secondary states, but will return error if any fails
func setSecondaryStates(networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	sessionIDToIMSI := map[string]string{}
	teidoHwId := map[string]string{}
	stateErrors := state_types.StateErrors{}
	for id, st := range states {
		params, err := extractRecordParameters(id, st)
		if err != nil {
			stateErrors[id] = err
			continue
		}
		if params.sessionId != "" {
			sessionIDToIMSI[params.sessionId] = params.imsi
		}

		for _, teid := range params.teids {
			teidoHwId[teid] = params.hwid
		}
	}

	if len(sessionIDToIMSI) == 0 && len(teidoHwId) == 0 {
		return stateErrors, nil
	}

	multiError := multierrors.NewMulti()
	if len(sessionIDToIMSI) != 0 {
		err := directoryd.MapSessionIDsToIMSIs(networkID, sessionIDToIMSI)
		multiError = multiError.AddFmt(err, "failed to update directoryd mapping of session IDs to IMSIs %+v", sessionIDToIMSI)
	}
	if len(teidoHwId) != 0 {
		err := directoryd.MapSgwCTeidToHWID(networkID, teidoHwId)
		multiError = multiError.AddFmt(err, "failed to update directoryd mapping of teid To HwID %+v", sessionIDToIMSI)
	}
	// multiError will only be nil if both updates succeeded
	return stateErrors, multiError.AsError()
}

// extractRecordParameters extracts IMSI, SessionID, TEID and HwID from directory record
// Returns error any error is found. No partial updates are allowed.
func extractRecordParameters(id state_types.ID, st state_types.State) (*directorydRecordParameters, error) {
	imsi := id.DeviceID
	record, ok := st.ReportedState.(*directoryd_types.DirectoryRecord)
	if !ok {
		return nil, fmt.Errorf(
			"convert reported state (id: <%+v>, state: <%+v>) to type %s",
			id, st, orc8r.DirectoryRecordType,
		)
	}
	sessionID, err := record.GetSessionID()
	if err != nil {
		return nil, err
	}
	teids, hwid, err := getTeidToHwIdPair(record)
	if err != nil {
		return nil, err
	}
	// log an error in case blank sessionId and no teid
	if sessionID == "" && len(teids) == 0 {
		glog.V(2).Infof("Session ID not found for record from %s", imsi)
	}

	return &directorydRecordParameters{
		imsi:      imsi,
		sessionId: sessionID,
		teids:     teids,
		hwid:      hwid,
	}, nil
}

// getTeidToHwIdPair will return all the TEIDs for that IMSI and its current location (HwID)
func getTeidToHwIdPair(record *directoryd_types.DirectoryRecord) ([]string, string, error) {
	teids, err := record.GetSgwCTeids()
	if err != nil {
		return nil, "", err
	}
	if len(teids) == 0 {
		return nil, "", nil
	}

	// GetLocationHistory will always return
	hwid, err := record.GetCurrentLocation()
	if err != nil {
		return nil, "", err
	}
	return teids, hwid, nil
}
