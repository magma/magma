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

package policydb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"magma/lte/cloud/go/protos"
)

const (
	// action direction proto from src to dst
	descriptorRegExp = "(.+)\\s+(.+)\\s+(.+)\\s+from\\s+(.+)\\s+to\\s+(.+)"
)

type address struct {
	ip      string
	version protos.IPAddress_IPVersion
	port    uint32
}

// GetFlowDescriptionFromFlowString returns a proto.FlowDescription from a IPFilterRule string
// passed in the Flow-Description AVP. This AVP can have many variations, but follows
// the format:
//    action direction proto from src to dst
// e.g.:
//    permit out ip from 1.2.3.0/24 to any
func GetFlowDescriptionFromFlowString(descriptorStr string) (*protos.FlowDescription, error) {
	split, err := tokenizeString(descriptorStr)
	if err != nil {
		return nil, err
	}
	return getFlowDescriptionFromSplit(split)
}

func tokenizeString(descriptorStr string) ([]string, error) {
	re, err := regexp.Compile(descriptorRegExp)
	if err != nil {
		return nil, err
	}
	matches := re.FindStringSubmatch(descriptorStr)
	if len(matches) == 0 {
		return nil, fmt.Errorf("Invalid format for descriptor %s", descriptorStr)
	}
	return matches[1:], nil
}

func getFlowDescriptionFromSplit(matches []string) (*protos.FlowDescription, error) {
	action, err := parseAction(matches[0])
	if err != nil {
		return nil, err
	}
	direction, err := parseDirection(matches[1])
	if err != nil {
		return nil, err
	}
	proto, err := parseProto(matches[2])
	if err != nil {
		return nil, err
	}
	src, err := parseAddress(matches[3])
	if err != nil {
		return nil, err
	}
	dst, err := parseAddress(matches[4])
	if err != nil {
		return nil, err
	}
	rule := &protos.FlowDescription{
		Action: action,
		Match: &protos.FlowMatch{
			Direction: direction,
			IpProto:   proto,
		},
	}
	if src.ip != "" {
		rule.Match.IpSrc = &protos.IPAddress{
			Version: src.version,
			Address: []byte(src.ip),
		}
	}
	if dst.ip != "" {
		rule.Match.IpDst = &protos.IPAddress{
			Version: dst.version,
			Address: []byte(dst.ip),
		}
	}
	if proto == protos.FlowMatch_IPPROTO_TCP {
		rule.Match.TcpSrc = src.port
		rule.Match.TcpDst = dst.port
	} else if proto == protos.FlowMatch_IPPROTO_UDP {
		rule.Match.UdpSrc = src.port
		rule.Match.UdpDst = dst.port
	}

	return rule, nil
}

// action can be "permit" or "deny"
func parseAction(action string) (protos.FlowDescription_Action, error) {
	if action == "permit" {
		return protos.FlowDescription_PERMIT, nil
	}
	if action == "deny" {
		return protos.FlowDescription_DENY, nil
	}
	return protos.FlowDescription_PERMIT, fmt.Errorf("Unable to parse action %s", action)
}

// direction can be "in" (Uplink) or "out" (Downlink)
func parseDirection(direction string) (protos.FlowMatch_Direction, error) {
	if direction == "in" {
		return protos.FlowMatch_UPLINK, nil
	} else if direction == "out" {
		return protos.FlowMatch_DOWNLINK, nil
	}
	return protos.FlowMatch_UPLINK, fmt.Errorf("Unable to parse direction %s", direction)
}

// proto can be "ip" or a proto number like "6" (TCP)
func parseProto(proto string) (protos.FlowMatch_IPProto, error) {
	if proto == "ip" {
		return protos.FlowMatch_IPPROTO_IP, nil
	}
	protoInt, err := strconv.Atoi(proto)
	if err != nil {
		return protos.FlowMatch_IPPROTO_IP, err
	}
	_, ok := protos.FlowMatch_IPProto_name[int32(protoInt)]
	if !ok {
		return protos.FlowMatch_IPPROTO_IP, err
	}
	return protos.FlowMatch_IPProto(protoInt), nil
}

// address looks like "any" or "1.1.1.1/32 4444"
func parseAddress(addr string) (*address, error) {
	matches := strings.Split(addr, " ")
	if len(matches) < 1 {
		return nil, fmt.Errorf("Invalid format for address %s", addr)
	}
	ipAddr := matches[0]
	version := protos.IPAddress_IPV4
	if ipAddr == "any" {
		return &address{ip: "", version: version, port: 0}, nil
	}

	// Easy check if its an ipv6 addr
	if strings.Count(ipAddr, ":") >= 2 {
		version = protos.IPAddress_IPV6
	}

	if len(matches) < 2 {
		return &address{ip: ipAddr, version: version, port: 0}, nil
	}

	// Don't support port ranges for now
	portInt, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, err
	}
	return &address{ip: ipAddr, version: version, port: uint32(portInt)}, nil
}
