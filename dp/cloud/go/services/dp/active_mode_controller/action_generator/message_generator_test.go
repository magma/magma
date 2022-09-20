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

package action_generator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator"
	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/action"
	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas"
	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas/frequency"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestGenerateMessages(t *testing.T) {
	const timeout = 100 * time.Second
	now := time.Unix(currentTimestamp, 0)
	data := []struct {
		name     string
		cbsd     *storage.DetailedCbsd
		expected []action_generator.Action
	}{{
		name: "Should do nothing for unregistered non active cbsd",
		cbsd: NewCbsdBuilder().
			Inactive().
			WithState(unregistered).
			Build(),
	}, {
		name: "Should do nothing when inactive cbsd has no grants",
		cbsd: NewCbsdBuilder().
			Inactive().
			Build(),
	}, {
		name: "Should generate deregistration request for non active registered cbsd if desired",
		cbsd: NewCbsdBuilder().
			Inactive().
			WithDesiredState(unregistered).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Deregistration, &sas.DeregistrationRequest{
				CbsdId: cbsdId,
			}),
		},
	}, {
		name: "Should generate registration request for active non registered cbsd",
		cbsd: NewCbsdBuilder().
			WithState(unregistered).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Registration, &sas.RegistrationRequest{
				UserId:           "some_user_id",
				FccId:            "some_fcc_id",
				CbsdSerialNumber: "some_serial_number",
			}),
		},
	}, {
		name: "Should generate spectrum inquiry request when there are no available channels",
		cbsd: NewCbsdBuilder().
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.SpectrumInquiry, &sas.SpectrumInquiryRequest{
				CbsdId: cbsdId,
				InquiredSpectrum: []*sas.FrequencyRange{{
					LowFrequency:  frequency.LowestHz,
					HighFrequency: frequency.HighestHz,
				}},
			}),
		},
	}, {
		name: "Should set available frequencies when they are nil but there are channels",
		cbsd: NewCbsdBuilder().
			WithChannel(storage.Channel{
				LowFrequencyHz:  3590e6,
				HighFrequencyHz: 3610e6,
				MaxEirp:         37,
			}).
			Build(),
		expected: []action_generator.Action{
			&action.Update{
				Data: &storage.DBCbsd{
					Id: db.MakeInt(dbId),
					AvailableFrequencies: []uint32{
						1<<9 | 1<<10 | 1<<11,
						1<<9 | 1<<10 | 1<<11,
						1 << 10,
						1 << 10,
					},
				},
				Mask: db.NewIncludeMask("available_frequencies"),
			},
		},
	}, {
		name: "Should generate spectrum inquiry request when no suitable available frequencies",
		cbsd: NewCbsdBuilder().
			WithChannel(someChannel).
			WithAvailableFrequencies([]uint32{0, 0, 0, 0}).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.SpectrumInquiry, &sas.SpectrumInquiryRequest{
				CbsdId: cbsdId,
				InquiredSpectrum: []*sas.FrequencyRange{{
					LowFrequency:  frequency.LowestHz,
					HighFrequency: frequency.HighestHz,
				}},
			}),
		},
	}, {
		name: "Should generate grant request when there are available frequencies and channels",
		cbsd: NewCbsdBuilder().
			WithChannel(someChannel).
			WithAvailableFrequencies([]uint32{0, 1 << 15, 0, 0}).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Grant, &sas.GrantRequest{
				CbsdId: cbsdId,
				OperationParam: &sas.OperationParam{
					MaxEirp: 35,
					OperationFrequencyRange: &sas.FrequencyRange{
						LowFrequency:  3620e6,
						HighFrequency: 3630e6,
					},
				},
			}),
		},
	}, {
		name: "Should request two grants in carrier aggregation mode",
		cbsd: NewCbsdBuilder().
			WithChannel(someChannel).
			WithAvailableFrequencies([]uint32{0, 0, 0, 1<<10 | 1<<20}).
			WithCarrierAggregation().
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Grant, &sas.GrantRequest{
				CbsdId: cbsdId,
				OperationParam: &sas.OperationParam{
					MaxEirp: 31,
					OperationFrequencyRange: &sas.FrequencyRange{
						LowFrequency:  3590e6,
						HighFrequency: 3610e6,
					},
				},
			}),
			makeRequest(sas.Grant, &sas.GrantRequest{
				CbsdId: cbsdId,
				OperationParam: &sas.OperationParam{
					MaxEirp: 31,
					OperationFrequencyRange: &sas.FrequencyRange{
						LowFrequency:  3640e6,
						HighFrequency: 3660e6,
					},
				},
			}),
		},
	}, {
		name: "Should send heartbeat message for grant in granted state",
		cbsd: NewCbsdBuilder().
			WithChannel(someChannel).
			WithAvailableFrequencies(noAvailableFrequencies).
			WithGrant(granted, someGrant).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Heartbeat, &sas.HeartbeatRequest{
				CbsdId:         cbsdId,
				GrantId:        grantId,
				OperationState: "GRANTED",
			}),
		},
	}, {
		name: "Should not send anything if heartbeat is not needed yet",
		cbsd: NewCbsdBuilder().
			WithChannel(someChannel).
			WithAvailableFrequencies(noAvailableFrequencies).
			WithGrant(authorized, &storage.DBGrant{
				GrantId:                  db.MakeString(grantId),
				HeartbeatIntervalSec:     db.MakeInt(int64(timeout/time.Second) + 1),
				LastHeartbeatRequestTime: db.MakeTime(now),
				LowFrequencyHz:           db.MakeInt(3590e6),
				HighFrequencyHz:          db.MakeInt(3610e6),
			}).
			Build(),
		expected: nil,
	}, {
		name: "Should send heartbeat request if necessary",
		cbsd: NewCbsdBuilder().
			WithChannel(someChannel).
			WithAvailableFrequencies(noAvailableFrequencies).
			WithGrant(authorized, &storage.DBGrant{
				GrantId:                  db.MakeString(grantId),
				HeartbeatIntervalSec:     db.MakeInt(int64(timeout / time.Second)),
				LastHeartbeatRequestTime: db.MakeTime(now),
				LowFrequencyHz:           db.MakeInt(3590e6),
				HighFrequencyHz:          db.MakeInt(3610e6),
			}).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Heartbeat, &sas.HeartbeatRequest{
				CbsdId:         cbsdId,
				GrantId:        grantId,
				OperationState: "AUTHORIZED",
			}),
		},
	}, {
		name: "Should send relinquish message for unsync grant",
		cbsd: NewCbsdBuilder().
			WithChannel(someChannel).
			WithAvailableFrequencies(noAvailableFrequencies).
			WithGrant(unsync, someGrant).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Relinquishment, &sas.RelinquishmentRequest{
				CbsdId:  cbsdId,
				GrantId: grantId,
			}),
		},
	}, {
		name: "Should send relinquish message when inactive for too long",
		cbsd: NewCbsdBuilder().
			Inactive().
			WithGrant(authorized, someGrant).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Relinquishment, &sas.RelinquishmentRequest{
				CbsdId:  cbsdId,
				GrantId: grantId,
			}),
		},
	}, {
		name: "Should send relinquish message when requested",
		cbsd: NewCbsdBuilder().
			ForRelinquish().
			WithGrant(authorized, someGrant).
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Relinquishment, &sas.RelinquishmentRequest{
				CbsdId:  cbsdId,
				GrantId: grantId,
			}),
		},
	}, {
		name: "Should deregister deleted cbsd",
		cbsd: NewCbsdBuilder().
			Deleted().
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Deregistration, &sas.DeregistrationRequest{
				CbsdId: cbsdId,
			}),
		},
	}, {
		name: "Should delete unregistered cbsd marked as deleted",
		cbsd: NewCbsdBuilder().
			WithState(unregistered).
			Deleted().
			Build(),
		expected: []action_generator.Action{
			&action.Delete{Id: dbId},
		},
	}, {
		name: "Should deregister updated cbsd",
		cbsd: NewCbsdBuilder().
			ForDeregistration().
			Build(),
		expected: []action_generator.Action{
			makeRequest(sas.Deregistration, &sas.DeregistrationRequest{
				CbsdId: cbsdId,
			}),
		},
	}, {
		name: "Should acknowledge update of unregistered cbsd marked as updated",
		cbsd: NewCbsdBuilder().
			WithState(unregistered).
			ForDeregistration().
			Build(),
		expected: []action_generator.Action{
			&action.Update{
				Data: &storage.DBCbsd{
					Id:               db.MakeInt(dbId),
					ShouldDeregister: db.MakeBool(false),
				},
				Mask: db.NewIncludeMask("should_deregister"),
			},
		},
	}, {
		name: "Should acknowledge relinquish when there are no grants",
		cbsd: NewCbsdBuilder().
			ForRelinquish().
			Build(),
		expected: []action_generator.Action{
			&action.Update{
				Data: &storage.DBCbsd{
					Id:               db.MakeInt(dbId),
					ShouldRelinquish: db.MakeBool(false),
				},
				Mask: db.NewIncludeMask("should_relinquish"),
			},
		},
	}}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			g := &action_generator.ActionGenerator{
				HeartbeatTimeout:  timeout,
				InactivityTimeout: timeout,
				Rng:               &stubRNG{},
			}

			cbsds := []*storage.DetailedCbsd{tt.cbsd}
			actual := g.GenerateActions(cbsds, now)

			require.Len(t, actual, len(tt.expected))
			for i := range tt.expected {
				assert.Equal(t, tt.expected[i], actual[i])
			}
		})
	}
}

