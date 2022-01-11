package message_generator_test

import (
	"testing"
	"time"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"

	"github.com/stretchr/testify/assert"
)

func TestFilterMessages(t *testing.T) {
	data := []struct {
		name            string
		pendingRequests []string
		expected        []*requests.RequestPayload
	}{
		{
			name:            "Should filter request if pending",
			pendingRequests: []string{`{"cbsdId":"some"}`},
			expected:        []*requests.RequestPayload{},
		},
		{
			name:            "Should not filter request if not pending",
			pendingRequests: nil,
			expected: []*requests.RequestPayload{{
				Payload: `{"deregistrationRequest":[{"cbsdId":"some"}]}`,
			}},
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			g := message_generator.NewMessageGenerator(0, 0)
			state := &active_mode.State{
				ActiveModeConfigs: []*active_mode.ActiveModeConfig{{
					DesiredState: active_mode.CbsdState_Unregistered,
					Cbsd: &active_mode.Cbsd{
						Id:              "some",
						State:           active_mode.CbsdState_Registered,
						PendingRequests: tt.pendingRequests,
					},
				}},
			}
			actual := g.GenerateMessages(state, time.Time{})
			assert.Equal(t, tt.expected, actual)
		})
	}
}
