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

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

type RelinquishmentRequestGenerator struct{}

func (*RelinquishmentRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
	grants := cbsd.Grants
	cbsdId := cbsd.CbsdId
	reqs := make([]*Request, 0, len(grants))
	for _, grant := range grants {
		req := &relinquishmentRequest{
			CbsdId:  cbsdId,
			GrantId: grant.Id,
		}
		reqs = append(reqs, asRequest(Relinquishment, req))
	}
	return reqs
}

type relinquishmentRequest struct {
	CbsdId  string `json:"cbsdId"`
	GrantId string `json:"grantId"`
}

type RelinquishmentProcessor struct {
	CbsdId string
	Grants map[int64]*active_mode.Grant
}

func (r *RelinquishmentProcessor) ProcessGrant(frequency int64, _ int64) *Request {
	req := &relinquishmentRequest{
		CbsdId:  r.CbsdId,
		GrantId: r.Grants[frequency].Id,
	}
	return asRequest(Relinquishment, req)
}
