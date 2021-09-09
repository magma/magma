package message_generator_test

import (
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHeartbeatMessages(t *testing.T) {
	data := []struct {
		name     string
		grants   []*active_mode.Grant
		expected []*requests.RequestPayload
	}{
		{
			name: "Should generate heartbeat request for authorized grant",
			grants: []*active_mode.Grant{{
				Id:                     "some_grant_id",
				State:                  active_mode.GrantState_Authorized,
				HeartbeatIntervalSec:   10,
				LastHeartbeatTimestamp: now().Add(-20 * time.Second).Unix(),
			}},
			expected: []*requests.RequestPayload{{
				Payload: `{
	"heartbeatRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"grantId": "some_grant_id",
			"operationState": "AUTHORIZED"
		}
	]
}`,
			}},
		},
		{
			name: "Should not generate heartbeat request when timeout has not expired yet",
			grants: []*active_mode.Grant{{
				Id:                     "some_grant_id",
				State:                  active_mode.GrantState_Authorized,
				HeartbeatIntervalSec:   10,
				LastHeartbeatTimestamp: now().Add(-5 * time.Second).Unix(),
			}},
			expected: nil,
		},
		{
			name: "Should generate heartbeat requests for multiple grants",
			grants: []*active_mode.Grant{{
				Id:    "some_grant_id",
				State: active_mode.GrantState_Granted,
			}, {
				Id:                     "other_grant_id",
				State:                  active_mode.GrantState_Authorized,
				HeartbeatIntervalSec:   10,
				LastHeartbeatTimestamp: now().Add(-20 * time.Second).Unix(),
			}},
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
			}, {
				Payload: `{
	"heartbeatRequest": [
		{
			"cbsdId": "some_cbsd_id",
			"grantId": "other_grant_id",
			"operationState": "AUTHORIZED"
		}
	]
}`,
			}},
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			state := &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Registered,
					Cbsd: &active_mode.Cbsd{
						Id:     "some_cbsd_id",
						State:  active_mode.CbsdState_Registered,
						Grants: tt.grants,
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
