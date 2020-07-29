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
	"log"
	"regexp"

	"magma/feg/cloud/go/protos/mconfig"
)

func GetVirtualApnRule(apnFilter, apnOverwrite string) (*VirtualApnRule, error) {
	reg, err := regexp.Compile(apnFilter)
	if err != nil {
		return nil, fmt.Errorf("Failed to compile apn filter: %s", err)
	}
	return &VirtualApnRule{
		ApnFilter:    reg,
		ApnOverwrite: apnOverwrite,
	}, nil
}

func GenerateVirtualApnRules(config []*mconfig.VirtualApnRule) []*VirtualApnRule {
	virtualApnConfigs := []*VirtualApnRule{}
	for _, virtualApnlCfg := range config {
		apnRule, err := GetVirtualApnRule(
			virtualApnlCfg.GetApnFilter(),
			virtualApnlCfg.GetApnOverwrite(),
		)
		if err != nil {
			log.Printf("%s Managed Gx Virtual APN Rule Config Load Error: %v", virtualApnlCfg.GetApnFilter(), err)
			continue
		}
		virtualApnConfigs = append(virtualApnConfigs, apnRule)
		log.Printf("Virtual APN Rule Activated filter: %s, Overwrite: %s", virtualApnlCfg.GetApnFilter(), virtualApnlCfg.GetApnOverwrite())
	}
	return virtualApnConfigs
}

func MatchAndGetOverwriteApn(apn string, rules []*VirtualApnRule) string {
	for _, rule := range rules {
		if len(rule.ApnOverwrite) > 0 && rule.ApnFilter.MatchString(apn) {
			return rule.ApnOverwrite
		}
	}
	return apn
}
