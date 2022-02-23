package sas

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type relinquishmentRequestGenerator struct{}

func NewRelinquishmentRequestGenerator() *relinquishmentRequestGenerator {
	return &relinquishmentRequestGenerator{}
}

func (*relinquishmentRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
	grants := cbsd.GetGrants()
	cbsdId := cbsd.GetId()
	reqs := make([]*Request, 0, len(grants))
	for _, grant := range grants {
		req := &relinquishmentRequest{
			CbsdId:  cbsdId,
			GrantId: grant.GetId(),
		}
		reqs = append(reqs, asRequest(Relinquishment, req))
	}
	return reqs
}

type relinquishmentRequest struct {
	CbsdId  string `json:"cbsdId"`
	GrantId string `json:"grantId"`
}
