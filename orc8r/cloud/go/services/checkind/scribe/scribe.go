/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package scribe

import (
	"encoding/json"
	"fmt"
	"strings"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/logger"

	"github.com/golang/glog"
)

const (
	GATEWAY_STATUS_SCRIBE_CATEGORY = "perfpipe_magma_gateway_status"
)

func LogGatewayStatusToScribe(status *protos.GatewayStatus, networkId string, logicalId string) {
	if status == nil {
		glog.Errorf("GatewayStatus for %v:%v is nil\n", networkId, logicalId)
		return
	}
	normalMsg, intMsg, err := FormatScribeGwStatusMessage(status, networkId, logicalId)
	if err != nil {
		glog.Errorf("Failed to convert gateway status into scribe message: %v\n", err)
		return
	}
	logEntries := []*protos.LogEntry{{
		Category:  GATEWAY_STATUS_SCRIBE_CATEGORY,
		NormalMap: normalMsg,
		IntMap:    intMsg,
		Time:      int64(status.Time / 1000),
	}}
	err = logger.LogToScribeWithSamplingRate(logEntries, 1)
	if err != nil {
		glog.Errorf("Failed to log gateway status to Scribe: %v\n", err)
	}
}

func FormatScribeGwStatusMessage(
	status *protos.GatewayStatus,
	networkId string,
	logicalId string,
) (map[string]string, map[string]int64, error) {
	if status.Checkin == nil {
		return nil, nil, fmt.Errorf("Checkin status is nil")
	}
	normalMsg := formatScribeGwStatusNormalMessage(status.Checkin.Status, status.Checkin, networkId, logicalId)
	intMsg := formatScribeGwStatusIntMessage(status.Checkin.SystemStatus, status.Checkin.GetMachineInfo(), status.CertExpirationTime)

	return normalMsg, intMsg, nil
}

func formatScribeGwStatusNormalMessage(
	serviceStatus *protos.ServiceStatus,
	checkin *protos.CheckinRequest,
	networkId string,
	logicalId string,
) map[string]string {
	normalMsg := map[string]string{}
	normalMsg["network_id"] = networkId
	normalMsg["gateway_id"] = logicalId
	normalMsg["hardware_id"] = checkin.GatewayId
	normalMsg["enodeb_configured"] = getValueFromMetaMap("enodeb_configured", serviceStatus)
	normalMsg["enodeb_connected"] = getValueFromMetaMap("enodeb_connected", serviceStatus)
	normalMsg["gps_connected"] = getValueFromMetaMap("gps_connected", serviceStatus)
	normalMsg["gps_latitude"] = getValueFromMetaMap("gps_latitude", serviceStatus)
	normalMsg["gps_longitude"] = getValueFromMetaMap("gps_longitude", serviceStatus)
	normalMsg["mme_connected"] = getValueFromMetaMap("mme_connected", serviceStatus)
	normalMsg["opstate_enabled"] = getValueFromMetaMap("opstate_enabled", serviceStatus)
	normalMsg["ptp_connected"] = getValueFromMetaMap("ptp_connected", serviceStatus)
	normalMsg["rf_tx_on"] = getValueFromMetaMap("rf_tx_on", serviceStatus)
	normalMsg["magma_package_version"] = getMagmaPkgVersion(checkin)
	normalMsg["vpn_ip"] = getVpnIP(checkin)
	normalMsg["kernel_version"] = getKernelVersion(checkin)
	normalMsg["kernel_versions_installed"] = strings.Join(getKernelVersionsInstalled(checkin), ",")

	if checkin.SystemStatus != nil {
		addSystemStatusNormalMessage(normalMsg, checkin.SystemStatus)
	}
	if checkin.PlatformInfo != nil {
		addPlatformInfoNormalMessage(normalMsg, checkin.PlatformInfo)
	}
	if checkin.MachineInfo != nil {
		addMachineInfoNormalMessage(normalMsg, checkin.MachineInfo)
	}
	return normalMsg
}

func getValueFromMetaMap(key string, serviceStatus *protos.ServiceStatus) string {
	ret := ""
	if serviceStatus == nil {
		return ret
	}
	meta := serviceStatus.Meta
	if meta == nil {
		return ret
	}
	if val, ok := meta[key]; ok {
		return val
	} else {
		return ret
	}
}

func getMagmaPkgVersion(checkin *protos.CheckinRequest) string {
	packages := checkin.GetPlatformInfo().GetPackages()
	if packages != nil {
		for _, pkg := range packages {
			if pkg.Name == "magma" {
				return pkg.Version
			}
		}
	}
	return checkin.MagmaPkgVersion
}

