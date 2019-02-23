/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
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
	assert.Equal(t, flow1.Ipv4Src, "1.1.1.0/28")
	assert.Equal(t, flow1.Ipv4Dst, "")

	flow2Desc, err := policydb.GetFlowDescriptionFromFlowString("permit in 17 from 1.1.1.0/28 to 1.1.2.0/32 8000")
	assert.NoError(t, err)
	flow2 := flow2Desc.Match
	assert.Equal(t, flow2.Ipv4Src, "1.1.1.0/28")
	assert.Equal(t, flow2.UdpSrc, uint32(0))
	assert.Equal(t, flow2.Ipv4Dst, "1.1.2.0/32")
	assert.Equal(t, flow2.UdpDst, uint32(8000))
	assert.Equal(t, flow2.TcpDst, uint32(0))

	flow3Desc, err := policydb.GetFlowDescriptionFromFlowString("permit in 6 from 1.1.1.0/28 8000 to 1.1.2.0/32")
	assert.NoError(t, err)
	flow3 := flow3Desc.Match
	assert.Equal(t, flow3.Ipv4Src, "1.1.1.0/28")
	assert.Equal(t, flow3.TcpSrc, uint32(8000))
	assert.Equal(t, flow3.UdpSrc, uint32(0))
	assert.Equal(t, flow3.Ipv4Dst, "1.1.2.0/32")
	assert.Equal(t, flow3.TcpDst, uint32(0))
}

func TestAll(t *testing.T) {
	flow1Desc, err := policydb.GetFlowDescriptionFromFlowString("deny out 17 from any to 1.1.2.0/32 8000")
	assert.NoError(t, err)
	flow1 := flow1Desc.Match
	assert.Equal(t, flow1Desc.Action, protos.FlowDescription_DENY)
	assert.Equal(t, flow1.Direction, protos.FlowMatch_DOWNLINK)
	assert.Equal(t, flow1.IpProto, protos.FlowMatch_IPPROTO_UDP)
	assert.Equal(t, flow1.Ipv4Src, "")
	assert.Equal(t, flow1.TcpSrc, uint32(0))
	assert.Equal(t, flow1.UdpSrc, uint32(0))
	assert.Equal(t, flow1.Ipv4Dst, "1.1.2.0/32")
	assert.Equal(t, flow1.TcpDst, uint32(0))
	assert.Equal(t, flow1.UdpDst, uint32(8000))
}
