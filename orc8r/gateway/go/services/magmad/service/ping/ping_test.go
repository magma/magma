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

package ping

import (
	"testing"

	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

var (
	pingResult1 = []byte(
		`PING google.com (172.217.5.110): 56 data bytes
64 bytes from 172.217.5.110: icmp_seq=0 ttl=55 time=15.070 ms
64 bytes from 172.217.5.110: icmp_seq=1 ttl=55 time=36.611 ms
64 bytes from 172.217.5.110: icmp_seq=2 ttl=55 time=17.490 ms
64 bytes from 172.217.5.110: icmp_seq=3 ttl=55 time=19.834 ms
64 bytes from 172.217.5.110: icmp_seq=4 ttl=55 time=65.760 ms
64 bytes from 172.217.5.110: icmp_seq=5 ttl=55 time=18.883 ms
64 bytes from 172.217.5.110: icmp_seq=6 ttl=55 time=21.936 ms
64 bytes from 172.217.5.110: icmp_seq=7 ttl=55 time=17.459 ms
64 bytes from 172.217.5.110: icmp_seq=8 ttl=55 time=15.802 ms
64 bytes from 172.217.5.110: icmp_seq=9 ttl=55 time=16.294 ms

--- google.com ping statistics ---
10 packets transmitted, 10 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 15.070/24.514/65.760/14.965 ms
`)
	pingResult2 = []byte(
		`PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.079 ms
64 bytes from 127.0.0.1: icmp_seq=2 ttl=64 time=0.113 ms
64 bytes from 127.0.0.1: icmp_seq=3 ttl=64 time=0.065 ms
64 bytes from 127.0.0.1: icmp_seq=4 ttl=64 time=0.254 ms
64 bytes from 127.0.0.1: icmp_seq=5 ttl=64 time=0.450 ms
64 bytes from 127.0.0.1: icmp_seq=6 ttl=64 time=0.242 ms
64 bytes from 127.0.0.1: icmp_seq=7 ttl=64 time=0.113 ms
64 bytes from 127.0.0.1: icmp_seq=8 ttl=64 time=0.117 ms
64 bytes from 127.0.0.1: icmp_seq=9 ttl=64 time=0.075 ms
64 bytes from 127.0.0.1: icmp_seq=10 ttl=64 time=0.114 ms

--- 127.0.0.1 ping statistics ---
10 packets transmitted, 10 received, 0% packet loss, time 9210ms
rtt min/avg/max/mdev = 0.065/0.162/0.450/0.114 ms
`)
	pingResult3 = []byte(
		`PING google.com (172.217.164.110): 56 data bytes
64 bytes from 172.217.164.110: seq=0 ttl=116 time=12.605 ms
64 bytes from 172.217.164.110: seq=1 ttl=116 time=11.947 ms
64 bytes from 172.217.164.110: seq=2 ttl=116 time=15.435 ms
64 bytes from 172.217.164.110: seq=3 ttl=116 time=15.160 ms

--- google.com ping statistics ---
4 packets transmitted, 4 packets received, 0% packet loss
round-trip min/avg/max = 11.947/13.786/15.435 ms
`)
	pingFailure = []byte(
		`PING google.com (216.58.194.174) 56(84) bytes of data.

--- google.com ping statistics ---
5 packets transmitted, 0 received, 100% packet loss, time 4118ms

`)
)

func TestPingResultParsing(t *testing.T) {
	res := &protos.PingResult{}
	ParseResult(pingResult1, res)
	assert.Empty(t, res.Error)
	assert.Equal(t, int32(10), res.PacketsReceived)
	assert.Equal(t, int32(10), res.PacketsTransmitted)
	assert.Equal(t, float32(24.514), res.AvgResponseMs)
	assert.Equal(t, "google.com", res.HostOrIp)
	ParseResult(pingResult2, res)
	assert.Empty(t, res.Error)
	assert.Equal(t, int32(10), res.PacketsReceived)
	assert.Equal(t, int32(10), res.PacketsTransmitted)
	assert.Equal(t, float32(0.162), res.AvgResponseMs)
	assert.Equal(t, "127.0.0.1", res.HostOrIp)
	ParseResult(pingResult3, res) // OpenWRT
	assert.Empty(t, res.Error)
	assert.Equal(t, int32(4), res.PacketsReceived)
	assert.Equal(t, int32(4), res.PacketsTransmitted)
	assert.Equal(t, float32(13.786), res.AvgResponseMs)
	assert.Equal(t, "google.com", res.HostOrIp)
	ParseResult(pingFailure, res)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, int32(0), res.PacketsReceived)
	assert.Equal(t, int32(5), res.PacketsTransmitted)
	assert.Equal(t, float32(0.), res.AvgResponseMs)
	assert.Equal(t, "google.com", res.HostOrIp)
}
