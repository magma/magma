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

package test_utils

import (
	"magma/feg/cloud/go/services/feg/obsidian/models"

	"github.com/go-openapi/swag"
)

func NewDefaultNetworkConfig() *models.NetworkFederationConfigs {
	// GyInitMethod_PER_SESSION
	gyInitMethodPerSession := uint32(1)

	return &models.NetworkFederationConfigs{
		S6a: &models.S6a{
			Server: &models.DiameterClientConfigs{
				Protocol:         "sctp",
				Retransmits:      3,
				WatchdogInterval: 1,
				RetryCount:       5,
				ProductName:      "magma",
				Host:             "magma-fedgw.magma.com",
				Realm:            "magma.com",
			},
			PlmnIds: []string{},
		},
		Gx: &models.Gx{
			DisableGx: swag.Bool(false),
			Servers: []*models.DiameterClientConfigs{
				{
					Protocol:         "tcp",
					Retransmits:      3,
					WatchdogInterval: 1,
					RetryCount:       5,
					ProductName:      "magma",
					Host:             "magma-fedgw.magma.com",
					Realm:            "magma.com",
				},
			},
		},
		Gy: &models.Gy{
			DisableGy: swag.Bool(false),
			Servers: []*models.DiameterClientConfigs{
				{
					Protocol:         "tcp",
					Retransmits:      3,
					WatchdogInterval: 1,
					RetryCount:       5,
					ProductName:      "magma",
					Host:             "magma-fedgw.magma.com",
					Realm:            "magma.com",
				},
			},
			InitMethod: &gyInitMethodPerSession,
		},
		Hss: &models.Hss{
			Server: &models.DiameterServerConfigs{
				Protocol:  "tcp",
				DestHost:  "magma.com",
				DestRealm: "magma.com",
			},
			LteAuthOp:  []byte("EREREREREREREREREREREQ=="),
			LteAuthAmf: []byte("gA"),
			DefaultSubProfile: &models.SubscriptionProfile{
				MaxUlBitRate: 100000000, // 100 Mbps
				MaxDlBitRate: 200000000, // 200 Mbps
			},
			SubProfiles:       make(map[string]models.SubscriptionProfile),
			StreamSubscribers: false,
		},
		Swx: &models.Swx{
			Servers: []*models.DiameterClientConfigs{
				{
					Protocol:         "sctp",
					Retransmits:      3,
					WatchdogInterval: 1,
					RetryCount:       5,
					ProductName:      "magma",
					Host:             "magma-fedgw.magma.com",
					Realm:            "magma.com",
				},
			},
			VerifyAuthorization: false,
			CacheTTLSeconds:     10800,
		},
		EapAka: &models.EapAka{
			Timeout: &models.EapAkaTimeouts{
				ChallengeMs:            20000,
				ErrorNotificationMs:    10000,
				SessionMs:              43200000,
				SessionAuthenticatedMs: 5000,
			},
			PlmnIds: []string{},
		},
		AaaServer: &models.AaaServer{
			IdleSessionTimeoutMs: 21600000,
			AccountingEnabled:    false,
			CreateSessionOnAuth:  false,
		},
		ServedNetworkIds: []string{},
		Health: &models.Health{
			HealthServices:           []string{"S6A_PROXY", "SESSION_PROXY"},
			UpdateIntervalSecs:       10,
			CloudDisablePeriodSecs:   10,
			LocalDisablePeriodSecs:   1,
			UpdateFailureThreshold:   3,
			RequestFailureThreshold:  0.50,
			MinimumRequestThreshold:  1,
			CPUUtilizationThreshold:  0.90,
			MemoryAvailableThreshold: 0.90,
		},

		Csfb: &models.Csfb{
			Client: &models.SctpClientConfigs{
				ServerAddress: "",
				LocalAddress:  "",
			},
		},
	}
}

func NewDefaultGatewayConfig() *models.GatewayFederationConfigs {
	return (*models.GatewayFederationConfigs)(NewDefaultNetworkConfig())
}
