package models_test

import (
	"context"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"

	b "magma/dp/cloud/go/services/dp/builders"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
)

func TestMutableCbsd_Validate(t *testing.T) {
	testData := []struct {
		name          string
		data          *models.MutableCbsd
		expectedError string
	}{{
		name:          "Should validate fcc id on create",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithFccId("").Payload,
		expectedError: "fcc_id in body is required",
	}, {
		name:          "Should validate serial number on create",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithSerialNumber("").Payload,
		expectedError: "serial_number in body is required",
	}, {
		name:          "Should validate user id on create",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithUserId("").Payload,
		expectedError: "user_id in body is required",
	}, {
		name:          "Should validate bandwidth",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithBandwidth(0).Payload,
		expectedError: "bandwidth_mhz in body is required",
	}, {
		name:          "Should validate frequencies",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithFrequencies(nil).Payload,
		expectedError: "frequencies_mhz in body is required",
	}, {
		name:          "Should validate max power",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithMaxPower(nil).Payload,
		expectedError: "max_power in body is required",
	}, {
		name:          "Should validate min power",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithMinPower(nil).Payload,
		expectedError: "min_power in body is required",
	}, {
		name:          "Should validate number of antennas",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithNumberOfAntennas(0).Payload,
		expectedError: "number_of_antennas in body is required",
	}, {
		name:          "Should validate incorrect bandwidth",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithBandwidth(12).Payload,
		expectedError: "bandwidth_mhz in body should be one of [5 10 15 20]",
	}, {
		name:          "Should validate too low frequency",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithFrequencies([]int64{123}).Payload,
		expectedError: "frequencies_mhz.0 in body should be greater than or equal to 3555",
	}, {
		name:          "Should validate too high frequency",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithFrequencies([]int64{12345}).Payload,
		expectedError: "frequencies_mhz.0 in body should be less than or equal to 3695",
	}, {
		name:          "Should validate single step enabled",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithSingleStepEnabled(nil).Payload,
		expectedError: "single_step_enabled in body is required",
	}, {
		name:          "Should validate cbsd category",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithCbsdCategory("").Payload,
		expectedError: "cbsd_category in body is required",
	}, {
		name:          "Should validate cbsd category value",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithCbsdCategory("c").Payload,
		expectedError: "cbsd_category in body should be one of [a b]",
	}, {
		name:          "Should validate carrier aggregation enabled",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithCarrierAggregationEnabled(nil).Payload,
		expectedError: "carrier_aggregation_enabled in body is required",
	}, {
		name:          "Should validate grant redundancy",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithGrantRedundancy(nil).Payload,
		expectedError: "grant_redundancy in body is required",
	}, {
		name:          "Should validate max ibw mhz",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithMaxIbwMhz(0).Payload,
		expectedError: "max_ibw_mhz in body is required",
	}, {
		name:          "Should validate max ibw mhz too high",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithMaxIbwMhz(151).Payload,
		expectedError: "max_ibw_mhz in body should be less than or equal to 150",
	}, {
		name:          "Should validate max ibw mhz not multiple of 5",
		data:          b.NewMutableCbsdModelPayloadBuilder().WithMaxIbwMhz(7).Payload,
		expectedError: "max_ibw_mhz in body should be a multiple of 5",
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate(strfmt.Default)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestMutableCbsd_ValidateModel(t *testing.T) {
	testData := []struct {
		name          string
		data          *models.MutableCbsd
		expectedError string
	}{{
		name: "Should validate grant_redundancy false with carrier aggregation enabled",
		data: b.NewMutableCbsdModelPayloadBuilder().
			WithGrantRedundancy(to_pointer.Bool(false)).
			WithCarrierAggregationEnabled(to_pointer.Bool(true)).
			Payload,
		expectedError: "grant_redundancy cannot be set to false when carrier_aggregation_enabled is enabled",
	}, {
		name: "Should validate max ibw mhz lesser than bandwidth mhz",
		data: b.NewMutableCbsdModelPayloadBuilder().
			WithMaxIbwMhz(5).WithBandwidth(10).Payload,
		expectedError: "max_ibw_mhz cannot be less than bandwidth_mhz",
	}}
	c := context.TODO()
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.ValidateModel(c)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}
