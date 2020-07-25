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

package protos_test

import (
	"fmt"
	"testing"

	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestSidProto(t *testing.T) {
	str := "IMSI12345"
	pb, err := protos.SidProto(str)
	assert.NoError(t, err)
	assert.True(t, pb.Id == "12345" && pb.Type == protos.SubscriberID_IMSI)

	_, err = protos.SidProto("BLAH12345")
	assert.Error(t, err)
}

func TestSidString(t *testing.T) {
	str := "IMSI12345"
	pb := protos.SubscriberID{Id: "12345"}
	out := protos.SidString(&pb)
	assert.Equal(t, out, str)
}

func TestIMSIandSessionIdParsers(t *testing.T) {
	randomSid := "99999"
	IMSI := "123456789012345"
	IMSInumeric := uint64(123456789012345)
	prefixedIMSI := fmt.Sprintf("IMSI%s", IMSI)
	magmaSessionId := fmt.Sprintf("%s-%s", prefixedIMSI, randomSid)

	// test GetIMSIFromSessionId
	resultIMSINoprefix, err := protos.GetIMSIFromSessionId(magmaSessionId)
	assert.NoError(t, err)
	assert.Equal(t, resultIMSINoprefix, IMSI)

	// test GetIMSIwithPrefixFromSessionId
	resultIMSIWithprefix, err := protos.GetIMSIwithPrefixFromSessionId(magmaSessionId)
	assert.NoError(t, err)
	assert.Equal(t, resultIMSIWithprefix, prefixedIMSI)

	// test StripPrefixFromIMSIandFormat
	resultIMSIstr, resultIMSInumeric, err := protos.StripPrefixFromIMSIandFormat(prefixedIMSI)
	assert.NoError(t, err)
	assert.Equal(t, resultIMSIstr, IMSI)
	assert.Equal(t, resultIMSInumeric, IMSInumeric)
}
