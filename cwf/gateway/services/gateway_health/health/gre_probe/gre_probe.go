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

package gre_probe

// GREProbe defines an interface to begin a probe of GRE endpoints
// and fetch that status at a later point.
type GREProbe interface {
	// Start begins the probe of the GRE endpoint(s).
	Start() error

	// Stop stops the probe of the GRE endpoint(s).
	Stop()

	// GetStatus fetches the status of the GRE probe. The GREProbeStatus
	// returned contains slices of reachable and unreachable endpoint IPs.
	GetStatus() *GREProbeStatus
}

type GREProbeStatus struct {
	Reachable   []string
	Unreachable []string
}

type GREEndpointStatus uint

const (
	EndpointReachable   GREEndpointStatus = 0
	EndpointUnreachable GREEndpointStatus = 1
)

type DummyGREProbe struct{}

func (d *DummyGREProbe) Start() error {
	return nil
}

func (d *DummyGREProbe) GetStatus() *GREProbeStatus {
	return &GREProbeStatus{
		Reachable:   []string{},
		Unreachable: []string{},
	}
}
