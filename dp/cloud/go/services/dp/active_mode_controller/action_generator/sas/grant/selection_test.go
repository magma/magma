/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grant_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas/grant"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestProcessGrants(t *testing.T) {
	testData := []struct {
		name     string
		cbsd     *storage.DBCbsd
		grants   []*storage.DBGrant
		expected []grantData
	}{{
		name: "Should do nothing when no grants or available frequencies",
		cbsd: &storage.DBCbsd{
			PreferredBandwidthMHz: db.MakeInt(20),
			MaxIbwMhx:             db.MakeInt(150),
			AvailableFrequencies:  []uint32{0, 0, 0, 0},
		},
		expected: []grantData{},
	}, {
		name: "Should select only one grant in no redundancy mode",
		cbsd: &storage.DBCbsd{
			PreferredBandwidthMHz: db.MakeInt(20),
			MaxIbwMhx:             db.MakeInt(150),
			AvailableFrequencies:  allAvailable,
		},
		grants: nil,
		expected: []grantData{{
			action:    add,
			frequency: 3560e6,
			bandwidth: 20e6,
		}},
	}, {
		name: "Should select keep existing grant in no redundancy mode",
		cbsd: &storage.DBCbsd{
			PreferredBandwidthMHz: db.MakeInt(20),
			MaxIbwMhx:             db.MakeInt(150),
			AvailableFrequencies:  allAvailable,
		},
		grants: []*storage.DBGrant{{
			LowFrequencyHz:  db.MakeInt(3590e6),
			HighFrequencyHz: db.MakeInt(3610e6),
		}},
		expected: []grantData{{
			action:    keep,
			frequency: 3600e6,
			bandwidth: 20e6,
		}},
	}, {
		name: "Should select grants for redundancy",
		cbsd: &storage.DBCbsd{
			GrantRedundancy:       db.MakeBool(true),
			PreferredBandwidthMHz: db.MakeInt(20),
			MaxIbwMhx:             db.MakeInt(150),
			AvailableFrequencies:  allAvailable,
		},
		grants: nil,
		expected: []grantData{{
			action:    add,
			frequency: 3560e6,
			bandwidth: 20e6,
		}, {
			action:    add,
			frequency: 3580e6,
			bandwidth: 20e6,
		}},
	}, {
		name: "Should use custom ordering in carrier aggregation mode",
		cbsd: &storage.DBCbsd{
			GrantRedundancy:           db.MakeBool(true),
			CarrierAggregationEnabled: db.MakeBool(true),
			PreferredBandwidthMHz:     db.MakeInt(15),
			MaxIbwMhx:                 db.MakeInt(150),
			AvailableFrequencies:      allAvailable,
		},
		grants: nil,
		expected: []grantData{{
			action:    add,
			frequency: 3555e6,
			bandwidth: 10e6,
		}, {
			action:    add,
			frequency: 3565e6,
			bandwidth: 10e6,
		}},
	}, {
		name: "Should use frequency and bandwidth preferences",
		cbsd: &storage.DBCbsd{
			GrantRedundancy:         db.MakeBool(true),
			PreferredBandwidthMHz:   db.MakeInt(15),
			PreferredFrequenciesMHz: []int64{3570},
			MaxIbwMhx:               db.MakeInt(150),
			AvailableFrequencies:    allAvailable,
		},
		grants: nil,
		expected: []grantData{{
			action:    add,
			frequency: 3570e6,
			bandwidth: 15e6,
		}, {
			action:    add,
			frequency: 3585e6,
			bandwidth: 15e6,
		}},
	}, {
		name: "Should add only one grant if only available in standard redundancy",
		cbsd: &storage.DBCbsd{
			GrantRedundancy:       db.MakeBool(true),
			PreferredBandwidthMHz: db.MakeInt(20),
			MaxIbwMhx:             db.MakeInt(30),
			AvailableFrequencies:  allAvailable,
		},
		grants: nil,
		expected: []grantData{{
			action:    add,
			frequency: 3560e6,
			bandwidth: 20e6,
		}},
	}, {
		name: "Should go to next bandwidth if only one available in carrier aggregation mode",
		cbsd: &storage.DBCbsd{
			GrantRedundancy:           db.MakeBool(true),
			CarrierAggregationEnabled: db.MakeBool(true),
			PreferredBandwidthMHz:     db.MakeInt(20),
			MaxIbwMhx:                 db.MakeInt(150),
			AvailableFrequencies:      []uint32{0, 1 << 10, 1 << 10, 0},
		},
		grants: nil,
		expected: []grantData{{
			action:    add,
			frequency: 3600e6,
			bandwidth: 15e6,
		}},
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			p := grant.Processors[grantData]{
				Keep: &stubGrantProcessor{action: keep},
				Del:  &stubGrantProcessor{action: del},
				Add:  &stubGrantProcessor{action: add},
			}
			actual := grant.ProcessGrants[grantData](tt.cbsd, tt.grants, p, 0)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

var allAvailable = []uint32{
	1<<30 - 1<<1,
	1<<30 - 1<<1,
	1<<29 - 1<<2,
	1<<29 - 1<<2,
}

type stubGrantProcessor struct {
	action action
}

type action uint8

const (
	keep action = iota
	del
	add
)

func (s *stubGrantProcessor) ProcessGrant(frequency int64, bandwidth int64) grantData {
	return grantData{
		action:    s.action,
		frequency: frequency,
		bandwidth: bandwidth,
	}
}

type grantData struct {
	action    action
	frequency int64
	bandwidth int64
}
