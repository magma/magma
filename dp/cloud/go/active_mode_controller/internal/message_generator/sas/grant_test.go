package sas_test

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestGrantRequestGenerator(t *testing.T) {
	data := []struct {
		name          string
		capabilities  *active_mode.EirpCapabilities
		channels      []*active_mode.Channel
		grantAttempts int
		preferences   active_mode.FrequencyPreferences
		expected      *grantParams
	}{{
		name:         "Should generate grant request with default max eirp",
		capabilities: getDefaultCapabilities(),
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3620 * 1e6,
			HighFrequencyHz: 3630 * 1e6,
		}},
		expected: &grantParams{
			maxEirp:       37,
			lowFrequency:  3620 * 1e6,
			highFrequency: 3630 * 1e6,
		},
	}, {
		name:         "Should generate grant request with max eirp from channels",
		capabilities: getDefaultCapabilities(),
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3625 * 1e6,
			HighFrequencyHz: 3635 * 1e6,
			MaxEirp:         wrapperspb.Float(15),
		}},
		expected: &grantParams{
			maxEirp:       15,
			lowFrequency:  3625 * 1e6,
			highFrequency: 3635 * 1e6,
		},
	}, {
		name: "Should generate grant request based on capabilities and bandwidth",
		capabilities: &active_mode.EirpCapabilities{
			MaxPower:      20,
			AntennaGain:   15,
			NumberOfPorts: 2,
		},
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3625 * 1e6,
			HighFrequencyHz: 3635 * 1e6,
		}},
		expected: &grantParams{
			maxEirp:       28,
			lowFrequency:  3625 * 1e6,
			highFrequency: 3635 * 1e6,
		},
	}, {
		name:         "Should use merged channels",
		capabilities: getDefaultCapabilities(),
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3550 * 1e6,
			HighFrequencyHz: 3560 * 1e6,
		}, {
			LowFrequencyHz:  3560 * 1e6,
			HighFrequencyHz: 3570 * 1e6,
		}},
		expected: &grantParams{
			maxEirp:       37,
			lowFrequency:  3550 * 1e6,
			highFrequency: 3570 * 1e6,
		},
	}, {
		name:         "Should not generate anything if there are no suitable channels",
		capabilities: getDefaultCapabilities(),
		channels: []*active_mode.Channel{{
			LowFrequencyHz:  3550 * 1e6,
			HighFrequencyHz: 3553 * 1e6,
		}},
		expected: nil,
	}, {
		name:         "Should not generate anything if there are no channels",
		capabilities: getDefaultCapabilities(),
		expected:     nil,
	}}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			capabilities := proto.Clone(tt.capabilities).(*active_mode.EirpCapabilities)
			preferences := proto.Clone(&tt.preferences).(*active_mode.FrequencyPreferences)
			var channels []*active_mode.Channel
			for _, ch := range tt.channels {
				channels = append(channels, proto.Clone(ch).(*active_mode.Channel))
			}
			cbsd := &active_mode.Cbsd{
				Id:               "some_cbsd_id",
				Channels:         channels,
				EirpCapabilities: capabilities,
				GrantAttempts:    int32(tt.grantAttempts),
				Preferences:      preferences,
			}
			g := sas.NewGrantRequestGenerator(stubRNG{})
			actual := g.GenerateRequests(cbsd)
			expected := toRequest(tt.expected)
			assertRequestsEqual(t, expected, actual)
		})
	}
}

func TestGrantSelectionOrder(t *testing.T) {
	data := []struct {
		grantAttempts int32
		expected      *grantParams
	}{{
		grantAttempts: 0,
		expected: &grantParams{
			lowFrequency:  3652.5 * 1e6,
			highFrequency: 3657.5 * 1e6,
			maxEirp:       28,
		},
	}, {
		grantAttempts: 1,
		expected: &grantParams{
			lowFrequency:  3572.5 * 1e6,
			highFrequency: 3587.5 * 1e6,
			maxEirp:       5,
		},
	}, {
		grantAttempts: 2,
		expected: &grantParams{
			lowFrequency:  3575 * 1e6,
			highFrequency: 3585 * 1e6,
			maxEirp:       5,
		},
	}, {
		grantAttempts: 3,
		expected: &grantParams{
			lowFrequency:  3550 * 1e6,
			highFrequency: 3560 * 1e6,
			maxEirp:       10,
		},
	}, {
		grantAttempts: 4,
		expected: &grantParams{
			lowFrequency:  3552.5 * 1e6,
			highFrequency: 3557.5 * 1e6,
			maxEirp:       10,
		},
	}, {
		grantAttempts: 5,
		expected: &grantParams{
			lowFrequency:  3670 * 1e6,
			highFrequency: 3680 * 1e6,
			maxEirp:       25,
		},
	}, {
		grantAttempts: 6,
		expected: &grantParams{
			lowFrequency:  3557.5 * 1e6,
			highFrequency: 3562.5 * 1e6,
			maxEirp:       10,
		},
	}, {
		grantAttempts: 7,
		expected:      nil,
	}}
	g := sas.NewGrantRequestGenerator(stubRNG{})
	cbsd := &active_mode.Cbsd{
		Id: "some_cbsd_id",
		Channels: []*active_mode.Channel{{
			LowFrequencyHz:  3652.5 * 1e6,
			HighFrequencyHz: 3657.5 * 1e6,
		}, {
			LowFrequencyHz:  3572.5 * 1e6,
			HighFrequencyHz: 3587.5 * 1e6,
			MaxEirp:         wrapperspb.Float(5),
		}, {
			LowFrequencyHz:  3550 * 1e6,
			HighFrequencyHz: 3564 * 1e6,
			MaxEirp:         wrapperspb.Float(10),
		}, {
			LowFrequencyHz:  3670 * 1e6,
			HighFrequencyHz: 3680 * 1e6,
		}},
		EirpCapabilities: &active_mode.EirpCapabilities{
			MinPower:      0,
			MaxPower:      20,
			AntennaGain:   15,
			NumberOfPorts: 1,
		},
		Preferences: &active_mode.FrequencyPreferences{
			BandwidthMhz:   15,
			FrequenciesMhz: []int32{3655, 3580, 3555},
		},
	}
	for _, tt := range data {
		t.Run(fmt.Sprintf("Attempt: %d", tt.grantAttempts), func(t *testing.T) {
			cbsd.GrantAttempts = tt.grantAttempts
			actual := g.GenerateRequests(cbsd)
			expected := toRequest(tt.expected)
			assertRequestsEqual(t, expected, actual)
		})
	}
}

type stubRNG struct{}

func (stubRNG) Int() int {
	return 0
}

func getDefaultCapabilities() *active_mode.EirpCapabilities {
	return &active_mode.EirpCapabilities{
		MinPower:      -1000,
		MaxPower:      1000,
		AntennaGain:   0,
		NumberOfPorts: 1,
	}
}

type grantParams struct {
	lowFrequency  int
	highFrequency int
	maxEirp       float32
}

func toRequest(g *grantParams) []*request {
	if g == nil {
		return nil
	}
	const requestTemplate = `{
	"cbsdId": "some_cbsd_id",
	"operationParam": {
		"maxEirp": %v,
		"operationFrequencyRange": {
			"lowFrequency": %d,
			"highFrequency": %d
		}
	}
}`
	payload := fmt.Sprintf(requestTemplate, g.maxEirp, g.lowFrequency, g.highFrequency)
	return []*request{{
		requestType: "grantRequest",
		data:        payload,
	}}
}
