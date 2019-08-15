/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package collection

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/coreos/go-systemd/dbus"
	"github.com/golang/glog"
	"github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs"
)

// DiskUsageMetricCollector is a MetricCollector which return a pair of metric
// families representing the total available disk space on the system and the
// total disk space used, respectively
type DiskUsageMetricCollector struct{}

func (*DiskUsageMetricCollector) GetMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs("/", &fs)
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, fmt.Errorf("Failed to collect disk usage statistics: %s", err)
	}

	all := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	used := all - free

	return []*io_prometheus_client.MetricFamily{
		makeTotalDiskSpaceMetric(all),
		makeUsedDiskSpaceMetric(used),
	}, nil
}

// SystemdStatusMetricCollector is a MetricCollector which queries systemd for
// the status of a given service, returning a single metric family with a gauge
// that has a value of 1 if the service is up and 0 otherwise. The gauge label
// will replace all instances of the "@" symbol with "-".
type SystemdStatusMetricCollector struct {
	ServiceNames []string
}

func (s *SystemdStatusMetricCollector) GetMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	conn, err := dbus.NewSystemdConnection()
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, err
	}
	defer conn.Close()

	systemdStatuses, err := conn.ListUnits()
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, err
	}

	fam := gatherSystemdStatusMetrics(
		s.ServiceNames,
		getSystemdStatusesByUnitName(systemdStatuses),
	)
	return []*io_prometheus_client.MetricFamily{fam}, nil
}

// ProcMetricsCollector is a MetricCollector which queries /proc for
// the number of open file descriptors across all processes running on the
// machine, returning a single metric family for the count.
type ProcMetricsCollector struct{}

func (s *ProcMetricsCollector) GetMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	procs, err := procfs.AllProcs()
	if err != nil {
		return []*io_prometheus_client.MetricFamily{}, err
	}
	totalFds := 0
	for _, proc := range procs {
		numFds, err := proc.FileDescriptorsLen()
		if err != nil {
			return []*io_prometheus_client.MetricFamily{}, err
		}
		totalFds = totalFds + numFds
	}
	return []*io_prometheus_client.MetricFamily{
		makeOpenFileDescriptorsMetric(uint64(totalFds)),
	}, nil
}

// makeTotalDiskSpaceMetric returns a prometheus MetricFamily with a single
// gauge value that indicates how much total disk space (in bytes) the current
// host has.
func makeTotalDiskSpaceMetric(availableSpaceBytes uint64) *io_prometheus_client.MetricFamily {
	name := "disk_total"
	help := "Total disk space on the machine"

	gaugeValue := float64(availableSpaceBytes)
	return MakeSingleGaugeFamily(name, help, nil, gaugeValue)
}

// makeUsedDiskSpaceMetric returns a prometheus MetricFamily with a single
// gauge value that indicates how much total disk space (in bytes) the current
// host has used.
func makeUsedDiskSpaceMetric(usedSpaceBytes uint64) *io_prometheus_client.MetricFamily {
	name := "disk_used"
	help := "Disk space used"

	gaugeValue := float64(usedSpaceBytes)
	return MakeSingleGaugeFamily(name, help, nil, gaugeValue)
}

func getSystemdStatusesByUnitName(systemdStatuses []dbus.UnitStatus) map[string]dbus.UnitStatus {
	ret := map[string]dbus.UnitStatus{}
	for _, stat := range systemdStatuses {
		ret[stat.Name] = stat
	}
	return ret
}

func gatherSystemdStatusMetrics(
	serviceNames []string,
	statusesByName map[string]dbus.UnitStatus,
) *io_prometheus_client.MetricFamily {
	metrics := make(map[MetricLabel]float64, len(serviceNames))

	for _, serviceName := range serviceNames {
		systemdKey := fmt.Sprintf("%s.service", serviceName)
		status, ok := statusesByName[systemdKey]
		if !ok {
			glog.V(5).Infof("Did not get status for unit %s from systemd", serviceName)
			continue
		}

		labelValue := strings.Replace(serviceName, "@", "-", -1)
		metricLabel := MetricLabel{Name: "service_name", Value: labelValue}
		if strings.ToLower(status.ActiveState) == "active" {
			metrics[metricLabel] = 1
		} else {
			metrics[metricLabel] = 0
		}
	}

	name := "systemd_status"
	help := "Status of a systemd service"
	return MakeMultiGaugeFamily(name, help, metrics)
}

// makeOpenFileDescriptorsMetric returns a prometheus MetricFamily with a
// single gauge value that indicates how many file descriptors are currently
// open across all processes on the current host
func makeOpenFileDescriptorsMetric(numFileDescriptors uint64) *io_prometheus_client.MetricFamily {
	name := "num_file_descriptors"
	help := "Total open file descriptors on the machine"

	gaugeValue := float64(numFileDescriptors)
	return MakeSingleGaugeFamily(name, help, nil, gaugeValue)
}
