package message_generator_test

import (
	"fmt"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGrantMessages(t *testing.T) {
	data := []struct {
		name           string
		eirpCapability *float32
		channel        *active_mode.Channel
		expected       []*requests.RequestPayload
	}{
		{
			name:     "Should generate grant request with default max eirp",
			channel:  &active_mode.Channel{},
			expected: makeGrantRequest(37),
		},
		{
			name: "Should generate grant request with max eirp from channel",
			channel: &active_mode.Channel{
				MaxEirp: makeOptionalFloat(15),
			},
			expected: makeGrantRequest(15),
		},
		{
			name:           "Should generate grant request based on eirp capability",
			eirpCapability: makeOptionalFloat(25),
			channel: &active_mode.Channel{
				MaxEirp: makeOptionalFloat(30),
			},
			expected: makeGrantRequest(15),
		},
		{
			name:           "Should generate grant request based on last max eirp",
			eirpCapability: makeOptionalFloat(25),
			channel: &active_mode.Channel{
				MaxEirp:  makeOptionalFloat(30),
				LastEirp: makeOptionalFloat(11),
			},
			expected: makeGrantRequest(10),
		},
		{
			name:           "Should not generate grant request if eirp 0 or less",
			eirpCapability: makeOptionalFloat(25),
			channel: &active_mode.Channel{
				MaxEirp:  makeOptionalFloat(30),
				LastEirp: makeOptionalFloat(1),
			},
			expected: nil,
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			tt.channel.FrequencyRange = &active_mode.FrequencyRange{
				Low:  10,
				High: 20,
			}
			state := &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						Id:             "some_cbsd_id",
						State:          active_mode.CbsdState_Registered,
						Channels:       []*active_mode.Channel{tt.channel},
						EirpCapability: tt.eirpCapability,
					},
				}},
			}
			actual := message_generator.GenerateMessages(now, state)
			assert.Len(t, actual, len(tt.expected))
			for i := range tt.expected {
				assert.JSONEq(t, tt.expected[i].Payload, actual[i].Payload)
			}
		})
	}
}

func makeGrantRequest(maxEirp float32) []*requests.RequestPayload {
	const requestTemplate = `{
	"grantRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"operationParam": {
				"maxEirp": %v,
				"operationFrequencyRange": {
					"lowFrequency": 10,
					"highFrequency": 20
				}
			}
		}
	]
}`
	payload := fmt.Sprintf(requestTemplate, maxEirp)
	return []*requests.RequestPayload{{Payload: payload}}
}
