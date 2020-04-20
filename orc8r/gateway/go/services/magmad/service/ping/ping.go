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

// package ping implements magmad ping execution functionality
package ping

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"magma/orc8r/lib/go/protos"
)

// regexps for parsing
const decReStr = `\d+(?:\.\d+)?`

var (
	packetReStr = fmt.Sprintf(
		`^(\d+) packets transmitted, (\d+) (?:packets )?received, (%s)%% packet loss`, decReStr)
	rttReStr = fmt.Sprintf(
		`^(?:rtt|round-trip) min/avg/max(?:/mdev|/stddev|) = (%[1]s)/(%[1]s)/(%[1]s).* ms`, decReStr)
	packetLineRe = regexp.MustCompile(packetReStr)
	rttLineRe    = regexp.MustCompile(rttReStr)
	statsRe      = regexp.MustCompile(`^--- ([\S]*)\s?ping statistics ---$`)
)

// ParseResult parses the output of ping command & updates res
func ParseResult(pingOut []byte, res *protos.PingResult) {
	if res == nil {
		return
	}
	res.PacketsTransmitted, res.PacketsReceived, res.AvgResponseMs = 0, 0, 0
	lines := bytes.Split(pingOut, []byte("\n"))
	end := len(lines) - 1
	for n, ln := range lines {
		if matches := statsRe.FindSubmatch(ln); len(matches) > 1 {
			res.HostOrIp = string(matches[1])
			if n == end {
				if len(res.Error) == 0 {
					res.Error = "Not enough output lines in ping output. The ping may have timed out."
				}
				return
			}
			matches = packetLineRe.FindSubmatch(lines[n+1])
			if len(matches) < 3 {
				res.Error = fmt.Sprintf("invalid packets summary line: '%s'", string(lines[n+1]))
				return
			}
			iVal, _ := strconv.Atoi(string(matches[1]))
			res.PacketsTransmitted = int32(iVal)
			iVal, _ = strconv.Atoi(string(matches[2]))
			res.PacketsReceived = int32(iVal)
			if n+1 == end {
				if len(res.Error) == 0 {
					res.Error = "No latency stats in ping output. The ping may have timed out."
				}
				return
			}
			matches = rttLineRe.FindSubmatch(lines[n+2])
			if len(matches) < 3 {
				if len(res.Error) == 0 {
					res.Error = fmt.Sprintf("invalid latency stats line: '%s'", string(lines[n+2]))
				}
				return
			}
			fVal, _ := strconv.ParseFloat(string(matches[2]), 32)
			res.AvgResponseMs = float32(fVal)
			return
		}
	}
	res.Error = "Could not find statistics header in ping output"
}
