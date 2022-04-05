package sas_helpers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas_helpers"
)

func TestBuild(t *testing.T) {
	const someDeregistrationRequest = `{"cbsdId":"someId"}`
	const otherDeregistrationRequest = `{"cbsdId":"otherId"}`
	const someHeartbeatRequest = `{"cbsdId":"someId","grantId":"grantId"}`
	const someRegistrationRequest = `{"key":"value"}`
	requests := []*sas.Request{{
		Type: sas.Deregistration,
		Data: []byte(someDeregistrationRequest),
	}, {
		Type: sas.Heartbeat,
		Data: []byte(someHeartbeatRequest),
	}, {
		Type: sas.Deregistration,
		Data: []byte(otherDeregistrationRequest),
	}, {
		Type: sas.Registration,
		Data: []byte(someRegistrationRequest),
	}}
	actual := sas_helpers.Build(requests)
	expected := []string{
		fmt.Sprintf(`{"%s":[%s]}`, sas.Registration, someRegistrationRequest),
		fmt.Sprintf(`{"%s":[%s]}`, sas.Heartbeat, someHeartbeatRequest),
		fmt.Sprintf(`{"%s":[%s,%s]}`, sas.Deregistration,
			someDeregistrationRequest, otherDeregistrationRequest),
	}
	assert.Equal(t, expected, actual)
}
