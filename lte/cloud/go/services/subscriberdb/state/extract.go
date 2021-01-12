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

package state

import (
	"encoding/base64"
	"net"
	"regexp"

	"magma/orc8r/cloud/go/services/state"

	"github.com/pkg/errors"
)

const (
	imsiAPNExpectedMatchCount = 3
)

var (
	// imsiAPNRegex matches the expected IMSI column value of IMSI.APN
	imsiAPNRegex = regexp.MustCompile(`^(IMSI\d+)\.(.+)$`)
)

// GetAssignedIPAddress extracts the IP address from a mobilityd state element.
// We expect something along the lines of:
// {
//   "state": "ALLOCATED",
//   "sid": {"id": "IMSI001010000000001.magma.ipv4"},
//   "ipBlock": {"netAddress": "wKiAAA==", "prefixLen": 24},
//   "ip": {"address": "wKiArg=="}
//  }
// The IP addresses are base64 encoded versions of the packed bytes
func GetAssignedIPAddress(mobilitydState state.ArbitraryJSON) (string, error) {
	ipField, ipExists := mobilitydState["ip"]
	if !ipExists {
		return "", errors.New("no ip field found in mobilityd state")
	}
	ipFieldAsMap, castOK := ipField.(map[string]interface{})
	if !castOK {
		return "", errors.New("could not cast ip field of mobilityd state to arbitrary JSON map type")
	}
	ipAddress, addrExists := ipFieldAsMap["address"]
	if !addrExists {
		return "", errors.New("no IP address found in mobilityd state")
	}
	ipAddressAsString, castOK := ipAddress.(string)
	if !castOK {
		return "", errors.New("encoded IP address is not a string as expected")
	}

	return base64DecodeIPAddress(ipAddressAsString)
}

// GetIMSIAndAPNFromMobilitydStateKey returns the IMSI and APN from the
// mobilityd state key of IMSI.APN.
func GetIMSIAndAPNFromMobilitydStateKey(key string) (string, string, error) {
	matches := imsiAPNRegex.FindStringSubmatch(key)
	if len(matches) != imsiAPNExpectedMatchCount {
		return "", "", errors.Errorf("mobilityd state key %s did not match regex", key)
	}
	imsi, apn := matches[1], matches[2]
	return imsi, apn, nil
}

func base64DecodeIPAddress(encodedIP string) (string, error) {
	ipBytes, err := base64.StdEncoding.DecodeString(encodedIP)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode mobilityd IP address")
	}
	if len(ipBytes) != 4 {
		return "", errors.Errorf("expected IP address to decode to 4 bytes, got %d", len(ipBytes))
	}
	return net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3]).String(), nil
}
