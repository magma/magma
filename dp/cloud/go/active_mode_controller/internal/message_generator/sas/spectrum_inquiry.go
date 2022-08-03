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

type SpectrumInquiryRequestGenerator struct{}

func (*SpectrumInquiryRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
	req := &spectrumInquiryRequest{
		CbsdId: cbsd.GetCbsdId(),
		InquiredSpectrum: []*frequencyRange{{
			LowFrequency:  lowestFrequencyHz,
			HighFrequency: highestFrequencyHz,
		}},
	}
	return []*Request{asRequest(SpectrumInquiry, req)}
}

const (
	lowestFrequencyHz  int64 = 3550 * 1e6
	highestFrequencyHz int64 = 3700 * 1e6
)

type spectrumInquiryRequest struct {
	CbsdId           string            `json:"cbsdId"`
	InquiredSpectrum []*frequencyRange `json:"inquiredSpectrum"`
}

type frequencyRange struct {
	LowFrequency  int64 `json:"lowFrequency"`
	HighFrequency int64 `json:"highFrequency"`
}
