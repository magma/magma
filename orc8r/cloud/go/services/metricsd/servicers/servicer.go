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

	metricContext := &exporters.PushedMetricContext{
		NetworkID: in.NetworkId,
	}

	for _, e := range metricsExporters {
		metrics := preprocessPushedMetrics(in)
		err := e.Submit(metrics, metricContext)
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

	metricsToSubmit := preprocessGatewayMetrics(in, networkID, gatewayID)
	metricContext := &exporters.GatewayMetricContext{
		NetworkID: networkID,
		GatewayID: gatewayID,
	}
	metricsExporters, err := metricsd.GetMetricsExporters()
	if err != nil {
		return &protos.Void{}, err
	}
	for _, e := range metricsExporters {
		err := e.Submit(metricsToSubmit, metricContext)
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
		metricContext := &exporters.CloudMetricContext{
			CloudHost: hostName,
		}
		metricsExporters, err := metricsd.GetMetricsExporters()
		if err != nil {
			glog.Error(err)
			continue
		}
		for _, e := range metricsExporters {
			err := e.Submit([]*prom_proto.MetricFamily{metricsToSubmit}, metricContext)
			if err != nil {
				glog.Error(err)
			}
		}
	}
	glog.Error("Consume cloud metrics channel unexpectedly closed")
}

func preprocessCloudMetrics(family *prom_proto.MetricFamily, hostName string) *prom_proto.MetricFamily {
	for _, metric := range family.Metric {
		metric.Label = protos.GetDecodedLabel(metric)
		addLabel(metric, metrics.CloudHostLabelName, hostName)
	}
	return family
}

func (srv *MetricsControllerServer) RegisterExporter(e exporters.Exporter) []exporters.Exporter {
	srv.exporters = append(srv.exporters, e)
	return srv.exporters
}

func preprocessGatewayMetrics(in *protos.MetricsContainer, networkID, gatewayID string) []*prom_proto.MetricFamily {
	ret := make([]*prom_proto.MetricFamily, 0, len(in.Family))
	for _, fam := range in.Family {
		for _, metric := range fam.Metric {
			metric.Label = protos.GetDecodedLabel(metric)
			addLabel(metric, metrics.NetworkLabelName, networkID)
			addLabel(metric, metrics.GatewayLabelName, gatewayID)
		}
		ret = append(ret, fam)
	}
	return ret
}

func preprocessPushedMetrics(in *protos.PushedMetricsContainer) []*prom_proto.MetricFamily {
	ret := make([]*prom_proto.MetricFamily, 0, len(in.Metrics))
	now := time.Now().Unix() * 1000
	for _, metric := range in.Metrics {
		ts := metric.TimestampMS
		if ts == 0 {
			ts = now
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
		ret = append(ret, fam)
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
		metric.Label = append(metric.Label, &prom_proto.LabelPair{Name: &labelName, Value: &labelValue})
	}
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
