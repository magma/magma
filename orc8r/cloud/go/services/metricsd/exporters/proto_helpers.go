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

package exporters

import (
	"magma/orc8r/cloud/go/services/metricsd/protos"
)

// MakeProtoMetrics converts native contextualized metrics to protos.
func MakeProtoMetrics(metrics []MetricAndContext) []*protos.ContextualizedMetric {
	var ret []*protos.ContextualizedMetric
	for _, m := range metrics {
		ret = append(ret, MakeProtoMetric(m))
	}
	return ret
}

// MakeProtoMetric converts native contextualized metric to proto.
func MakeProtoMetric(m MetricAndContext) *protos.ContextualizedMetric {
	p := &protos.ContextualizedMetric{
		Family:  m.Family,
		Context: &protos.Context{MetricName: m.Context.MetricName},
	}

	switch ctx := m.Context.AdditionalContext.(type) {
	case *CloudMetricContext:
		p.Context.OriginContext = &protos.Context_CloudMetric{
			CloudMetric: &protos.CloudContext{CloudHost: ctx.CloudHost},
		}
	case *GatewayMetricContext:
		p.Context.OriginContext = &protos.Context_GatewayMetric{
			GatewayMetric: &protos.GatewayContext{NetworkId: ctx.NetworkID, GatewayId: ctx.GatewayID},
		}
	case *PushedMetricContext:
		p.Context.OriginContext = &protos.Context_PushedMetric{
			PushedMetric: &protos.PushedContext{NetworkId: ctx.NetworkID},
		}
	}

	return p
}

// MakeNativeMetrics converts protos to native contextualized metrics.
func MakeNativeMetrics(protoMetrics []*protos.ContextualizedMetric) []MetricAndContext {
	var ret []MetricAndContext
	for _, p := range protoMetrics {
		ret = append(ret, MakeNativeMetric(p))
	}
	return ret
}

// MakeNativeMetric converts proto to native contextualized metric.
func MakeNativeMetric(p *protos.ContextualizedMetric) MetricAndContext {
	m := MetricAndContext{
		Family:  p.Family,
		Context: MetricContext{MetricName: p.Context.MetricName},
	}

	switch ctx := p.Context.OriginContext.(type) {
	case *protos.Context_CloudMetric:
		m.Context.AdditionalContext = &CloudMetricContext{CloudHost: ctx.CloudMetric.CloudHost}
	case *protos.Context_GatewayMetric:
		m.Context.AdditionalContext = &GatewayMetricContext{NetworkID: ctx.GatewayMetric.NetworkId, GatewayID: ctx.GatewayMetric.GatewayId}
	case *protos.Context_PushedMetric:
		m.Context.AdditionalContext = &PushedMetricContext{NetworkID: ctx.PushedMetric.NetworkId}
	}

	return m
}
