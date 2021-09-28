package message_generator_test

import (
	"testing"
	"time"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMessages(t *testing.T) {
	data := []struct {
		name     string
		state    *active_mode.State
		expected []*requests.RequestPayload
	}{
		{
			name: "Should do nothing for unregistered non active cbsd",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Unregistered,
					Cbsd: &active_mode.Cbsd{
						State: active_mode.CbsdState_Unregistered,
					},
				}},
			},
			expected: nil,
		},
		{
			name: "Should generate deregistration request for non active registered cbsd",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Unregistered,
					Cbsd: &active_mode.Cbsd{
						Id:    "some_cbsd_id",
						State: active_mode.CbsdState_Registered,
					},
				}},
			},
			expected: []*requests.RequestPayload{{
				Payload: `{
	"deregistrationRequest": [
		{
			"cbsdId": "some_cbsd_id"
		}
	]
}`,
			}},
		},
		{
			name: "Should generate registration request for active non registered cbsd",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						UserId:       "some_user_id",
						FccId:        "some_fcc_id",
						SerialNumber: "some_serial_number",
						State:        active_mode.CbsdState_Unregistered,
					},
				}},
			},
			expected: []*requests.RequestPayload{{
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
		},
		{
			name: "Should generate spectrum inquiry request when there are no available channels",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						Id:    "some_cbsd_id",
						State: active_mode.CbsdState_Registered,
					},
				}},
			},
			expected: []*requests.RequestPayload{{
				Payload: `{
	"spectrumInquiryRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"inquiredSpectrum": {
				"lowFrequency": 3550000000,
				"highFrequency": 3700000000
			}
		}
	]
}`,
			}},
		},
		{
			name: "Should generate grant request when there are channels",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						Id:    "some_cbsd_id",
						State: active_mode.CbsdState_Registered,
						Channels: []*active_mode.Channel{{
							FrequencyRange: &active_mode.FrequencyRange{
								Low:  10,
								High: 20,
							},
							MaxEirp: makeOptionalFloat(15),
						}},
					},
				}},
			},
			expected: []*requests.RequestPayload{{
				Payload: `{
	"grantRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"operationParam": {
				"maxEirp": 15,
				"operationFrequencyRange": {
					"lowFrequency": 10,
					"highFrequency": 20
				}
			}
		}
	]
}`,
			}},
		},
		{
			name: "Should send heartbeat message for grant in granted state",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						Id:    "some_cbsd_id",
						State: active_mode.CbsdState_Registered,
						Grants: []*active_mode.Grant{{
							Id:    "some_grant_id",
							State: active_mode.GrantState_Granted,
						}},
					},
				}},
			},
			expected: []*requests.RequestPayload{{
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
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			actual := message_generator.GenerateMessages(now, tt.state)
			assert.Len(t, actual, len(tt.expected))
			for i := range tt.expected {
				assert.JSONEq(t, tt.expected[i].Payload, actual[i].Payload)
			}
		})
	}
}

func makeOptionalFloat(v float32) *float32 {
	return &v
}

func now() time.Time {
	return time.Unix(1000, 0)
}
