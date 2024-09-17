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

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type MetricAggregate interface {
	Update(metricFamily *dto.MetricFamily)
	GetCollector() prometheus.Collector
}

type GaugeMetricAggregate struct {
	gauge *prometheus.Gauge
}

type GaugeVecMetricAggregate struct {
	gaugeVec *prometheus.GaugeVec
}

// CounterMetricAggregate uses a Gauge underneath to avoid dealing with
// counter values which may occasionally decrement.
type CounterMetricAggregate struct {
	gauge *prometheus.Gauge
}

// CounterVecMetricAggregate uses a GaugeVec underneath to avoid dealing with
// counter values which may occasionally decrement.
type CounterVecMetricAggregate struct {
	gaugeVec *prometheus.GaugeVec
}

// CreateMetricAggregate builds a new MetricAggregate if the metric is a
// prometheus Gauge, GaugeVec, Counter, or CounterVec.
func CreateMetricAggregate(metricName string, metricFamily *dto.MetricFamily) (MetricAggregate, error) {
	metricType := metricFamily.GetType()
	labelArr, err := getLabelNames(metricFamily)
	if err != nil {
		return nil, fmt.Errorf("Failed to build metric collector: %s\n", err)
	}
	metricHelp := metricFamily.GetHelp()

	if metricType == dto.MetricType_COUNTER {
		if len(labelArr) == 0 {
			gauge := prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: metricName,
					Help: metricHelp,
				},
			)
			return &CounterMetricAggregate{&gauge}, nil
		} else {
			gaugeVec := prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: metricName,
					Help: metricHelp,
				},
				labelArr,
			)
			return &CounterVecMetricAggregate{gaugeVec}, nil
		}
	} else if metricType == dto.MetricType_GAUGE {
		if len(labelArr) == 0 {
			gauge := prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: metricName,
					Help: metricHelp,
				},
			)
			return &GaugeMetricAggregate{&gauge}, nil
		} else {
			gaugeVec := prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: metricName,
					Help: metricHelp,
				},
				labelArr,
			)
			return &GaugeVecMetricAggregate{gaugeVec}, nil
		}
	}
	// Do not parse other metric types
	return nil, fmt.Errorf("Not building MetricAggregate for metric type %s", metricType)
}

func (ma *GaugeMetricAggregate) Update(metricFamily *dto.MetricFamily) {
	metricArr := metricFamily.GetMetric()
	for _, metricSample := range metricArr {
		inp := metricSample.GetGauge()
		val := inp.GetValue()
		(*ma.gauge).Set(val)
	}
}

func (ma GaugeMetricAggregate) GetCollector() prometheus.Collector {
	return *ma.gauge
}

func (ma *GaugeVecMetricAggregate) Update(metricFamily *dto.MetricFamily) {
	metricArr := metricFamily.GetMetric()
	for _, metricSample := range metricArr {
		inp := metricSample.GetGauge()
		val := inp.GetValue()
		labelValArr := getLabelVals(metricSample)
		(*ma.gaugeVec).WithLabelValues(labelValArr...).Set(val)
	}
}

func (ma GaugeVecMetricAggregate) GetCollector() prometheus.Collector {
	return *ma.gaugeVec
}

func (ma *CounterMetricAggregate) Update(metricFamily *dto.MetricFamily) {
	metricArr := metricFamily.GetMetric()
	for _, metricSample := range metricArr {
		inp := metricSample.GetCounter()
		val := inp.GetValue()
		(*ma.gauge).Set(val)
	}
}

func (ma CounterMetricAggregate) GetCollector() prometheus.Collector {
	return *ma.gauge
}

func (ma *CounterVecMetricAggregate) Update(metricFamily *dto.MetricFamily) {
	metricArr := metricFamily.GetMetric()
	for _, metricSample := range metricArr {
		inp := metricSample.GetCounter()
		val := inp.GetValue()
		labelValArr := getLabelVals(metricSample)
		(*ma.gaugeVec).WithLabelValues(labelValArr...).Set(val)
	}
}

func (ma CounterVecMetricAggregate) GetCollector() prometheus.Collector {
	return *ma.gaugeVec
}
