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
import logging
import time
from typing import Callable, Dict, List, NamedTuple, Optional, Union

import metrics_pb2
import prometheus_client.core
import requests
import snowflake
from magma.common.service_registry import ServiceRegistry
from orc8r.protos import metricsd_pb2
from orc8r.protos.common_pb2 import Void
from orc8r.protos.metricsd_pb2 import MetricsContainer
from orc8r.protos.metricsd_pb2_grpc import MetricsControllerStub
from orc8r.protos.service303_pb2_grpc import Service303Stub
from prometheus_client.parser import text_string_to_metric_families

# ScrapeTarget Holds information required to scrape and process metrics from a
# prometheus target
ScrapeTarget = NamedTuple(
    'ScrapeTarget', [
        ('url', str), ('name', str),
        ('interval', int),
    ],
)


class MetricsCollector(object):
    """
    Polls magma services periodicaly for metrics and posts them to cloud
    """
    _services = []

    def __init__(
        self, services: List[str],
        collect_interval: int,
        sync_interval: int,
        grpc_timeout: int,
        grpc_max_msg_size_mb: Union[int, float],
        loop: Optional[asyncio.AbstractEventLoop] = None,
        post_processing_fn: Optional[Callable] = None,
        scrape_targets: [ScrapeTarget] = None,
    ):
        self.sync_interval = sync_interval
        self.collect_interval = collect_interval
        self.grpc_timeout = grpc_timeout
        self.grpc_max_msg_size_bytes = grpc_max_msg_size_mb * 1024 * 1024
        self._services = services
        self._loop = loop if loop else asyncio.get_event_loop()
        self._samples_for_service = {}
        for s in self._services:
            self._samples_for_service[s] = []
        self._grpc_options = _get_metrics_chan_grpc_options(
            grpc_max_msg_size_mb,
        )
        self.scrape_targets = scrape_targets if scrape_targets else []
        # @see example_metrics_postprocessor_fn
        self.post_processing_fn = post_processing_fn

    def run(self):
        """
        Starts the service metric collection loop, the cloud sync loop, and the
        non-magma service prometheus scrape loop.
        """
        logging.info("Starting collector...")
        for s in self._services:
            self._loop.call_later(self.sync_interval, self.sync, s)
            self._loop.call_soon(self.collect, s)

        for target in self.scrape_targets:
            self._loop.call_later(
                target.interval,
                self.scrape_prometheus_target,
                target,
            )

    def sync(self, service_name):
        """
        Synchronizes sample queue for specific service to cloud and reschedules
        sync loop
        """
        if service_name in self._samples_for_service and \
           self._samples_for_service[service_name]:
            chan = ServiceRegistry.get_rpc_channel(
                'metricsd',
                ServiceRegistry.CLOUD,
                grpc_options=self._grpc_options,
            )
            client = MetricsControllerStub(chan)
            if self.post_processing_fn:
                # If services wants to, let it run a postprocessing function
                # If we throw an exception here, we'll have no idea whether
                # something was postprocessed or not, so I guess try and make it
                # idempotent?  #m sevchicken
                self.post_processing_fn(
                    self._samples_for_service[service_name],
                )

            samples = self._samples_for_service[service_name]
            sample_chunks = self._chunk_samples(samples)
            for idx, chunk in enumerate(sample_chunks):
                metrics_container = MetricsContainer(
                    gatewayId=snowflake.snowflake(),
                    family=chunk,
                )
                future = client.Collect.future(
                    metrics_container,
                    self.grpc_timeout,
                )
                future.add_done_callback(
                    self._make_sync_done_func(
                        service_name, idx,
                    ),
                )
            self._samples_for_service[service_name].clear()
        self._loop.call_later(self.sync_interval, self.sync, service_name)

    def sync_done(self, service_name, chunk, collect_future):
        """
        Sync callback to handle exceptions
        """
        err = collect_future.exception()
        if err:
            logging.error(
                "Metrics upload error for service %s (chunk %d)! "
                "[%s] %s", service_name, chunk, err.code(),
                err.details(),
            )
        else:
            logging.debug(
                "Metrics upload success for service %s (chunk %d)",
                service_name, chunk,
            )

    def collect(self, service_name):
        """
        Calls into Service303 to get service metrics samples and
        reschedule collection.
        """
        chan = ServiceRegistry.get_rpc_channel(
            service_name,
            ServiceRegistry.LOCAL,
        )
        client = Service303Stub(chan)
        future = client.GetMetrics.future(Void(), self.grpc_timeout)
        future.add_done_callback(
            lambda future:
            self._loop.call_soon_threadsafe(
                self.collect_done, service_name, future,
            ),
        )
        self._loop.call_later(
            self.collect_interval, self.collect,
            service_name,
        )

    def collect_done(self, service_name, get_metrics_future):
        """
        Collect callback to add sample results to queue or handle exceptions
        """
        err = get_metrics_future.exception()
        if err:
            logging.warning(
                "Collect %s Error! [%s] %s",
                service_name, err.code(), err.details(),
            )
            self._samples_for_service[service_name].append(
                _get_collect_success_metric(service_name, False),
            )
        else:
            container = get_metrics_future.result()
            logging.debug(
                "Collected %d from %s...",
                len(container.family), service_name,
            )
            for family in container.family:
                for sample in family.metric:
                    sample.label.add(name="service", value=service_name)
                self._samples_for_service[service_name].append(family)
                if _is_start_time_metric(family):
                    self._add_uptime_metric(service_name, family)
            self._samples_for_service[service_name].append(
                _get_collect_success_metric(service_name, True),
            )

    def _add_uptime_metric(self, service_name, family):
        if (
            not family.metric
            or len(family.metric) == 0
            or not family.metric[0].gauge.value
        ):
            logging.error("Could not parse start time metric: %s", family)
            return
        start_time = family.metric[0].gauge.value
        uptime = _get_uptime_metric(service_name, start_time)
        if uptime is not None:
            self._samples_for_service[service_name].append(uptime)

    def _make_sync_done_func(self, service_name, chunk):
        return lambda future: self._loop.call_soon_threadsafe(
            self.sync_done, service_name, chunk,
            future,
        )

    def _chunk_samples(self, samples):
        # Add 1kiB for gRPC overhead
        max_msg_bytes = self.grpc_max_msg_size_bytes - 1000

        chunked_samples = []
        chunked_samples_size = 0
        for s in samples:
            if chunked_samples_size + s.ByteSize() <= max_msg_bytes:
                chunked_samples.append(s)
                chunked_samples_size += s.ByteSize()
            else:
                yield chunked_samples
                chunked_samples = [s]
                chunked_samples_size = s.ByteSize()
        # Send leftover samples
        if chunked_samples_size > 0:
            yield chunked_samples

    def scrape_prometheus_target(self, target: ScrapeTarget) -> None:
        """
        Scrape a prometheus metrics target, convert to protobuf, send results
        to cloud. Reschedule collection if error.
        """
        try:
            r = requests.get(target.url)
            metrics = _parse_metrics_response(r.text)
            _add_scrape_label_to_metrics(metrics, target.name)
            self._package_and_send_metrics(metrics, target)

        except Exception as e:  # pylint: disable=broad-except
            logging.error(
                "Error scraping prometheus target: %s", str(e),
            )
            self._loop.call_later(
                target.interval,
                self.scrape_prometheus_target, target,
            )

    def _package_and_send_metrics(
            self, metrics: [metrics_pb2.MetricFamily],
            target: ScrapeTarget,
    ) -> None:
        """
        Send parsed and protobuf-converted metrics to cloud.
        """
        chan = ServiceRegistry.get_rpc_channel(
            'metricsd',
            ServiceRegistry.CLOUD,
            grpc_options=self._grpc_options,
        )

        client = MetricsControllerStub(chan)
        for chunk in self._chunk_samples(metrics):
            metrics_container = MetricsContainer(
                gatewayId=snowflake.snowflake(),
                family=chunk,
            )
            future = client.Collect.future(
                metrics_container,
                self.grpc_timeout,
            )
            future.add_done_callback(
                lambda future:
                self._loop.call_soon_threadsafe(
                    self.scrape_done, future, target,
                ),
            )

        self._loop.call_later(
            target.interval,
            self.scrape_prometheus_target, target,
        )

    def scrape_done(self, collect_future, target):
        """
        Log error if send fails, otherwise reschedule collection
        """
        err = collect_future.exception()
        if err:
            logging.error(
                "Prometheus Target Metrics upload error! [%s] %s",
                err.code(), err.details(),
            )
        else:
            logging.debug(
                "Prometheus Target Metrics upload success: %s",
                target.name,
            )


