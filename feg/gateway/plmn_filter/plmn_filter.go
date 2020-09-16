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
package plmn_filter

import (
	"github.com/golang/glog"
)

type PlmnIdVals map[string]PlmnIdVal

type PlmnIdVal struct {
	l5 bool
	b6 byte
}

func GetPlmnVals(plmns []string) PlmnIdVals {
	plmnIds := PlmnIdVals{}
	for _, plmnid := range plmns {
		glog.Infof("Adding PLMN ID: %s", plmnid)
		switch len(plmnid) {
		case 5:
			plmnIds[plmnid] = PlmnIdVal{l5: true}
		case 6:
			plmnid5 := plmnid[:5]
			val, _ := plmnIds[plmnid5]
			val.b6 = plmnid[5]
			plmnIds[plmnid5] = val
		default:
			glog.Warningf("Invalid HLR PLMN ID: %s", plmnid)
		}
	}
	return plmnIds
}

// CheckImsiOnPlmnIdListIfAny returns true either if there is no PLMN ID filters (allowlist) configured or
// one the configured PLMN IDs matches passed IMSI
func CheckImsiOnPlmnIdListIfAny(imsi string, plmnIds PlmnIdVals) bool {
	if len(plmnIds) == 0 {
		return true
	}
	val, ok := plmnIds[imsi[:5]]
	return ok && (val.l5 || (len(imsi) > 5 && val.b6 == imsi[5]))
}
