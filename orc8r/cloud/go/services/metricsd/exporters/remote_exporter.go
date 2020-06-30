package exporters

import (
	"context"

	"magma/orc8r/cloud/go/services/metricsd/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// remoteExporter identifies a remote metrics exporter.
type remoteExporter struct {
	// service name of the exporter
	// should always be uppercase to match service registry convention
	service string
}

func NewRemoteExporter(serviceName string) Exporter {
	return &remoteExporter{service: serviceName}
}

func (r *remoteExporter) Submit(metrics []MetricAndContext) error {
	c, err := r.getExporterClient()
	if err != nil {
		return err
	}
	_, err = c.Submit(context.Background(), &protos.SubmitMetricsRequest{Metrics: MakeProtoMetrics(metrics)})
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

// MakeProtoMetrics converts native contextualized metrics to protobuf.
func MakeProtoMetrics(metrics []MetricAndContext) []*protos.ContextualizedMetric {
	var p []*protos.ContextualizedMetric
	for _, m := range metrics {
		p = append(p, MakeProtoMetric(m))
	}
	return p
}

// MakeProtoMetric converts native contextualized metric to protobuf.
func MakeProtoMetric(m MetricAndContext) *protos.ContextualizedMetric {
	p := &protos.ContextualizedMetric{
		Family:  m.Family,
		Context: &protos.Context{MetricName: m.Context.MetricName},
	}
	switch additionalCtx := m.Context.AdditionalContext.(type) {
	case *CloudMetricContext:
		p.Context.OriginContext = &protos.Context_CloudMetric{
			CloudMetric: &protos.CloudContext{CloudHost: additionalCtx.CloudHost},
		}
	case *GatewayMetricContext:
		p.Context.OriginContext = &protos.Context_GatewayMetric{
			GatewayMetric: &protos.GatewayContext{NetworkId: additionalCtx.NetworkID, GatewayId: additionalCtx.GatewayID},
		}
	case *PushedMetricContext:
		p.Context.OriginContext = &protos.Context_PushedMetric{
			PushedMetric: &protos.PushedContext{NetworkId: additionalCtx.NetworkID},
		}
	}
	return p
}
