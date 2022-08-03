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
	"math"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/vishvananda/netlink"

	"magma/gateway/config"
	"magma/gateway/mconfig"
	"magma/orc8r/lib/go/security/cert"
)

const GW_CERT_CHECK_INTERVAL = time.Minute * 30

var (
	mu                sync.RWMutex
	certExpirationMs  int64
	nextCertCheckTime time.Time
)

// GetCertExpirationTime returns current GW certificate expiration time in milliseconds
// GetCertExpirationTime refreshes cached cert expiration value every GW_CERT_CHECK_INTERVAL
func GetCertExpirationTime() int64 {
	now := time.Now()
	mu.RLock()
	if now.Before(nextCertCheckTime) {
		defer mu.RUnlock()
		return certExpirationMs
	}
	mu.RUnlock()

	expirationMs := certExpirationMs
	crt, err := cert.LoadCert(config.GetControlProxyConfigs().GwCertFile)
	if err == nil {
		expirationMs = UnixMs(crt.NotAfter)
		mu.Lock()
		certExpirationMs = expirationMs
		nextCertCheckTime = now.Add(GW_CERT_CHECK_INTERVAL)
		mu.Unlock()
	}
	return expirationMs
}

// GetConfigInfo returns mconfig file information
func GetConfigInfo() *ConfigInfo {
	path, fi := mconfig.Info()
	var modTime uint64
	if fi != nil {
		modTime = uint64(UnixMs(fi.ModTime()))
	}
	return &ConfigInfo{
		MconfigCreatedAt: modTime,
		MconfigPath:      path,
	}
}

// GetCpuInfo returns legacy Magma GW CPU Info, please use GetCpusInfo for more accurate information
func GetCpuInfo() *CpuInfo {
	return &CpuInfo{
		Architecture:   runtime.GOARCH,
		CoreCount:      uint64(cpuInfo.Cores),
		ModelName:      cpuInfo.ModelName,
		ThreadsPerCore: 0,
	}
}

// GetCpuInfo returns legacy Magma GW CPU Info, please use GetCpusInfo for more accurate information
func GetCpusInfo() *CpusInfo {
	info := &CpusInfo{
		Architecture: runtime.GOARCH,
		Cpus:         make([]CPU, len(cpusInfo)),
	}
	for i, ci := range cpusInfo {
		cpuInf := &(info.Cpus[i])
		cpuInf.CpuNumber = ci.CPU
		cpuInf.CoreCount = ci.Cores
		cpuInf.ModelName = ci.ModelName
		cpuInf.Mhz = ci.Mhz
		cpuInf.CacheSize = ci.CacheSize
	}
	return info
}

// GetPlatformInfo
func GetPlatformInfo() *PlatformInfo {
	return &PlatformInfo{
		ConfigInfo:    GetConfigInfo(),
		KernelVersion: hostInfo.KernelVersion,
		Packages:      []*Package{{Name: "magma"}},
		VpnIp:         vpnIp,
	}
}

// GetNetworkInfo
func GetNetworkInfo() *NetworkInfo {
	interfaces := make([]*NetworkInterface, len(netInterfaces))
	for i, ni := range netInterfaces {
		// dumb down, interface status
		status := "UNKNOWN"
	statConvertLoop:
		for _, f := range ni.Flags {
			switch strings.ToLower(f) {
			case "up":
				status = "UP"
				break statConvertLoop
			case "down", "disabled":
				status = "DOWN"
				break statConvertLoop
			}
		}
		netIface := &NetworkInterface{
			IpAddresses:        []string{},
			IpV6Addresses:      []string{},
			MacAddress:         ni.HardwareAddr,
			NetworkInterfaceID: ni.Name,
			Status:             status,
		}
		for _, addr := range ni.Addrs {
			if ip, _, err := net.ParseCIDR(addr.Addr); err == nil {
				if ip.To4() == nil {
					netIface.IpV6Addresses = append(netIface.IpV6Addresses, addr.Addr)
					continue
				}
			}
			netIface.IpAddresses = append(netIface.IpAddresses, addr.Addr)
		}
		interfaces[i] = netIface
	}

	return &NetworkInfo{
		NetworkInterfaces: interfaces,
		RoutingTable:      getRoutingTable(hostRoutes),
	}
}

// getRoutingTable resolves the routing table based on the hostRoutes. In the first
// call to getRoutes, we obtain routes that can have an associated source IPs.
// We use the remaining outputs to infer the source IPs of the other hostRoutes
// and extract the desired format with another call to getRoutes.
func getRoutingTable(hostRoutes []netlink.Route) []*Route {
	routes, unlinkedHostRoutes, interfaceToIP := getRoutes(hostRoutes)
	unlinkedHostRoutes = linkRoutes(unlinkedHostRoutes, interfaceToIP)
	additionalRoutes, stillUnlinked, _ := getRoutes(unlinkedHostRoutes)
	if stillUnlinked != nil {
		glog.Warningf("There are entries in the route list which have no resolvable source IP: %+v", stillUnlinked)
	}

	routes = append(routes, additionalRoutes...)
	return routes
}