def _parse_metrics_response(response_text: str) -> [metrics_pb2.MetricFamily]:
    parsed_families = list(text_string_to_metric_families(response_text))
    metrics_to_send = list(map(core_metric_to_proto, parsed_families))
    # Flatten list
    return [item for sublist in metrics_to_send for item in sublist]


def _add_scrape_label_to_metrics(
        metrics: [metrics_pb2.MetricFamily],
        scrape_label: str,
) -> None:
    for family in metrics:
        for sample in family.metric:
            sample.label.add(name="scrape_target", value=scrape_label)


def _get_collect_success_metric(service_name, gw_up):
    """
    Get a the service_metrics_collected metric for a service which is either 0
    if the collection was unsuccessful or 1 if it was successful
    """
    family_proto = metrics_pb2.MetricFamily(
        type=metrics_pb2.GAUGE,
        name=str(metricsd_pb2.service_metrics_collected),
    )
    metric_proto = metrics_pb2.Metric(timestamp_ms=int(time.time() * 1000))
    metric_proto.gauge.value = 1 if gw_up else 0
    metric_proto.label.add(name="service", value=service_name)
    family_proto.metric.extend([metric_proto])
    return family_proto


def _is_start_time_metric(family):
    return family.name == str(metricsd_pb2.process_start_time_seconds)


