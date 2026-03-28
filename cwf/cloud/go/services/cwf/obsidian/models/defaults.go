/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import fegmodels "magma/feg/cloud/go/services/feg/obsidian/models"

func NewDefaultNetworkCarrierWifiConfigs() *NetworkCarrierWifiConfigs {
	defaultRuleID := "default_rule_1"
	return &NetworkCarrierWifiConfigs{
		AaaServer: &fegmodels.AaaServer{
			AccountingEnabled:    false,
			CreateSessionOnAuth:  false,
			IdleSessionTimeoutMs: 21600000,
		},
		DefaultRuleID: &defaultRuleID,
		EapAka: &fegmodels.EapAka{
			PlmnIds: []string{},
			Timeout: &fegmodels.EapAkaTimeouts{
				ChallengeMs:            20000,
				ErrorNotificationMs:    10000,
				SessionAuthenticatedMs: 5000,
				SessionMs:              43200000,
			},
		},
		EapSim: &fegmodels.EapSim{
			PlmnIds: []string{},
			Timeout: &fegmodels.EapSimTimeouts{
				ChallengeMs:            20000,
				ErrorNotificationMs:    10000,
				SessionAuthenticatedMs: 5000,
				SessionMs:              43200000,
			},
		},
		NetworkServices: []string{"dpi", "policy_enforcement"},
	}
}

func NewDefaultModifiedNetworkCarrierWifiConfigs() *NetworkCarrierWifiConfigs {
	configs := NewDefaultNetworkCarrierWifiConfigs()
	configs.AaaServer.AccountingEnabled = true
	return configs
}