func makeRequest(requestType string, payload any) action_generator.Action {
	req := &storage.MutableRequest{
		Request: &storage.DBRequest{
			CbsdId:  db.MakeInt(dbId),
			Payload: payload,
		},
		RequestType: &storage.DBRequestType{
			Name: db.MakeString(requestType),
		},
	}
	return &action.Request{Data: req}
}

type stubRNG struct{}

func (s *stubRNG) Int() int {
	return 0
}

const (
	currentTimestamp = 1000
	dbId             = 123
	cbsdId           = "some_cbsd_id"
	grantId          = "some_grant_id"

	registered   = "registered"
	unregistered = "unregistered"

	granted    = "granted"
	authorized = "authorized"
	unsync     = "unsync"
)

var (
	someChannel = storage.Channel{
		LowFrequencyHz:  3550e6,
		HighFrequencyHz: 3700e6,
		MaxEirp:         37,
	}
	someGrant = &storage.DBGrant{
		GrantId:         db.MakeString(grantId),
		LowFrequencyHz:  db.MakeInt(3590e6),
		HighFrequencyHz: db.MakeInt(3610e6),
	}
	noAvailableFrequencies = []uint32{0, 0, 0, 0}
)

type cbsdBuilder struct {
	cbsd *storage.DetailedCbsd
}

