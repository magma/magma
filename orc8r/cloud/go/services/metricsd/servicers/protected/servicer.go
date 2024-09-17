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

package protected

import (
	"context"
	"time"

	"github.com/golang/glog"
	prom_proto "github.com/prometheus/client_model/go"

	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"
)

// CloudMetricsControllerServer implements a handler to the gRPC server run by the
// Metrics Controller. It can register instances of the Exporter interface for
// writing to storage
type CloudMetricsControllerServer struct {
	exporters []exporters.Exporter
}

func NewCloudMetricsControllerServer() *CloudMetricsControllerServer {
	return &CloudMetricsControllerServer{}
}

func (srv *CloudMetricsControllerServer) Push(ctx context.Context, in *protos.PushedMetricsContainer) (*protos.Void, error) {
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

func (srv *CloudMetricsControllerServer) PushRaw(ctx context.Context, in *protos.RawMetricsContainer) (*protos.Void, error) {
	for _, family := range in.Families {
		for _, metric := range family.Metric {
			addLabel(metric, "service", in.Service)
		}
		processMetricFamily(family, in.HostName)
	}
	return new(protos.Void), nil
}

// ConsumeCloudMetrics pulls metrics off the given input channel and sends
// them to all exporters after some preprocessing.
// Returns only when inputChan closed, which should never happen.
func (srv *CloudMetricsControllerServer) ConsumeCloudMetrics(inputChan chan *prom_proto.MetricFamily, hostName string) {
	for family := range inputChan {
		processMetricFamily(family, hostName)
	}
	glog.Error("Consume cloud metrics channel unexpectedly closed")
}

func processMetricFamily(family *prom_proto.MetricFamily, hostName string) {
	metricsToSubmit := preprocessCloudMetrics(family, hostName)
	metricsExporters, err := metricsd.GetMetricsExporters()
	if err != nil {
		glog.Error(err)
		return
	}
	for _, e := range metricsExporters {
		err := e.Submit([]exporters.MetricAndContext{metricsToSubmit})
		if err != nil {
			glog.Error(err)
		}
	}
}

func preprocessCloudMetrics(family *prom_proto.MetricFamily, hostName string) exporters.MetricAndContext {
	ctx := exporters.MetricContext{
		MetricName: family.GetName(),
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
