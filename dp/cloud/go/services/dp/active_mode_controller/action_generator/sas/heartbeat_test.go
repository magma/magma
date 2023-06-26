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
	"fmt"
	"testing"
	"time"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

const (
	someCbsdId  = "some_cbsd_id"
	someGrantId = "some_grant_id"

	granted    = "GRANTED"
	authorized = "AUTHORIZED"

	nextSend          = 1000
	heartbeatInterval = 250

	frequency int64 = 3600e6
	bandwidth int64 = 20e6
)

func TestHeartbeatRequestGenerator(t *testing.T) {
	data := []struct {
		name     string
		grant    *storage.DetailedGrant
		expected *request
	}{{
		name: "Should generate heartbeat immediately when grant is not authorized yet",
		grant: &storage.DetailedGrant{
			Grant: &storage.DBGrant{
				GrantId:                  db.MakeString(someGrantId),
				HeartbeatIntervalSec:     db.MakeInt(heartbeatInterval),
				LastHeartbeatRequestTime: db.MakeTime(time.Unix(nextSend, 0)),
			},
			GrantState: &storage.DBGrantState{
				Name: db.MakeString("granted"),
			},
		},
		expected: getHeartbeatRequest(granted),
	}, {
		name: "Should generate heartbeat when timeout has expired",
		grant: &storage.DetailedGrant{
			Grant: &storage.DBGrant{
				GrantId:                  db.MakeString(someGrantId),
				HeartbeatIntervalSec:     db.MakeInt(heartbeatInterval),
				LastHeartbeatRequestTime: db.MakeTime(time.Unix(nextSend-heartbeatInterval, 0)),
			},
			GrantState: &storage.DBGrantState{
				Name: db.MakeString("authorized"),
			},
		},
		expected: getHeartbeatRequest(authorized),
	}, {
		name: "Should not generate heartbeat request when timeout has not expired yet",
		grant: &storage.DetailedGrant{
			Grant: &storage.DBGrant{
				GrantId:                  db.MakeString(someGrantId),
				HeartbeatIntervalSec:     db.MakeInt(heartbeatInterval),
				LastHeartbeatRequestTime: db.MakeTime(time.Unix(nextSend-heartbeatInterval+1, 0)),
			},
			GrantState: &storage.DBGrantState{
				Name: db.MakeString("authorized"),
			},
		},
		expected: nil,
	}, {
		name: "Should generate relinquish request for unsync grant",
		grant: &storage.DetailedGrant{
			Grant: &storage.DBGrant{
				GrantId: db.MakeString(someGrantId),
			},
			GrantState: &storage.DBGrantState{
				Name: db.MakeString("unsync"),
			},
		},
		expected: getRelinquishmentRequest(),
	}}
	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			p := sas.HeartbeatProcessor{
				NextSendTimestamp: nextSend,
				CbsdId:            someCbsdId,
				Grants: map[int64]*storage.DetailedGrant{
					frequency: tt.grant,
				},
			}
			actual := p.ProcessGrant(frequency, bandwidth)
			assertRequestEqual(t, tt.expected, actual)
		})
	}
}

func getHeartbeatRequest(state string) *request {
	const requestTemplate = `{
	"cbsdId": "%s",
	"grantId": "%s",
	"operationState": "%s"
}`
	payload := fmt.Sprintf(requestTemplate, someCbsdId, someGrantId, state)
	return &request{
		requestType: "heartbeatRequest",
		data:        payload,
	}
}