func NewCbsdBuilder() *cbsdBuilder {
	return &cbsdBuilder{
		cbsd: &storage.DetailedCbsd{
			Cbsd: &storage.DBCbsd{
				Id:                    db.MakeInt(dbId),
				CbsdId:                db.MakeString(cbsdId),
				UserId:                db.MakeString("some_user_id"),
				FccId:                 db.MakeString("some_fcc_id"),
				CbsdSerialNumber:      db.MakeString("some_serial_number"),
				LastSeen:              db.MakeTime(time.Unix(currentTimestamp, 0)),
				PreferredBandwidthMHz: db.MakeInt(20),
				MinPower:              db.MakeFloat(0),
				MaxPower:              db.MakeFloat(30),
				AntennaGainDbi:        db.MakeFloat(15),
				NumberOfPorts:         db.MakeInt(1),
				SingleStepEnabled:     db.MakeBool(false),
				CbsdCategory:          db.MakeString("A"),
				MaxIbwMhx:             db.MakeInt(150),
			},
			CbsdState: &storage.DBCbsdState{
				Name: db.MakeString(registered),
			},
			DesiredState: &storage.DBCbsdState{
				Name: db.MakeString(registered),
			},
		},
	}
}

func (c *cbsdBuilder) Build() *storage.DetailedCbsd {
	return c.cbsd
}

func (c *cbsdBuilder) Inactive() *cbsdBuilder {
	c.cbsd.Cbsd.LastSeen = db.MakeTime(time.Unix(0, 0))
	return c
}

func (c *cbsdBuilder) WithState(state string) *cbsdBuilder {
	c.cbsd.CbsdState.Name = db.MakeString(state)
	return c
}

func (c *cbsdBuilder) WithDesiredState(state string) *cbsdBuilder {
	c.cbsd.DesiredState.Name = db.MakeString(state)
	return c
}

func (c *cbsdBuilder) Deleted() *cbsdBuilder {
	c.cbsd.Cbsd.IsDeleted = db.MakeBool(true)
	return c
}

func (c *cbsdBuilder) ForDeregistration() *cbsdBuilder {
	c.cbsd.Cbsd.ShouldDeregister = db.MakeBool(true)
	return c
}

func (c *cbsdBuilder) ForRelinquish() *cbsdBuilder {
	c.cbsd.Cbsd.ShouldRelinquish = db.MakeBool(true)
	return c
}

func (c *cbsdBuilder) WithChannel(channel storage.Channel) *cbsdBuilder {
	c.cbsd.Cbsd.Channels = append(c.cbsd.Cbsd.Channels, channel)
	return c
}

func (c *cbsdBuilder) WithGrant(state string, grant *storage.DBGrant) *cbsdBuilder {
	g := &storage.DetailedGrant{
		Grant: grant,
		GrantState: &storage.DBGrantState{
			Name: db.MakeString(state),
		},
	}
	c.cbsd.Grants = append(c.cbsd.Grants, g)
	return c
}

func (c *cbsdBuilder) WithAvailableFrequencies(frequencies []uint32) *cbsdBuilder {
	c.cbsd.Cbsd.AvailableFrequencies = frequencies
	return c
}

func (c *cbsdBuilder) WithCarrierAggregation() *cbsdBuilder {
	c.cbsd.Cbsd.GrantRedundancy = db.MakeBool(true)
	c.cbsd.Cbsd.CarrierAggregationEnabled = db.MakeBool(true)
	return c
}

func (c *cbsdBuilder) WithName(name string) *cbsdBuilder {
	c.cbsd.Cbsd.CbsdId = db.MakeString(name)
	c.cbsd.Cbsd.CbsdSerialNumber = db.MakeString(name)
	c.cbsd.Cbsd.FccId = db.MakeString(name)
	c.cbsd.Cbsd.UserId = db.MakeString(name)
	return c
}
