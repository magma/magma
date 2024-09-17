package exporters

import (
	"context"
	"strings"

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/lib/go/merrors"
	lib_protos "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
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

func (r *remoteExporter) Submit(metrics []MetricAndContext) error {
	c, err := r.getExporterClient()
	if err != nil {
		return err
	}
	_, err = c.Submit(context.Background(), &protos.SubmitMetricsRequest{Metrics: MakeProtoMetrics(metrics)})
	return err
}

func (r *remoteExporter) getExporterClient() (protos.MetricsExporterClient, error) {
	conn, err := registry.GetConnection(r.service, lib_protos.ServiceType_PROTECTED)
	if err != nil {
		initErr := merrors.NewInitError(err, r.service)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewMetricsExporterClient(conn), nil
}
