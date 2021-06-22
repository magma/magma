"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import unittest
import unittest.mock

import metrics_pb2
from magma.common import metrics_export
from orc8r.protos import metricsd_pb2
from prometheus_client import (
    CollectorRegistry,
    Counter,
    Gauge,
    Histogram,
    Summary,
)


class Service303MetricTests(unittest.TestCase):
    """
    Tests for the Service303 metrics interface
    """

    def setUp(self):
        self.registry = CollectorRegistry()
        self.maxDiff = None

    def test_counter(self):
        """Test that we can track counters in Service303"""
        # Add a counter with a label to the regisry
        c = Counter(
            'process_max_fds', 'A counter', ['result'],
            registry=self.registry,
        )

        # Create two series for value1 and value2
        c.labels('success').inc(1.23)
        c.labels('failure').inc(2.34)

        # Build proto outputs
        counter1 = metrics_pb2.Counter(value=1.23)
        counter2 = metrics_pb2.Counter(value=2.34)
        metric1 = metrics_pb2.Metric(
            counter=counter1,
            timestamp_ms=1234000,
        )
        metric2 = metrics_pb2.Metric(
            counter=counter2,
            timestamp_ms=1234000,
        )
        family = metrics_pb2.MetricFamily(
            name=str(metricsd_pb2.process_max_fds),
            type=metrics_pb2.COUNTER,
        )
        metric1.label.add(
            name=str(metricsd_pb2.result),
            value='success',
        )
        metric2.label.add(
            name=str(metricsd_pb2.result),
            value='failure',
        )
        family.metric.extend([metric1, metric2])

        with unittest.mock.patch('time.time') as mock_time:
            mock_time.side_effect = lambda: 1234
            self.assertCountEqual(
                list(metrics_export.get_metrics(self.registry))[0].metric,
                family.metric,
            )

    def test_gauge(self):
        """Test that we can track gauges in Service303"""
        # Add a gauge with a label to the regisry
        c = Gauge(
            'process_max_fds', 'A gauge', ['result'],
            registry=self.registry,
        )

        # Create two series for value1 and value2
        c.labels('success').inc(1.23)
        c.labels('failure').inc(2.34)

        # Build proto outputs
        gauge1 = metrics_pb2.Gauge(value=1.23)
        gauge2 = metrics_pb2.Gauge(value=2.34)
        metric1 = metrics_pb2.Metric(
            gauge=gauge1,
            timestamp_ms=1234000,
        )
        metric2 = metrics_pb2.Metric(
            gauge=gauge2,
            timestamp_ms=1234000,
        )
        family = metrics_pb2.MetricFamily(
            name=str(metricsd_pb2.process_max_fds),
            type=metrics_pb2.GAUGE,
        )
        metric1.label.add(
            name=str(metricsd_pb2.result),
            value='success',
        )
        metric2.label.add(
            name=str(metricsd_pb2.result),
            value='failure',
        )
        family.metric.extend([metric1, metric2])

        with unittest.mock.patch('time.time') as mock_time:
            mock_time.side_effect = lambda: 1234
            self.assertCountEqual(
                list(metrics_export.get_metrics(self.registry))[0].metric,
                family.metric,
            )

    def test_summary(self):
        """Test that we can track summaries in Service303"""
        # Add a summary with a label to the regisry
        c = Summary(
            'process_max_fds', 'A summary', [
                'result',
            ], registry=self.registry,
        )
        c.labels('success').observe(1.23)
        c.labels('failure').observe(2.34)

        # Build proto outputs
        summary1 = metrics_pb2.Summary(sample_count=1, sample_sum=1.23)
        summary2 = metrics_pb2.Summary(sample_count=1, sample_sum=2.34)
        metric1 = metrics_pb2.Metric(
            summary=summary1,
            timestamp_ms=1234000,
        )
        metric2 = metrics_pb2.Metric(
            summary=summary2,
            timestamp_ms=1234000,
        )
        family = metrics_pb2.MetricFamily(
            name=str(metricsd_pb2.process_max_fds),
            type=metrics_pb2.SUMMARY,
        )
        metric1.label.add(
            name=str(metricsd_pb2.result),
            value='success',
        )
        metric2.label.add(
            name=str(metricsd_pb2.result),
            value='failure',
        )
        family.metric.extend([metric1, metric2])

        with unittest.mock.patch('time.time') as mock_time:
            mock_time.side_effect = lambda: 1234
            self.assertCountEqual(
                list(metrics_export.get_metrics(self.registry))[0].metric,
                family.metric,
            )

    def test_histogram(self):
        """Test that we can track histogram in Service303"""
        # Add a histogram with a label to the regisry
        c = Histogram(
            'process_max_fds', 'A summary', ['result'],
            registry=self.registry, buckets=[0, 2, float('inf')],
        )
        c.labels('success').observe(1.23)
        c.labels('failure').observe(2.34)

        # Build proto outputs
        histogram1 = metrics_pb2.Histogram(sample_count=1, sample_sum=1.23)
        histogram1.bucket.add(upper_bound=0, cumulative_count=0)
        histogram1.bucket.add(upper_bound=2, cumulative_count=1)
        histogram1.bucket.add(upper_bound=float('inf'), cumulative_count=1)
        histogram2 = metrics_pb2.Histogram(sample_count=1, sample_sum=2.34)
        histogram2.bucket.add(upper_bound=0, cumulative_count=0)
        histogram2.bucket.add(upper_bound=2, cumulative_count=0)
        histogram2.bucket.add(upper_bound=float('inf'), cumulative_count=1)
        metric1 = metrics_pb2.Metric(
            histogram=histogram1,
            timestamp_ms=1234000,
        )
        metric2 = metrics_pb2.Metric(
            histogram=histogram2,
            timestamp_ms=1234000,
        )
        family = metrics_pb2.MetricFamily(
            name=str(metricsd_pb2.process_max_fds),
            type=metrics_pb2.HISTOGRAM,
        )
        metric1.label.add(
            name=str(metricsd_pb2.result),
            value='success',
        )
        metric2.label.add(
            name=str(metricsd_pb2.result),
            value='failure',
        )
        family.metric.extend([metric1, metric2])

        with unittest.mock.patch('time.time') as mock_time:
            mock_time.side_effect = lambda: 1234
            self.assertCountEqual(
                list(metrics_export.get_metrics(self.registry))[0].metric,
                family.metric,
            )

    def test_converted_enums(self):
        """ Test that metric names and labels are auto converted """
        # enum values (from metricsd.proto):
        # mme_new_association => 500, result => 0
        c = Counter(
            'mme_new_association', 'A counter', ['result'],
            registry=self.registry,
        )

        c.labels('success').inc(1.23)

        metric_family = list(metrics_export.get_metrics(self.registry))[0]

        self.assertEqual(
            metric_family.name,
            str(metricsd_pb2.mme_new_association),
        )
        metric_labels = metric_family.metric[0].label
        # Order not guaranteed=
        self.assertEqual(metric_labels[0].name, str(metricsd_pb2.result))
        self.assertEqual(metric_labels[0].value, 'success')


if __name__ == "__main__":
    unittest.main()
