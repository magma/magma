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

package collection

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const RadiusMetricPrefix = "radius_"

type MetricAggregateRegistry struct {
	metricAggregateByName map[string]MetricAggregate
}

func NewMetricAggregateRegistry() MetricAggregateRegistry {
	return MetricAggregateRegistry{make(map[string]MetricAggregate)}
}

// Update reads in exported metric families from the radius server's metrics.
// The metrics will be mirrored in the prometheus default registry.
// Gauge and counter metrics are registered in the default registry, and
// metric values are collected and updated via this function.
//
// NOTE: This prepends the metric names before processing them to specify
// that they are coming from the radius server
func (r *MetricAggregateRegistry) Update(metricFamilies map[string]*dto.MetricFamily) {
	for metricName, metricFamily := range metricFamilies {
		// We want to mark the metric as coming from the radius server
		// Some metrics already have this prefix, but others need it added
		if !strings.HasPrefix(metricName, RadiusMetricPrefix) {
			metricName = fmt.Sprintf("%s%s", RadiusMetricPrefix, metricName)
		}
		r.register(metricName, metricFamily)
		r.update(metricName, metricFamily)
	}
}

// registerFamilyIfNotRegistered will register the metric to both the
// MetricAggregateRegistry as well as the default prometheus registry
func (r *MetricAggregateRegistry) register(metricName string, metricFamily *dto.MetricFamily) {
	_, ok := r.get(metricName)
	if ok {
		return
	}

	aggregator, err := CreateMetricAggregate(metricName, metricFamily)
	if err != nil {
		glog.Infof("Ignoring metric %s: unsupported metric type by radiusd, err: %s", metricName, err)
		return // Just ignore if we have an unsupported metric
	}

	prometheus.MustRegister(aggregator.GetCollector())
	r.metricAggregateByName[metricName] = aggregator
}

func (r MetricAggregateRegistry) get(metricName string) (MetricAggregate, bool) {
	aggregator, ok := r.metricAggregateByName[metricName]
	return aggregator, ok
}

func (r *MetricAggregateRegistry) update(metricName string, metricFamily *dto.MetricFamily) {
	aggregator, ok := r.get(metricName)
	if !ok {
		// Warning is logged in 'register' method
		return // Just ignore if we have an unsupported metric
	}
	aggregator.Update(metricFamily)
}
