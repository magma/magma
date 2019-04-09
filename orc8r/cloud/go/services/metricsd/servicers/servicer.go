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
	promo_exp "magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"

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
			ctx := exporters.NewMetricsContext(family, promo_exp.CloudMetricID, promo_exp.CloudMetricID, hostName)
			err := e.Submit([]exporters.MetricAndContext{{Family: family, Context: ctx}})
			if err != nil {
				glog.Error(err)
			}
		}
	}
	return nil
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
		ctx := exporters.NewMetricsContext(fam, networkID, hardwareID, gatewayID)
		ret = append(ret, exporters.MetricAndContext{Family: fam, Context: ctx})
	}
	return ret
}