def _get_uptime_metric(service_name, start_time):
    """
    Returns a metric for service uptime using the prometheus exposed
    process_start_time_seconds
    """
    # Metrics collection should never fail, so only log exceptions
    curr_time = calendar.timegm(time.gmtime())
    uptime = curr_time - start_time
    family_proto = metrics_pb2.MetricFamily(
        type=metrics_pb2.GAUGE,
        name=str(metricsd_pb2.process_uptime_seconds),
    )
    metric_proto = metrics_pb2.Metric(timestamp_ms=int(time.time() * 1000))
    metric_proto.gauge.value = uptime
    metric_proto.label.add(name="service", value=service_name)
    family_proto.metric.extend([metric_proto])
    return family_proto


def _get_metrics_chan_grpc_options(msg_size_mb: int):
    """
    Returns a list of gRPC options for metricsd cloud grpc channel
    :param msg_size_mb: msg size in MBs
    :return: list of tuples containing grpc options for channel
    """
    grpc_max_msg_size_bytes = msg_size_mb * 1024 * 1024
    logging.debug(
        'Setting metricsd gRPC chan Max Message Size to: %s bytes',
        grpc_max_msg_size_bytes,
    )
    return [('grpc.max_send_message_length', grpc_max_msg_size_bytes)]


def example_metrics_postprocessor_fn(
        samples: List[metrics_pb2.MetricFamily],
) -> None:
    """
    An example metrics postprocessor function for MetricsCollector

    A metrics postprocessor function can mutate samples before they are sent out
    to the metricsd cloud service. The purpose of this is usually to add labels
    to the metrics, though it is also possible to add, remove, or change the
    value of samples (though you probably shouldn't).

    Uncaught exceptions will crash the server, so if you are doing anything
    non-trivial, consider wrapping in a try/catch and figuring out whether a
    failure is fatal
    (whether you are willing to accept malformed/unprocessed stats).

    You are guaranteed that samples will only be run through this function once
    (though retries can cause delays between when this is run on samples and
    when it makes it the cloud).
    """
    failed = 0
    for family in samples:
        for sample in family.metric:
            try:
                sample.label.add(name="new_label", value="foo")
            except Exception:  # pylint: disable=broad-except
                # This operation is trivial enough that it probably shouldn't
                # be caught, but this is for example purposes. It would be a
                # bad idea to log per sample, because you could have thousands
                failed += 1
    if failed:
        logging.error("Failed to add label to %d samples!", failed)


def do_nothing_metrics_postprocessor(
        _samples: List[metrics_pb2.MetricFamily],
) -> None:
    """This metrics post processor does nothing for config examples"""


_METRIC_TYPES = ('counter', 'gauge', 'summary', 'histogram', 'untyped')


def core_metric_to_proto(
        metric: prometheus_client.core.Metric,
) -> [metrics_pb2.MetricFamily]:
    """
    Converts metrics from the prometheus client parser format to protobuf for
    sending to cloud
    """
    typ = metric.type
    if typ == 'counter':
        return [_counter_to_proto(metric)]
    elif typ == 'gauge':
        return [_gauge_to_proto(metric)]
    elif typ == 'summary':
        return _summary_to_proto(metric)
    elif typ == 'histogram':
        return _histogram_to_proto(metric)
    else:  # untyped
        return [_untyped_to_proto(metric)]


def _counter_to_proto(
        metric: prometheus_client.core.Metric,
) -> metrics_pb2.MetricFamily:
    ret = metrics_pb2.MetricFamily(name=metric.name, type=metrics_pb2.COUNTER)
    for sample in metric.samples:
        counter = metrics_pb2.Counter(value=sample[2])
        met = metrics_pb2.Metric(counter=counter)
        for key in sample[1]:
            met.label.add(name=key, value=sample[1][key])
        ret.metric.extend([met])
    return ret


