package message_generator

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type registrationRequest struct {
	UserId           string `json:"userId"`
	FccId            string `json:"fccId"`
	CbsdSerialNumber string `json:"cbsdSerialNumber"`
}

func (*registrationRequest) name() string {
	return "registration"
}

type registerMessageGenerator struct{}

func (*registerMessageGenerator) generateMessages(config *active_mode.ActiveModeConfig) []message {
	cbsd := config.GetCbsd()
	req := &registrationRequest{
		UserId:           cbsd.GetUserId(),
		FccId:            cbsd.GetFccId(),
		CbsdSerialNumber: cbsd.GetSerialNumber(),
	}
	return []message{req}
}
