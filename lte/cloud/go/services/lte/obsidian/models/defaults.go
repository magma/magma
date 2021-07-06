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

import "github.com/go-openapi/swag"

// Move default config function to test_models package
func NewDefaultTDDNetworkConfig() *NetworkCellularConfigs {
	return &NetworkCellularConfigs{
		Ran: &NetworkRanConfigs{
			BandwidthMhz: 20,
			TddConfig: &NetworkRanConfigsTddConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
		},
		Epc: &NetworkEpcConfigs{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf: []byte("\x80\x00"),

			HssRelayEnabled:          swag.Bool(false),
			GxGyRelayEnabled:         swag.Bool(false),
			CloudSubscriberdbEnabled: false,
			CongestionControlEnabled: swag.Bool(true),
			DefaultRuleID:            "",
		},
	}
}

// Move default config function to test_models package
func NewDefaultFDDNetworkConfig() *NetworkCellularConfigs {
	return &NetworkCellularConfigs{
		Ran: &NetworkRanConfigs{
			BandwidthMhz: 20,
			FddConfig: &NetworkRanConfigsFddConfig{
				Earfcndl: 1,
				Earfcnul: 18001,
			},
		},
		Epc: &NetworkEpcConfigs{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:                []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:               []byte("\x80\x00"),
			HssRelayEnabled:          swag.Bool(false),
			GxGyRelayEnabled:         swag.Bool(false),
			CloudSubscriberdbEnabled: false,
			DefaultRuleID:            "",
		},
	}
}

// Move default config function to test_models package
func NewDefaultEnodebStatus() *EnodebState {
	return &EnodebState{
		EnodebConfigured: swag.Bool(true),
		EnodebConnected:  swag.Bool(true),
		GpsConnected:     swag.Bool(true),
		GpsLatitude:      swag.String("1.1"),
		GpsLongitude:     swag.String("2.2"),
		OpstateEnabled:   swag.Bool(true),
		RfTxOn:           swag.Bool(true),
		RfTxDesired:      swag.Bool(false),
		PtpConnected:     swag.Bool(false),
		MmeConnected:     swag.Bool(true),
		FsmState:         swag.String("TEST"),
		IPAddress:        "192.168.0.1",
		UesConnected:     5,
	}
}
