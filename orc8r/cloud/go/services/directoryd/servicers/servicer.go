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
	"sort"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/types"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type directoryLookupServicer struct {
	store storage.DirectorydStorage
}

func NewDirectoryLookupServicer(store storage.DirectorydStorage) (protos.DirectoryLookupServer, error) {
	srv := &directoryLookupServicer{store: store}
	return srv, nil
}

func (d *directoryLookupServicer) GetHostnameForHWID(
	ctx context.Context, req *protos.GetHostnameForHWIDRequest,
) (*protos.GetHostnameForHWIDResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	hostname, err := d.store.GetHostnameForHWID(req.Hwid)
	res := &protos.GetHostnameForHWIDResponse{Hostname: hostname}

	return res, err
}

func (d *directoryLookupServicer) MapHWIDsToHostnames(ctx context.Context, req *protos.MapHWIDToHostnameRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	err = d.store.MapHWIDsToHostnames(req.HwidToHostname)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) GetIMSIForSessionID(
	ctx context.Context, req *protos.GetIMSIForSessionIDRequest,
) (*protos.GetIMSIForSessionIDResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	imsi, err := d.store.GetIMSIForSessionID(req.NetworkID, req.SessionID)
	res := &protos.GetIMSIForSessionIDResponse{Imsi: imsi}

	return res, err
}

func (d *directoryLookupServicer) MapSessionIDsToIMSIs(ctx context.Context, req *protos.MapSessionIDToIMSIRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	err = d.store.MapSessionIDsToIMSIs(req.NetworkID, req.SessionIDToIMSI)

	return &protos.Void{}, err
}

// MapHWIDToDirectoryRecordIDs maps {hwid -> directory record IDs}. If existing
// record IDs already exist for a given hwid, new record IDs will be appended
// to the existing state.
func (d *directoryLookupServicer) MapHWIDToDirectoryRecordIDs(ctx context.Context, req *protos.MapHWIDToDirectoryRecordIDsRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}
	for hwid, newRecords := range req.HwidToRecordIDs {
		existingRecords, err := d.store.GetDirectoryRecordIDsForHWID(req.NetworkID, hwid)
		// if this is the first record reported for a given hwid then we will
		// receive an error for non-existent state.
		if err == merrors.ErrNotFound {
			continue
		}
		if err != nil {
			return &protos.Void{}, errors.Wrap(err, "failed to lookup existing directory record IDs")
		}
		allRecordIds := getUnionOfIDs(existingRecords.Ids, newRecords.Ids)
		req.HwidToRecordIDs[hwid].Ids = allRecordIds
	}
	err = d.store.MapHWIDToDirectoryRecordIDs(req.NetworkID, req.HwidToRecordIDs)

	return &protos.Void{}, err
}

// GetHWIDToDirectoryRecordIDs returns the directory record IDs mapped to a
// given hwid. The returned state is pruned of all stale directory records.
func (d *directoryLookupServicer) GetDirectoryRecordIDsForHWID(ctx context.Context, req *protos.GetDirectoryRecordIDsForHWIDRequest) (*protos.GetDirectoryRecordIDsForHWIDResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}
	directoryRecords, err := d.store.GetDirectoryRecordIDsForHWID(req.NetworkID, req.Hwid)
	if err != nil {
		return nil, err
	}
	prunedIDs, err := pruneStaleDirectoryRecordIDs(req.NetworkID, directoryRecords.Ids)
	if err != nil {
		return &protos.GetDirectoryRecordIDsForHWIDResponse{}, errors.Wrap(err, "prune stale directory records")
	}
	res := &protos.GetDirectoryRecordIDsForHWIDResponse{Ids: &protos.DirectoryRecordIDs{Ids: prunedIDs}}

	return res, nil
}

func pruneStaleDirectoryRecordIDs(networkID string, recordIDs []string) ([]string, error) {
	// since state indexing does not handle deletion of stale state, prune
	// non-existent directory records before returning
	prunedIDs := []string{}
	stateIDs := types.MakeIDs(orc8r.DirectoryRecordType, recordIDs...)
	statesByID, err := state.GetStates(networkID, stateIDs, serdes.State)
	if err != nil {
		return []string{}, err
	}
	for _, stateID := range stateIDs {
		_, ok := statesByID[stateID]
		if !ok {
			continue
		}
		prunedIDs = append(prunedIDs, stateID.DeviceID)
	}
	return prunedIDs, nil
}

func getUnionOfIDs(existingIDs []string, newIDs []string) []string {
	idSet := map[string]bool{}
	for _, id := range existingIDs {
		idSet[id] = true
	}
	for _, id := range newIDs {
		idSet[id] = true
	}
	idKeys := funk.Keys(idSet).([]string)
	sort.Strings(idKeys)
	return idKeys
}
