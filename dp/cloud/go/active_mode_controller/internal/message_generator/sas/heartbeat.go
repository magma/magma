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

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type HeartbeatProcessor struct {
	NextSendTimestamp int64
	CbsdId            string
	Grants            map[int64]*active_mode.Grant
}

func (h *HeartbeatProcessor) ProcessGrant(frequency int64, _ int64) *Request {
	grant := h.Grants[frequency]
	if grant.State == active_mode.GrantState_Unsync {
		req := &relinquishmentRequest{
			CbsdId:  h.CbsdId,
			GrantId: grant.Id,
		}
		return asRequest(Relinquishment, req)
	}
	if grant.State == active_mode.GrantState_Granted ||
		shouldSendNow(grant, h.NextSendTimestamp) {
		req := &heartbeatRequest{
			CbsdId:         h.CbsdId,
			GrantId:        grant.Id,
			OperationState: strings.ToUpper(grant.State.String()),
		}
		return asRequest(Heartbeat, req)
	}
	return nil
}

type heartbeatRequest struct {
	CbsdId         string `json:"cbsdId"`
	GrantId        string `json:"grantId"`
	OperationState string `json:"operationState"`
}

func shouldSendNow(grant *active_mode.Grant, nextSendTimestamp int64) bool {
	deadline := grant.HeartbeatIntervalSec + grant.LastHeartbeatTimestamp
	return deadline <= nextSendTimestamp
}
