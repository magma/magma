/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import (
	"time"

	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	prom_proto "github.com/prometheus/client_model/go"
	"golang.org/x/net/context"
)

// MetricsControllerServer implements a handler to the gRPC server run by the
// Metrics Controller. It can register instances of the Exporter interface for
// writing to storage
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

	metricsExporters, err := metricsd.GetMetricsExporters()
	if err != nil {
		return &protos.Void{}, err
	}
	for _, e := range metricsExporters {
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
	checkID, err := protos.GetGatewayIdentity(ctx)
	if err != nil {
		return new(protos.Void), err
	}
	if len(checkID.HardwareId) > 0 && checkID.HardwareId != hardwareID {
		glog.Errorf("Expected %s, but found %s as Hardware ID", checkID.HardwareId, hardwareID)
		hardwareID = checkID.HardwareId
	}
	networkID, gatewayID, err := getNetworkAndEntityIDForPhysicalID(hardwareID)
	if err != nil {
		return new(protos.Void), err
	}
	glog.V(2).Infof("collecting %v metrics from gateway %v\n", len(in.Family), in.GatewayId)

	metricsToSubmit := metricsContainerToMetricAndContexts(in, networkID, gatewayID)
	metricsExporters, err := metricsd.GetMetricsExporters()
	if err != nil {
		return &protos.Void{}, err
	}
	for _, e := range metricsExporters {
		err := e.Submit(metricsToSubmit)
		if err != nil {
			glog.Error(err)
		}
	}
	return new(protos.Void), nil
}

// ConsumeCloudMetrics pulls metrics off the given input channel and sends
// them to all exporters after some preprocessing.
// Returns only when inputChan closed, which should never happen.
func (srv *MetricsControllerServer) ConsumeCloudMetrics(inputChan chan *prom_proto.MetricFamily, hostName string) {
	for family := range inputChan {
		metricsToSubmit := preprocessCloudMetrics(family, hostName)
		metricsExporters, err := metricsd.GetMetricsExporters()
		if err != nil {
			glog.Error(err)
			continue
		}
		for _, e := range metricsExporters {
			err := e.Submit([]exporters.MetricAndContext{metricsToSubmit})
			if err != nil {
				glog.Error(err)
			}
		}
	}
	glog.Error("Consume cloud metrics channel unexpectedly closed")
}

func preprocessCloudMetrics(family *prom_proto.MetricFamily, hostName string) exporters.MetricAndContext {
	ctx := exporters.MetricContext{
		MetricName: protos.GetDecodedName(family),
		AdditionalContext: &exporters.CloudMetricContext{
			CloudHost: hostName,
		},
	}
	for _, metric := range family.Metric {
		metric.Label = protos.GetDecodedLabel(metric)
		addLabel(metric, metrics.CloudHostLabelName, hostName)
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
		ctx := exporters.MetricContext{
			MetricName: protos.GetDecodedName(fam),
			AdditionalContext: &exporters.GatewayMetricContext{
				NetworkID: networkID,
				GatewayID: gatewayID,
			},
		}
		for _, metric := range fam.Metric {
			metric.Label = protos.GetDecodedLabel(metric)
			addLabel(metric, metrics.NetworkLabelName, networkID)
			addLabel(metric, metrics.GatewayLabelName, gatewayID)
		}
		ret = append(ret, exporters.MetricAndContext{Family: fam, Context: ctx})
	}
	return ret
}

func pushedMetricsToMetricsAndContext(in *protos.PushedMetricsContainer) []exporters.MetricAndContext {
	ret := make([]exporters.MetricAndContext, 0, len(in.Metrics))
	for _, metric := range in.Metrics {
		ctx := exporters.MetricContext{
			MetricName: metric.MetricName,
			AdditionalContext: &exporters.PushedMetricContext{
				NetworkID: in.NetworkId,
			},
		}

		ts := metric.TimestampMS
		if ts == 0 {
			ts = time.Now().Unix() * 1000
		}

		prometheusLabels := make([]*prom_proto.LabelPair, 0, len(metric.Labels))
		for _, label := range metric.Labels {
			prometheusLabels = append(prometheusLabels, &prom_proto.LabelPair{Name: &label.Name, Value: &label.Value})
		}
		promoMetric := &prom_proto.Metric{
			Label: prometheusLabels,
			Gauge: &prom_proto.Gauge{
				Value: &metric.Value,
			},
			TimestampMs: &ts,
		}
		addLabel(promoMetric, metrics.NetworkLabelName, in.NetworkId)

		gaugeType := prom_proto.MetricType_GAUGE
		fam := &prom_proto.MetricFamily{
			Name:   &metric.MetricName,
			Type:   &gaugeType,
			Metric: []*prom_proto.Metric{promoMetric},
		}
		ret = append(ret, exporters.MetricAndContext{Family: fam, Context: ctx})
	}
	return ret
}

// addLabel ensures that the desired name-value pairing is present in the
// metric's labels.
func addLabel(metric *prom_proto.Metric, labelName, labelValue string) {
	labelAdded := false
	for _, label := range metric.Label {
		if label.GetName() == labelName {
			label.Value = &labelValue
			labelAdded = true
		}
	}
	if !labelAdded {
		metric.Label = append(metric.Label, &prom_proto.LabelPair{Name: strPtr(labelName), Value: &labelValue})
	}
}

func strPtr(s string) *string {
	return &s
}

func getNetworkAndEntityIDForPhysicalID(physicalID string) (string, string, error) {
	if len(physicalID) == 0 {
		return "", "", errors.New("Empty Hardware ID")
	}
	entity, err := configurator.LoadEntityForPhysicalID(physicalID, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return "", "", err
	}
	return entity.NetworkID, entity.Key, nil
}
