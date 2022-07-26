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
	"github.com/stretchr/testify/require"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestRegistrationRequestGenerator(t *testing.T) {
	data := []struct {
		name     string
		cbsd     *active_mode.Cbsd
		expected string
	}{{
		name: "Should generate multi step registration request",
		cbsd: &active_mode.Cbsd{
			SasSettings: &active_mode.SasSettings{
				UserId:       "some_user_id",
				FccId:        "some_fcc_id",
				SerialNumber: "some_serial_number",
			},
		},
		expected: `{
	"userId": "some_user_id",
	"fccId": "some_fcc_id",
	"cbsdSerialNumber": "some_serial_number"
}`,
	}, {
		name: "Should generate cpi less single step registration request",
		cbsd: &active_mode.Cbsd{
			SasSettings: &active_mode.SasSettings{
				SingleStepEnabled: true,
				CbsdCategory:      "a",
				SerialNumber:      "some_serial_number",
				FccId:             "some_fcc_id",
				UserId:            "some_user_id",
			},
			InstallationParams: &active_mode.InstallationParams{
				LatitudeDeg:      12,
				LongitudeDeg:     34,
				HeightM:          5,
				HeightType:       "agl",
				IndoorDeployment: true,
				AntennaGainDbi:   15,
			},
		},
		expected: `{
	"userId": "some_user_id",
	"fccId": "some_fcc_id",
	"cbsdSerialNumber": "some_serial_number",
	"cbsdCategory": "A",
	"airInterface": {
		"radioTechnology": "E_UTRA"
	},
	"installationParam": {
		"latitude": 12,
		"longitude": 34,
		"height": 5,
		"heightType": "AGL",
		"indoorDeployment": true,
		"antennaGain": 15
	},
	"measCapability": []
}`,
	}}
	g := &sas.RegistrationRequestGenerator{}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			actual := g.GenerateRequests(tt.cbsd)
			expected := []*request{{
				requestType: "registrationRequest",
				data:        tt.expected,
			}}
			assertRequestsEqual(t, expected, actual)
		})
	}
}

type request struct {
	requestType string
	data        string
}

func assertRequestsEqual(t *testing.T, expected []*request, actual []*sas.Request) {
	require.Len(t, actual, len(expected))
	for i := range actual {
		args := []any{"at %d", i}
		assertRequestEqual(t, expected[i], actual[i], args...)
	}
}

func assertRequestEqual(t *testing.T, expected *request, actual *sas.Request, args ...any) {
	if expected == nil {
		assert.Nil(t, actual, args...)
		return
	}
	assert.Equal(t, expected.requestType, actual.Type.String(), args...)
	assert.JSONEq(t, expected.data, string(actual.Data), args...)
}
