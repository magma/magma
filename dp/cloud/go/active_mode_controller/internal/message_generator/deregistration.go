package message_generator

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type deregistrationMessageGenerator struct{}

func (*deregistrationMessageGenerator) generateMessages(config *active_mode.ActiveModeConfig) []message {
	req := &DeregistrationRequest{
		CbsdId: config.GetCbsd().GetId(),
	}
	return []message{req}
}

type DeregistrationRequest struct {
	CbsdId string `json:"cbsdId"`
}

func (*DeregistrationRequest) name() string {
	return "deregistration"
}
