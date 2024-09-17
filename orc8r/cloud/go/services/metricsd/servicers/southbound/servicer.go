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

package southbound

import (
	"context"
	"errors"

	"github.com/golang/glog"
	prom_proto "github.com/prometheus/client_model/go"

	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/lib/go/metrics"
	"magma/orc8r/lib/go/protos"
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
	networkID, gatewayID, err := getNetworkAndEntityIDForPhysicalID(ctx, hardwareID)
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

func getNetworkAndEntityIDForPhysicalID(ctx context.Context, physicalID string) (string, string, error) {
	if len(physicalID) == 0 {
		return "", "", errors.New("Empty Hardware ID")
	}
	entity, err := configurator.LoadEntityForPhysicalID(ctx, physicalID, configurator.EntityLoadCriteria{}, serdes.Entity)
	if err != nil {
		return "", "", err
	}
	return entity.NetworkID, entity.Key, nil
}

func metricsContainerToMetricAndContexts(
	in *protos.MetricsContainer,
	networkID, gatewayID string,
) []exporters.MetricAndContext {
	ret := make([]exporters.MetricAndContext, 0, len(in.Family))
	for _, fam := range in.Family {
		ctx := exporters.MetricContext{
			MetricName: fam.GetName(),
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
