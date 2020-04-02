/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// MetricsControllerServer implements a handler to the gRPC server run by the
// Metrics Controller. It can register instances of the Exporter interface for
// writing to storage
package servicers

import (
	"time"

	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	prometheusProto "github.com/prometheus/client_model/go"
	"golang.org/x/net/context"
)

type MetricsControllerServer struct {
	exporters []exporters.Exporter
}

func NewMetricsControllerServer() *MetricsControllerServer {
	return &MetricsControllerServer{}
}

func (srv *MetricsControllerServer) Push(ctx context.Context, in *protos.PushedMetricsContainer) (*protos.Void, error) {
	if in.Metrics == nil || len(in.Metrics) == 0 {
		return new(protos.Void), nil
	}

	for _, e := range srv.exporters {
		metricsToSubmit := pushedMetricsToMetricsAndContext(in)
		err := e.Submit(metricsToSubmit)
		if err != nil {
			glog.Error(err)
		}
	}
	return new(protos.Void), nil
}

func (srv *MetricsControllerServer) Collect(ctx context.Context, in *protos.MetricsContainer) (*protos.Void, error) {
	if in.Family == nil || len(in.Family) == 0 {
		return new(protos.Void), nil
	}

	hardwareID := in.GetGatewayId()
	networkID, gatewayID, err := configurator.GetNetworkAndEntityIDForPhysicalID(hardwareID)
	if err != nil {
		return new(protos.Void), err
	}
	glog.V(2).Infof("collecting %v metrics from gateway %v\n", len(in.Family), in.GatewayId)

	metricsToSubmit := metricsContainerToMetricAndContexts(in, networkID, gatewayID)
	for _, e := range srv.exporters {
		err := e.Submit(metricsToSubmit)
		if err != nil {
			glog.Error(err)
		}
	}
	return new(protos.Void), nil
}

// Pulls metrics off the given input channel and sends them to all exporters
// after some preprocessing. Should be run in a goroutine as this blocks
// forever.
func (srv *MetricsControllerServer) ConsumeCloudMetrics(inputChan chan *prometheusProto.MetricFamily, hostName string) error {
	for family := range inputChan {
		metricsToSubmit := preprocessCloudMetrics(family, hostName)
		for _, e := range srv.exporters {
			err := e.Submit([]exporters.MetricAndContext{metricsToSubmit})
			if err != nil {
				glog.Error(err)
			}
		}
	}
	return nil
}

func preprocessCloudMetrics(family *prometheusProto.MetricFamily, hostName string) exporters.MetricAndContext {
	ctx := exporters.MetricsContext{
		MetricName: protos.GetDecodedName(family),
		AdditionalContext: &exporters.CloudMetricContext{
			CloudHost: hostName,
		},
	}
	for _, metric := range family.Metric {
		metric.Label = protos.GetDecodedLabel(metric)
		addRequiredLabelToMetric(metric, metrics.CloudHostLabelName, hostName)
	}
	return exporters.MetricAndContext{Family: family, Context: ctx}
}

func (srv *MetricsControllerServer) RegisterExporter(e exporters.Exporter) []exporters.Exporter {
	srv.exporters = append(srv.exporters, e)
	return srv.exporters
}

func metricsContainerToMetricAndContexts(
	in *protos.MetricsContainer,
	networkID, gatewayID string,
) []exporters.MetricAndContext {
	ret := make([]exporters.MetricAndContext, 0, len(in.Family))
	for _, fam := range in.Family {
		ctx := exporters.MetricsContext{
			MetricName: protos.GetDecodedName(fam),
			AdditionalContext: &exporters.GatewayMetricContext{
				NetworkID: networkID,
				GatewayID: gatewayID,
			},
		}
		for _, metric := range fam.Metric {
			metric.Label = protos.GetDecodedLabel(metric)
			addRequiredLabelToMetric(metric, metrics.NetworkLabelName, networkID)
			addRequiredLabelToMetric(metric, metrics.GatewayLabelName, gatewayID)
		}
		ret = append(ret, exporters.MetricAndContext{Family: fam, Context: ctx})
	}
	return ret
}

func pushedMetricsToMetricsAndContext(in *protos.PushedMetricsContainer) []exporters.MetricAndContext {
	ret := make([]exporters.MetricAndContext, 0, len(in.Metrics))
	for _, metric := range in.Metrics {
		ctx := exporters.MetricsContext{
			MetricName: metric.MetricName,
			AdditionalContext: &exporters.PushedMetricContext{
				NetworkID: in.NetworkId,
			},
		}

		ts := metric.TimestampMS
		if ts == 0 {
			ts = time.Now().Unix() * 1000
		}

		prometheusLabels := make([]*prometheusProto.LabelPair, 0, len(metric.Labels))
		for _, label := range metric.Labels {
			prometheusLabels = append(prometheusLabels, &prometheusProto.LabelPair{Name: &label.Name, Value: &label.Value})
		}
		promoMetric := &prometheusProto.Metric{
			Label: prometheusLabels,
			Gauge: &prometheusProto.Gauge{
				Value: &metric.Value,
			},
			TimestampMs: &ts,
		}
		addRequiredLabelToMetric(promoMetric, metrics.NetworkLabelName, in.NetworkId)

		gaugeType := prometheusProto.MetricType_GAUGE
		fam := &prometheusProto.MetricFamily{
			Name:   &metric.MetricName,
			Type:   &gaugeType,
			Metric: []*prometheusProto.Metric{promoMetric},
		}
		ret = append(ret, exporters.MetricAndContext{Family: fam, Context: ctx})
	}
	return ret
}

func addRequiredLabelToMetric(metric *prometheusProto.Metric, labelName, labelValue string) {
	labelAdded := false
	for _, label := range metric.Label {
		if label.GetName() == labelName {
			label.Value = &labelValue
			labelAdded = true
		}
	}
	if !labelAdded {
		metric.Label = append(metric.Label, &prometheusProto.LabelPair{Name: makeStringPointer(labelName), Value: &labelValue})
	}
}

func makeStringPointer(s string) *string {
	return &s
}
