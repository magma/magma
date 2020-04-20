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

package servicers_test

import (
	"testing"

	"magma/feg/cloud/go/protos/mconfig"
	definitions "magma/feg/gateway/services/s6a_proxy/servicers"
	hss "magma/feg/gateway/services/testcore/hss/servicers"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/stretchr/testify/assert"
)

func TestConstructPermanentFailureAnswer(t *testing.T) {
	msg := diam.NewMessage(diam.AuthenticationInformation, diam.RequestFlag, diam.TGPP_S6A_APP_ID, 1, 2, dict.Default)
	serverCfg := &mconfig.DiamServerConfig{
		DestHost:  "magma_host",
		DestRealm: "magma_realm",
	}
	response := hss.ConstructFailureAnswer(msg, datatype.UTF8String("magma"), serverCfg, 1000)

	assert.Equal(t, msg.Header.CommandCode, response.Header.CommandCode)
	assert.Equal(t, uint8(0), response.Header.CommandFlags)
	assert.Equal(t, uint32(diam.TGPP_S6A_APP_ID), response.Header.ApplicationID)
	assert.Equal(t, uint32(1), response.Header.HopByHopID)
	assert.Equal(t, uint32(2), response.Header.EndToEndID)

	_, err := response.FindAVP(avp.ExperimentalResult, dict.UndefinedVendorID)
	assert.NoError(t, err)

	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1000), aia.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, datatype.DiameterIdentity("magma_host"), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma_realm"), aia.OriginRealm)
	assert.Equal(t, "magma", aia.SessionID)
}

func TestConstructSuccessAnswer(t *testing.T) {
	msg := diam.NewMessage(diam.AuthenticationInformation, diam.RequestFlag, diam.TGPP_S6A_APP_ID, 1, 2, dict.Default)
	serverCfg := &mconfig.DiamServerConfig{
		DestHost:  "magma_host",
		DestRealm: "magma_realm",
	}
	response := hss.ConstructSuccessAnswer(msg, datatype.UTF8String("magma"), serverCfg, diam.TGPP_S6A_APP_ID)

	var aia definitions.AIA
	err := response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, uint32(diam.Success), aia.ResultCode)
	assert.Equal(t, datatype.DiameterIdentity("magma_host"), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma_realm"), aia.OriginRealm)
	assert.Equal(t, "magma", aia.SessionID)
	assert.Equal(t, int32(1), aia.AuthSessionState)
}

func TestAddStandardAnswerAVPS(t *testing.T) {
	msg := diam.NewMessage(diam.AuthenticationInformation, 0, diam.TGPP_S6A_APP_ID, 1, 2, dict.Default)
	serverCfg := &mconfig.DiamServerConfig{
		DestHost:  "magma_host",
		DestRealm: "magma_realm",
	}
	hss.AddStandardAnswerAVPS(msg, "magma", serverCfg, diam.Success)

	var aia definitions.AIA
	err := msg.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, datatype.DiameterIdentity("magma_host"), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma_realm"), aia.OriginRealm)
	assert.Equal(t, "magma", aia.SessionID)
}
