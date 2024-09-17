/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package credit_control

import (
	"strings"

	"magma/lte/cloud/go/protos"
)

func (gsu *GrantedServiceUnit) ToProto() *protos.GrantedUnits {
	if gsu == nil {
		return &protos.GrantedUnits{
			Total: &protos.CreditUnit{IsValid: false},
			Tx:    &protos.CreditUnit{IsValid: false},
			Rx:    &protos.CreditUnit{IsValid: false},
		}
	}
	return &protos.GrantedUnits{
		Total: getCreditUnit(gsu.TotalOctets),
		Tx:    getCreditUnit(gsu.InputOctets),  // Input == Tx == Uplink
		Rx:    getCreditUnit(gsu.OutputOctets), // Output == Rx == Downlink
	}
}

func (gsu *GrantedServiceUnit) IsEmpty() bool {
	return gsu.TotalOctets == nil && gsu.InputOctets == nil && gsu.OutputOctets == nil
}

func getCreditUnit(volume *uint64) *protos.CreditUnit {
	if volume == nil {
		return &protos.CreditUnit{IsValid: false}
	}
	return &protos.CreditUnit{IsValid: true, Volume: *volume}
}

func RemoveIMSIPrefix(imsi string) string {
	return strings.TrimPrefix(imsi, "IMSI")
}

func AddIMSIPrefix(imsi string) string {
	return "IMSI" + imsi
}
