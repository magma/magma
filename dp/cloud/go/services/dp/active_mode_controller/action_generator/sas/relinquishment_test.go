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

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestRelinquishmentRequestGenerator(t *testing.T) {
	data := &storage.DetailedCbsd{
		Cbsd: &storage.DBCbsd{
			CbsdId: db.MakeString("some_cbsd_id"),
		},
		Grants: []*storage.DetailedGrant{{
			Grant: &storage.DBGrant{
				GrantId: db.MakeString("some_grant_id"),
			},
		}},
	}
	g := &sas.RelinquishmentRequestGenerator{}
	actual := g.GenerateRequests(data)
	expected := []*request{getRelinquishmentRequest()}
	assertRequestsEqual(t, expected, actual)
}

func TestRelinquishmentProcessor(t *testing.T) {
	const frequency = 3600e6
	p := &sas.RelinquishmentProcessor{
		CbsdId: "some_cbsd_id",
		Grants: map[int64]*storage.DetailedGrant{
			frequency: {
				Grant: &storage.DBGrant{
					GrantId: db.MakeString("some_grant_id"),
				},
			},
		},
	}
	actual := p.ProcessGrant(frequency, 20e6)
	expected := getRelinquishmentRequest()
	assertRequestEqual(t, expected, actual)
}

func getRelinquishmentRequest() *request {
	return &request{
		requestType: "relinquishmentRequest",
		data: `{
	"cbsdId": "some_cbsd_id",
	"grantId": "some_grant_id"
}`,
	}
}
