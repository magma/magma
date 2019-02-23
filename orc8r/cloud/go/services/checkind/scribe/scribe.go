/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package scribe

import (
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
	intMsg := formatScribeGwStatusIntMessage(status.Checkin.SystemStatus, status.CertExpirationTime)

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
	normalMsg["magma_package_version"] = checkin.MagmaPkgVersion
	normalMsg["enodeb_configured"] = getValueFromMetaMap("enodeb_configured", serviceStatus)
	normalMsg["enodeb_connected"] = getValueFromMetaMap("enodeb_connected", serviceStatus)
	normalMsg["gps_connected"] = getValueFromMetaMap("gps_connected", serviceStatus)
	normalMsg["gps_latitude"] = getValueFromMetaMap("gps_latitude", serviceStatus)
	normalMsg["gps_longitude"] = getValueFromMetaMap("gps_longitude", serviceStatus)
	normalMsg["mme_connected"] = getValueFromMetaMap("mme_connected", serviceStatus)
	normalMsg["opstate_enabled"] = getValueFromMetaMap("opstate_enabled", serviceStatus)
	normalMsg["ptp_connected"] = getValueFromMetaMap("ptp_connected", serviceStatus)
	normalMsg["rf_tx_on"] = getValueFromMetaMap("rf_tx_on", serviceStatus)
	normalMsg["vpn_ip"] = checkin.VpnIp
	normalMsg["kernel_version"] = checkin.KernelVersion
	normalMsg["kernel_versions_installed"] = strings.Join(checkin.KernelVersionsInstalled, ",")
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

func formatScribeGwStatusIntMessage(systemStatus *protos.SystemStatus, certExpTime int64) map[string]int64 {
	if systemStatus == nil {
		return map[string]int64{"cert_expiration_time": certExpTime}
	}
	return map[string]int64{
		"cpu_idle":             int64(systemStatus.CpuIdle),
		"cpu_system":           int64(systemStatus.CpuSystem),
		"cpu_user":             int64(systemStatus.CpuUser),
		"mem_available":        int64(systemStatus.MemAvailable),
		"mem_free":             int64(systemStatus.MemFree),
		"mem_total":            int64(systemStatus.MemTotal),
		"mem_used":             int64(systemStatus.MemUsed),
		"uptime_secs":          int64(systemStatus.UptimeSecs),
		"cert_expiration_time": certExpTime,
	}
}
