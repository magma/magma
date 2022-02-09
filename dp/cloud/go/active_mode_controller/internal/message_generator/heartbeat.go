package message_generator

import (
	"strings"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type heartbeatRequest struct {
	CbsdId         string `json:"cbsdId"`
	GrantId        string `json:"grantId"`
	OperationState string `json:"operationState"`
}

func (*heartbeatRequest) name() string {
	return "heartbeat"
}

type heartbeatOrRelinquishMessageGenerator struct {
	deadlineUnix int64
}

func (h *heartbeatOrRelinquishMessageGenerator) generateMessages(config *active_mode.ActiveModeConfig) []message {
	var req message
	cbsd := config.GetCbsd()
	cbsdId := cbsd.GetId()
	grants := cbsd.GetGrants()
	messages := make([]message, 0, len(grants))
	for _, grant := range grants {
		grantState := grant.GetState()
		grantId := grant.GetId()
		if grantState == active_mode.GrantState_Authorized && !shouldSendNow(grant, h.deadlineUnix) {
			continue
		} else if grantState == active_mode.GrantState_Unsync {
			req = &relinquishmentRequest{
				CbsdId:  cbsdId,
				GrantId: grantId,
			}
		} else {
			req = &heartbeatRequest{
				CbsdId:         cbsdId,
				GrantId:        grantId,
				OperationState: strings.ToUpper(grantState.String()),
			}
		}
		messages = append(messages, req)
	}
	return messages
}

func shouldSendNow(grant *active_mode.Grant, deadline int64) bool {
	return grant.GetHeartbeatIntervalSec()+grant.GetLastHeartbeatTimestamp() <= deadline
}
