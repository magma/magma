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

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/directoryd"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/protos"
	state_types "magma/orc8r/cloud/go/services/state/types"
)

const (
	indexerVersion indexer.Version = 1
)

type teidType int

const (
	// types of TEIDs
	controlPlaneTeid teidType = iota
	userPlaneTeid
)

var (
	indexerTypes = []string{orc8r.DirectoryRecordType}
)

type directorydRecordParameters struct {
	imsi      string
	sessionId string
	cTeids    []string
	uTeids    []string
	hwid      string
}

type indexerServicer struct{}

// NewIndexerServicer returns the state indexer for directoryd.
//
// The directoryd indexer performs the following indexing functions:
//   - sidToIMSI: map session ID to IMSI
//
// sidToIMSI
//
// Directoryd records are reported as {IMSI -> {records...}}. The sidToIMSI
// function is an online generation of the derived reverse map, producing {session ID -> IMSI}.
// NOTE: the indexer provides a best-effort generation of the session ID -> IMSI mapping, meaning
//   - a {session ID -> IMSI} mapping may be missing even though the IMSI has a session ID record
//   - a {session ID -> IMSI} mapping may be stale
func NewIndexerServicer() protos.IndexerServer {
	return &indexerServicer{}
}

func (i *indexerServicer) Index(ctx context.Context, req *protos.IndexRequest) (*protos.IndexResponse, error) {
	states, err := state_types.MakeStatesByID(req.States, serdes.State)
	if err != nil {
		return nil, err
	}
	stErrs, err := setSecondaryStates(ctx, req.NetworkId, states)
	if err != nil {
		return nil, err
	}
	res := &protos.IndexResponse{StateErrors: state_types.MakeProtoStateErrors(stErrs)}
	return res, nil
}

