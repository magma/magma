package sas_test

import (
	"testing"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestRelinquishmentRequestGenerator(t *testing.T) {
	cbsd := &active_mode.Cbsd{
		Id: "some_id",
		Grants: []*active_mode.Grant{{
			Id: "some_grant_id",
		}},
	}
	g := sas.NewRelinquishmentRequestGenerator()
	actual := g.GenerateRequests(cbsd)
	expected := []*request{{
		requestType: "relinquishmentRequest",
		data: `{
	"cbsdId": "some_id",
	"grantId": "some_grant_id"
}`,
	}}
	assertRequestsEqual(t, expected, actual)
}
