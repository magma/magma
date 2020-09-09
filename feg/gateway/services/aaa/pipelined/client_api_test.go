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

package pipelined_test

import (
	"strings"
	"testing"

	"magma/feg/gateway/services/aaa/pipelined"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/test/mock_pipelined"
	lte_protos "magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

const (
	IMSI1      = "123456789012345"
	SESSIONID1 = "sessionid0001"
	SESSIONID2 = "sessionid9999"
	imsiPrefix = "IMSI"
)

// TestAccountingBadMACaddress uses a wrong macaddres.
// This shouldn't cause an error.
func TestAAAPipelinedClientWithGoodAPN(t *testing.T) {
	mockPipelined := mock_pipelined.NewRunningPipelined(t)
	aaaCtx := getAAAcontext(SESSIONID1, IMSI1)
	subscriberId := makeSID(SESSIONID1)

	// good apn
	expectedMac := (strings.Split(aaaCtx.GetApn(), ":"))[0]
	expectedApName := (strings.Split(aaaCtx.GetApn(), ":"))[1]

	err := pipelined.AddUeMacFlow(subscriberId, aaaCtx)
	assert.NoError(t, err)
	mock_pipelined.AssertMacFlowInstall(t, mockPipelined)
	mock_pipelined.AssertReceivedApMacAndAddress(t, mockPipelined, expectedMac, expectedApName)

	err = pipelined.DeleteUeMacFlow(subscriberId, aaaCtx)
	assert.NoError(t, err)
	mock_pipelined.AssertIDeleteMacFlow(t, mockPipelined)
	mock_pipelined.AssertReceivedApMacAndAddress(t, mockPipelined, expectedMac, expectedApName)
}

func TestAAAPipelinedClientWithBadAPN(t *testing.T) {
	mockPipelined := mock_pipelined.NewRunningPipelined(t)
	aaaCtx := getAAAcontext(SESSIONID1, IMSI1)
	subscriberId := makeSID(SESSIONID1)

	// bad apn
	aaaCtx.Apn = "98-76-54-AA-BB-C:Wifi-Offload-hotspot20"
	expectedMac := ""
	expectedApName := "98-76-54-AA-BB-C:Wifi-Offload-hotspot20"

	err := pipelined.AddUeMacFlow(subscriberId, aaaCtx)
	assert.NoError(t, err)
	mock_pipelined.AssertReceivedApMacAndAddress(t, mockPipelined, expectedMac, expectedApName)
	mock_pipelined.AssertMacFlowInstall(t, mockPipelined)

	err = pipelined.DeleteUeMacFlow(subscriberId, aaaCtx)
	assert.NoError(t, err)
	mock_pipelined.AssertIDeleteMacFlow(t, mockPipelined)
	mock_pipelined.AssertReceivedApMacAndAddress(t, mockPipelined, expectedMac, expectedApName)
}

func getAAAcontext(sessionId, IMSI string) *protos.Context {
	return &protos.Context{
		SessionId: sessionId,
		Imsi:      IMSI,
		Msisdn:    "0015551234567",
		Apn:       "98-76-54-AA-BB-CC:Wifi-Offload-hotspot20",
		MacAddr:   "12-34-AB-CD-EF-FF",
	}
}

func makeSID(imsi string) *lte_protos.SubscriberID {
	return &lte_protos.SubscriberID{Id: imsiPrefix + imsi, Type: lte_protos.SubscriberID_IMSI}
}
