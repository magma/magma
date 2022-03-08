package sas_test

import (
	"fmt"
	"testing"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

const (
	someGrantId  = "some_grant_id"
	otherGrantId = "other_grant_id"

	granted    = "GRANTED"
	authorized = "AUTHORIZED"

	nextSend          = 1000
	heartbeatInterval = 250
)

func TestHeartbeatRequestGenerator(t *testing.T) {
	data := []struct {
		name     string
		grants   []*active_mode.Grant
		expected []*request
	}{
		{
			name: "Should generate heartbeat immediately when grant is not authorized yet",
			grants: []*active_mode.Grant{{
				Id:                     "some_grant_id",
				State:                  active_mode.GrantState_Granted,
				HeartbeatIntervalSec:   heartbeatInterval,
				LastHeartbeatTimestamp: nextSend,
			}},
			expected: []*request{
				newHeartbeatParams(withState(granted)).toRequest(),
			},
		},
		{
			name: "Should generate heartbeat when timeout has expired",
			grants: []*active_mode.Grant{{
				Id:                     "some_grant_id",
				State:                  active_mode.GrantState_Authorized,
				HeartbeatIntervalSec:   heartbeatInterval,
				LastHeartbeatTimestamp: nextSend - heartbeatInterval,
			}},
			expected: []*request{
				newHeartbeatParams(withState(authorized)).toRequest(),
			},
		},
		{
			name: "Should not generate heartbeat request when timeout has not expired yet",
			grants: []*active_mode.Grant{{
				Id:                     "some_grant_id",
				State:                  active_mode.GrantState_Authorized,
				HeartbeatIntervalSec:   heartbeatInterval,
				LastHeartbeatTimestamp: nextSend - heartbeatInterval + 1,
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
				HeartbeatIntervalSec:   heartbeatInterval,
				LastHeartbeatTimestamp: nextSend - heartbeatInterval,
			}},
			expected: []*request{
				newHeartbeatParams(
					withGrantId(someGrantId),
					withState(granted),
				).toRequest(),
				newHeartbeatParams(
					withGrantId(otherGrantId),
					withState(authorized),
				).toRequest(),
			},
		},
		{
			name: "Should generate relinquish request for unsync grant",
			grants: []*active_mode.Grant{{
				Id:    "some_grant_id",
				State: active_mode.GrantState_Unsync,
			}},
			expected: []*request{{
				requestType: "relinquishmentRequest",
				data: `{
	"cbsdId": "some_cbsd_id",
	"grantId": "some_grant_id"
}`,
			}},
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			cbsd := &active_mode.Cbsd{
				Id:     "some_cbsd_id",
				Grants: tt.grants,
			}
			g := sas.NewHeartbeatRequestGenerator(nextSend)
			actual := g.GenerateRequests(cbsd)
			assertRequestsEqual(t, tt.expected, actual)
		})
	}
}

type heartbeatParams struct {
	grantId string
	state   string
}

type heartbeatOption func(*heartbeatParams)

func withGrantId(grantId string) heartbeatOption {
	return func(h *heartbeatParams) {
		h.grantId = grantId
	}
}

func withState(state string) heartbeatOption {
	return func(h *heartbeatParams) {
		h.state = state
	}
}

func newHeartbeatParams(options ...heartbeatOption) *heartbeatParams {
	h := &heartbeatParams{
		grantId: someGrantId,
	}
	for _, o := range options {
		o(h)
	}
	return h
}

func (h *heartbeatParams) toRequest() *request {
	const requestTemplate = `{
	"cbsdId": "some_cbsd_id",
	"grantId": "%s",
	"operationState": "%s"
}`
	payload := fmt.Sprintf(requestTemplate, h.grantId, h.state)
	return &request{
		requestType: "heartbeatRequest",
		data:        payload,
	}
}
