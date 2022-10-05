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

package adaptruckus

import (
	"encoding/binary"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2866"
	"fbc/lib/go/radius/ruckus"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestVerifyConversion(t *testing.T) {
	// Arrange
	logger, _ := zap.NewDevelopment()
	require.NotNil(t, logger)
	packet := radius.Packet{
		Attributes: radius.Attributes{},
	}
	acctCtrs := []ruckus.RuckusTCAcctCtrs{
		{
			RuckusAcctCtrsTCName:        "xwf",
			RuckusAcctCtrsInputOctets:   100,
			RuckusAcctCtrsOutputOctets:  100,
			RuckusAcctCtrsInputPackets:  100,
			RuckusAcctCtrsOutputPackets: 100,
		},
		{
			RuckusAcctCtrsTCName:        "internet",
			RuckusAcctCtrsInputOctets:   1234,
			RuckusAcctCtrsOutputOctets:  4321,
			RuckusAcctCtrsInputPackets:  1423,
			RuckusAcctCtrsOutputPackets: 1432,
		},
	}
	ruckus.RuckusTCAcctCtrs_Set(&packet, acctCtrs)

	// Act
	_, err := Handle(
		nil,
		&modules.RequestContext{
			Logger: logger,
		}, &radius.Request{Packet: &packet},
		func(c *modules.RequestContext, r *radius.Request) (*modules.Response, error) {
			r.Get(rfc2866.AcctOutputOctets_Type)
			expect(t, r, rfc2866.AcctInputOctets_Type, 1234)
			expect(t, r, rfc2866.AcctOutputOctets_Type, 4321)
			expect(t, r, rfc2866.AcctInputPackets_Type, 1423)
			expect(t, r, rfc2866.AcctOutputPackets_Type, 1432)
			return nil, nil
		})

	// Assert
	require.Nil(t, err)
}

func expect(t *testing.T, r *radius.Request, attrType radius.Type, expected uint64) {
	// Get the attribute
	attr, ok := r.Lookup(attrType)
	require.True(t, ok)

	// Deserialize to uint64
	val := binary.BigEndian.Uint64(attr)
	require.Equal(t, expected, val)
}
