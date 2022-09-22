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
	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas/frequency"
	"magma/dp/cloud/go/services/dp/storage"
)

type SpectrumInquiryRequestGenerator struct{}

func (*SpectrumInquiryRequestGenerator) GenerateRequests(cbsd *storage.DetailedCbsd) []*storage.MutableRequest {
	payload := &SpectrumInquiryRequest{
		CbsdId: cbsd.Cbsd.CbsdId.String,
		InquiredSpectrum: []*FrequencyRange{{
			LowFrequency:  frequency.LowestHz,
			HighFrequency: frequency.HighestHz,
		}},
	}
	req := makeRequest(SpectrumInquiry, payload)
	return []*storage.MutableRequest{req}
}

type SpectrumInquiryRequest struct {
	CbsdId           string            `json:"cbsdId"`
	InquiredSpectrum []*FrequencyRange `json:"inquiredSpectrum"`
}

type FrequencyRange struct {
	LowFrequency  int64 `json:"lowFrequency"`
	HighFrequency int64 `json:"highFrequency"`
}
