package sas

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type registrationRequestGenerator struct{}

func NewRegistrationRequestGenerator() *registrationRequestGenerator {
	return &registrationRequestGenerator{}
}

func (*registrationRequestGenerator) GenerateRequests(config *active_mode.ActiveModeConfig) []*Request {
	cbsd := config.GetCbsd()
	req := &registrationRequest{
		UserId:           cbsd.GetUserId(),
		FccId:            cbsd.GetFccId(),
		CbsdSerialNumber: cbsd.GetSerialNumber(),
	}
	return []*Request{asRequest(Registration, req)}
}

type registrationRequest struct {
	UserId           string `json:"userId"`
	FccId            string `json:"fccId"`
	CbsdSerialNumber string `json:"cbsdSerialNumber"`
}
