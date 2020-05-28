/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexers_test

import (
	"strings"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/directoryd/indexers"
	directoryd_test "magma/orc8r/cloud/go/services/directoryd/test_init"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/stretchr/testify/assert"
)

const (
	imsi0 = "some_imsi_0"
	imsi1 = "some_imsi_1"
	nid0  = "some_network_id_0"
	sid0  = "some_session_id_0"
	sid1  = "some_session_id_1"
)

func TestIndexerSessionID(t *testing.T) {
	indexer := indexers.NewSessionIDToIMSI()
	directoryd_test.StartTestService(t)

	record := &directoryd.DirectoryRecord{
		Identifiers: map[string]interface{}{
			directoryd.RecordKeySessionID: sid0, // imsi0->sid0
		},
	}

	id := state_types.ID{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi0,
	}
	st := state_types.State{
		Type:               orc8r.DirectoryRecordType,
		ReportedState:      record,
		Version:            44,
		TimeMs:             42,
		CertExpirationTime: 43,
	}

	// Indexer ID has prefix directoryd
	assert.True(t, strings.HasPrefix(indexer.GetID(), strings.ToLower(directoryd.ServiceName)))

	// Indexer subscription matches directory records
	assert.True(t, len(indexer.GetSubscriptions()) > 0)
	assert.True(t, indexer.GetSubscriptions()[0].Match(id))

	// Index the imsi0->sid0 state, result is sid0->imsi0 reverse mapping
	errs, err := indexer.Index(nid0, state_types.StatesByID{id: st})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err := directoryd.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, imsi)

	// Update sid -- index imsi0->sid1, result is sid1->imsi0 reverse mapping
	// Note that we specifically don't test for the presence of {sid0 -> ?}, as we allow stale derived state to persist.
	st.ReportedState.(*directoryd.DirectoryRecord).Identifiers[directoryd.RecordKeySessionID] = sid1
	errs, err = indexer.Index(nid0, state_types.StatesByID{id: st})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, imsi)

	// Update imsi -- index imsi1->sid1, result is sid1->imsi1 reverse mapping
	// Note that we specifically don't test for the presence of {sid0 -> ?}, as we allow stale derived state to persist.
	id.DeviceID = imsi1
	errs, err = indexer.Index(nid0, state_types.StatesByID{id: st})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, imsi)

	// Errs contains an err when e.g. reported state is wrong type -- and sid1->imsi1 still intact
	st.ReportedState = 42
	errs, err = indexer.Index(nid0, state_types.StatesByID{id: st})
	assert.NoError(t, err)
	assert.Error(t, errs[id])
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, imsi)
}
