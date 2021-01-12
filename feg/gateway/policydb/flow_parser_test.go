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

package policydb_test

import (
	"testing"

	"magma/feg/gateway/policydb"
	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestFlowAction(t *testing.T) {
	permit, err := policydb.GetFlowDescriptionFromFlowString("permit in ip from any to any")
	assert.NoError(t, err)
	assert.Equal(t, permit.Action, protos.FlowDescription_PERMIT)
	deny, err := policydb.GetFlowDescriptionFromFlowString("deny in ip from any to any")
	assert.NoError(t, err)
	assert.Equal(t, deny.Action, protos.FlowDescription_DENY)
}

func TestFlowDirection(t *testing.T) {
	in, err := policydb.GetFlowDescriptionFromFlowString("permit in ip from any to any")
	assert.NoError(t, err)
	assert.Equal(t, in.Match.Direction, protos.FlowMatch_UPLINK)

	out, err := policydb.GetFlowDescriptionFromFlowString("permit out ip from any to any")
	assert.Equal(t, out.Match.Direction, protos.FlowMatch_DOWNLINK)
	assert.NoError(t, err)
}

func TestFlowProto(t *testing.T) {
	ipProto, err := policydb.GetFlowDescriptionFromFlowString("permit in ip from any to any")
	assert.NoError(t, err)
	assert.Equal(t, ipProto.Match.IpProto, protos.FlowMatch_IPPROTO_IP)

	ipProto2, err := policydb.GetFlowDescriptionFromFlowString("permit in 0 from any to any")
	assert.Equal(t, ipProto2.Match.IpProto, protos.FlowMatch_IPPROTO_IP)
	assert.NoError(t, err)

	udp, err := policydb.GetFlowDescriptionFromFlowString("permit in 17 from any to any")
	assert.NoError(t, err)
	assert.Equal(t, udp.Match.IpProto, protos.FlowMatch_IPPROTO_UDP)
}

func TestFlowAddresses(t *testing.T) {
	flow1Desc, err := policydb.GetFlowDescriptionFromFlowString("permit in ip from 1.1.1.0/28 to any")
	assert.NoError(t, err)
	flow1 := flow1Desc.Match
	assert.Equal(t, flow1.IpSrc.Address, []byte("1.1.1.0/28"))
	assert.Equal(t, flow1.IpDst, (*protos.IPAddress)(nil))

	flow2Desc, err := policydb.GetFlowDescriptionFromFlowString("permit in 17 from 1.1.1.0/28 to 1.1.2.0/32 8000")
	assert.NoError(t, err)
	flow2 := flow2Desc.Match
	assert.Equal(t, flow2.IpSrc.Address, []byte("1.1.1.0/28"))
	assert.Equal(t, flow2.UdpSrc, uint32(0))
	assert.Equal(t, flow2.IpDst.Address, []byte("1.1.2.0/32"))
	assert.Equal(t, flow2.UdpDst, uint32(8000))
	assert.Equal(t, flow2.TcpDst, uint32(0))

	flow3Desc, err := policydb.GetFlowDescriptionFromFlowString("permit in 6 from 1.1.1.0/28 8000 to 1.1.2.0/32")
	assert.NoError(t, err)
	flow3 := flow3Desc.Match
	assert.Equal(t, flow3.IpSrc.Address, []byte("1.1.1.0/28"))
	assert.Equal(t, flow3.IpSrc.Version, protos.IPAddress_IPV4)
	assert.Equal(t, flow3.TcpSrc, uint32(8000))
	assert.Equal(t, flow3.UdpSrc, uint32(0))
	assert.Equal(t, flow3.IpDst.Address, []byte("1.1.2.0/32"))
	assert.Equal(t, flow3.IpDst.Version, protos.IPAddress_IPV4)
	assert.Equal(t, flow3.TcpDst, uint32(0))

	flow4Desc, err := policydb.GetFlowDescriptionFromFlowString("permit in 6 from 2001:db8::/32 8000 to 8522:44e5:595a::c523")
	assert.NoError(t, err)
	flow4 := flow4Desc.Match
	assert.Equal(t, flow4.IpSrc.Address, []byte("2001:db8::/32"))
	assert.Equal(t, flow4.IpSrc.Version, protos.IPAddress_IPV6)
	assert.Equal(t, flow4.TcpSrc, uint32(8000))
	assert.Equal(t, flow4.UdpSrc, uint32(0))
	assert.Equal(t, flow4.IpDst.Address, []byte("8522:44e5:595a::c523"))
	assert.Equal(t, flow4.IpDst.Version, protos.IPAddress_IPV6)
	assert.Equal(t, flow4.TcpDst, uint32(0))

	flow5Desc, err := policydb.GetFlowDescriptionFromFlowString("permit in 6 from b522::10 92 to any")
	assert.NoError(t, err)
	flow5 := flow5Desc.Match
	assert.Equal(t, flow5.IpSrc.Address, []byte("b522::10"))
	assert.Equal(t, flow5.IpSrc.Version, protos.IPAddress_IPV6)
	assert.Equal(t, flow5.TcpSrc, uint32(92))
	assert.Equal(t, flow5.UdpSrc, uint32(0))
	assert.Equal(t, flow5.IpDst, (*protos.IPAddress)(nil))
	assert.Equal(t, flow5.TcpDst, uint32(0))
}

func TestAll(t *testing.T) {
	flow1Desc, err := policydb.GetFlowDescriptionFromFlowString("deny out 17 from any to 1.1.2.0/32 8000")
	assert.NoError(t, err)
	flow1 := flow1Desc.Match
	assert.Equal(t, flow1Desc.Action, protos.FlowDescription_DENY)
	assert.Equal(t, flow1.Direction, protos.FlowMatch_DOWNLINK)
	assert.Equal(t, flow1.IpProto, protos.FlowMatch_IPPROTO_UDP)
	assert.Equal(t, flow1.IpSrc, (*protos.IPAddress)(nil))
	assert.Equal(t, flow1.TcpSrc, uint32(0))
	assert.Equal(t, flow1.UdpSrc, uint32(0))
	assert.Equal(t, flow1.IpDst.Address, []byte("1.1.2.0/32"))
	assert.Equal(t, flow1.TcpDst, uint32(0))
	assert.Equal(t, flow1.UdpDst, uint32(8000))
}
