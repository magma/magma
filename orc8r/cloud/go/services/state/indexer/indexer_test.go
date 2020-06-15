/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package indexer_test

import (
	"sort"
	"testing"

	"magma/orc8r/cloud/go/services/state/indexer"
	state_types "magma/orc8r/cloud/go/services/state/types"

	assert "github.com/stretchr/testify/require"
)

// Test filter IDs and filter states
func TestFilter(t *testing.T) {
	const (
		did0 = "some_deviceid_0"
		did1 = "some_deviceid_1"
		did2 = "some_deviceid_2"

		type0 = "some_type_0"
		type1 = "some_type_1"
		type2 = "some_type_2"
	)

	id0 := state_types.ID{Type: type0, DeviceID: did0}
	id1 := state_types.ID{Type: type1, DeviceID: did1}
	id2 := state_types.ID{Type: type2, DeviceID: did2}

	st0 := state_types.State{ReportedState: 42, Type: type0}
	st1 := state_types.State{ReportedState: 42, Type: type1}
	st2 := state_types.State{ReportedState: 42, Type: type2}

	type args struct {
		types  []string
		states state_types.StatesByID
	}
	tests := []struct {
		name string
		args args
		want state_types.StatesByID
	}{
		{
			name: "one state one sub",
			args: args{
				types:  []string{type0},
				states: state_types.StatesByID{id0: st0},
			},
			want: state_types.StatesByID{id0: st0},
		},
		{
			name: "one state zero sub",
			args: args{
				types:  nil,
				states: state_types.StatesByID{id0: st0},
			},
			want: state_types.StatesByID{},
		},
		{
			name: "zero state one sub",
			args: args{
				types:  []string{type0},
				states: state_types.StatesByID{},
			},
			want: state_types.StatesByID{},
		},
		{
			name: "wrong type",
			args: args{
				types:  []string{type1},
				states: state_types.StatesByID{id0: st0},
			},
			want: state_types.StatesByID{},
		},
		{
			name: "multi state multi sub",
			args: args{
				types: []string{type0, type1},
				states: state_types.StatesByID{
					id0: st0,
					id1: st1,
					id2: st2,
				},
			},
			want: state_types.StatesByID{
				id0: st0,
				id1: st1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStates := indexer.FilterStates(tt.args.types, tt.args.states)
			assert.Equal(t, tt.want, gotStates)
			gotIDs := indexer.FilterIDs(tt.args.types, statesToIDs(tt.args.states))
			assert.Equal(t, statesToIDs(tt.want), gotIDs)
		})
	}
}

func statesToIDs(states state_types.StatesByID) []state_types.ID {
	var ret []state_types.ID
	for id := range states {
		ret = append(ret, id)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Type+ret[i].DeviceID < ret[j].Type+ret[j].DeviceID }) // make deterministic
	return ret
}
