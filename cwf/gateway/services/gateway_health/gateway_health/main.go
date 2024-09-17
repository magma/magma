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

package main

import (
	"flag"
	"net"
	"strings"

	mconfigprotos "magma/cwf/cloud/go/protos/mconfig"
	"magma/cwf/gateway/registry"
	"magma/cwf/gateway/services/gateway_health/health/gre_probe"
	"magma/cwf/gateway/services/gateway_health/health/service_health"
	"magma/cwf/gateway/services/gateway_health/health/system_health"
	"magma/cwf/gateway/services/gateway_health/servicers"
	fegprotos "magma/feg/cloud/go/protos"
	"magma/gateway/mconfig"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

const (
	defaultMemUtilPct       = 0.9
	defaultCpuUtilPct       = 0.9
	defaultGREProbeInterval = 10
	defaultICMPPktCount     = 3
	defaultInterface        = "eth1"
)

func main() {
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.GatewayHealth)
	if err != nil {
		glog.Fatalf("Error creating %s service: %s", registry.GatewayHealth, err)
	}
	cfg := getHealthMconfig()
	probe := gre_probe.NewICMPProbe(cfg.GrePeers, cfg.GreProbeInterval, int(cfg.IcmpProbePktCount))

	transportVIP := cfg.GetClusterVirtualIp()
	if len(transportVIP) == 0 {
		glog.Warningf("No transport VIP has been specified. If running with HA enabled, this is a critical error!")
	} else if _, _, err := net.ParseCIDR(transportVIP); err != nil {
		glog.Fatalf("Transport VIP must be specified with a subnet (i.e. 10.10.10.10/24)")
	}
	systemHealth, err := system_health.NewCWAGSystemHealthProvider(defaultInterface, transportVIP)
	if err != nil {
		glog.Fatalf("Error creating CWAGServiceHealthProvider: %s", err)
	}
	dockerHealth, err := service_health.NewDockerServiceHealthProvider()
	if err != nil {
		glog.Fatalf("Error creating DockerServiceHealthProvider: %s", err)
	}
	servicer := servicers.NewGatewayHealthServicer(cfg, probe, dockerHealth, systemHealth)
	fegprotos.RegisterServiceHealthServer(srv.GrpcServer, servicer)

	// Start GRE probe
	err = probe.Start()
	if err != nil {
		glog.Fatalf("Error running GRE health probe: %s", err)
	}
	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running %s service: %s", registry.GatewayHealth, err)
	}
}

func getHealthMconfig() *mconfigprotos.CwfGatewayHealthConfig {
	ret := &mconfigprotos.CwfGatewayHealthConfig{}
	err := mconfig.GetServiceConfigs(strings.ToLower(registry.GatewayHealth), ret)
	if err != nil {
		ret.CpuUtilThresholdPct = defaultCpuUtilPct
		ret.MemUtilThresholdPct = defaultMemUtilPct
		ret.GreProbeInterval = defaultGREProbeInterval
		ret.IcmpProbePktCount = defaultICMPPktCount
		glog.Errorf("Could not load mconfig. Using defaults: %v", ret)
		return ret
	}
	if ret.CpuUtilThresholdPct == 0 {
		ret.CpuUtilThresholdPct = defaultCpuUtilPct
	}
	if ret.MemUtilThresholdPct == 0 {
		ret.MemUtilThresholdPct = defaultMemUtilPct
	}
	if ret.GreProbeInterval == 0 {
		ret.GreProbeInterval = defaultGREProbeInterval
	}
	ret.GrePeers = removeSubnetEndpoints(ret.GrePeers)
	glog.Infof("Using config: %v", ret)
	return ret
}

// removeSubnetEndpoints removes all endpoints that do not correspond to a
// singular IP address
func removeSubnetEndpoints(endpoints []*mconfigprotos.CwfGatewayHealthConfigGrePeer) []*mconfigprotos.CwfGatewayHealthConfigGrePeer {
	ret := []*mconfigprotos.CwfGatewayHealthConfigGrePeer{}
	for _, endpoint := range endpoints {
		parsedIP, _, err := net.ParseCIDR(endpoint.Ip)
		if err != nil {
			ret = append(ret, endpoint)
			continue
		}
		if !strings.HasSuffix(parsedIP.String(), ".0") {
			parsedGrePeer := &mconfigprotos.CwfGatewayHealthConfigGrePeer{
				Ip: parsedIP.String(),
			}
			ret = append(ret, parsedGrePeer)
			continue
		}
		glog.Infof("Not monitoring GRE peer: %s. Health service only supports monitoring specific (non-subnet) endpoints", endpoint.Ip)
	}
	return ret
}
