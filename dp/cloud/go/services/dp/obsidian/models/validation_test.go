package models_test

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"

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
		data:          newMutableCbsd(withFccId("")),
		expectedError: "fcc_id in body is required",
	}, {
		name:          "Should validate serial number on create",
		data:          newMutableCbsd(withSerialNumber("")),
		expectedError: "serial_number in body is required",
	}, {
		name:          "Should validate user id on create",
		data:          newMutableCbsd(withUserId("")),
		expectedError: "user_id in body is required",
	}, {
		name:          "Should validate bandwidth",
		data:          newMutableCbsd(withBandwidth(0)),
		expectedError: "bandwidth_mhz in body is required",
	}, {
		name:          "Should validate frequencies",
		data:          newMutableCbsd(withFrequencies(nil)),
		expectedError: "frequencies_mhz in body is required",
	}, {
		name:          "Should validate antenna gain",
		data:          newMutableCbsd(withAntennaGain(nil)),
		expectedError: "antenna_gain in body is required",
	}, {
		name:          "Should validate max power",
		data:          newMutableCbsd(withMaxPower(nil)),
		expectedError: "max_power in body is required",
	}, {
		name:          "Should validate min power",
		data:          newMutableCbsd(withMinPower(nil)),
		expectedError: "min_power in body is required",
	}, {
		name:          "Should validate number of antennas",
		data:          newMutableCbsd(withNumberOfAntennas(0)),
		expectedError: "number_of_antennas in body is required",
	}, {
		name:          "Should validate incorrect bandwidth",
		data:          newMutableCbsd(withBandwidth(12)),
		expectedError: "bandwidth_mhz in body should be one of [5 10 15 20]",
	}, {
		name:          "Should validate too low frequency",
		data:          newMutableCbsd(withFrequencies([]int64{123})),
		expectedError: "frequencies_mhz.0 in body should be greater than or equal to 3555",
	}, {
		name:          "Should validate too high frequency",
		data:          newMutableCbsd(withFrequencies([]int64{12345})),
		expectedError: "frequencies_mhz.0 in body should be less than or equal to 3695",
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate(strfmt.Default)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func newMutableCbsd(options ...mutableCbsdOption) *models.MutableCbsd {
	m := &models.MutableCbsd{
		Capabilities: models.Capabilities{
			AntennaGain:      to_pointer.Float(1),
			MaxPower:         to_pointer.Float(24),
			MinPower:         to_pointer.Float(0),
			NumberOfAntennas: 1,
		},
		FrequencyPreferences: models.FrequencyPreferences{
			BandwidthMhz:   10,
			FrequenciesMhz: []int64{3600},
		},
		FccID:        "someFCCId",
		SerialNumber: "someSerialNumber",
		UserID:       "someUserId",
	}
	for _, o := range options {
		o(m)
	}
	return m
}

type mutableCbsdOption func(cbsd *models.MutableCbsd)

func withBandwidth(bandwidth int64) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.FrequencyPreferences.BandwidthMhz = bandwidth
	}
}

func withFrequencies(frequencies []int64) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.FrequencyPreferences.FrequenciesMhz = frequencies
	}
}

func withMaxPower(maxPower *float64) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.Capabilities.MaxPower = maxPower
	}
}

func withMinPower(minPower *float64) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.Capabilities.MinPower = minPower
	}
}

func withAntennaGain(antennaGain *float64) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.Capabilities.AntennaGain = antennaGain
	}
}

func withNumberOfAntennas(numberOfAntennas int64) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.Capabilities.NumberOfAntennas = numberOfAntennas
	}
}

func withFccId(fccId string) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.FccID = fccId
	}
}

func withSerialNumber(serialNumber string) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.SerialNumber = serialNumber
	}
}

func withUserId(userId string) mutableCbsdOption {
	return func(m *models.MutableCbsd) {
		m.UserID = userId
	}
}
