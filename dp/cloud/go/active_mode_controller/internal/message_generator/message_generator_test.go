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

package message_generator_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator"
	"magma/dp/cloud/go/active_mode_controller/internal/test_utils/builders"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
)

func TestGenerateMessages(t *testing.T) {
	const timeout = 100 * time.Second
	now := time.Unix(builders.Now, 0)
	data := []struct {
		name             string
		cbsd             *active_mode.Cbsd
		expectedRequests []*requests.RequestPayload
		expectedActions  []any
	}{{
		name: "Should do nothing for unregistered non active cbsd",
		cbsd: builders.NewCbsdBuilder().
			Inactive().
			WithState(active_mode.CbsdState_Unregistered).
			Build(),
	}, {
		name: "Should do nothing when inactive cbsd has no grants",
		cbsd: builders.NewCbsdBuilder().
			Inactive().
			Build(),
	}, {
		name: "Should generate deregistration request for non active registered cbsd if desired",
		cbsd: builders.NewCbsdBuilder().
			Inactive().
			WithDesiredState(active_mode.CbsdState_Unregistered).
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"deregistrationRequest": [
		{
			"cbsdId": "some_cbsd_id"
		}
	]
}`,
		}},
	}, {
		name: "Should generate registration request for active non registered cbsd",
		cbsd: builders.NewCbsdBuilder().
			WithState(active_mode.CbsdState_Unregistered).
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"registrationRequest": [
		{
			"userId": "some_user_id",
			"fccId": "some_fcc_id",
			"cbsdSerialNumber": "some_serial_number"
		}
]
}`,
		}},
	}, {
		name: "Should generate spectrum inquiry request when there are no available channels",
		cbsd: builders.NewCbsdBuilder().
			Build(),
		expectedRequests: getSpectrumInquiryRequest(),
	}, {
		name: "Should set available frequencies when they are nil but there are channels",
		cbsd: builders.NewCbsdBuilder().
			WithChannel(&active_mode.Channel{
				LowFrequencyHz:  3590e6,
				HighFrequencyHz: 3610e6,
			}).
			Build(),
		expectedActions: []any{
			&active_mode.StoreAvailableFrequenciesRequest{
				Id: builders.DbId,
				AvailableFrequencies: []uint32{
					1<<9 | 1<<10 | 1<<11,
					1<<9 | 1<<10 | 1<<11,
					1 << 10,
					1 << 10,
				},
			},
		},
	}, {
		name: "Should generate spectrum inquiry request when no suitable available frequencies",
		cbsd: builders.NewCbsdBuilder().
			WithChannel(builders.SomeChannel).
			WithAvailableFrequencies([]uint32{0, 0, 0, 0}).
			Build(),
		expectedRequests: getSpectrumInquiryRequest(),
	}, {
		name: "Should generate grant request when there are available frequencies and channels",
		cbsd: builders.NewCbsdBuilder().
			WithChannel(builders.SomeChannel).
			WithAvailableFrequencies([]uint32{0, 1 << 15, 0, 0}).
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"grantRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"operationParam": {
				"maxEirp": 35,
				"operationFrequencyRange": {
					"lowFrequency": 3620000000,
					"highFrequency": 3630000000
				}
			}
		}
	]
}`,
		}},
	}, {
		name: "Should request two grants in carrier aggregation mode",
		cbsd: builders.NewCbsdBuilder().
			WithChannel(builders.SomeChannel).
			WithAvailableFrequencies([]uint32{0, 0, 0, 1<<10 | 1<<20}).
			WithCarrierAggregation().
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"grantRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"operationParam": {
				"maxEirp": 31,
				"operationFrequencyRange": {
					"lowFrequency": 3590000000,
					"highFrequency": 3610000000
				}
			}
		},
		{
			"cbsdId": "some_cbsd_id",
			"operationParam": {
				"maxEirp": 31,
				"operationFrequencyRange": {
					"lowFrequency": 3640000000,
					"highFrequency": 3660000000
				}
			}
		}
	]
}`,
		}},
	}, {
		name: "Should send heartbeat message for grant in granted state",
		cbsd: builders.NewCbsdBuilder().
			WithChannel(builders.SomeChannel).
			WithAvailableFrequencies(builders.NoAvailableFrequencies).
			WithGrant(&active_mode.Grant{
				Id:              builders.GrantId,
				State:           active_mode.GrantState_Granted,
				LowFrequencyHz:  3590e6,
				HighFrequencyHz: 3610e6,
			}).
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"heartbeatRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"grantId": "some_grant_id",
			"operationState": "GRANTED"
		}
	]
}`,
		}},
	}, {
		name: "Should send relinquish message for unsync grant",
		cbsd: builders.NewCbsdBuilder().
			WithChannel(builders.SomeChannel).
			WithAvailableFrequencies(builders.NoAvailableFrequencies).
			WithGrant(&active_mode.Grant{
				Id:              builders.GrantId,
				State:           active_mode.GrantState_Unsync,
				LowFrequencyHz:  3590e6,
				HighFrequencyHz: 3610e6,
			}).
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"relinquishmentRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"grantId": "some_grant_id"
		}
	]
}`,
		}},
	}, {
		name: "Should send relinquish message when inactive for too long",
		cbsd: builders.NewCbsdBuilder().
			Inactive().
			WithGrant(&active_mode.Grant{
				Id:              builders.GrantId,
				State:           active_mode.GrantState_Authorized,
				LowFrequencyHz:  3590e6,
				HighFrequencyHz: 3610e6,
			}).
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"relinquishmentRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"grantId": "some_grant_id"
		}
	]
}`,
		}},
	}, {
		name: "Should deregister deleted cbsd",
		cbsd: builders.NewCbsdBuilder().
			Deleted().
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"deregistrationRequest": [
		{
			"cbsdId": "some_cbsd_id"
		}
	]
}`,
		}},
	}, {
		name: "Should delete unregistered cbsd marked as deleted",
		cbsd: builders.NewCbsdBuilder().
			WithState(active_mode.CbsdState_Unregistered).
			Deleted().
			Build(),
		expectedActions: []any{
			&active_mode.DeleteCbsdRequest{Id: 123},
		},
	}, {
		name: "Should deregister updated cbsd",
		cbsd: builders.NewCbsdBuilder().
			ForDeregistration().
			Build(),
		expectedRequests: []*requests.RequestPayload{{
			Payload: `{
	"deregistrationRequest": [
		{
			"cbsdId": "some_cbsd_id"
		}
	]
}`,
		}},
	}, {
		name: "Should acknowledge update of unregistered cbsd marked as updated",
		cbsd: builders.NewCbsdBuilder().
			WithState(active_mode.CbsdState_Unregistered).
			ForDeregistration().
			Build(),
		expectedActions: []interface{}{
			&active_mode.AcknowledgeCbsdUpdateRequest{Id: 123},
		},
	}}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			g := message_generator.NewMessageGenerator(0, timeout, &stubRNG{})
			state := &active_mode.State{Cbsds: []*active_mode.Cbsd{tt.cbsd}}
			msgs := g.GenerateMessages(state, now)
			p := &stubProvider{}
			for _, msg := range msgs {
				_ = msg.Send(context.Background(), p)
			}
			require.Len(t, p.requests, len(tt.expectedRequests))
			for i := range tt.expectedRequests {
				assert.JSONEq(t, tt.expectedRequests[i].Payload, p.requests[i].Payload)
			}
			require.Len(t, p.actions, len(tt.expectedActions))
			for i := range tt.expectedActions {
				assert.Equal(t, tt.expectedActions[i], p.actions[i])
			}
		})
	}
}

