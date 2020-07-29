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
// package status - definition & implementation of a gateway status API
package status

import (
	"net"

	"github.com/emakeev/snowflake"
	"github.com/moriyoshi/routewrapper"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	gopsutil_net "github.com/shirou/gopsutil/net"
)

var (
	hwId          string
	cpusInfo      []cpu.InfoStat
	cpuInfo       cpu.InfoStat
	hostInfo      *host.InfoStat
	netInterfaces []gopsutil_net.InterfaceStat
	disksInfo     []disk.PartitionStat
	hostRoutes    []routewrapper.Route
	bootTime      uint64
	vpnIp         string
	machineInfo   *MachineInfo
	platformInfo  *PlatformInfo
)

func init() {
	cpusInfo, _ = cpu.Info()
	if len(cpusInfo) > 0 {
		cpuInfo = cpusInfo[0]
		// aggregate multiple CPU cores into one (support for legasy reporting)
		cpuInfo.Cores = 0
		for _, ci := range cpusInfo {
			cpuInfo.Cores += ci.Cores
		}
	}
	uuid, _ := snowflake.Get()
	hwId = uuid.String()
	bootTime, _ = host.BootTime()
	hostInfo, _ = host.Info()
	if hostInfo == nil {
		hostInfo = &host.InfoStat{HostID: hwId}
	}
	netInterfaces, _ = gopsutil_net.Interfaces()
	for _, iface := range netInterfaces {
		if iface.Name == "tun0" {
			ifIpv6 := ""
			for _, addr := range iface.Addrs {
				if ip, _, err := net.ParseCIDR(addr.Addr); err == nil {
					if len(ip) <= net.IPv4len {
						vpnIp = ip.String()
						break
					}
					if len(ifIpv6) == 0 {
						ifIpv6 = ip.String()
					}
				}
			}
			if len(vpnIp) == 0 {
				vpnIp = ifIpv6
			}
			break
		}
	}
	disksInfo, _ = disk.Partitions(true)

	w, err := routewrapper.NewRouteWrapper()
	if err == nil {
		hostRoutes, _ = w.Routes()
	}
	machineInfo = GetMachineInfo()
	platformInfo = GetPlatformInfo()
}

func GetHwId() string {
	return hwId
}
