package sas_test

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

const mega = 1e6

func TestGrantRequestGenerator(t *testing.T) {
	data := []struct {
		name         string
		capabilities *active_mode.EirpCapabilities
		channels     []*active_mode.Channel
		expected     []*request
	}{
		{
			name:         "Should generate grant request with default max eirp",
			capabilities: getDefaultCapabilities(),
			channels: []*active_mode.Channel{{
				FrequencyRange: getDefaultFrequencyRange(),
			}},
			expected: newGrantParams().toRequest(),
		},
		{
			name:         "Should generate grant request with max eirp from channels",
			capabilities: getDefaultCapabilities(),
			channels: []*active_mode.Channel{{
				FrequencyRange: getDefaultFrequencyRange(),
				MaxEirp:        wrapperspb.Float(15),
			}},
			expected: newGrantParams(withMaxEirp(15)).toRequest(),
		},
		{
			name:         "Should generate grant request based on last max eirp",
			capabilities: getDefaultCapabilities(),
			channels: []*active_mode.Channel{{
				FrequencyRange: getDefaultFrequencyRange(),
				MaxEirp:        wrapperspb.Float(30),
				LastEirp:       wrapperspb.Float(11),
			}},
			expected: newGrantParams(withMaxEirp(10)).toRequest(),
		},
		{
			name: "Should generate grant request based on capabilities and bandwidth",
			capabilities: &active_mode.EirpCapabilities{
				MaxPower:      20,
				AntennaGain:   15,
				NumberOfPorts: 2,
			},
			channels: []*active_mode.Channel{{
				FrequencyRange: getDefaultFrequencyRange(),
			}},
			expected: newGrantParams(withMaxEirp(28)).toRequest(),
		},
		{
			name: "Should not generate grant request if eirp 0 or less",
			capabilities: &active_mode.EirpCapabilities{
				MinPower:      0,
				AntennaGain:   0,
				NumberOfPorts: 1,
			},
			channels: []*active_mode.Channel{{
				FrequencyRange: getDefaultFrequencyRange(),
				LastEirp:       wrapperspb.Float(-10),
			}},
			expected: nil,
		},
		{
			name: "Should switch to another channel if current is unusable",
			capabilities: &active_mode.EirpCapabilities{
				MinPower:      0,
				MaxPower:      10,
				AntennaGain:   15,
				NumberOfPorts: 1,
			},
			channels: []*active_mode.Channel{{
				FrequencyRange: getDefaultFrequencyRange(),
				LastEirp:       wrapperspb.Float(5),
			}, {
				FrequencyRange: &active_mode.FrequencyRange{
					Low:  3550 * mega,
					High: 3560 * mega,
				},
				MaxEirp: wrapperspb.Float(6),
			}},
			expected: newGrantParams(
				withMaxEirp(6),
				withFrequencyMHz(3550*mega, 3560*mega),
			).toRequest(),
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			config := &active_mode.ActiveModeConfig{
				Cbsd: &active_mode.Cbsd{
					Id:               "some_cbsd_id",
					Channels:         tt.channels,
					EirpCapabilities: tt.capabilities,
				},
			}
			g := sas.NewGrantRequestGenerator()
			actual := g.GenerateRequests(config)
			assertRequestsEqual(t, tt.expected, actual)
		})
	}
}

func getDefaultFrequencyRange() *active_mode.FrequencyRange {
	return &active_mode.FrequencyRange{
		Low:  3.62e9,
		High: 3.63e9,
	}
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
	maxEirp      float32
	minFrequency int
	maxFrequency int
}

type grantOption func(*grantParams)

func withFrequencyMHz(low int, high int) grantOption {
	return func(g *grantParams) {
		g.minFrequency = low
		g.maxFrequency = high
	}
}

func withMaxEirp(eirp float32) grantOption {
	return func(g *grantParams) {
		g.maxEirp = eirp
	}
}

func newGrantParams(options ...grantOption) *grantParams {
	g := &grantParams{
		maxEirp:      37,
		minFrequency: 3620 * mega,
		maxFrequency: 3630 * mega,
	}
	for _, o := range options {
		o(g)
	}
	return g
}

func (g *grantParams) toRequest() []*request {
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
	payload := fmt.Sprintf(requestTemplate, g.maxEirp, g.minFrequency, g.maxFrequency)
	return []*request{{
		requestType: "grantRequest",
		data:        payload,
	}}
}
