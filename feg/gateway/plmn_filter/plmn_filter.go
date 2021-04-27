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
	"strings"

	"github.com/golang/glog"
)

type PlmnIdVals map[string]PlmnIdVal

type PlmnIdVal struct {
	l5 bool
	b6 byte
}

func GetPlmnVals(plmnids []string, plmnidModuleName ...string) PlmnIdVals {
	var moduleName string
	plmnIds := PlmnIdVals{}
	if len(plmnidModuleName) > 0 {
		moduleName = strings.TrimSpace(plmnidModuleName[0])
		if len(moduleName) > 0 {
			moduleName += " "
		}
	}
	for _, plmnid := range plmnids {
		glog.Infof("Adding %sPLMN ID: %s", moduleName, plmnid)
		switch len(plmnid) {
		case 5:
			plmnIds[plmnid] = PlmnIdVal{l5: true}
		case 6:
			plmnid5 := plmnid[:5]
			val := plmnIds[plmnid5]
			val.b6 = plmnid[5]
			plmnIds[plmnid5] = val
		default:
			glog.Warningf("Invalid %sPLMN ID: %s", moduleName, plmnid)
		}
	}
	return plmnIds
}

// Check returns true when either the plmnIdFilerTable is empty (no PLMN ID filtering configured)
// or one of the configured PLMN IDs matches passed IMSI
func (plmnIdFilerTable PlmnIdVals) Check(imsi string) bool {
	if len(plmnIdFilerTable) == 0 {
		return true
	}
	val, ok := plmnIdFilerTable[imsi[:5]]
	return ok && (val.l5 || (len(imsi) > 5 && val.b6 == imsi[5]))
}

// CheckImsiOnPlmnIdListIfAny returns true when either the plmnIdFilerTable is empty (no PLMN ID filtering configured)
// or one of the configured PLMN IDs matches passed IMSI
// CheckImsiOnPlmnIdListIfAny is a functional alias to Check()
func CheckImsiOnPlmnIdListIfAny(imsi string, plmnIds PlmnIdVals) bool {
	return plmnIds.Check(imsi)
}
