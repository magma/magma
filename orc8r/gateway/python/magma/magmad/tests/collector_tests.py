"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
import calendar
import time
import unittest
import unittest.mock

from magma.common.service_registry import ServiceRegistry
from magma.magmad.metrics_collector import MetricsCollector
from metrics_pb2 import Metric, MetricFamily
from orc8r.protos import metricsd_pb2
from orc8r.protos.metricsd_pb2 import MetricsContainer


# Allow access to protected variables for unit testing
# pylint: disable=protected-access


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
        ServiceRegistry._PROXY_CONFIG = {'local_port': 1234,
                'cloud_address': 'test',
                'proxy_cloud_connections': True}

        self._services = ['test']
        self.gateway_id = "2876171d-bf38-4254-b4da-71a713952904"
        self.queue_length = 5
        self.timeout = 1
        self._collector = MetricsCollector(self._services, 5, 10,
                                           self.timeout,
                                           queue_length=self.queue_length,
                                           loop=asyncio.new_event_loop())

    @unittest.mock.patch('magma.magmad.metrics_collector.MetricsControllerStub')
    def test_sync(self, controller_mock):
        """
        Test if the collector syncs our sample.
        """
        # Mock out Collect.future
        mock = unittest.mock.Mock()
        mock.Collect.future.side_effect = [unittest.mock.Mock()]
        controller_mock.side_effect = [mock]

        # Call with no samples
        self._collector.sync()
        controller_mock.Collect.future.assert_not_called()
        self._collector._loop.stop()

        # Call with new samples to send and some to retry
        samples = [MetricFamily(name="1234")]
        self._collector._samples.extend(samples)
        self._collector._retry_queue.extend(samples)
        with unittest.mock.patch('snowflake.snowflake') as mock_snowflake:
            mock_snowflake.side_effect = lambda: self.gateway_id
            self._collector.sync()
        mock.Collect.future.assert_called_once_with(
            MetricsContainer(
                gatewayId=self.gateway_id,
                family=samples * 2),
            self.timeout)
        self.assertCountEqual(self._collector._samples, [])
        self.assertCountEqual(self._collector._retry_queue, [])

    def test_sync_queue(self):
        """
        Test if the sync queues items on failure
        """
        # We should retry sending the newest samples
        samples = [MetricFamily(name=str(i))
                   for i in range(self.queue_length + 1)]
        mock_future = MockFuture(is_error=True)
        self._collector.sync_done(samples, mock_future)
        self.assertCountEqual(self._collector._samples, [])
        self.assertCountEqual(self._collector._retry_queue,
                              samples[-self.queue_length:])

        # On success don't retry to send
        self._collector._retry_queue.clear()
        mock_future = MockFuture(is_error=False)
        self._collector.sync_done(samples, mock_future)
        self.assertCountEqual(self._collector._samples, [])
        self.assertCountEqual(self._collector._retry_queue, [])

    def test_collect(self):
        """
        Test if the collector syncs our sample.
        """
        mock = unittest.mock.MagicMock()
        samples = [MetricFamily(name="2345")]
        self._collector._samples.clear()
        self._collector._samples.extend(samples)
        mock.result.side_effect = [MetricsContainer(family=samples)]
        mock.exception.side_effect = [False]

        self._collector.collect_done('test', mock)
        # Should dequeue sample from the left, and enqueue on right
        # collector should add one more metric for collection success/failure
        self.assertEqual(len(self._collector._samples), len(samples * 2) + 1)

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
        self._collector._samples.clear()
        mock.result.side_effect = [MetricsContainer(family=samples)]
        mock.exception.side_effect = [False]

        self._collector.collect_done('test', mock)

        # should have uptime, start time, and collection success
        self.assertEqual(len(self._collector._samples), 3)
        uptime_list = [fam for fam in self._collector._samples
                       if fam.name == str(metricsd_pb2.process_uptime_seconds)]
        self.assertEqual(len(uptime_list), 1)
        self.assertEqual(len(uptime_list[0].metric), 1)
        self.assertGreater(uptime_list[0].metric[0].gauge.value, 0)

        # ensure no exceptions with empty metric
        empty = MetricFamily(name=str(metricsd_pb2.process_start_time_seconds))
        samples = [empty]
        self._collector._samples.clear()
        mock.result.side_effect = [MetricsContainer(family=samples)]
        mock.exception.side_effect = [False]
        try:
            self._collector.collect_done('test', mock)
        except Exception:   # pylint: disable=broad-except
            self.fail("Collection with empty metric should not have failed")


if __name__ == "__main__":
    unittest.main()
