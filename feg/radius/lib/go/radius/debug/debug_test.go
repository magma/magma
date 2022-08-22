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

package debug_test

import (
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/debug"
	. "fbc/lib/go/radius/rfc2865"
	. "fbc/lib/go/radius/rfc2866"
	. "fbc/lib/go/radius/rfc2869"
	. "fbc/lib/go/radius/rfc3162"
)

var secret = []byte(`1234567`)

func TestDumpPacket(t *testing.T) {
	tests := []*struct {
		Packet func() *radius.Packet
		Output []string
	}{
		{
			func() *radius.Packet {
				p := &radius.Packet{
					Code:       radius.CodeAccessRequest,
					Identifier: 33,
					Secret:     secret,
					Attributes: make(radius.Attributes),
				}
				p.Authenticator[0] = 0x01

				UserName_SetString(p, "Tim")
				UserPassword_SetString(p, "12345")
				NASIPAddress_Set(p, net.IPv4(10, 0, 2, 5))
				AcctStatusType_Add(p, 3) // Alive, exists in dictionary file
				AcctStatusType_Add(p, AcctStatusType_Value_InterimUpdate)
				AcctLinkCount_Set(p, 2)
				EventTimestamp_Set(p, time.Date(2018, 5, 13, 11, 55, 10, 0, time.UTC))
				NASIPv6Address_Set(p, net.ParseIP("::1"))
				mac, _ := net.ParseMAC("01:02:03:04:05:06:ff:ff")
				FramedInterfaceID_Set(p, mac)

				return p
			},
			[]string{
				`Access-Request Id 33`,
				`  User-Name = "Tim"`,
				`  User-Password = "12345"`,
				`  NAS-IP-Address = 10.0.2.5`,
				`  Acct-Status-Type = Alive / Interim-Update`,
				`  Acct-Status-Type = Alive / Interim-Update`,
				`  Acct-Link-Count = 2`,
				`  Event-Timestamp = 2018-05-13T11:55:10Z`,
				`  NAS-IPv6-Address = ::1`,
				`  Framed-Interface-Id = 01:02:03:04:05:06:ff:ff`,
			},
		},
	}

	config := &debug.Config{
		Dictionary: debug.IncludedDictionary,
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := tt.Packet()
			result := debug.DumpString(config, p)
			outputStr := strings.Join(tt.Output, "\n")
			if result != outputStr {
				t.Fatalf("\nexpected:\n%s\ngot:\n%s", outputStr, result)
			}
		})
	}
}

func TestDumpRequest(t *testing.T) {
	tests := []*struct {
		Request func() *radius.Request
		Output  []string
	}{
		{
			func() *radius.Request {
				local, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1812")
				remote, _ := net.ResolveUDPAddr("udp4", "10.0.10.3:34521")

				req := &radius.Request{
					LocalAddr:  local,
					RemoteAddr: remote,
					Packet: &radius.Packet{
						Code:       radius.CodeAccessRequest,
						Identifier: 5,
						Secret:     secret,
						Attributes: make(radius.Attributes),
					},
				}

				UserName_SetString(req.Packet, "Tim")
				UserPassword_SetString(req.Packet, "12345")
				NASIPAddress_Set(req.Packet, net.IPv4(10, 0, 2, 5))

				return req
			},
			[]string{
				`Access-Request Id 5 from 10.0.10.3:34521 to 127.0.0.1:1812`,
				`  User-Name = "Tim"`,
				`  User-Password = "12345"`,
				`  NAS-IP-Address = 10.0.2.5`,
			},
		},
	}

	config := &debug.Config{
		Dictionary: debug.IncludedDictionary,
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			req := tt.Request()
			result := debug.DumpRequestString(config, req)
			outputStr := strings.Join(tt.Output, "\n")
			if result != outputStr {
				t.Fatalf("\nexpected:\n%s\ngot:\n%s", outputStr, result)
			}
		})
	}
}
