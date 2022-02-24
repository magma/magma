package sas

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type deregistrationRequestGenerator struct{}

func NewDeregistrationRequestGenerator() *deregistrationRequestGenerator {
	return &deregistrationRequestGenerator{}
}

func (*deregistrationRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
	req := &deregistrationRequest{
		CbsdId: cbsd.GetId(),
	}
	return []*Request{asRequest(Deregistration, req)}
}

type deregistrationRequest struct {
	CbsdId string `json:"cbsdId"`
}