// getRoutes resolves the routing table based on the hostRoutes. The function
// uses the source IP of the hostRoute to find the matching NetworkInterfaceId.
func getRoutes(hostRoutes []netlink.Route) ([]*Route, []netlink.Route, map[string]string) {
	interfaceToIP := make(map[string]string)
	var unlinkedHostRoutes []netlink.Route
	var routes []*Route

	for _, hostRoute := range hostRoutes {
		src := getSourceIP(hostRoute)

		if src == "" {
			// No source IP is stored in hostRoute, we try to resolve this later
			unlinkedHostRoutes = append(unlinkedHostRoutes, hostRoute)
			continue
		}

		netInterfaceID := getNetInterfaceID(hostRoute.LinkIndex)
		route := &Route{
			DestinationIp:      getDestinationIP(hostRoute),
			GatewayIp:          getGatewayIP(hostRoute),
			Genmask:            getMaskStr(hostRoute),
			NetworkInterfaceId: netInterfaceID,
		}
		routes = append(routes, route)
		interfaceToIP[netInterfaceID] = src
	}
	return routes, unlinkedHostRoutes, interfaceToIP

}

// linkRoutes links routes to their source IP based on the LinkIndex
func linkRoutes(unlinkedHostRoutes []netlink.Route, interfaceToIP map[string]string) []netlink.Route {
	linkIndexToIP := make(map[int]string)
	linkList, _ := netlink.LinkList()
	for _, l := range linkList {
		linkIndexToIP[l.Attrs().Index] = interfaceToIP[l.Attrs().Name]
	}
	for i, uhr := range unlinkedHostRoutes {
		unlinkedHostRoutes[i].Src = net.ParseIP(linkIndexToIP[uhr.LinkIndex])
	}
	return unlinkedHostRoutes
}

// GetMachineInfo
func GetMachineInfo() *MachineInfo {
	return &MachineInfo{
		CpuInfo:     GetCpuInfo(),
		NetworkInfo: GetNetworkInfo(),
	}
}

// GetDiskPartitions
func GetDiskPartitions() []*DiskPartition {
	partitions := make([]*DiskPartition, len(disksInfo))
	for i, dp := range disksInfo {
		part := &DiskPartition{
			Device:     dp.Device,
			MountPoint: dp.Mountpoint,
		}
		if usage, err := disk.Usage(part.MountPoint); err == nil && usage != nil {
			part.Free = usage.Free
			part.Total = usage.Total
			part.Used = usage.Used
			part.UsedPercent = usage.UsedPercent
		}
		partitions[i] = part
	}
	return partitions
}

// GetSystemStatus
func GetSystemStatus() *SystemStatus {
	now := time.Now()
	stat := &SystemStatus{
		DiskPartitions: GetDiskPartitions(),
		Time:           uint64(UnixMs(now)),
		UptimeSecs:     uint64(now.Unix()) - bootTime,
	}
	times, _ := cpu.Times(false)
	if len(times) > 0 {
		stat.CpuIdle, stat.CpuSystem, stat.CpuUser =
			uint64(math.Round(times[0].Idle)), uint64(math.Round(times[0].System)), uint64(math.Round(times[0].User))
	}
	m, _ := mem.VirtualMemory()
	if m != nil {
		stat.MemAvailable, stat.MemFree, stat.MemTotal, stat.MemUsed, stat.SwapFree, stat.SwapTotal, stat.SwapUsed =
			m.Available, m.Free, m.Total, m.Used, m.SwapFree, m.SwapTotal, m.SwapTotal-m.SwapFree
	}
	return stat
}

// GetGatewayStatus
func GetGatewayStatus() *GatewayStatus {
	return &GatewayStatus{
		CertExpirationTime: GetCertExpirationTime(),
		HardwareID:         hwId,
		MachineInfo:        machineInfo,
		PlatformInfo:       platformInfo,
		SystemStatus:       GetSystemStatus(),
	}
}

// UnixMs returns Unix time in milliseconds
func UnixMs(t time.Time) int64 {
	return t.Unix() + int64(t.Nanosecond())/int64(time.Millisecond)
}

func getNetInterfaceID(index int) string {
	for _, ni := range netInterfaces {
		if ni.Index == index {
			return ni.Name
		}
	}
	return ""
}

// getSourceIP returns an empty string as default which getRoutes checks.
func getSourceIP(hostRoute netlink.Route) string {
	if hostRoute.Src == nil {
		return ""
	} else {
		return hostRoute.Src.To4().String()
	}
}

// getDestinationIP defaults to "0.0.0.0/0" if nothing is found in the host route.
func getDestinationIP(hostRoute netlink.Route) string {
	return getDestinationIPNet(hostRoute).IP.To4().String()
}

func getDestinationIPNet(hostRoute netlink.Route) net.IPNet {
	if hostRoute.Dst == nil {
		_, dest, _ := net.ParseCIDR("0.0.0.0/0")
		return *dest
	} else {
		return *hostRoute.Dst
	}
}

// getGatewayIP defaults to "0.0.0.0" if no IP can be resolved.
func getGatewayIP(hostRoute netlink.Route) string {
	gw := hostRoute.Gw.To4()
	if gw == nil {
		gw = hostRoute.Gw
		if len(gw) == 0 {
			gw = []byte{0, 0, 0, 0}
		}
	}
	return gw.String()
}

// getMaskStr Subnet Mask.
func getMaskStr(hostRoute netlink.Route) string {
	maskV4 := net.IP(getDestinationIPNet(hostRoute).Mask).To4()
	if maskV4 != nil {
		return maskV4.String()
	}
	return getDestinationIPNet(hostRoute).Mask.String()
}
