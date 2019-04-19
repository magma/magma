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
	"errors"
	"strings"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/metricsd/exporters"

	"github.com/golang/glog"
	prometheus_proto "github.com/prometheus/client_model/go"
	"golang.org/x/net/context"
)

const (
	CloudNetworkLabelName = "networkId"
	CloudGatewayLabelName = "gatewayId"
)

type MetricsControllerServer struct {
	exporters []exporters.Exporter
}

func NewMetricsControllerServer() *MetricsControllerServer {
	return &MetricsControllerServer{}
}

func (srv *MetricsControllerServer) Collect(ctx context.Context, in *protos.MetricsContainer) (*protos.Void, error) {
	if in.Family == nil || len(in.Family) == 0 {
		return new(protos.Void), nil
	}

	hardwareID := in.GetGatewayId()
	networkID, gatewayID, err := srv.getNetworkAndGatewayID(hardwareID)
	if err != nil {
		return new(protos.Void), err
	}
	glog.V(2).Infof("collecting %v metrics from gateway %v\n", len(in.Family), in.GatewayId)

	metricsToSubmit := metricsContainerToMetricAndContexts(in, networkID, hardwareID, gatewayID)
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
func (srv *MetricsControllerServer) ConsumeCloudMetrics(inputChan chan *prometheus_proto.MetricFamily, hostName string) error {
	for family := range inputChan {
		for _, e := range srv.exporters {
			decodedName := protos.GetDecodedName(family)
			networkID, gatewayID := determineCloudNetworkAndGatewayID(family, hostName)
			ctx := exporters.MetricsContext{
				MetricName:        removeCloudMetricLabels(decodedName),
				NetworkID:         networkID,
				GatewayID:         gatewayID,
				OriginatingEntity: networkID + "." + gatewayID,
				DecodedName:       decodedName,
			}
			err := e.Submit([]exporters.MetricAndContext{{Family: family, Context: ctx}})
			if err != nil {
				glog.Error(err)
			}
		}
	}
	return nil
}

func determineCloudNetworkAndGatewayID(family *prometheus_proto.MetricFamily, defaultGateway string) (string, string) {
	// First try to unpack the name from labels in the family name
	networkID, gatewayID := unpackCloudMetricName(family.GetName())
	if networkID != "" && gatewayID != "" {
		return networkID, gatewayID
	}
	// If not found in name, look at the labels of the family's metrics
	labeledNetwork, labeledGateway := getLabeledNetworkAndGatewayID(family)
	if networkID == "" {
		networkID = labeledNetwork
	}
	if gatewayID == "" {
		gatewayID = labeledGateway
	}
	// Finally use default of network: 'cloud' and provided default gateway
	if networkID == "" {
		networkID = exporters.CloudMetricID
	}
	if gatewayID == "" {
		gatewayID = defaultGateway
	}
	return networkID, gatewayID
}

func getLabeledNetworkAndGatewayID(family *prometheus_proto.MetricFamily) (string, string) {
	var networkID, gatewayID string
	for _, metric := range family.GetMetric() {
		for _, label := range metric.GetLabel() {
			if label.GetName() == CloudNetworkLabelName && networkID == "" {
				networkID = label.GetValue()
			}
			if label.GetName() == CloudGatewayLabelName && gatewayID == "" {
				gatewayID = label.GetValue()
			}
		}
	}
	return networkID, gatewayID
}

func (srv *MetricsControllerServer) RegisterExporter(e exporters.Exporter) []exporters.Exporter {
	srv.exporters = append(srv.exporters, e)
	return srv.exporters
}

func (srv *MetricsControllerServer) getNetworkAndGatewayID(hardwareID string) (string, string, error) {
	if len(hardwareID) == 0 {
		return "", "", errors.New("Empty Hardware ID")
	}
	networkID, err := magmad.FindGatewayNetworkId(hardwareID)
	if err != nil {
		return "", "", err
	}
	gatewayID, err := magmad.FindGatewayId(networkID, hardwareID)
	return networkID, gatewayID, err
}

func metricsContainerToMetricAndContexts(
	in *protos.MetricsContainer,
	networkID string, hardwareID string, gatewayID string,
) []exporters.MetricAndContext {
	ret := make([]exporters.MetricAndContext, 0, len(in.Family))
	for _, fam := range in.Family {
		ctx := exporters.MetricsContext{
			MetricName:        protos.GetDecodedName(fam),
			NetworkID:         networkID,
			HardwareID:        hardwareID,
			GatewayID:         gatewayID,
			OriginatingEntity: networkID + "." + gatewayID,
			DecodedName:       protos.GetDecodedName(fam),
		}
		ret = append(ret, exporters.MetricAndContext{Family: fam, Context: ctx})
	}
	return ret
}

// unpackCloudMetricName takes a "cloud" metric name and attempts to parse out
// the networkID and gatewayID from the name. Returns an error if either do not
// exist.
func unpackCloudMetricName(metricName string) (string, string) {
	var networkID, gatewayID string

	networkLabelIndex := strings.Index(metricName, CloudNetworkLabelName)
	gatewayLabelIndex := strings.Index(metricName, CloudGatewayLabelName)
	if gatewayLabelIndex == -1 {
		if networkLabelIndex == -1 {
			return "", ""
		}
		networkStart := networkLabelIndex + len(CloudNetworkLabelName) + 1
		networkID = metricName[networkStart:]
		return networkID, ""
	}

	networkStart := networkLabelIndex + len(CloudNetworkLabelName) + 1
	gatewayStart := gatewayLabelIndex + len(CloudGatewayLabelName) + 1

	gatewayID = metricName[gatewayStart : networkLabelIndex-1]
	networkID = metricName[networkStart:]

	return networkID, gatewayID
}

// removeCloudMetricLabels takes a cloud metric name and removes the networkID
// and gatewayID labels from the name if they exist
func removeCloudMetricLabels(metricName string) string {
	networkLabelIndex := strings.Index(metricName, CloudNetworkLabelName)
	gatewayLabelIndex := strings.Index(metricName, CloudGatewayLabelName)
	if gatewayLabelIndex == -1 {
		if networkLabelIndex == -1 {
			return metricName
		}
		networkStart := networkLabelIndex - 1
		return metricName[:networkStart]
	}

	gatewayStart := gatewayLabelIndex - 1

	return metricName[:gatewayStart]
}