def _gauge_to_proto(
        metric: prometheus_client.core.Metric,
) -> metrics_pb2.MetricFamily:
    ret = metrics_pb2.MetricFamily(name=metric.name, type=metrics_pb2.GAUGE)
    for sample in metric.samples:
        (_, labels, value, *_) = sample
        met = metrics_pb2.Metric(gauge=metrics_pb2.Gauge(value=value))
        for key in labels:
            met.label.add(name=key, value=labels[key])
        ret.metric.extend([met])
    return ret


def _summary_to_proto(
        metric: prometheus_client.core.Metric,
) -> [metrics_pb2.MetricFamily]:
    """
    1. Get metrics by unique labelset ignoring quantile
    2. convert to proto separately for each one
    """

    family_by_labelset = {}

    for sample in metric.samples:
        (name, labels, value, *_) = sample
        # get real family by checking labels (ignoring quantile)
        distinct_labels = frozenset(_remove_label(labels, 'quantile').items())
        if distinct_labels not in family_by_labelset:
            fam = metrics_pb2.MetricFamily(
                name=metric.name,
                type=metrics_pb2.SUMMARY,
            )
            summ = metrics_pb2.Summary(sample_count=0, sample_sum=0)
            fam.metric.extend([metrics_pb2.Metric(summary=summ)])
            family_by_labelset[distinct_labels] = fam

        unique_family = family_by_labelset[distinct_labels]

        if str.endswith(name, "_sum"):
            unique_family.metric[0].summary.sample_sum = value
        elif str.endswith(name, "_count"):
            unique_family.metric[0].summary.sample_count = int(value)
        elif 'quantile' in labels:
            unique_family.metric[0].summary.quantile.extend([
                metrics_pb2.Quantile(
                    quantile=float(labels['quantile']),
                    value=value,
                ),
            ])

    # Add non-quantile labels to all metrics
    for labelset in family_by_labelset.keys():
        for label in labelset:
            family_by_labelset[labelset].metric[0].label.add(
                name=label[0],
                value=label[1],
            )

    return list(family_by_labelset.values())


def _histogram_to_proto(
        metric: prometheus_client.core.Metric,
) -> [metrics_pb2.MetricFamily]:
    """
    1. Get metrics by unique labelset ignoring quantile
    2. convert to proto separately for each one
    """

    family_by_labelset = {}

    for sample in metric.samples:
        (name, labels, value, *_) = sample
        # get real family by checking labels (ignoring le)
        distinct_labels = frozenset(_remove_label(labels, 'le').items())
        if distinct_labels not in family_by_labelset:
            fam = metrics_pb2.MetricFamily(
                name=metric.name,
                type=metrics_pb2.HISTOGRAM,
            )
            hist = metrics_pb2.Histogram(sample_count=0, sample_sum=0)
            fam.metric.extend([metrics_pb2.Metric(histogram=hist)])
            family_by_labelset[distinct_labels] = fam

        unique_family = family_by_labelset[distinct_labels]

        if str.endswith(name, "_sum"):
            unique_family.metric[0].histogram.sample_sum = value
        elif str.endswith(name, "_count"):
            unique_family.metric[0].histogram.sample_count = int(value)
        elif 'le' in labels:
            unique_family.metric[0].histogram.bucket.extend([
                metrics_pb2.Bucket(
                    upper_bound=float(labels['le']),
                    cumulative_count=value,
                ),
            ])

    # Add non-quantile labels to all metrics
    for labelset in family_by_labelset.keys():
        for label in labelset:
            family_by_labelset[labelset].metric[0].label.add(
                name=str(label[0]), value=str(label[1]),
            )

    return list(family_by_labelset.values())


def _untyped_to_proto(
        metric: prometheus_client.core.Metric,
) -> metrics_pb2.MetricFamily:
    ret = metrics_pb2.MetricFamily(name=metric.name, type=metrics_pb2.UNTYPED)
    for sample in metric.samples:
        (_, labels, value, *_) = sample
        new_untyped = metrics_pb2.Untyped(value=value)
        met = metrics_pb2.Metric(untyped=new_untyped)
        for key in labels:
            met.label.add(name=key, value=labels[key])
        ret.metric.extend([met])
    return ret


def _remove_label(
        labels: Dict[str, str],
        label_to_remove: str,
) -> Dict[str, str]:
    ret = labels.copy()
    if label_to_remove in ret:
        del ret[label_to_remove]
    return ret
