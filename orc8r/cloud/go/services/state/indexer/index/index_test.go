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
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"

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
	type1 = "some_type_1"
	type2 = "some_type_2"

	iid0 = "some_indexerid_0"
	iid1 = "some_indexerid_1"
	iid2 = "some_indexerid_2"
)

var (
	someErr = errors.New("some_error")
)

func TestFilterStates(t *testing.T) {
	initTest()

	id0 := state.ID{Type: type0, DeviceID: did0}
	id1 := state.ID{Type: type1, DeviceID: did1}
	id2 := state.ID{Type: type2, DeviceID: did2}

	st0 := state.State{ReportedState: 42, Type: type0}
	st1 := state.State{ReportedState: 42, Type: type1}
	st2 := state.State{ReportedState: 42, Type: type2}

	type args struct {
		idx    indexer.Indexer
		states state.StatesByID
	}
	tests := []struct {
		name string
		args args
		want state.StatesByID
	}{
		{
			name: "one state one sub",
			args: args{
				idx:    getIndexer("", []indexer.Subscription{{Type: type0, KeyMatcher: indexer.MatchAll}}),
				states: state.StatesByID{id0: st0},
			},
			want: state.StatesByID{id0: st0},
		},
		{
			name: "one state zero sub",
			args: args{
				idx:    getIndexer("", nil),
				states: state.StatesByID{id0: st0},
			},
			want: state.StatesByID{},
		},
		{
			name: "zero state one sub",
			args: args{
				idx:    getIndexer("", []indexer.Subscription{{Type: type0, KeyMatcher: indexer.MatchAll}}),
				states: state.StatesByID{},
			},
			want: state.StatesByID{},
		},
		{
			name: "wrong type",
			args: args{
				idx:    getIndexer("", []indexer.Subscription{{Type: type1, KeyMatcher: indexer.MatchAll}}),
				states: state.StatesByID{id0: st0},
			},
			want: state.StatesByID{},
		},
		{
			name: "wrong device ID",
			args: args{
				idx:    getIndexer("", []indexer.Subscription{{Type: type0, KeyMatcher: indexer.NewMatchExact("0xdeadbeef")}}),
				states: state.StatesByID{id0: st0},
			},
			want: state.StatesByID{},
		},
		{
			name: "multi state multi sub",
			args: args{
				idx: getIndexer("", []indexer.Subscription{
					{Type: type0, KeyMatcher: indexer.MatchAll},
					{Type: type1, KeyMatcher: indexer.NewMatchPrefix(id1.DeviceID[0:3])},
					{Type: type2, KeyMatcher: indexer.NewMatchExact(id2.DeviceID[0:3])},
				}),
				states: state.StatesByID{
					id0: st0,
					id1: st1,
					id2: st2,
				},
			},
			want: state.StatesByID{
				id0: st0,
				id1: st1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterStates(tt.args.idx, tt.args.states)
			assert.Equal(t, tt.want, got, tt.want)
		})
	}
}

func TestIndexImpl_HappyPath(t *testing.T) {
	initTest()
	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	id00 := state.ID{Type: orc8r.DirectoryRecordType, DeviceID: did0}
	id01 := state.ID{Type: orc8r.DirectoryRecordType, DeviceID: did1}

	id12 := state.ID{Type: orc8r.AccessGatewayRecordType, DeviceID: did2}
	id13 := state.ID{Type: orc8r.AccessGatewayRecordType, DeviceID: did3}

	reported0 := &directoryd.DirectoryRecord{LocationHistory: []string{"rec0_location_history"}}
	reported1 := &directoryd.DirectoryRecord{LocationHistory: []string{"rec1_location_history"}}
	st00 := state.State{ReportedState: reported0, Type: orc8r.DirectoryRecordType}
	st01 := state.State{ReportedState: reported1, Type: orc8r.DirectoryRecordType}

	reported2 := &models.GatewayDevice{HardwareID: "42"}
	reported3 := &models.GatewayDevice{HardwareID: "43"}
	st12 := state.State{ReportedState: reported2, Type: orc8r.AccessGatewayRecordType}
	st13 := state.State{ReportedState: reported3, Type: orc8r.AccessGatewayRecordType}

	index0 := state.StatesByID{
		id00: st00,
		id01: st01,
	}
	index1 := state.StatesByID{
		id12: st12,
	}
	index2 := state.StatesByID{
		id13: st13,
	}

	in := state.StatesByID{
		id00: st00,
		id01: st01,
		id12: st12,
		id13: st13,
	}

	idx0 := getIndexerBasic(iid0, []indexer.Subscription{{Type: orc8r.DirectoryRecordType, KeyMatcher: indexer.MatchAll}})
	idx1 := getIndexerBasic(iid1, []indexer.Subscription{{Type: orc8r.AccessGatewayRecordType, KeyMatcher: indexer.NewMatchExact(did2)}})
	idx2 := getIndexerBasic(iid2, []indexer.Subscription{{Type: orc8r.AccessGatewayRecordType, KeyMatcher: indexer.NewMatchExact(did3)}})
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
	initTest()
	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	id00 := state.ID{Type: orc8r.DirectoryRecordType, DeviceID: did0}
	id01 := state.ID{Type: orc8r.DirectoryRecordType, DeviceID: did1}

	id12 := state.ID{Type: orc8r.AccessGatewayRecordType, DeviceID: did2}
	id13 := state.ID{Type: orc8r.AccessGatewayRecordType, DeviceID: did3}

	reported0 := &directoryd.DirectoryRecord{LocationHistory: []string{"rec0_location_history"}}
	reported1 := &directoryd.DirectoryRecord{LocationHistory: []string{"rec1_location_history"}}
	st00 := state.State{ReportedState: reported0, Type: orc8r.DirectoryRecordType}
	st01 := state.State{ReportedState: reported1, Type: orc8r.DirectoryRecordType}

	reported2 := &models.GatewayDevice{HardwareID: "42"}
	reported3 := &models.GatewayDevice{HardwareID: "43"}
	st12 := state.State{ReportedState: reported2, Type: orc8r.AccessGatewayRecordType}
	st13 := state.State{ReportedState: reported3, Type: orc8r.AccessGatewayRecordType}

	in := state.StatesByID{
		id00: st00,
		id01: st01,
		id12: st12,
		id13: st13,
	}

	// All states get filtered -> no err
	idx0 := getIndexerBasic(iid0, []indexer.Subscription{{Type: orc8r.DirectoryRecordType, KeyMatcher: indexer.NewMatchPrefix("0xdeadbeef")}})
	idx1 := getIndexerBasic(iid1, []indexer.Subscription{{Type: type0, KeyMatcher: indexer.MatchAll}})
	idx0.On("Index", nid0, mock.Anything).Return(nil, nil)
	idx1.On("Index", nid0, mock.Anything).Return(nil, nil)
	indexer.DeregisterAllForTest(t)
	assert.NoError(t, indexer.RegisterAll(idx0, idx1))
	actual := indexImpl(nid0, in)
	assert.Empty(t, actual)
	idx0.AssertNotCalled(t, "Index", mock.Anything, mock.Anything)
	idx1.AssertNotCalled(t, "Index", mock.Anything, mock.Anything)
}

func initTest() {
	// Uncomment below to view reindex queue logs during test
	//_ = flag.Set("alsologtostderr", "true")
}

func getIndexerBasic(id string, subs []indexer.Subscription) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetSubscriptions").Return(subs)
	return idx
}

func getIndexer(id string, subs []indexer.Subscription) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetSubscriptions").Return(subs)
	idx.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Once()
	return idx
}
