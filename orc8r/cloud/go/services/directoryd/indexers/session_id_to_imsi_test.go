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
	directorydTest "magma/orc8r/cloud/go/services/directoryd/test_init"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"

	"github.com/stretchr/testify/assert"
)

const (
	hwid0 = "some_hardware_id_0"
	imsi0 = "some_imsi_0"
	imsi1 = "some_imsi_1"
	nid0  = "some_network_id_0"
	sid0  = "some_session_id_0"
	sid1  = "some_session_id_1"
)

func TestIndexerSessionID(t *testing.T) {
	indexer := indexers.NewSessionIDToIMSI()
	directorydTest.StartTestService(t)

	record := &directoryd.DirectoryRecord{
		LocationHistory: []string{hwid0}, // imsi0->hwid0
		Identifiers: map[string]interface{}{
			directoryd.RecordKeySessionID: sid0, // imsi0->sid0
		},
	}

	st := state.State{
		ReporterID:         imsi0,
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
	assert.True(t, indexer.GetSubscriptions()[0].Match(st))

	// Index the imsi0->sid0 state, result is sid0->imsi0 reverse mapping
	errs, err := indexer.Index(nid0, hwid0, []state.State{st})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err := directoryd.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, imsi)

	// Update sid -- index imsi0->sid1, result is sid1->imsi0 reverse mapping
	// Note that we specifically don't test for the presence of {sid0 -> ?}, as we allow stale derived state to persist.
	st.ReportedState.(*directoryd.DirectoryRecord).Identifiers[directoryd.RecordKeySessionID] = sid1
	errs, err = indexer.Index(nid0, hwid0, []state.State{st})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, imsi)

	// Update imsi -- index imsi1->sid1, result is sid1->imsi1 reverse mapping
	// Note that we specifically don't test for the presence of {sid0 -> ?}, as we allow stale derived state to persist.
	st.ReporterID = imsi1
	errs, err = indexer.Index(nid0, hwid0, []state.State{st})
	assert.NoError(t, err)
	assert.Empty(t, errs)
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, imsi)

	// Errs contains an err when e.g. reported state is wrong type -- and sid1->imsi1 still intact
	st.ReportedState = 42
	errs, err = indexer.Index(nid0, hwid0, []state.State{st})
	tk := storage.TypeAndKey{Type: orc8r.DirectoryRecordType, Key: st.ReporterID}
	assert.NoError(t, err)
	assert.Error(t, errs[tk])
	imsi, err = directoryd.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, imsi)
}
