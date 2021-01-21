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

package servicers

import (
	"context"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/s8_proxy/servicers/mock_pgw"

	"github.com/stretchr/testify/assert"
)

const (
	//port 0 means golang will choose the port. Selected port will be injected on getDefaultConfig
	s8proxyAddrs = "127.0.0.1:0" // equivalent to sgwAddrs
	pgwAddrs     = "127.0.0.1:0"
	IMSI1        = "123456789012345"
)

func TestS8Proxy(t *testing.T) {
	// Create and run PGW
	mockPgw, err := mock_pgw.NewStarted(nil, s8proxyAddrs, pgwAddrs)
	if err != nil {
		t.Fatalf("Error creating mock PGW: +%s", err)
		return
	}
	defer mockPgw.Close()
	t.Logf("Running PGW at %s\n", mockPgw.LocalAddr().String())

	// Run S8_proxy
	config := getDefaultConfig(mockPgw.LocalAddr().String())
	s8p, err := NewS8Proxy(config)
	if err != nil {
		t.Fatalf("Error creating S8 proxy +%s", err)
		return
	}

	// Create Session Request message
	csReq := &protos.CreateSessionRequestPgw{
		Imsi:   IMSI1,
		Msisdn: "00111",
		Mei:    "111",
		ServingNetwork: &protos.ServingNetwork{
			Mcc: "222",
			Mnc: "333",
		},
		RatType: 0,
		BearerContext: &protos.BearerContext{
			Id: 5,
			AgwUserPlaneFteid: &protos.Fteid{
				Ipv4Address: "127.0.0.10",
				Ipv6Address: "",
				Teid:        11,
			},
			Qos: &protos.QosInformation{
				Pci:                     0,
				PriorityLevel:           0,
				PreemptionCapability:    0,
				PreemptionVulnerability: 0,
				Qci:                     0,
				Gbr: &protos.Ambr{
					BrUl: 123,
					BrDl: 234,
				},
				Mbr: &protos.Ambr{
					BrUl: 567,
					BrDl: 890,
				},
			},
		},
		PdnType: protos.PDNType_IPV4,
		Paa: &protos.PdnAddressAllocation{
			Ipv4Address: "10.0.0.10",
			Ipv6Address: "",
			Ipv6Prefix:  0,
		},

		Apn:            "internet.com",
		SelectionMode:  "",
		ApnRestriction: 0,
		Ambr: &protos.Ambr{
			BrUl: 999,
			BrDl: 888,
		},
		Uli: &protos.UserLocationInformation{
			Lac:    1,
			Ci:     2,
			Sac:    3,
			Rac:    4,
			Tac:    5,
			Eci:    6,
			MeNbi:  7,
			EMeNbi: 8,
		},
		IndicationFlag: nil,
	}

	// Send and receive Create Session Request
	csRes, err := s8p.CreateSession(context.Background(), csReq)
	assert.NoError(t, err)
	assert.NotEmpty(t, csRes)

	// check fteid was received properly
	assert.NotEqual(t, 0, csRes.PgwFteidU.Teid)
	assert.NotEmpty(t, csRes.PgwFteidU.Ipv4Address)
	assert.Empty(t, csRes.PgwFteidU.Ipv6Address)
}

func getDefaultConfig(pgwActualAddrs string) *S8ProxyConfig {
	return &S8ProxyConfig{
		ServerAddr: pgwActualAddrs,
	}
}
