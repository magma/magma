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

// Package - description of magma package
type Package struct {
	// name
	Name string `json:"name,omitempty"`
	// version
	Version string `json:"version,omitempty"`
}

// ConfigInfo config info
type ConfigInfo struct {
	// mconfig created at
	MconfigCreatedAt uint64 `json:"mconfig_created_at,omitempty"`
	// latest mconfig file path
	MconfigPath string `json:"mconfig_path,omitempty"`
}

// CpuInfo legacy machine CPU info
type CpuInfo struct {
	// architecture
	Architecture string `json:"architecture,omitempty"`
	// core count
	CoreCount uint64 `json:"core_count,omitempty"`
	// model name
	ModelName string `json:"model_name,omitempty"`
	// threads per core
	ThreadsPerCore uint64 `json:"threads_per_core,omitempty"`
}

// CPU - per CPU information
type CPU struct {
	CpuNumber int32 `json:"cpu_number"`
	// core count
	CoreCount int32 `json:"core_count,omitempty"`
	// model name
	ModelName string  `json:"model_name,omitempty"`
	Mhz       float64 `json:"mhz,omitempty"`
	CacheSize int32   `json:"cache_size,omitempty"`
}

// CpusInfo machine info CPU info
type CpusInfo struct {
	// architecture
	Architecture string `json:"architecture,omitempty"`
	Cpus         []CPU  `json:"cpus,omitempty"`
}

// PlatformInfo platform info
type PlatformInfo struct {
	// config info
	ConfigInfo *ConfigInfo `json:"config_info,omitempty"`
	// kernel version
	KernelVersion string `json:"kernel_version,omitempty"`
	// kernel versions installed
	KernelVersionsInstalled []string `json:"kernel_versions_installed,omitempty"`
	// packages
	Packages []*Package `json:"packages,omitempty"`
	// vpn ip
	VpnIp string `json:"vpn_ip,omitempty"`
}

// Route route
type Route struct {
	// destination ip
	DestinationIp string `json:"destination_ip,omitempty"`
	// gateway ip
	GatewayIp string `json:"gateway_ip,omitempty"`
	// genmask
	Genmask string `json:"genmask,omitempty"`
	// network interface id
	NetworkInterfaceId string `json:"network_interface_id,omitempty"`
}

// NetworkInterface network interface
type NetworkInterface struct {
	// ip addresses
	IpAddresses []string `json:"ip_addresses,omitempty"`
	// ipv6 addresses
	IpV6Addresses []string `json:"ipv6_addresses,omitempty"`
	// mac address
	MacAddress string `json:"mac_address,omitempty"`
	// network interface id
	NetworkInterfaceID string `json:"network_interface_id,omitempty" magma_alt_name:"NetworkInterfaceId"`
	// status
	// [UP DOWN UNKNOWN]
	Status string `json:"status,omitempty"`
}

// NetworkInfo machine info network info
type NetworkInfo struct {
	// network interfaces
	NetworkInterfaces []*NetworkInterface `json:"network_interfaces,omitempty"`
	// routing table
	RoutingTable []*Route `json:"routing_table,omitempty"`
}

// MachineInfo machine info
type MachineInfo struct {
	// cpu info
	CpuInfo *CpuInfo `json:"cpu_info,omitempty"`
	// network info
	NetworkInfo *NetworkInfo `json:"network_info,omitempty"`
}

// DiskPartition disk partition
type DiskPartition struct {
	// Name of the device
	Device string `json:"device,omitempty"`
	// Free disk space of the device in bytes
	Free uint64 `json:"free,omitempty"`
	// Mount point of the device
	MountPoint string `json:"mount_point,omitempty"`
	// Total disk space of the device in bytes
	Total uint64 `json:"total,omitempty"`
	// Used disk space of the device in bytes
	Used        uint64  `json:"used,omitempty"`
	UsedPercent float64 `json:"used_percent,omitempty"`
}

// SystemStatus system status of a gateway
type SystemStatus struct {
	// cpu idle
	CpuIdle uint64 `json:"cpu_idle,omitempty"`
	// cpu system
	CpuSystem uint64 `json:"cpu_system,omitempty"`
	// cpu user
	CpuUser uint64 `json:"cpu_user,omitempty"`
	// disk partitions
	DiskPartitions []*DiskPartition `json:"disk_partitions,omitempty"`
	// mem available
	MemAvailable uint64 `json:"mem_available,omitempty"`
	// mem free
	MemFree uint64 `json:"mem_free,omitempty"`
	// mem total
	MemTotal uint64 `json:"mem_total,omitempty"`
	// mem used
	MemUsed uint64 `json:"mem_used,omitempty"`
	// swap free
	SwapFree uint64 `json:"swap_free,omitempty"`
	// swap total
	SwapTotal uint64 `json:"swap_total,omitempty"`
	// swap used
	SwapUsed uint64 `json:"swap_used,omitempty"`
	// time
	Time uint64 `json:"time,omitempty"`
	// uptime secs
	UptimeSecs uint64 `json:"uptime_secs,omitempty"`
}

// GatewayStatus - gateway status definition
type GatewayStatus struct {
	// cert expiration time milliseconds
	CertExpirationTime int64 `json:"cert_expiration_time,omitempty"`
	// hardware id
	HardwareID string `json:"hardware_id,omitempty"`
	// machine info
	MachineInfo *MachineInfo `json:"machine_info,omitempty"`
	// platform info
	PlatformInfo *PlatformInfo `json:"platform_info,omitempty"`
	// system status
	SystemStatus *SystemStatus `json:"system_status,omitempty"`
	// meta
	Meta map[string]string `json:"meta,omitempty"`
}
