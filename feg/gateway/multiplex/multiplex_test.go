/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package multiplex_test

import (
	"magma/feg/gateway/multiplex"
	"testing"

	"github.com/stretchr/testify/assert"
)

type muxSelectorsScenario struct {
	sessionId     string
	imsiStr       string
	imsiNumeric   uint64
	conversionErr bool
}

type multiplexorScenario struct {
	sessionId    string
	numServers   int
	outputServer int
	outputErr    bool
}

var (
	// Scenarios to test multiplex.Context IMSI parer
	muxSelectorTestScenarios = []muxSelectorsScenario{
		{"IMSI123456789012345-54321", "IMSI123456789012345", 123456789012345, false},
		{"IMSI999999999999999-54321", "IMSI999999999999999", 999999999999999, false},
		{"IMSI000000000000010-54321", "IMSI000000000000010", 10, false},
		{"IMSI1234AAAAA012345-54321", "IMSI1234AAAAA012345", 0, true},
	}
	// Scenarios to test StaticMultiplexByIMSI
	staticMultiplexByIMSIScenario = []multiplexorScenario{
		{"IMSI123456789012345-54321", 5, 0, false},
		{"IMSI999999999999999-54321", 7, 5, false},
		{"IMSI000000000000010-54321", 6, 4, false},
		{"IMSI1234AAAAA012345-54321", -1, -1, true},
	}
)

func TestMultiplexContextSessionIdConversions(t *testing.T) {
	for _, scenario := range muxSelectorTestScenarios {
		assertIMSIConversion(t, multiplex.NewContext().WithSessionId(scenario.sessionId), scenario)
	}
}

func TestMultiplexContextIMSIstrConversions(t *testing.T) {
	for _, scenario := range muxSelectorTestScenarios {
		assertIMSIConversion(t, multiplex.NewContext().WithIMSI(scenario.imsiStr), scenario)
	}
}

func assertIMSIConversion(t *testing.T, muxCtx *multiplex.Context, scenario muxSelectorsScenario) {
	if scenario.conversionErr {
		assert.Errorf(t, muxCtx.GetError(), "Conversion should have failed on scenario: %+v", scenario)
	} else {
		imsi, err := muxCtx.GetIMSI()
		assert.NoErrorf(t, err, "on scenario: %+v", scenario)
		assert.Equalf(t, scenario.imsiNumeric, imsi, "on scenario: %+v", scenario)
	}
}

func TestStaticMultiplexByIMSI(t *testing.T) {
	for _, scenario := range staticMultiplexByIMSIScenario {
		mux, err := multiplex.NewStaticMultiplexByIMSI(scenario.numServers)
		if scenario.outputErr {
			assert.Errorf(t, err, "StaticMultiplexByIMSI should have returned an error on scenario: %+v", scenario)
		} else {
			assert.NoErrorf(t, err, "on scenario: %+v", scenario)
			server, err := mux.GetIndex(multiplex.NewContext().WithSessionId(scenario.sessionId))
			assert.NoError(t, err)
			assert.Equalf(t, scenario.outputServer, server, "on scenario: %+v", scenario)
		}
	}

}
