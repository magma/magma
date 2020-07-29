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

package protos

import "magma/feg/cloud/go/protos/mconfig"

// ToMconfig creates new mconfig.DiamServerConfig, copies controller diameter
// server config proto to a managed config proto & returns the new mconfig.DiamServerConfig
func (config *DiamServerConfig) ToMconfig() *mconfig.DiamServerConfig {
	return &mconfig.DiamServerConfig{
		Protocol:     config.GetProtocol(),
		Address:      config.GetAddress(),
		LocalAddress: config.GetLocalAddress(),
		DestRealm:    config.GetDestRealm(),
		DestHost:     config.GetDestHost(),
	}
}

// ToMconfig copies diameter client config controller proto to a managed config proto & returns it
func (config *DiamClientConfig) ToMconfig() *mconfig.DiamClientConfig {
	return &mconfig.DiamClientConfig{
		Protocol:         config.GetProtocol(),
		Address:          config.GetAddress(),
		Retransmits:      config.GetRetransmits(),
		WatchdogInterval: config.GetWatchdogInterval(),
		RetryCount:       config.GetRetryCount(),
		LocalAddress:     config.GetLocalAddress(),
		ProductName:      config.GetProductName(),
		Realm:            config.GetRealm(),
		Host:             config.GetHost(),
		DestRealm:        config.GetDestRealm(),
		DestHost:         config.GetDestHost(),
	}
}

// ToMconfig copies controller subscription profile proto to a a new managed config proto & returns it
func (profile *HSSConfig_SubscriptionProfile) ToMconfig() *mconfig.HSSConfig_SubscriptionProfile {
	return &mconfig.HSSConfig_SubscriptionProfile{
		MaxUlBitRate: profile.GetMaxUlBitRate(),
		MaxDlBitRate: profile.GetMaxDlBitRate(),
	}
}

// ToMconfig creates new mconfig.EapAkaConfig_Timeouts, copies config proto to a managed config proto & returns
// the new mconfig.EapAkaConfig_Timeouts
func (config *EapAkaConfig_Timeouts) ToMconfig() *mconfig.EapAkaConfig_Timeouts {
	return &mconfig.EapAkaConfig_Timeouts{
		ChallengeMs:            config.GetChallengeMs(),
		ErrorNotificationMs:    config.GetErrorNotificationMs(),
		SessionMs:              config.GetSessionMs(),
		SessionAuthenticatedMs: config.GetSessionAuthenticatedMs(),
	}
}
