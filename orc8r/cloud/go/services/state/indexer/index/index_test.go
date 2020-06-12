/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// index_test.go tests indexing with local indexers.

package index

import (
	"testing"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/pkg/errors"
	assert "github.com/stretchr/testify/require"
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestIndexImpl_HappyPath(t *testing.T) {
	const (
		nid0 = "some_networkid_0"

		iid0 = "some_indexerid_0"
		iid1 = "some_indexerid_1"
		iid2 = "some_indexerid_2"
		iid3 = "some_indexerid_3"
	)
	var someErr = errors.New("some_error")

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	id0 := state_types.ID{Type: orc8r.GatewayStateType}
	id1 := state_types.ID{Type: orc8r.AccessGatewayRecordType}
	reported0 := &models.GatewayStatus{Meta: map[string]string{"foo": "bar"}}
	reported1 := &models.GatewayDevice{HardwareID: "42"}
	st0 := state_types.State{ReportedState: reported0, Type: orc8r.GatewayStateType}
	st1 := state_types.State{ReportedState: reported1, Type: orc8r.AccessGatewayRecordType}

	indexTwo := state_types.StatesByID{id0: st0, id1: st1}
	indexOne := state_types.StatesByID{id1: st1}
	in := state_types.StatesByID{id0: st0, id1: st1}

	idx0 := getIndexer(iid0, []string{orc8r.GatewayStateType, orc8r.AccessGatewayRecordType})
	idx1 := getIndexer(iid1, []string{orc8r.AccessGatewayRecordType})
	idx2 := getIndexer(iid2, []string{"type_with_no_reported_states"})
	idx3 := getIndexer(iid3, []string{})

	idx0.On("Index", nid0, indexTwo).Return(state_types.StateErrors{id0: someErr}, nil).Once()
	idx1.On("Index", nid0, indexOne).Return(nil, someErr).Times(maxRetry)
	idx0.On("GetVersion").Return(indexer.Version(42))
	idx1.On("GetVersion").Return(indexer.Version(42))

	// All indexing occurs as expected
	indexer.DeregisterAllForTest(t)
	assert.NoError(t, indexer.RegisterIndexers(idx0, idx1, idx2, idx3))
	actual := indexImpl(nid0, in)
	assert.Len(t, actual, 1) // from idx1's overarching err return
	e := actual[0].Error()
	assert.Contains(t, e, iid1)
	assert.Contains(t, e, ErrIndex)
	assert.Contains(t, e, someErr.Error())
	idx0.AssertExpectations(t)
	idx1.AssertExpectations(t)
	idx2.AssertExpectations(t)
	idx3.AssertExpectations(t)
}

func getIndexer(id string, types []string) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetTypes").Return(types)
	return idx
}
