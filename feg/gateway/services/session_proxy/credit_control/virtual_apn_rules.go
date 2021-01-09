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
	"fmt"
	"regexp"

	"magma/feg/cloud/go/protos/mconfig"

	"github.com/golang/glog"
)

type VirtualApnRule struct {
	ApnFilter                     *regexp.Regexp
	ChargingCharacteristicsFilter *regexp.Regexp
	ApnOverwrite                  string
}

func (v *VirtualApnRule) FromMconfig(pVirtualApnConfig *mconfig.VirtualApnRule) error {
	apnRegex, err := regexp.Compile(pVirtualApnConfig.GetApnFilter())
	if err != nil {
		return fmt.Errorf("Failed to compile apn filter: %s", err)
	}
	ccRegex, err := regexp.Compile(pVirtualApnConfig.GetChargingCharacteristicsFilter())
	if err != nil {
		return fmt.Errorf("Failed to compile charging characteristics filter: %s", err)
	}
	v.ApnFilter = apnRegex
	v.ChargingCharacteristicsFilter = ccRegex
	v.ApnOverwrite = pVirtualApnConfig.GetApnOverwrite()
	return nil
}

func GenerateVirtualApnRules(pConfigs []*mconfig.VirtualApnRule) []*VirtualApnRule {
	virtualApnConfigs := []*VirtualApnRule{}
	for _, pVirtualApnCfg := range pConfigs {
		apnRule := &VirtualApnRule{}
		err := apnRule.FromMconfig(pVirtualApnCfg)
		if err != nil {
			glog.Errorf("%s Managed Virtual APN Rule Config Load Error: %v", pVirtualApnCfg.GetApnFilter(), err)
			continue
		}
		virtualApnConfigs = append(virtualApnConfigs, apnRule)
		glog.Infof("Virtual APN Rule Activated APN filter: %s, Charging Characteristics filter: %s, Overwrite: %s",
			pVirtualApnCfg.GetApnFilter(), pVirtualApnCfg.GetChargingCharacteristicsFilter(), pVirtualApnCfg.GetApnOverwrite())
	}
	return virtualApnConfigs
}

func MatchAndGetOverwriteApn(apn, chargingCharacteristics string, rules []*VirtualApnRule) string {
	for _, rule := range rules {
		if len(rule.ApnOverwrite) > 0 &&
			rule.ApnFilter.MatchString(apn) &&
			rule.ChargingCharacteristicsFilter.MatchString(chargingCharacteristics) {
			glog.Infof("VirtualAPN match found! Mapping apn %s->%s", apn, rule.ApnOverwrite)
			return rule.ApnOverwrite
		}
	}
	return apn
}