func (i *indexerServicer) DeIndex(ctx context.Context, req *protos.DeIndexRequest) (*protos.DeIndexResponse, error) {
	states, err := state_types.MakeStatesByID(req.States, serdes.State)
	if err != nil {
		return &protos.DeIndexResponse{}, err
	}
	stErrs, err := unsetSecondaryStates(ctx, req.NetworkId, states)
	if err != nil {
		return nil, err
	}
	res := &protos.DeIndexResponse{StateErrors: state_types.MakeProtoStateErrors(stErrs)}
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

// setSecondaryStates maps {sessionID -> IMSI} and {TEID -> HWID}
// Will attempt to update all secondary states, but will return error if any fails
func setSecondaryStates(ctx context.Context, networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	sessionIDToIMSI, cTeidToHwId, uTeidToHwId, stateErrors := getMappings(states)
	if len(sessionIDToIMSI) == 0 && len(cTeidToHwId) == 0 {
		return stateErrors, nil
	}
	errs := &multierror.Error{}
	if len(sessionIDToIMSI) != 0 {
		err := directoryd.MapSessionIDsToIMSIs(ctx, networkID, sessionIDToIMSI)
		errs = errsAppend(errs, "failed to update directoryd mapping of session IDs to IMSIs %+v %v", sessionIDToIMSI, err)
	}
	if len(cTeidToHwId) != 0 {
		err := directoryd.MapSgwCTeidToHWID(ctx, networkID, cTeidToHwId)
		errs = errsAppend(errs, "failed to update directoryd mapping of control plane teid To HwID %+v %v", sessionIDToIMSI, err)
	}
	if len(uTeidToHwId) != 0 {
		err := directoryd.MapSgwUTeidToHWID(ctx, networkID, uTeidToHwId)
		errs = errsAppend(errs, "failed to update directoryd mapping of user plane teid To HwID %+v %v", sessionIDToIMSI, err)
	}

	// errs will only be nil if both updates succeeded
	return stateErrors, errs.ErrorOrNil()
}

// unsetSecondaryStates removes {sessionID -> IMSI} and {TEID -> HWID} mappings
func unsetSecondaryStates(ctx context.Context, networkID string, states state_types.StatesByID) (state_types.StateErrors, error) {
	sessionIDToIMSI, cTeidToHwId, uTeidToHwId, stateErrors := getMappings(states)
	errs := &multierror.Error{}

	err := unsetTeids(ctx, controlPlaneTeid, networkID, cTeidToHwId)
	errs = multierror.Append(errs, err)
	err = unsetTeids(ctx, userPlaneTeid, networkID, uTeidToHwId)
	errs = multierror.Append(errs, err)
	err = unsetSessionIDs(ctx, networkID, sessionIDToIMSI)
	errs = multierror.Append(errs, err)
	return stateErrors, errs.ErrorOrNil()
}

// unsetTeids removes teidType {TEID -> TEID} mappings
func unsetTeids(ctx context.Context, tType teidType, networkID string, cTeidToHwId map[string]string) error {
	var teids []string
	for teid := range cTeidToHwId {
		teids = append(teids, teid)
	}
	var err error
	if len(teids) != 0 {
		switch tType {
		case controlPlaneTeid:
			err = directoryd.UnmapSgwCTeidToHWID(ctx, networkID, teids)
		case userPlaneTeid:
			err = directoryd.UnmapSgwUTeidToHWID(ctx, networkID, teids)
		default:
			err = fmt.Errorf("unsetTeids: TeidType not found")
		}
		if err != nil {
			err = fmt.Errorf("Unmap TeidToHWID failed (teidType: %d): %s", tType, err)
			glog.Error(err)
		}
	}
	return err
}

// unsetSessionIDs removes {sessionID -> IMSI} mappings
func unsetSessionIDs(ctx context.Context, networkID string, sessionIDToIMSI map[string]string) error {
	var sessionIDs []string
	for sessionID := range sessionIDToIMSI {
		sessionIDs = append(sessionIDs, sessionID)
	}
	var err error
	if len(sessionIDs) != 0 {
		err = directoryd.UnmapSessionIDsToIMSIs(ctx, networkID, sessionIDs)
		if err != nil {
			err = fmt.Errorf("UnmapSessionIDsToIMSIs failed: %s", err)
			glog.Error(err)
		}
	}
	return err
}

// getMappings builds SessionID to IMSI and TEIDs to HWID maps from state
func getMappings(states state_types.StatesByID) (
	sessionIDToIMSI map[string]string,
	cTeidToHwId map[string]string,
	uTeidToHwId map[string]string,
	stateErrors state_types.StateErrors,
) {
	sessionIDToIMSI = map[string]string{}
	cTeidToHwId = map[string]string{}
	uTeidToHwId = map[string]string{}
	stateErrors = state_types.StateErrors{}
	for id, st := range states {
		params, err := extractRecordParameters(id, st)
		if err != nil {
			stateErrors[id] = err
			continue
		}
		if params.sessionId != "" {
			sessionIDToIMSI[params.sessionId] = params.imsi
		}

		for _, cTeid := range params.cTeids {
			cTeidToHwId[cTeid] = params.hwid
		}
		for _, uTeid := range params.uTeids {
			uTeidToHwId[uTeid] = params.hwid
		}
	}
	return
}

// extractRecordParameters extracts IMSI, SessionID, TEID and HWID from directory record
// Returns error if any error is found. No partial updates are allowed.
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
	cTeids, hwid, err := getTeidToHwIdPair(controlPlaneTeid, record)
	if err != nil {
		return nil, err
	}
	uTeids, _, err := getTeidToHwIdPair(userPlaneTeid, record)
	if err != nil {
		return nil, err
	}
	// log an error in case blank sessionId and no TEID
	if sessionID == "" && len(cTeids) == 0 {
		glog.V(2).Infof("Session ID not found for IMSI %s in record %v", imsi, record)
	}
	return &directorydRecordParameters{
		imsi:      imsi,
		sessionId: sessionID,
		cTeids:    cTeids,
		uTeids:    uTeids,
		hwid:      hwid,
	}, nil
}

// getTeidToHwIdPair will return all TEIDs of a type for that IMSI and its current location (HWID)
func getTeidToHwIdPair(tType teidType, record *directoryd_types.DirectoryRecord) ([]string, string, error) {
	// GetLocationHistory will always return
	hwid, err := record.GetCurrentLocation()
	if err != nil {
		return nil, "", err
	}
	// get either cTeid or uTeid
	var teids []string
	switch tType {
	case controlPlaneTeid:
		teids, err = record.GetSgwCTeids()
	case userPlaneTeid:
		teids, err = record.GetSgwUTeids()
	default:
		err = fmt.Errorf("getTeidToHwIdPair: TeidType not found")
	}
	return teids, hwid, err
}

// errsAppend concatenates message, sessionIDToIMSI and err into a new error message and appends to errs
func errsAppend(errs *multierror.Error, message string, sessionIDToIMSI map[string]string, err error) *multierror.Error {
	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf(message, sessionIDToIMSI, err))
	}
	return errs
}
