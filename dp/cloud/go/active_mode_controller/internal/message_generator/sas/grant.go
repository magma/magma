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
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type GrantProcessor struct {
	CbsdId   string
	Calc     eirpCalculator
	Channels []*active_mode.Channel
}

type eirpCalculator interface {
	CalcUpperBoundForRange([]*active_mode.Channel, int64, int64) float64
}

func (g *GrantProcessor) ProcessGrant(frequency int64, bandwidth int64) *Request {
	low := frequency - bandwidth/2
	high := frequency + bandwidth/2
	maxEirp := g.Calc.CalcUpperBoundForRange(g.Channels, low, high)
	req := &grantRequest{
		CbsdId: g.CbsdId,
		OperationParam: &operationParam{
			MaxEirp: maxEirp,
			OperationFrequencyRange: &frequencyRange{
				LowFrequency:  low,
				HighFrequency: high,
			},
		},
	}
	return asRequest(Grant, req)
}

type grantRequest struct {
	CbsdId         string          `json:"cbsdId"`
	OperationParam *operationParam `json:"operationParam"`
}

type operationParam struct {
	MaxEirp                 float64         `json:"maxEirp"`
	OperationFrequencyRange *frequencyRange `json:"operationFrequencyRange"`
}
