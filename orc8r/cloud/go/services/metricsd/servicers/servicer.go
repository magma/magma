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
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
	"magma/orc8r/cloud/go/services/metricsd/exporters"

	"github.com/golang/glog"
	prometheus_proto "github.com/prometheus/client_model/go"
	"golang.org/x/net/context"
)

type MetricsControllerServer struct {
	exporters []exporters.Exporter
}

func NewMetricsControllerServer() *MetricsControllerServer {
	return &MetricsControllerServer{}
}

func (srv *MetricsControllerServer) getNetworkAndGatewayID(
	hardwareID string) (string, string, error) {
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

func (srv *MetricsControllerServer) Collect(
	ctx context.Context,
	in *protos.MetricsContainer) (*protos.Void, error) {

	hardwareID := in.GetGatewayId()
	networkID, gatewayID, err := srv.getNetworkAndGatewayID(hardwareID)
	if err != nil {
		return new(protos.Void), err
	}
	if len(in.Family) != 0 {
		glog.V(2).Infof("collecting %v metrics from gateway %v\n", len(in.Family), in.GatewayId)
	}
	for _, family := range in.GetFamily() {
		for _, e := range srv.exporters {
			context := exporters.MetricsContext{
				MetricName:        protos.GetDecodedName(family),
				NetworkID:         networkID,
				HardwareID:        hardwareID,
				GatewayID:         gatewayID,
				OriginatingEntity: networkID + "." + gatewayID,
				DecodedName:       protos.GetDecodedName(family),
			}
			err := e.Submit(family, context)
			if err != nil {
				glog.Error(err)
			}
		}
	}
	return new(protos.Void), nil
}

// Pulls metrics off the given input channel and sends them to all exporters
// after some preprocessing. Should be run in a goroutine as this blocks
// forever.
func (srv *MetricsControllerServer) ConsumeCloudMetrics(
	inputChan chan *prometheus_proto.MetricFamily,
	hostName string,
) error {
	for family := range inputChan {
		for _, e := range srv.exporters {
			context := exporters.MetricsContext{
				MetricName:        protos.GetDecodedName(family),
				NetworkID:         "cloud",
				GatewayID:         hostName,
				HardwareID:        "cloud",
				OriginatingEntity: "cloud." + hostName,
				DecodedName:       protos.GetDecodedName(family),
			}
			err := e.Submit(family, context)
			if err != nil {
				glog.Errorf("Error submitting metric family to exporter: %s", err)
			}
		}
	}
	return nil
}

func (srv *MetricsControllerServer) RegisterExporter(
	e exporters.Exporter) []exporters.Exporter {
	srv.exporters = append(srv.exporters, e)
	return srv.exporters
}
