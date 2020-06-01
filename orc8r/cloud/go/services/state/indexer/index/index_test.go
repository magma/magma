/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package index

import (
	"testing"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	assert "github.com/stretchr/testify/require"
)

const (
	nid0 = "some_networkid_0"

	did0 = "some_deviceid_0"
	did1 = "some_deviceid_1"
	did2 = "some_deviceid_2"
	did3 = "some_deviceid_3"

	type0 = "some_type_0"

	iid0 = "some_indexerid_0"
	iid1 = "some_indexerid_1"
	iid2 = "some_indexerid_2"
)

var (
	someErr = errors.New("some_error")
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestIndexImpl_HappyPath(t *testing.T) {
	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	id00 := state_types.ID{Type: orc8r.DirectoryRecordType, DeviceID: did0}
	id01 := state_types.ID{Type: orc8r.DirectoryRecordType, DeviceID: did1}

	id12 := state_types.ID{Type: orc8r.AccessGatewayRecordType, DeviceID: did2}
	id13 := state_types.ID{Type: orc8r.AccessGatewayRecordType, DeviceID: did3}

	reported0 := &directoryd.DirectoryRecord{LocationHistory: []string{"rec0_location_history"}}
	reported1 := &directoryd.DirectoryRecord{LocationHistory: []string{"rec1_location_history"}}
	st00 := state_types.State{ReportedState: reported0, Type: orc8r.DirectoryRecordType}
	st01 := state_types.State{ReportedState: reported1, Type: orc8r.DirectoryRecordType}

	reported2 := &models.GatewayDevice{HardwareID: "42"}
	reported3 := &models.GatewayDevice{HardwareID: "43"}
	st12 := state_types.State{ReportedState: reported2, Type: orc8r.AccessGatewayRecordType}
	st13 := state_types.State{ReportedState: reported3, Type: orc8r.AccessGatewayRecordType}

	index0 := state_types.StatesByID{
		id00: st00,
		id01: st01,
	}
	index1 := state_types.StatesByID{
		id12: st12,
	}
	index2 := state_types.StatesByID{
		id13: st13,
	}

	in := state_types.StatesByID{
		id00: st00,
		id01: st01,
		id12: st12,
		id13: st13,
	}

	idx0 := getIndexerWithVersion(iid0, []indexer.Subscription{{Type: orc8r.DirectoryRecordType, KeyMatcher: indexer.MatchAll}})
	idx1 := getIndexerWithVersion(iid1, []indexer.Subscription{{Type: orc8r.AccessGatewayRecordType, KeyMatcher: indexer.NewMatchExact(did2)}})
	idx2 := getIndexerWithVersion(iid2, []indexer.Subscription{{Type: orc8r.AccessGatewayRecordType, KeyMatcher: indexer.NewMatchExact(did3)}})
	idx0.On("Index", nid0, index0).Return(indexer.StateErrors{id00: someErr}, nil).Once()
	idx1.On("Index", nid0, index1).Return(nil, nil).Once()
	idx2.On("Index", nid0, index2).Return(nil, someErr).Times(maxRetry)

	// All indexing occurs as expected
	indexer.DeregisterAllForTest(t)
	assert.NoError(t, indexer.RegisterAll(idx0, idx1, idx2))
	actual := indexImpl(nid0, in)
	assert.Len(t, actual, 1)
	e := actual[0].Error()
	assert.Contains(t, e, iid2)
	assert.Contains(t, e, ErrIndex)
	assert.Contains(t, e, someErr.Error())
	idx0.AssertExpectations(t)
	idx1.AssertExpectations(t)
	idx2.AssertExpectations(t)
}

func TestIndexImpl_AllStatesFiltered(t *testing.T) {
	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	id00 := state_types.ID{Type: orc8r.DirectoryRecordType, DeviceID: did0}
	id01 := state_types.ID{Type: orc8r.DirectoryRecordType, DeviceID: did1}

	id12 := state_types.ID{Type: orc8r.AccessGatewayRecordType, DeviceID: did2}
	id13 := state_types.ID{Type: orc8r.AccessGatewayRecordType, DeviceID: did3}

	reported0 := &directoryd.DirectoryRecord{LocationHistory: []string{"rec0_location_history"}}
	reported1 := &directoryd.DirectoryRecord{LocationHistory: []string{"rec1_location_history"}}
	st00 := state_types.State{ReportedState: reported0, Type: orc8r.DirectoryRecordType}
	st01 := state_types.State{ReportedState: reported1, Type: orc8r.DirectoryRecordType}

	reported2 := &models.GatewayDevice{HardwareID: "42"}
	reported3 := &models.GatewayDevice{HardwareID: "43"}
	st12 := state_types.State{ReportedState: reported2, Type: orc8r.AccessGatewayRecordType}
	st13 := state_types.State{ReportedState: reported3, Type: orc8r.AccessGatewayRecordType}

	in := state_types.StatesByID{
		id00: st00,
		id01: st01,
		id12: st12,
		id13: st13,
	}

	// All states get filtered -> no err
	idx0 := getIndexer(iid0, []indexer.Subscription{{Type: orc8r.DirectoryRecordType, KeyMatcher: indexer.NewMatchPrefix("0xdeadbeef")}})
	idx1 := getIndexer(iid1, []indexer.Subscription{{Type: type0, KeyMatcher: indexer.MatchAll}})
	idx0.On("Index", nid0, mock.Anything).Return(nil, nil)
	idx1.On("Index", nid0, mock.Anything).Return(nil, nil)
	indexer.DeregisterAllForTest(t)
	assert.NoError(t, indexer.RegisterAll(idx0, idx1))
	actual := indexImpl(nid0, in)
	assert.Empty(t, actual)
	idx0.AssertNotCalled(t, "Index", mock.Anything, mock.Anything)
	idx1.AssertNotCalled(t, "Index", mock.Anything, mock.Anything)
}

func getIndexerWithVersion(id string, subs []indexer.Subscription) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetSubscriptions").Return(subs)
	idx.On("GetVersion").Return(indexer.Version(42))
	return idx
}

func getIndexer(id string, subs []indexer.Subscription) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetSubscriptions").Return(subs)
	return idx
}
