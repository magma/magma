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

	"magma/orc8r/cloud/go/service/client"
	"magma/orc8r/lib/go/protos"

	"github.com/prometheus/client_model/go"
)

// CloudServiceMetricCollector is a MetricCollector which uses service303 to
// collect metrics from a specific cloud service.
type CloudServiceMetricCollector struct {
	service string
}

func NewCloudServiceMetricCollector(service string) MetricCollector {
	return &CloudServiceMetricCollector{service: service}
}

func (c *CloudServiceMetricCollector) GetMetrics() ([]*io_prometheus_client.MetricFamily, error) {
	container, err := client.Service303GetMetrics(c.service)
	serviceName := c.service
	if err != nil {
		return []*io_prometheus_client.MetricFamily{
			makeGetMetricsStatusMetric(serviceName, getMetricsStatusFailure),
		}, fmt.Errorf("Failed to get metrics from service %s: %v", serviceName, err)
	}

	ret := c.postprocessCloudServiceMetrics(container)
	return append(ret, makeGetMetricsStatusMetric(serviceName, getMetricsStatusSuccess)), nil
}

// Appends service name label to all samples
func (c *CloudServiceMetricCollector) postprocessCloudServiceMetrics(container *protos.MetricsContainer) []*io_prometheus_client.MetricFamily {
	for _, fam := range container.Family {
		for _, sample := range fam.Metric {
			labelName := "service"
			labelValue := strings.ToLower(c.service)
			sample.Label = append(
				sample.Label,
				&io_prometheus_client.LabelPair{
					Name:  &labelName,
					Value: &labelValue,
				},
			)
		}
	}
	return container.Family
}

type getMetricsStatus uint8

const (
	getMetricsStatusSuccess getMetricsStatus = 1
	getMetricsStatusFailure getMetricsStatus = 0
)

// makeGetMetricsStatusMetric returns a prometheus MetricFamily with a gauge
// value that indicates that a GetMetrics call to a specific service succeeded
// or failed.
func makeGetMetricsStatusMetric(serviceName string, status getMetricsStatus) *io_prometheus_client.MetricFamily {
	name := "get_metrics_status"
	help := "1 if get_metrics call to service succeeds, 0 if it fails."

	labelName := "serviceName"
	gaugeValue := float64(status)
	return MakeSingleGaugeFamily(name, help, &MetricLabel{Name: labelName, Value: serviceName}, gaugeValue)
}
