/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// package testutil provides utilities for integration tests
package testlib

import (
	"fmt"
	"strings"
)

func fetchHardwareID(hostname string, bastionIP string) (string, error) {
	out, err := runRemoteCommand(hostname, bastionIP, []string{
		"cat /etc/snowflake",
	})
	if err != nil {
		return "", err
	}

	if len(out) != 1 {
		return "", fmt.Errorf("Unexpected error retrieving hardwareID %s", hostname)
	}
	return strings.TrimSpace(out[0]), nil
}
