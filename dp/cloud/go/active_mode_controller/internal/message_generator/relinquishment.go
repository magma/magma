package message_generator

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type relinquishmentRequest struct {
	CbsdId  string `json:"cbsdId"`
	GrantId string `json:"grantId"`
}

func (*relinquishmentRequest) name() string {
	return "relinquishment"
}

type relinquishmentMessageGenerator struct{}

func (*relinquishmentMessageGenerator) generateMessages(config *active_mode.ActiveModeConfig) []message {
	cbsd := config.GetCbsd()
	grants := cbsd.GetGrants()
	messages := make([]message, 0, len(grants))
	for _, grant := range grants {
		req := &relinquishmentRequest{
			CbsdId:  cbsd.GetId(),
			GrantId: grant.GetId(),
		}
		messages = append(messages, req)
	}
	return messages
}
