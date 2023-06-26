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
	"magma/dp/cloud/go/services/dp/storage"
)

type RelinquishmentRequestGenerator struct{}

func (*RelinquishmentRequestGenerator) GenerateRequests(cbsd *storage.DetailedCbsd) []*storage.MutableRequest {
	reqs := make([]*storage.MutableRequest, len(cbsd.Grants))
	for i, grant := range cbsd.Grants {
		payload := &RelinquishmentRequest{
			CbsdId:  cbsd.Cbsd.CbsdId.String,
			GrantId: grant.Grant.GrantId.String,
		}
		reqs[i] = makeRequest(Relinquishment, payload)
	}
	return reqs
}

type RelinquishmentRequest struct {
	CbsdId  string `json:"cbsdId"`
	GrantId string `json:"grantId"`
}

type RelinquishmentProcessor struct {
	CbsdId string
	Grants map[int64]*storage.DetailedGrant
}

func (r *RelinquishmentProcessor) ProcessGrant(frequency int64, _ int64) *storage.MutableRequest {
	payload := &RelinquishmentRequest{
		CbsdId:  r.CbsdId,
		GrantId: r.Grants[frequency].Grant.GrantId.String,
	}
	return makeRequest(Relinquishment, payload)
}
