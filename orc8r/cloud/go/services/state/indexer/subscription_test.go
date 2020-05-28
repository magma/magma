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

	"github.com/stretchr/testify/assert"
)

const (
	imsi0 = "some_imsi_0"
	imsi1 = "some_imsi_1"

	did0 = "some_deviceid_0"
	did1 = "some_deviceid_1"
	did2 = "some_deviceid_2"

	type0 = "some_type_0"
	type1 = "some_type_1"
	type2 = "some_type_2"
)

var (
	emptyMatch []string
)

// Test filter IDs and filter states
func TestFilter(t *testing.T) {
	id0 := state_types.ID{Type: type0, DeviceID: did0}
	id1 := state_types.ID{Type: type1, DeviceID: did1}
	id2 := state_types.ID{Type: type2, DeviceID: did2}

	st0 := state_types.State{ReportedState: 42, Type: type0}
	st1 := state_types.State{ReportedState: 42, Type: type1}
	st2 := state_types.State{ReportedState: 42, Type: type2}

	type args struct {
		subs   []indexer.Subscription
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
				subs:   []indexer.Subscription{{Type: type0, KeyMatcher: indexer.MatchAll}},
				states: state_types.StatesByID{id0: st0},
			},
			want: state_types.StatesByID{id0: st0},
		},
		{
			name: "one state zero sub",
			args: args{
				subs:   nil,
				states: state_types.StatesByID{id0: st0},
			},
			want: state_types.StatesByID{},
		},
		{
			name: "zero state one sub",
			args: args{
				subs:   []indexer.Subscription{{Type: type0, KeyMatcher: indexer.MatchAll}},
				states: state_types.StatesByID{},
			},
			want: state_types.StatesByID{},
		},
		{
			name: "wrong type",
			args: args{
				subs:   []indexer.Subscription{{Type: type1, KeyMatcher: indexer.MatchAll}},
				states: state_types.StatesByID{id0: st0},
			},
			want: state_types.StatesByID{},
		},
		{
			name: "wrong device ID",
			args: args{
				subs:   []indexer.Subscription{{Type: type0, KeyMatcher: indexer.NewMatchExact("0xdeadbeef")}},
				states: state_types.StatesByID{id0: st0},
			},
			want: state_types.StatesByID{},
		},
		{
			name: "multi state multi sub",
			args: args{
				subs: []indexer.Subscription{
					{Type: type0, KeyMatcher: indexer.MatchAll},
					{Type: type1, KeyMatcher: indexer.NewMatchPrefix(id1.DeviceID[0:3])},
					{Type: type2, KeyMatcher: indexer.NewMatchExact(id2.DeviceID[0:3])},
				},
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
			gotStates := indexer.FilterStates(tt.args.subs, tt.args.states)
			assert.Equal(t, tt.want, gotStates)
			gotIDs := indexer.FilterIDs(tt.args.subs, statesToIDs(tt.args.states))
			assert.Equal(t, statesToIDs(tt.want), gotIDs)
		})
	}
}

func TestSubscription_Match(t *testing.T) {
	sub := indexer.Subscription{Type: type0}

	id00 := state_types.ID{DeviceID: imsi0, Type: type0}
	id10 := state_types.ID{DeviceID: imsi1, Type: type0}
	id01 := state_types.ID{DeviceID: imsi0, Type: type1}
	id11 := state_types.ID{DeviceID: imsi1, Type: type1}

	// Match all
	sub.KeyMatcher = indexer.MatchAll
	assert.True(t, sub.Match(id00))
	assert.True(t, sub.Match(id10))
	assert.False(t, sub.Match(id01))
	assert.False(t, sub.Match(id11))

	// Match exact
	sub.KeyMatcher = indexer.NewMatchExact(imsi0)
	assert.True(t, sub.Match(id00))
	assert.False(t, sub.Match(id10))
	assert.False(t, sub.Match(id01))
	assert.False(t, sub.Match(id11))
}

func TestMatchAll_Match(t *testing.T) {
	m := indexer.MatchAll
	yesMatch := []string{"", "f", "fo", "foo", "foobar", "foobarbaz", "foobar !@#$%^&*()_+"}
	testMatchImpl(t, m, yesMatch, emptyMatch)
}

func TestMatchExact_Match(t *testing.T) {
	m := indexer.NewMatchExact("foobar")
	yesMatch := []string{"foobar"}
	noMatch := []string{"", "f", "fo", "foo", "foobarbaz", "foobar !@#$%^&*()_+"}
	testMatchImpl(t, m, yesMatch, noMatch)
}

func TestMatchPrefix_Match(t *testing.T) {
	m := indexer.NewMatchPrefix("foo")
	yesMatch := []string{"foo", "foobar", "foobarbaz", "foobar !@#$%^&*()_+"}
	noMatch := []string{"", "f", "fo"}
	testMatchImpl(t, m, yesMatch, noMatch)
}

func testMatchImpl(t *testing.T, m indexer.KeyMatcher, yesMatch, noMatch []string) {
	for _, str := range yesMatch {
		assert.True(t, m.Match(str))
	}
	for _, str := range noMatch {
		assert.False(t, m.Match(str))
	}
}

func statesToIDs(states state_types.StatesByID) []state_types.ID {
	var ret []state_types.ID
	for id := range states {
		ret = append(ret, id)
	}

	// Sort for deterministic output
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Type+ret[i].DeviceID < ret[j].Type+ret[j].DeviceID
	})
	return ret
}