func getSpectrumInquiryRequest() []*requests.RequestPayload {
	return []*requests.RequestPayload{{
		Payload: `{
	"spectrumInquiryRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"inquiredSpectrum": [
				{
					"lowFrequency": 3550000000,
					"highFrequency": 3700000000
				}
			]
		}
	]
}`,
	}}
}

type stubRNG struct{}

func (s *stubRNG) Int() int {
	return 0
}

type stubProvider struct {
	requests []*requests.RequestPayload
	actions  []interface{}
}

func (s *stubProvider) GetRequestsClient() requests.RadioControllerClient {
	return &stubRadioControllerClient{requests: &s.requests}
}

func (s *stubProvider) GetActiveModeClient() active_mode.ActiveModeControllerClient {
	return &stubActiveModeControllerClient{actions: &s.actions}
}

type stubActiveModeControllerClient struct {
	actions *[]interface{}
}

func (s *stubActiveModeControllerClient) GetState(_ context.Context, _ *active_mode.GetStateRequest, _ ...grpc.CallOption) (*active_mode.State, error) {
	panic("not implemented")
}

func (s *stubActiveModeControllerClient) DeleteCbsd(_ context.Context, in *active_mode.DeleteCbsdRequest, _ ...grpc.CallOption) (*empty.Empty, error) {
	*s.actions = append(*s.actions, in)
	return nil, nil
}

func (s *stubActiveModeControllerClient) AcknowledgeCbsdUpdate(_ context.Context, in *active_mode.AcknowledgeCbsdUpdateRequest, _ ...grpc.CallOption) (*empty.Empty, error) {
	*s.actions = append(*s.actions, in)
	return nil, nil
}
func (s *stubActiveModeControllerClient) StoreAvailableFrequencies(_ context.Context, in *active_mode.StoreAvailableFrequenciesRequest, _ ...grpc.CallOption) (*empty.Empty, error) {
	*s.actions = append(*s.actions, in)
	return nil, nil
}

type stubRadioControllerClient struct {
	requests *[]*requests.RequestPayload
}

func (s *stubRadioControllerClient) UploadRequests(_ context.Context, in *requests.RequestPayload, _ ...grpc.CallOption) (*requests.RequestDbIds, error) {
	*s.requests = append(*s.requests, in)
	return nil, nil
}
