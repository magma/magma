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
	lteprotos "magma/lte/cloud/go/protos"

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
		Sid: &lteprotos.SubscriberID{
			Id:   IMSI1,
			Type: lteprotos.SubscriberID_IMSI,
		},
		MSISDN:               "00111",
		MEI:                  "111",
		MCC:                  "222",
		MNC:                  "333",
		RatType:              0,
		IndicationFlag:       nil,
		BearerId:             5,
		UserPlaneTeid:        0,
		S5S8Ip4UserPane:      "127.0.0.10",
		S5S8Ip6UserPane:      "",
		Apn:                  "internet.com",
		SelectionMode:        "",
		PdnType:              0,
		PdnAddressAllocation: "",
		ApnRestriction:       0,
		AmbrUp:               0,
		AmbrDown:             0,
		Uli:                  "",
	}

	// Send and receive Create Session Request
	_, err = s8p.CreateSession(context.Background(), csReq)
	assert.NoError(t, err)
}

func getDefaultConfig(pgwActualAddrs string) *S8ProxyConfig {
	return &S8ProxyConfig{
		ServerAddr: pgwActualAddrs,
	}
}