func getVpnIP(checkin *protos.CheckinRequest) string {
	if checkin.PlatformInfo != nil {
		return checkin.PlatformInfo.VpnIp
	}
	return checkin.VpnIp
}

func getKernelVersion(checkin *protos.CheckinRequest) string {
	if checkin.PlatformInfo != nil {
		return checkin.PlatformInfo.KernelVersion
	}
	return checkin.KernelVersion
}

func getKernelVersionsInstalled(checkin *protos.CheckinRequest) []string {
	if checkin.PlatformInfo != nil {
		return checkin.PlatformInfo.KernelVersionsInstalled
	}
	return checkin.KernelVersionsInstalled
}

func addSystemStatusNormalMessage(normalMsg map[string]string, systemStatus *protos.SystemStatus) {
	if systemStatus.GetDiskPartitions() != nil {
		tryAddProtoNormalMessage(normalMsg, "disk_partitions", systemStatus.GetDiskPartitions())
	}
}

func addPlatformInfoNormalMessage(normalMsg map[string]string, platformInfo *protos.PlatformInfo) {
	normalMsg["platform_info.vpn_ip"] = platformInfo.GetVpnIp()
	if platformInfo.GetPackages() != nil {
		tryAddProtoNormalMessage(normalMsg, "platform_info.packages", platformInfo.GetPackages())
	}
	normalMsg["platform_info.kernel_version"] = platformInfo.GetKernelVersion()
	if platformInfo.GetKernelVersionsInstalled() != nil {
		tryAddProtoNormalMessage(normalMsg, "platform_info.kernel_versions_installed", platformInfo.GetKernelVersionsInstalled())
	}
}

func addMachineInfoNormalMessage(normalMsg map[string]string, machineInfo *protos.MachineInfo) {
	normalMsg["machine_info.cpu_info.architecture"] = machineInfo.GetCpuInfo().GetArchitecture()
	normalMsg["machine_info.cpu_info.model_name"] = machineInfo.GetCpuInfo().GetModelName()
	if machineInfo.GetNetworkInfo().GetNetworkInterfaces() != nil {
		tryAddProtoNormalMessage(normalMsg, "machine_info.network_info.network_interfaces", machineInfo.GetNetworkInfo().GetNetworkInterfaces())
	}
	if machineInfo.GetNetworkInfo().GetRoutingTable() != nil {
		tryAddProtoNormalMessage(normalMsg, "machine_info.network_info.routing_table", machineInfo.GetNetworkInfo().GetRoutingTable())
	}
}

func tryAddProtoNormalMessage(normalMsg map[string]string, key string, message interface{}) {
	marshalledJSON, err := json.Marshal(message)
	if err == nil {
		normalMsg[key] = string(marshalledJSON)
	} else {
		glog.Errorf("Failed to marshal json for key %s: %s\n", key, err)
	}
}

func formatScribeGwStatusIntMessage(systemStatus *protos.SystemStatus, machineInfo *protos.MachineInfo, certExpTime int64) map[string]int64 {
	intMessage := map[string]int64{"cert_expiration_time": certExpTime}
	if systemStatus != nil {
		intMessage["cpu_idle"] = int64(systemStatus.CpuIdle)
		intMessage["cpu_system"] = int64(systemStatus.CpuSystem)
		intMessage["cpu_user"] = int64(systemStatus.CpuUser)
		intMessage["mem_available"] = int64(systemStatus.MemAvailable)
		intMessage["mem_free"] = int64(systemStatus.MemFree)
		intMessage["mem_total"] = int64(systemStatus.MemTotal)
		intMessage["mem_used"] = int64(systemStatus.MemUsed)
		intMessage["uptime_secs"] = int64(systemStatus.UptimeSecs)
		intMessage["swap_total"] = int64(systemStatus.SwapTotal)
		intMessage["swap_used"] = int64(systemStatus.SwapUsed)
		intMessage["swap_free"] = int64(systemStatus.SwapFree)
	}
	if machineInfo.GetCpuInfo() != nil {
		intMessage["machine_info.cpu_info.core_count"] = int64(machineInfo.GetCpuInfo().GetCoreCount())
		intMessage["machine_info.cpu_info.threads_per_core"] = int64(machineInfo.GetCpuInfo().GetThreadsPerCore())
	}
	return intMessage
}
