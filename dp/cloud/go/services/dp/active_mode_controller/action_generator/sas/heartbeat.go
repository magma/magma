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

package sas

import (
	"strings"

	"magma/dp/cloud/go/services/dp/storage"
)

type HeartbeatProcessor struct {
	NextSendTimestamp int64
	CbsdId            string
	Grants            map[int64]*storage.DetailedGrant
}

func (h *HeartbeatProcessor) ProcessGrant(frequency int64, _ int64) *storage.MutableRequest {
	grant := h.Grants[frequency]
	if grant.GrantState.Name.String == unsync {
		req := &RelinquishmentRequest{
			CbsdId:  h.CbsdId,
			GrantId: grant.Grant.GrantId.String,
		}
		return makeRequest(Relinquishment, req)
	}
	if grant.GrantState.Name.String == granted ||
		shouldSendNow(grant.Grant, h.NextSendTimestamp) {
		req := &HeartbeatRequest{
			CbsdId:         h.CbsdId,
			GrantId:        grant.Grant.GrantId.String,
			OperationState: strings.ToUpper(grant.GrantState.Name.String),
		}
		return makeRequest(Heartbeat, req)
	}
	return nil
}

type HeartbeatRequest struct {
	CbsdId         string `json:"cbsdId"`
	GrantId        string `json:"grantId"`
	OperationState string `json:"operationState"`
}

const (
	granted = "granted"
	unsync  = "unsync"
)

func shouldSendNow(grant *storage.DBGrant, nextSendTimestamp int64) bool {
	deadline := grant.HeartbeatIntervalSec.Int64 + grant.LastHeartbeatRequestTime.Time.Unix()
	return deadline <= nextSendTimestamp
}
