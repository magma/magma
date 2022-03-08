package sas

import (
	"strings"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type heartbeatRequestGenerator struct {
	nextSendTimestamp int64
}

func NewHeartbeatRequestGenerator(nextSendTimestamp int64) *heartbeatRequestGenerator {
	return &heartbeatRequestGenerator{
		nextSendTimestamp: nextSendTimestamp,
	}
}

func (h *heartbeatRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
	grants := cbsd.GetGrants()
	var reqs []*Request
	for _, grant := range grants {
		if grant.GetState() == active_mode.GrantState_Unsync {
			req := &relinquishmentRequest{
				CbsdId:  cbsd.GetId(),
				GrantId: grant.GetId(),
			}
			reqs = append(reqs, asRequest(Relinquishment, req))
			continue
		}
		if grant.GetState() == active_mode.GrantState_Authorized &&
			!shouldSendNow(grant, h.nextSendTimestamp) {
			continue
		}
		req := &heartbeatRequest{
			CbsdId:         cbsd.GetId(),
			GrantId:        grant.GetId(),
			OperationState: strings.ToUpper(grant.GetState().String()),
		}
		reqs = append(reqs, asRequest(Heartbeat, req))
	}
	return reqs
}

type heartbeatRequest struct {
	CbsdId         string `json:"cbsdId"`
	GrantId        string `json:"grantId"`
	OperationState string `json:"operationState"`
}

func shouldSendNow(grant *active_mode.Grant, nextSendTimestamp int64) bool {
	deadline := grant.GetHeartbeatIntervalSec() + grant.GetLastHeartbeatTimestamp()
	return deadline <= nextSendTimestamp
}
