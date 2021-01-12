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

import (
	"sync"
	"time"

	"magma/cwf/cloud/go/protos/mconfig"
	"magma/cwf/gateway/services/gateway_health/metrics"

	"github.com/sparrc/go-ping"
)

// ICMPProbe implements the GRE probe interface
// using ICMP over GRE.
type ICMPProbe struct {
	Endpoints []*mconfig.CwfGatewayHealthConfigGrePeer
	Interval  time.Duration
	PktCount  int

	// Maps endpoints IPs to their current GREEndpointStatus
	endpointStatus map[string]GREEndpointStatus
	sync.RWMutex   // R/W lock synchronizing map endpoint status access
	stop           chan bool
}

const (
	defaultPingTimeout = 5 * time.Second
)

// NewICMPProbe create a new ICMPProbe with the provided endpoints and
// probe interval.
func NewICMPProbe(endpoints []*mconfig.CwfGatewayHealthConfigGrePeer, interval uint32, pktCount int) *ICMPProbe {
	return &ICMPProbe{
		Endpoints:      endpoints,
		Interval:       time.Duration(interval) * time.Second,
		PktCount:       pktCount,
		endpointStatus: map[string]GREEndpointStatus{},
		stop:           make(chan bool),
	}
}

// Start begins the ICMP probes of the ICMPProbe's endpoints.
func (i *ICMPProbe) Start() error {
	var pingers []*ping.Pinger
	for _, endpoint := range i.Endpoints {
		p, err := ping.NewPinger(endpoint.Ip)
		if err != nil {
			return err
		}
		// Need privileged mode to work with docker
		p.SetPrivileged(true)
		pingers = append(pingers, p)
	}
	startProbe := func() {
		for {
			select {
			case <-i.stop:
				return
			default:
				time.Sleep(i.Interval)
				i.executeProbe(pingers)
			}
		}
	}

	go startProbe()
	return nil
}

// Stop stops the ICMP probes of the ICMPProbe's endpoints
func (i *ICMPProbe) Stop() {
	i.stop <- true
	i.Lock()
	i.endpointStatus = map[string]GREEndpointStatus{}
	i.Unlock()
}

// GetStatus returns the current GREEndpointStatus of each endpoint.
func (i *ICMPProbe) GetStatus() *GREProbeStatus {
	var reachable []string
	var unreachable []string
	i.RLock()
	defer i.RUnlock()
	for ip, status := range i.endpointStatus {
		if status == EndpointUnreachable {
			unreachable = append(unreachable, ip)
		} else if status == EndpointReachable {
			reachable = append(reachable, ip)
		}
	}
	return &GREProbeStatus{
		Reachable:   reachable,
		Unreachable: unreachable,
	}
}

func (i *ICMPProbe) executeProbe(pingers []*ping.Pinger) {
	for _, pinger := range pingers {
		pinger.Count = i.PktCount
		pinger.Timeout = defaultPingTimeout
		// reduce frequency of updates by only updating on finish
		pinger.OnFinish = func(stats *ping.Statistics) {
			i.Lock()
			defer i.Unlock()
			if stats.PacketsRecv == 0 {
				i.endpointStatus[stats.Addr] = EndpointUnreachable
				metrics.GreEndpointReachable.WithLabelValues(stats.Addr).Set(0)
			} else {
				i.endpointStatus[stats.Addr] = EndpointReachable
				metrics.GreEndpointReachable.WithLabelValues(stats.Addr).Set(1)
			}
		}
		pinger.Run()
	}
}
