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
	"log"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/s8_proxy/servicers/mock_pgw"

	"github.com/stretchr/testify/assert"
	"github.com/wmnsk/go-gtp/gtpv2"
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
	log.Printf("Running PGW at %s\n", mockPgw.LocalAddr().String())

	// Run S8_proxy
	config := getDefaultConfig(mockPgw.LocalAddr().String())
	s8p, err := NewS8Proxy(config)
	if err != nil {
		t.Fatalf("Error creating S8 proxy +%s", err)
		return
	}

	//------------------------
	//---- Create Session ----
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
			UserPlaneFteid: &protos.Fteid{
				Ipv4Address: "127.0.0.10",
				Ipv6Address: "",
				Teid:        11,
			},
			Qos: &protos.QosInformation{
				Pci:                     0,
				PriorityLevel:           0,
				PreemptionCapability:    0,
				PreemptionVulnerability: 0,
				Qci:                     9,
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

	// check User Plane FTEID was received properly
	assert.Equal(t, mockPgw.LastTEIDu, csRes.BearerContext.UserPlaneFteid.Teid)
	assert.NotEmpty(t, csRes.BearerContext.UserPlaneFteid.Ipv4Address)
	assert.Empty(t, csRes.BearerContext.UserPlaneFteid.Ipv6Address)

	// check Control Plane TEID
	session, err := s8p.gtpClient.GetSessionByIMSI(IMSI1)
	assert.NoError(t, err)
	sessionCteid, err := session.GetTEID(gtpv2.IFTypeS5S8PGWGTPC)
	assert.NoError(t, err)
	assert.Equal(t, mockPgw.LastTEIDc, sessionCteid)

	// check received QOS
	sentQos := csReq.BearerContext.Qos
	receivedQos := mockPgw.LastQos
	assert.Equal(t, sentQos.Gbr.BrDl, receivedQos.Gbr.BrDl)
	assert.Equal(t, sentQos.Gbr.BrUl, receivedQos.Gbr.BrUl)
	assert.Equal(t, sentQos.Mbr.BrDl, receivedQos.Mbr.BrDl)
	assert.Equal(t, sentQos.Mbr.BrUl, receivedQos.Mbr.BrUl)
	assert.Equal(t, sentQos.Qci, receivedQos.Qci)

	//------------------------
	//---- Delete Session ----
	cdReq := &protos.DeleteSessionRequestPgw{Imsi: IMSI1}
	_, err = s8p.DeleteSession(context.Background(), cdReq)
	assert.NoError(t, err)
	// session shouldnt exist anymore
	_, err = s8p.gtpClient.GetSessionByIMSI(IMSI1)
	assert.Error(t, err)

	//------------------------
	//---- Echo Request ----
	eReq := &protos.EchoRequest{}
	_, err = s8p.SendEcho(context.Background(), eReq)
	assert.NoError(t, err)
}

func getDefaultConfig(pgwActualAddrs string) *S8ProxyConfig {
	return &S8ProxyConfig{
		ServerAddr: pgwActualAddrs,
	}
}
