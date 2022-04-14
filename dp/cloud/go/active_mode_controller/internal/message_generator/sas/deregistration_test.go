package sas_test

import (
	"testing"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestDeregistrationRequestGenerator(t *testing.T) {
	cbsd := &active_mode.Cbsd{Id: "some_id"}
	g := sas.NewDeregistrationRequestGenerator()
	actual := g.GenerateRequests(cbsd)
	expected := []*request{{
		requestType: "deregistrationRequest",
		data: `{
	"cbsdId": "some_id"
}`,
	}}
	assertRequestsEqual(t, expected, actual)
}
