package message_generator

import (
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"strings"
)

type heartbeatRequest struct {
	CbsdId         string `json:"cbsdId"`
	GrantId        string `json:"grantId"`
	OperationState string `json:"operationState"`
}

func (*heartbeatRequest) name() string {
	return "heartbeat"
}

type heartbeatMessageGenerator struct {
	now clock
}

func (h *heartbeatMessageGenerator) generateMessages(config *active_mode.ActiveModeConfig) []message {
	now := h.now().Unix()
	cbsd := config.GetCbsd()
	grants := cbsd.GetGrants()
	messages := make([]message, 0, len(grants))
	for _, grant := range grants {
		if grant.GetState() == active_mode.GrantState_Authorized && !isExpired(grant, now) {
			continue
		}
		req := &heartbeatRequest{
			CbsdId:         cbsd.GetId(),
			GrantId:        grant.GetId(),
			OperationState: strings.ToUpper(grant.GetState().String()),
		}
		messages = append(messages, req)
	}
	return messages
}

func isExpired(grant *active_mode.Grant, unixNow int64) bool {
	return grant.GetHeartbeatIntervalSec()+grant.GetLastHeartbeatTimestamp() <= unixNow
}
