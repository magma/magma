""""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import collections
import logging
import time

import grpc
import metrics_pb2 as metrics_proto
import orc8r.protos.metricsd_pb2 as metricsd
from integ_tests.gateway.rpc import get_rpc_channel
from orc8r.protos.common_pb2 import Void
from orc8r.protos.service303_pb2 import ServiceInfo
from orc8r.protos.service303_pb2_grpc import Service303Stub


class MetricNotFoundException(Exception):
    """
    Throws when no metric can be found
    under certain metric_name, or
    a metric_name is not found under
    the given service.
    """
    pass


class GatewayServicesUtil(object):
    """
    Util to wrap the service 303 stubs for MME and SPGW so they can be
    contained and queried on together
    """

    def __init__(self):
        self._mme_service = Service303Util("mme")
        self._mobilityd_service = Service303Util("mobilityd")

    def wait_for_healthy_gateway(self, after_start_time=0):
        """
        Waits for both MME and SPGW to be healthy before returning
        Args:
            after_start_time: time since epoch in seconds, argument to wait for
                healthy gateway once it's started after a certain point. Can be
                used to detect a restart after a command. Defaults to 0 because
                all services by default start after 0 seconds since epoch
        """
        assert(self._mme_service.wait_for_healthy_service(after_start_time))
        # Mobilityd comes up after mme, so wait for this too
        assert(
            self._mobilityd_service.wait_for_healthy_service(
                after_start_time,
            )
        )

    def get_service_by_name(self, service_name):
        if service_name == 'mme':
            return self._mme_service
        if service_name == 'mobilityd':
            return self._mobilityd_service
        raise NotImplementedError("Unrecognized service_name " + service_name)

    def get_mme_service_util(self):
        return self._mme_service

    def get_mobilityd_service_util(self):
        return self._mobilityd_service


class Service303Util(object):
    """
    Util to query a single service, e.g. MME. Currently this supports getting
    the start time, uptime, and service info (status, version, health)
    """

    # Wait for total of 60 seconds
    MAX_HEALTHY_ITER = 30
    SLEEP_TIME = 4  # seconds

    def __init__(self, channel_name):
        self._service_stub = Service303Stub(get_rpc_channel(channel_name))
        self._service_name = channel_name

    @staticmethod
    def _does_metric_have_all_labels(metric, label_values):
        for label_name, label_value in label_values.items():
            if not any(
                    [
                        label_name == label.name and label_value == label.value
                        for label in metric.label
                    ],
            ):
                return False
        return True

    @staticmethod
    def _is_metric_type_supported(metric_type):
        return metric_type == metrics_proto.GAUGE \
            or metric_type == metrics_proto.COUNTER

    def get_metric_value(
        self,
        metric_name,
        expected_labels=None,
        default=None,
    ):
        """
        Get the value of a metric in the list of metrics by its name
        and its labels. Throws if metric name is not found, or if
        no metric under given metric_name matches the labels.

        Args:
             metric_name (str): encoded string name of metric.
             e.g.: "501" for ue_attach. Use str(protos.metricsd_pb2.ue_attach)
             to get "501".
             expected_labels ({str: str}): a dict of {label_name, label_value}
             that the metric should have. Optional
             e.g.: {'0', 'success'}
             default (float): default value to return if metric is not found.
             e.g.: 0.

        Returns:
            float: the gauge or counter value of the metric depending
            on its type.
        """
        expected_labels = {} if expected_labels is None else expected_labels
        metrics_list = self.get_metrics().family
        try:
            metric = next(m for m in metrics_list if m.name == metric_name)
        except StopIteration:
            if default is not None:
                return default
            else:
                raise MetricNotFoundException(
                    "No {metric_name} metric found."
                    .format(metric_name=metric_name),
                )
        metric_type = metric.type
        if not self._is_metric_type_supported(metric_type):
            raise NotImplementedError(
                "Metric value type not supported: " + metric.type,
            )
        for m in metric.metric:
            if self._does_metric_have_all_labels(m, expected_labels):
                if metric_type == metrics_proto.GAUGE:
                    return m.gauge.value
                elif metric_type == metrics_proto.COUNTER:
                    return m.counter.value
        if default is not None:
            return default
        else:
            raise MetricNotFoundException(
                "No metric under {metric_name} "
                "has all  of the label_values specified"
                .format(metric_name=metric_name),
            )

    def get_start_time(self):
        """ Get the start time of the running service """
        return self.get_metric_value(
            str(metricsd.process_start_time_seconds),
        )

    def get_uptime(self):
        """ Get the uptime (time since service start) of running service """
        return self.get_metric_value(str(metricsd.process_cpu_seconds_total))

    def get_metrics(self):
        """ Helper function to query for metrics over grpc """
        return self._service_stub.GetMetrics(Void())

    def get_service_info(self):
        """ Helper function to query for service info over grpc """
        return self._service_stub.GetServiceInfo(Void())

    def _service_started_after(self, begin_time):
        """
        Helper function to see if a service has started after the given time.
        If the given begin_time is 0, this always returns true.

        Args:
            begin_time: time since epoch in seconds. Function tests if the
                service started after this time

        Returns true if service started after begin_time, false otherwise
        """
        if (time == 0):
            return True
        uptime = self.get_uptime()
        curr_time = time.time()
        elapsed_time = curr_time - begin_time
        return elapsed_time >= uptime

    def _try_to_grpc(self, grpc_call):
        """
        Wrapper function to attempt to make a GRPC call, but safely return None
        if the GRPC channel is not ready or not up.

        Args:
            grpc_call: function to call that may throw an RPC error

        Returns result of GRPC call if available, None if it isn't ready or up
        """
        try:
            return grpc_call()
        except grpc.RpcError as error:
            err_code = error.exception().code()
            # Ignore errors caused if grpc isn't up or ready
            if (
                err_code in (
                    grpc.StatusCode.FAILED_PRECONDITION,
                    grpc.StatusCode.UNAVAILABLE,
                    grpc.StatusCode.UNKNOWN,
                    grpc.StatusCode.CANCELLED,
                )
            ):
                return None
            else:
                raise

    def wait_for_healthy_service(self, started_after_time=0):
        """
        Waits for the service to be healthy before returning
        Args:
            after_start_time: time since epoch in seconds, argument to wait for
                healthy gateway once it's started after a certain point. Can be
                used to detect a restart after a command. Defaults to 0 because
                all services by default start after 0 seconds since epoch
        """
        for _ in range(self.MAX_HEALTHY_ITER):
            info = self._try_to_grpc(self.get_service_info)

            if (
                info and info.health == ServiceInfo.APP_HEALTHY
                and info.state == ServiceInfo.ALIVE
            ):
                if self._service_started_after(started_after_time):
                    return True
                else:
                    print("%s hasn't restarted yet" % (self._service_name))
            else:
                print("%s not healthy, waiting..." % (self._service_name))
            time.sleep(self.SLEEP_TIME)
        print("max iterations hit, %s not healthy" % (self._service_name))
        return False


# Container for storing metric values
MetricValue = collections.namedtuple(
    'MetricValue', 'service name labels value',
)


def verify_gateway_metrics(test):
    """
    Decorator to verify metrics on an s1aptest. It does this by taking in the
    tested metrics and how much they should increment
    No args required, but the test case class needs to have the following
    instance variables defined:
        gateway_services: A GatewayServicesUtil instance to get metrics from
        TEST_METRICS: A list of MetricValue where the value is the expected
            increase in the gauge/counter
    """

    def wrapper(self):
        services = self.gateway_services
        initial = _get_metric_counts(services, self.TEST_METRICS)
        test(self)
        final = _get_metric_counts(services, self.TEST_METRICS)
        print("Verifying metrics")
        _verify_metric_differences(self, self.TEST_METRICS, initial, final)
        print("Metrics verified")
    # HACK to allow decorator on a unit test
    wrapper.__name__ = test.__name__
    return wrapper


def _get_metric_counts(gateway_services, metric_list):
    """
    Goes through all the specified metrics, gets their value and stores the
    value in the metric_list object with the key specified by store_name
    Args:
        service_util: A GatewayServicesUtil instance to get metrics from
        metric_list: A list of MetricValue specified as above
    Returns:
        list of MetricValue with the value as the current value
    """
    curr_values = []
    for metric in metric_list:
        service_util = gateway_services.get_service_by_name(metric.service)
        value = service_util.get_metric_value(
            metric.name,
            expected_labels=metric.labels,
            default=0,
        )
        curr_values.append(
            MetricValue(
                service=metric.service,
                name=metric.name,
                labels=metric.labels,
                value=value,
            ),
        )
    return curr_values


def _verify_metric_differences(
    test_class,
    expected_vals,
    initial_vals,
    final_vals,
):
    """
    Verifies the metrics changed by the expected amount. Asserts if not
    Args:
        test_class: the test passed into the decorator that extends
            unittest.TestCase
        expected_vals: list of MetricValue, where the value is the expected
            increase
        initial_vals: list of MetricValue, with the initial values
        final_vals: list of MetricValue, with the final values
    """
    # Since lists are all in same order, can just iterate through
    for i in range(len(expected_vals)):
        diff = final_vals[i].value - initial_vals[i].value
        expected = expected_vals[i].value
        logging.debug(
            'Metric %s %s increased by %d, expected %d',
            expected_vals[i].name,
            str(expected_vals[i].labels),
            diff,
            expected,
        )
        test_class.assertEqual(diff, expected)
    return
