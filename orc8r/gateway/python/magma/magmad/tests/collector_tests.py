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
import asyncio
import calendar
import time
import unittest
import unittest.mock
from random import randrange

import metrics_pb2
import prometheus_client
from magma.common.service_registry import ServiceRegistry
# Allow access to protected variables for unit testing
# pylint: disable=protected-access
from magma.magmad.metrics_collector import (
    MetricsCollector,
    _counter_to_proto,
    _gauge_to_proto,
    _histogram_to_proto,
    _summary_to_proto,
    _untyped_to_proto,
)
from metrics_pb2 import Metric, MetricFamily
from orc8r.protos import metricsd_pb2
from orc8r.protos.metricsd_pb2 import MetricsContainer


class MockFuture(object):
    is_error = True

    def __init__(self, is_error):
        self.is_error = is_error

    def exception(self):
        if self.is_error:
            return self.MockException()
        return None

    class MockException(object):
        def details(self):
            return ''

        def code(self):
            return 0


class MetricsCollectorTests(unittest.TestCase):
    """
    Tests for the MetricCollector collect and sync
    """

    @classmethod
    def setUpClass(cls):
        cls.queue_size = 5

    def setUp(self):
        ServiceRegistry.add_service('test', '0.0.0.0', 0)
        ServiceRegistry._PROXY_CONFIG = {
            'local_port': 1234,
            'cloud_address': 'test',
            'proxy_cloud_connections': True,
        }

        self._services = ['test']
        self.gateway_id = "2876171d-bf38-4254-b4da-71a713952904"
        self.timeout = 1
        self._collector = MetricsCollector(
            self._services, 5, 10,
            self.timeout,
            grpc_max_msg_size_mb=4,
            loop=asyncio.new_event_loop(),
        )

    @unittest.mock.patch('magma.magmad.metrics_collector.MetricsControllerStub')
    def test_sync(self, controller_mock):
        """
        Test if the collector syncs our sample.
        """
        # Mock out Collect.future
        mock = unittest.mock.Mock()
        mock.Collect.future.side_effect = [
            unittest.mock.Mock(),
            unittest.mock.Mock(),
            unittest.mock.Mock(),
        ]
        controller_mock.side_effect = [mock, mock, mock]

        # Call with no samples
        service_name = "test"
        self._collector.sync(service_name)
        controller_mock.Collect.future.assert_not_called()
        self._collector._loop.stop()

        # Call with new samples to send
        samples = [MetricFamily(name="1234")]
        self._collector._samples_for_service[service_name].extend(samples)
        with unittest.mock.patch('snowflake.snowflake') as mock_snowflake:
            mock_snowflake.side_effect = lambda: self.gateway_id
            self._collector.sync(service_name)
        mock.Collect.future.assert_called_once_with(
            MetricsContainer(
                gatewayId=self.gateway_id,
                family=samples,
            ),
            self.timeout,
        )
        self.assertCountEqual(
            self._collector._samples_for_service[service_name],
            [],
        )

        # Reduce max msg size to trigger msg chunking
        self._collector.grpc_max_msg_size_bytes = 10000
        samples = self._generate_samples(2000)
        self._collector._samples_for_service[service_name].extend(samples)
        chunked_samples = self._collector._chunk_samples(samples)
        chunk1 = next(chunked_samples)
        chunk2 = next(chunked_samples)

        with unittest.mock.patch('snowflake.snowflake') as mock_snowflake:
            mock_snowflake.side_effect = lambda: self.gateway_id
            self._collector.sync(service_name)
        mock.Collect.future.assert_any_call(
            MetricsContainer(
                gatewayId=self.gateway_id,
                family=chunk1,
            ),
            self.timeout,
        )
        mock.Collect.future.assert_any_call(
            MetricsContainer(
                gatewayId=self.gateway_id,
                family=chunk2,
            ),
            self.timeout,
        )
        self.assertCountEqual(
            self._collector._samples_for_service[service_name],
            [],
        )

    def test_collect(self):
        """
        Test if the collector syncs our sample.
        """
        mock = unittest.mock.MagicMock()
        service_name = "test"
        samples = [MetricFamily(name="2345")]
        self._collector._samples_for_service[service_name].clear()
        self._collector._samples_for_service[service_name].extend(samples)
        mock.result.side_effect = [MetricsContainer(family=samples)]
        mock.exception.side_effect = [False]

        self._collector.collect_done('test', mock)
        # Should dequeue sample from the left, and enqueue on right
        # collector should add one more metric for collection success/failure
        self.assertEqual(
            len(self._collector._samples_for_service[service_name]),
            len(samples * 2) + 1,
        )

    def test_collect_start_time(self):
        """
        Test if the collector syncs our sample.
        """
        mock = unittest.mock.MagicMock()
        start_metric = Metric()
        start_metric.gauge.value = calendar.timegm(time.gmtime()) - 1
        start_time = MetricFamily(
            name=str(metricsd_pb2.process_start_time_seconds),
            metric=[start_metric],
        )
        samples = [start_time]
        service_name = "test"
        self._collector._samples_for_service[service_name].clear()
        mock.result.side_effect = [MetricsContainer(family=samples)]
        mock.exception.side_effect = [False]

        self._collector.collect_done('test', mock)

        # should have uptime, start time, and collection success
        self.assertEqual(
            len(self._collector._samples_for_service[service_name]),
            3,
        )
        uptime_list = [
            fam for fam in self._collector._samples_for_service[service_name]
            if fam.name == str(metricsd_pb2.process_uptime_seconds)
        ]
        self.assertEqual(len(uptime_list), 1)
        self.assertEqual(len(uptime_list[0].metric), 1)
        self.assertGreater(uptime_list[0].metric[0].gauge.value, 0)

        # ensure no exceptions with empty metric
        empty = MetricFamily(name=str(metricsd_pb2.process_start_time_seconds))
        samples = [empty]
        self._collector._samples_for_service[service_name].clear()
        mock.result.side_effect = [MetricsContainer(family=samples)]
        mock.exception.side_effect = [False]
        try:
            self._collector.collect_done('test', mock)
        except Exception:  # pylint: disable=broad-except
            self.fail("Collection with empty metric should not have failed")

    def test_counter_to_proto(self):
        test_counter = prometheus_client.core.CounterMetricFamily(
            "test",
            "",
            labels=["testLabel"],
        )
        test_counter.add_metric(["val"], 1.23)
        test_counter.add_metric(["val2"], 2.34)

        proto = _counter_to_proto(test_counter)
        self.assertEqual(proto.name, test_counter.name)
        self.assertEqual(proto.type, metrics_pb2.COUNTER)

        self.assertEqual(2, len(proto.metric))
        self.assertEqual("val", proto.metric[0].label[0].value)
        self.assertEqual(1.23, proto.metric[0].counter.value)
        self.assertEqual("val2", proto.metric[1].label[0].value)
        self.assertEqual(2.34, proto.metric[1].counter.value)

    def test_gauge_to_proto(self):
        test_gauge = prometheus_client.core.GaugeMetricFamily(
            "test",
            "",
            labels=["testLabel"],
        )
        test_gauge.add_metric(["val"], 1.23)
        test_gauge.add_metric(["val2"], 2.34)

        proto = _gauge_to_proto(test_gauge)
        self.assertEqual(proto.name, test_gauge.name)
        self.assertEqual(proto.type, metrics_pb2.GAUGE)

        self.assertEqual(2, len(proto.metric))
        self.assertEqual("val", proto.metric[0].label[0].value)
        self.assertEqual(1.23, proto.metric[0].gauge.value)
        self.assertEqual("val2", proto.metric[1].label[0].value)
        self.assertEqual(2.34, proto.metric[1].gauge.value)

    def test_untyped_to_proto(self):
        test_untyped = prometheus_client.core.UntypedMetricFamily(
            "test",
            "",
            labels=["testLabel"],
        )
        test_untyped.add_metric(["val"], 1.23)
        test_untyped.add_metric(["val2"], 2.34)

        proto = _untyped_to_proto(test_untyped)
        self.assertEqual(proto.name, test_untyped.name)
        self.assertEqual(proto.type, metrics_pb2.UNTYPED)

        self.assertEqual(2, len(proto.metric))
        self.assertEqual("val", proto.metric[0].label[0].value)
        self.assertEqual(1.23, proto.metric[0].untyped.value)
        self.assertEqual("val2", proto.metric[1].label[0].value)
        self.assertEqual(2.34, proto.metric[1].untyped.value)

    def test_summary_to_proto(self):
        test_summary = prometheus_client.core.SummaryMetricFamily(
            "test",
            "",
            labels=["testLabel"],
        )
        # Add first unique labelset metrics
        test_summary.add_metric(["val1"], 10, 0.1)
        test_summary.add_sample(
            "test",
            {"quantile": "0.0", "testLabel": "val1"}, 0.01,
        )
        test_summary.add_sample(
            "test",
            {"quantile": "0.5", "testLabel": "val1"}, 0.02,
        )
        test_summary.add_sample(
            "test",
            {"quantile": "1.0", "testLabel": "val1"}, 0.03,
        )

        # Add second unique labelset metrics
        test_summary.add_metric(["val2"], 20, 0.2)
        test_summary.add_sample(
            "test",
            {"quantile": "0.0", "testLabel": "val2"}, 0.02,
        )
        test_summary.add_sample(
            "test",
            {"quantile": "0.5", "testLabel": "val2"}, 0.04,
        )
        test_summary.add_sample(
            "test",
            {"quantile": "1.0", "testLabel": "val2"}, 0.06,
        )

        protos = _summary_to_proto(test_summary)
        self.assertEqual(2, len(protos))

        for proto in protos:
            self.assertEqual(proto.name, test_summary.name)
            self.assertEqual(proto.type, metrics_pb2.SUMMARY)
            if proto.metric[0].label[0].value == "val1":
                self.assertEqual(1, len(proto.metric))
                self.assertEqual(10, proto.metric[0].summary.sample_count)
                self.assertEqual(0.1, proto.metric[0].summary.sample_sum)
                self.assertEqual(3, len(proto.metric[0].summary.quantile))
                self.assertEqual(
                    0.01,
                    proto.metric[0].summary.quantile[0].value,
                )
                self.assertEqual(
                    0.02,
                    proto.metric[0].summary.quantile[1].value,
                )
                self.assertEqual(
                    0.03,
                    proto.metric[0].summary.quantile[2].value,
                )
            else:
                self.assertEqual(1, len(proto.metric))
                self.assertEqual(20, proto.metric[0].summary.sample_count)
                self.assertEqual(0.2, proto.metric[0].summary.sample_sum)
                self.assertEqual(3, len(proto.metric[0].summary.quantile))
                self.assertEqual(
                    0.02,
                    proto.metric[0].summary.quantile[0].value,
                )
                self.assertEqual(
                    0.04,
                    proto.metric[0].summary.quantile[1].value,
                )
                self.assertEqual(
                    0.06,
                    proto.metric[0].summary.quantile[2].value,
                )

    def test_histogram_to_proto(self):
        test_hist = prometheus_client.core.HistogramMetricFamily(
            "test",
            "",
            labels=["testLabel"],
        )
        # Add first unique labelset metrics
        test_hist.add_metric(["val1"], [(1, 1), (10, 2), (100, 3)], 6)

        # Add second unique labelset metrics
        test_hist.add_metric(["val2"], [(1, 2), (10, 3), (100, 4)], 9)

        protos = _histogram_to_proto(test_hist)
        self.assertEqual(2, len(protos))

        for proto in protos:
            self.assertEqual(proto.name, test_hist.name)
            self.assertEqual(proto.type, metrics_pb2.HISTOGRAM)
            if proto.metric[0].label[0].value == "val1":
                self.assertEqual(1, len(proto.metric))
                self.assertEqual(3, proto.metric[0].histogram.sample_count)
                self.assertEqual(6, proto.metric[0].histogram.sample_sum)
                self.assertEqual(3, len(proto.metric[0].histogram.bucket))
                self.assertEqual(
                    1, proto.metric[0].histogram.bucket[
                        0
                    ].cumulative_count,
                )
                self.assertEqual(
                    2, proto.metric[0].histogram.bucket[
                        1
                    ].cumulative_count,
                )
                self.assertEqual(
                    3, proto.metric[0].histogram.bucket[
                        2
                    ].cumulative_count,
                )
            else:
                self.assertEqual(1, len(proto.metric))
                self.assertEqual(4, proto.metric[0].histogram.sample_count)
                self.assertEqual(9, proto.metric[0].histogram.sample_sum)
                self.assertEqual(3, len(proto.metric[0].histogram.bucket))
                self.assertEqual(
                    2, proto.metric[0].histogram.bucket[
                        0
                    ].cumulative_count,
                )
                self.assertEqual(
                    3, proto.metric[0].histogram.bucket[
                        1
                    ].cumulative_count,
                )
                self.assertEqual(
                    4, proto.metric[0].histogram.bucket[
                        2
                    ].cumulative_count,
                )

    def test_chunk_samples(self):
        # Generate 1.2 kbs of metric samples
        test_collector = MetricsCollector(
            self._services, 5, 10,
            self.timeout,
            grpc_max_msg_size_mb=0.01,
            loop=asyncio.new_event_loop(),
        )
        samples = self._generate_samples(2000)
        chunked_samples = test_collector._chunk_samples(samples)
        self.assertEqual(len(list(chunked_samples)), 2)

    def _generate_samples(self, number):
        samples = []
        for _ in range(number):
            sample_name = randrange(10000)
            samples.append(MetricFamily(name=str(sample_name)))
        return samples


if __name__ == "__main__":
    unittest.main()
