package sas_test

import (
	"testing"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestSpectrumInquiryRequestGenerator(t *testing.T) {
	config := &active_mode.ActiveModeConfig{
		Cbsd: &active_mode.Cbsd{
			Id: "some_id",
		},
	}
	g := sas.NewSpectrumInquiryRequestGenerator()
	actual := g.GenerateRequests(config)
	expected := []*request{{
		requestType: "spectrumInquiryRequest",
		data: `{
	"cbsdId": "some_id",
	"inquiredSpectrum": [{
		"lowFrequency": 3550000000,
		"highFrequency": 3700000000
	}]
}`,
	}}
	assertRequestsEqual(t, expected, actual)
}
