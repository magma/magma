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

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

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
		Packages:      []*Package{&Package{Name: "magma"}},
		VpnIp:         vpnIp,
	}
}

// GetNetworkInfo
func GetNetworkInfo() *NetworkInfo {
	interfaces := make([]*NetworkInterface, len(netInterfaces))
	routes := make([]*Route, len(hostRoutes))
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
	for i, rt := range hostRoutes {
		dest := rt.Destination.IP.To4()
		if dest == nil {
			dest = rt.Destination.IP
		}
		gw := rt.Gateway.To4()
		if gw == nil {
			gw = rt.Gateway
			if len(gw) == 0 {
				if len(dest) == net.IPv4len {
					gw = []byte{0, 0, 0, 0}
				} else {
					gw = net.IP([]byte{0, 0, 0, 0}).To16()
				}
			}
		}
		maskStr := rt.Destination.Mask.String()
		if len(dest) == net.IPv4len {
			maskV4 := net.IP(rt.Destination.Mask).To4()
			if maskV4 != nil {
				maskStr = maskV4.String()
			}
		}
		route := &Route{
			DestinationIp: dest.String(),
			GatewayIp:     gw.String(),
			Genmask:       maskStr,
		}
		if rt.Interface != nil {
			route.NetworkInterfaceId = rt.Interface.Name
		}
		routes[i] = route
	}
	return &NetworkInfo{
		NetworkInterfaces: interfaces,
		RoutingTable:      routes,
	}
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
