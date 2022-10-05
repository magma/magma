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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestRegistrationRequestGenerator(t *testing.T) {
	data := []struct {
		name     string
		cbsd     *storage.DBCbsd
		expected string
	}{{
		name: "Should generate multi step registration request",
		cbsd: &storage.DBCbsd{
			UserId:           db.MakeString("some_user_id"),
			FccId:            db.MakeString("some_fcc_id"),
			CbsdSerialNumber: db.MakeString("some_serial_number"),
		},
		expected: `{
	"userId": "some_user_id",
	"fccId": "some_fcc_id",
	"cbsdSerialNumber": "some_serial_number"
}`,
	}, {
		name: "Should generate cpi less single step registration request",
		cbsd: &storage.DBCbsd{
			SingleStepEnabled: db.MakeBool(true),
			CbsdCategory:      db.MakeString("a"),
			CbsdSerialNumber:  db.MakeString("some_serial_number"),
			FccId:             db.MakeString("some_fcc_id"),
			UserId:            db.MakeString("some_user_id"),
			LatitudeDeg:       db.MakeFloat(12),
			LongitudeDeg:      db.MakeFloat(34),
			HeightM:           db.MakeFloat(5),
			HeightType:        db.MakeString("agl"),
			IndoorDeployment:  db.MakeBool(true),
			AntennaGainDbi:    db.MakeFloat(15),
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
			data := &storage.DetailedCbsd{Cbsd: tt.cbsd}
			actual := g.GenerateRequests(data)
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

func assertRequestsEqual(t *testing.T, expected []*request, actual []*storage.MutableRequest) {
	require.Len(t, actual, len(expected))
	for i := range actual {
		args := []any{"at %d", i}
		assertRequestEqual(t, expected[i], actual[i], args...)
	}
}

func assertRequestEqual(t *testing.T, expected *request, actual *storage.MutableRequest, args ...any) {
	if expected == nil {
		assert.Nil(t, actual, args...)
		return
	}
	assert.Equal(t, expected.requestType, actual.RequestType.Name.String, args...)
	actualPayload, _ := json.Marshal(actual.Request.Payload)
	assert.JSONEq(t, expected.data, string(actualPayload), args...)
}
