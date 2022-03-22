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
		name          string
		capabilities  *active_mode.EirpCapabilities
		channels      []*active_mode.Channel
		grantAttempts int
		expected      *grantParams
	}{
		{
			name:         "Should generate grant request with default max eirp",
			capabilities: getDefaultCapabilities(),
			channels: []*active_mode.Channel{{
				FrequencyRange: &active_mode.FrequencyRange{
					Low:  3620 * mega,
					High: 3630 * mega,
				},
			}},
			expected: &grantParams{
				maxEirp:      37,
				minFrequency: 3620 * mega,
				maxFrequency: 3630 * mega,
			},
		},
		{
			name:         "Should generate grant request with max eirp from channels",
			capabilities: getDefaultCapabilities(),
			channels: []*active_mode.Channel{{
				FrequencyRange: &active_mode.FrequencyRange{
					Low:  3625 * mega,
					High: 3635 * mega,
				},
				MaxEirp: wrapperspb.Float(15),
			}},
			expected: &grantParams{
				maxEirp:      15,
				minFrequency: 3625 * mega,
				maxFrequency: 3635 * mega,
			},
		},
		{
			name: "Should generate grant request based on capabilities and bandwidth",
			capabilities: &active_mode.EirpCapabilities{
				MaxPower:      20,
				AntennaGain:   15,
				NumberOfPorts: 2,
			},
			channels: []*active_mode.Channel{{
				FrequencyRange: &active_mode.FrequencyRange{
					Low:  3625 * mega,
					High: 3635 * mega,
				},
			}},
			expected: &grantParams{
				maxEirp:      28,
				minFrequency: 3625 * mega,
				maxFrequency: 3635 * mega,
			},
		},
		{
			name:         "Should use merged channels",
			capabilities: getDefaultCapabilities(),
			channels: []*active_mode.Channel{{
				FrequencyRange: &active_mode.FrequencyRange{
					Low:  3550 * mega,
					High: 3560 * mega,
				},
			}, {
				FrequencyRange: &active_mode.FrequencyRange{
					Low:  3560 * mega,
					High: 3570 * mega,
				},
			}},
			expected: &grantParams{
				maxEirp:      37,
				minFrequency: 3550 * mega,
				maxFrequency: 3570 * mega,
			},
		},
		{
			name:         "Should not generate anything if there are no suitable channels",
			capabilities: getDefaultCapabilities(),
			channels: []*active_mode.Channel{{
				FrequencyRange: &active_mode.FrequencyRange{
					Low:  3550 * mega,
					High: 3553 * mega,
				},
			}},
			expected: nil,
		},
		{
			name:         "Should not generate anything if there are no channels",
			capabilities: getDefaultCapabilities(),
			expected:     nil,
		},
		{
			name:         "Should not generate anything if there are grant attempts",
			capabilities: getDefaultCapabilities(),
			channels: []*active_mode.Channel{{
				FrequencyRange: &active_mode.FrequencyRange{
					Low:  3550 * mega,
					High: 3700 * mega,
				},
			}},
			grantAttempts: 1,
		},
	}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			cbsd := &active_mode.Cbsd{
				Id:               "some_cbsd_id",
				Channels:         tt.channels,
				EirpCapabilities: tt.capabilities,
				GrantAttempts:    int32(tt.grantAttempts),
			}
			g := sas.NewGrantRequestGenerator(stubRNG{})
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
	maxEirp      float32
	minFrequency int
	maxFrequency int
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
	payload := fmt.Sprintf(requestTemplate, g.maxEirp, g.minFrequency, g.maxFrequency)
	return []*request{{
		requestType: "grantRequest",
		data:        payload,
	}}
}
