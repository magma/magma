package sas

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type deregistrationRequestGenerator struct{}

func NewDeregistrationRequestGenerator() *deregistrationRequestGenerator {
	return &deregistrationRequestGenerator{}
}

func (*deregistrationRequestGenerator) GenerateRequests(config *active_mode.ActiveModeConfig) []*Request {
	req := &deregistrationRequest{
		CbsdId: config.GetCbsd().GetId(),
	}
	return []*Request{asRequest(Deregistration, req)}
}

type deregistrationRequest struct {
	CbsdId string `json:"cbsdId"`
}
