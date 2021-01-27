package exporters

import (
	"context"
	"strings"

	io_prometheus_client "github.com/prometheus/client_model/go"

	"magma/orc8r/cloud/go/services/metricsd/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// remoteExporter identifies a remote metrics exporter.
type remoteExporter struct {
	// service name of the exporter
	// should always be lowercase to match service registry convention
	service string
}

func NewRemoteExporter(serviceName string) Exporter {
	return &remoteExporter{service: strings.ToLower(serviceName)}
}

func (r *remoteExporter) Submit(metrics []*io_prometheus_client.MetricFamily, ctx MetricContext) error {
	c, err := r.getExporterClient()
	if err != nil {
		return err
	}
	_, err = c.Submit(context.Background(), makeProtoSubmitRequest(metrics, ctx))
	return err
}

func (r *remoteExporter) getExporterClient() (protos.MetricsExporterClient, error) {
	conn, err := registry.GetConnection(r.service)
	if err != nil {
		initErr := merrors.NewInitError(err, r.service)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewMetricsExporterClient(conn), nil
}

func makeProtoSubmitRequest(metrics []*io_prometheus_client.MetricFamily, ctx MetricContext) *protos.SubmitMetricsRequest {
	switch metCtx := ctx.(type) {
	case *GatewayMetricContext:
		return &protos.SubmitMetricsRequest{Metrics: metrics, Context: &protos.SubmitMetricsRequest_GatewayContext{
			GatewayContext: &protos.GatewayContext{
				NetworkId: metCtx.NetworkID,
				GatewayId: metCtx.GatewayID,
			},
		}}
	case *CloudMetricContext:
		return &protos.SubmitMetricsRequest{Metrics: metrics, Context: &protos.SubmitMetricsRequest_CloudContext{
			CloudContext: &protos.CloudContext{
				CloudHost: metCtx.CloudHost,
			},
		}}
	case *PushedMetricContext:
		return &protos.SubmitMetricsRequest{Metrics: metrics, Context: &protos.SubmitMetricsRequest_PushedContext{
			PushedContext: &protos.PushedContext{
				NetworkId: metCtx.NetworkID,
			},
		}}
	default:
		return &protos.SubmitMetricsRequest{}
	}
}
