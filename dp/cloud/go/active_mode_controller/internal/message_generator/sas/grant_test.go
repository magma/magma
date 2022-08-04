/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sas_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestGrantProcessor(t *testing.T) {
	const (
		eirp            = 10
		frequency int64 = 3600e6
		bandwidth int64 = 20e6
	)
	calc := &stubEirpCalculator{
		eirp: eirp,
	}
	channels := []*active_mode.Channel{{
		LowFrequencyHz:  3550e6,
		HighFrequencyHz: 3700e6,
		MaxEirp:         wrapperspb.Float(eirp),
	}}
	p := &sas.GrantProcessor{
		CbsdId:   "some_id",
		Calc:     calc,
		Channels: channels,
	}
	actual := p.ProcessGrant(frequency, bandwidth)
	expected := &request{
		requestType: "grantRequest",
		data: `{
	"cbsdId": "some_id",
	"operationParam": {
		"maxEirp": 10,
		"operationFrequencyRange": {
			"lowFrequency": 3590000000,
			"highFrequency": 3610000000
		}
	}
}`,
	}
	assert.Equal(t, channels, calc.channels)
	assert.Equal(t, frequency-bandwidth/2, calc.low)
	assert.Equal(t, frequency+bandwidth/2, calc.high)
	assertRequestEqual(t, expected, actual)
}

type stubEirpCalculator struct {
	channels []*active_mode.Channel
	low      int64
	high     int64
	eirp     float64
}

func (s *stubEirpCalculator) CalcUpperBoundForRange(channels []*active_mode.Channel, low int64, high int64) float64 {
	s.channels = channels
	s.low = low
	s.high = high
	return s.eirp
}
