package sas_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestRegistrationRequestGenerator(t *testing.T) {
	cbsd := &active_mode.Cbsd{
		UserId:       "some_user_id",
		FccId:        "some_fcc_id",
		SerialNumber: "some_serial_number",
	}
	g := sas.NewRegistrationRequestGenerator()
	actual := g.GenerateRequests(cbsd)
	expected := []*request{{
		requestType: "registrationRequest",
		data: `{
	"userId": "some_user_id",
	"fccId": "some_fcc_id",
	"cbsdSerialNumber": "some_serial_number"
}`,
	}}
	assertRequestsEqual(t, expected, actual)
}

type request struct {
	requestType string
	data        string
}

func assertRequestsEqual(t *testing.T, expected []*request, actual []*sas.Request) {
	require.Len(t, actual, len(expected))
	for i := range actual {
		args := []interface{}{"at %d", i}
		x, y := expected[i], actual[i]
		assert.Equal(t, x.requestType, y.Type.String(), args...)
		assert.JSONEq(t, x.data, string(y.Data), args...)
	}
}
