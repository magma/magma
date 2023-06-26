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

type GrantProcessor struct {
	CbsdId   string
	Calc     eirpCalculator
	Channels []storage.Channel
}

type eirpCalculator interface {
	CalcUpperBoundForRange([]storage.Channel, int64, int64) float64
}

func (g *GrantProcessor) ProcessGrant(frequency int64, bandwidth int64) *storage.MutableRequest {
	low := frequency - bandwidth/2
	high := frequency + bandwidth/2
	maxEirp := g.Calc.CalcUpperBoundForRange(g.Channels, low, high)
	payload := &GrantRequest{
		CbsdId: g.CbsdId,
		OperationParam: &OperationParam{
			MaxEirp: maxEirp,
			OperationFrequencyRange: &FrequencyRange{
				LowFrequency:  low,
				HighFrequency: high,
			},
		},
	}
	return makeRequest(Grant, payload)
}

type GrantRequest struct {
	CbsdId         string          `json:"cbsdId"`
	OperationParam *OperationParam `json:"operationParam"`
}

type OperationParam struct {
	MaxEirp                 float64         `json:"maxEirp"`
	OperationFrequencyRange *FrequencyRange `json:"operationFrequencyRange"`
}
