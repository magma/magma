package message_generator_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
)

func TestGenerateMessages(t *testing.T) {
	const timeout = 100 * time.Second
	now := time.Unix(1000, 0)
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
						State:             active_mode.CbsdState_Unregistered,
						LastSeenTimestamp: now.Unix(),
					},
				}},
			},
			expected: nil,
		},
		{
			name: "Should do nothing when inactive cbsd has no grants",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						State:             active_mode.CbsdState_Unregistered,
						LastSeenTimestamp: 0,
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
						Id:                "some_cbsd_id",
						State:             active_mode.CbsdState_Registered,
						LastSeenTimestamp: now.Unix(),
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
						UserId:            "some_user_id",
						FccId:             "some_fcc_id",
						SerialNumber:      "some_serial_number",
						State:             active_mode.CbsdState_Unregistered,
						LastSeenTimestamp: now.Unix(),
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
						Id:                "some_cbsd_id",
						State:             active_mode.CbsdState_Registered,
						LastSeenTimestamp: now.Unix(),
					},
				}},
			},
			expected: getSpectrumInquiryRequest(),
		},
		{
			name: "Should generate spectrum inquiry request when all channels are unsuitable",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						Id:    "some_cbsd_id",
						State: active_mode.CbsdState_Registered,
						Channels: []*active_mode.Channel{{
							FrequencyRange: &active_mode.FrequencyRange{
								Low:  3.62e9,
								High: 3.63e9,
							},
							MaxEirp: wrapperspb.Float(4),
						}},
						EirpCapabilities: &active_mode.EirpCapabilities{
							MinPower:      0,
							MaxPower:      10,
							AntennaGain:   15,
							NumberOfPorts: 1,
						},
						LastSeenTimestamp: now.Unix(),
					},
				}},
			},
			expected: getSpectrumInquiryRequest(),
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
								Low:  3.62e9,
								High: 3.63e9,
							},
							MaxEirp: wrapperspb.Float(15),
						}},
						EirpCapabilities: &active_mode.EirpCapabilities{
							MinPower:      0,
							MaxPower:      100,
							AntennaGain:   0,
							NumberOfPorts: 1,
						},
						LastSeenTimestamp: now.Unix(),
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
					"lowFrequency": 3620000000,
					"highFrequency": 3630000000
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
						LastSeenTimestamp: now.Unix(),
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
		{
			name: "Should send both heartbeat and relinquish message for 2 grants when one is in Granted state and the other in Unsync",
			state: &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						Id:    "some_cbsd_id",
						State: active_mode.CbsdState_Registered,
						Grants: []*active_mode.Grant{
							{
								Id:    "some_grant_id",
								State: active_mode.GrantState_Granted,
							},
							{
								Id:    "some_other_grant_id",
								State: active_mode.GrantState_Unsync,
							},
						},
						LastSeenTimestamp: now.Unix(),
					},
				}},
			},
			expected: []*requests.RequestPayload{
				{
					Payload: `{
	"heartbeatRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"grantId": "some_grant_id",
			"operationState": "GRANTED"
		}
	]
}`,
				},
				{
					Payload: `{
	"relinquishmentRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"grantId": "some_other_grant_id"
		}
	]
}`,
				},
			},
		},
		{
			name: "Should send relinquish message when inactive for too long",
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
						LastSeenTimestamp: 0,
					},
				}},
			},
			expected: []*requests.RequestPayload{{
				Payload: `{
	"relinquishmentRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"grantId": "some_grant_id"
		}
	]
}`,
			}},
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			g := message_generator.NewMessageGenerator(0, timeout)
			actual := g.GenerateMessages(tt.state, now)
			require.Len(t, actual, len(tt.expected))
			for i := range tt.expected {
				assert.JSONEq(t, tt.expected[i].Payload, actual[i].Payload)
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
