"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
import calendar
import logging
import time
import typing

import snowflake
import metrics_pb2
from orc8r.protos import metricsd_pb2
from orc8r.protos.common_pb2 import Void
from orc8r.protos.metricsd_pb2 import MetricsContainer
from orc8r.protos.metricsd_pb2_grpc import MetricsControllerStub
from orc8r.protos.service303_pb2_grpc import Service303Stub

from magma.common.service_registry import ServiceRegistry


class MetricsCollector(object):
    """
    Polls magma services periodicaly for metrics and posts them to cloud
    """
    _services = []

    def __init__(self, services, collect_interval,
                 sync_interval, grpc_timeout, queue_length, loop=None,
                 post_processing_fn=None):
        self.sync_interval = sync_interval
        self.collect_interval = collect_interval
        self.grpc_timeout = grpc_timeout
        self.queue_length = queue_length
        self._services = services
        self._loop = loop if loop else asyncio.get_event_loop()
        self._retry_queue = []
        self._samples = []
        # @see example_metrics_postprocessor_fn
        self.post_processing_fn = post_processing_fn

    def run(self):
        """
        Starts the service metric collection loop, and the cloud sync loop.
        """
        logging.info("Starting collector...")
        self._loop.call_later(self.sync_interval, self.sync)
        for s in self._services:
            self._loop.call_soon(self.collect, s)

    def sync(self):
        """
        Synchronizes sample queue to cloud and reschedules sync loop
        """
        if self._samples:
            chan = ServiceRegistry.get_rpc_channel('metricsd',
                                                   ServiceRegistry.CLOUD)
            client = MetricsControllerStub(chan)
            if self.post_processing_fn:
                # If services wants to, let it run a postprocessing function
                # If we throw an exception here, we'll have no idea whether
                # something was postprocessed or not, so I guess try and make it
                # idempotent?  #m sevchicken
                self.post_processing_fn(self._samples)
            samples = self._retry_queue + self._samples
            metrics_container = MetricsContainer(
                gatewayId=snowflake.snowflake(),
                family=samples
            )
            future = client.Collect.future(metrics_container, self.grpc_timeout)
            future.add_done_callback(lambda future:
                self._loop.call_soon_threadsafe(
                    self.sync_done, samples, future))
            self._retry_queue.clear()
            self._samples.clear()
        self._loop.call_later(self.sync_interval, self.sync)

    def sync_done(self, samples, collect_future):
        """
        Sync callback to handle exceptions
        """
        err = collect_future.exception()
        if err:
            self._retry_queue = samples[-self.queue_length:]
            logging.error("Metrics upload error! [%s] %s",
                          err.code(), err.details())
        else:
            logging.debug("Metrics upload success")

    def collect(self, service_name):
        """
        Calls into Service303 to get service metrics samples and
        rescheudle collection.
        """
        chan = ServiceRegistry.get_rpc_channel(service_name,
                                               ServiceRegistry.LOCAL)
        client = Service303Stub(chan)
        future = client.GetMetrics.future(Void(), self.grpc_timeout)
        future.add_done_callback(lambda future:
            self._loop.call_soon_threadsafe(
                self.collect_done, service_name, future))
        self._loop.call_later(self.collect_interval, self.collect, service_name)

    def collect_done(self, service_name, get_metrics_future):
        """
        Collect callback to add sample results to queue or handle exceptions
        """
        err = get_metrics_future.exception()
        if err:
            logging.warning("Collect %s Error! [%s] %s",
                            service_name, err.code(), err.details())
            self._samples.append(
                _get_collect_success_metric(service_name, False))
        else:
            container = get_metrics_future.result()
            logging.debug("Collected %d from %s...",
                          len(container.family), service_name)
            for family in container.family:
                for sample in family.metric:
                    sample.label.add(name="service", value=service_name)
                self._samples.append(family)
                if _is_start_time_metric(family):
                    self._add_uptime_metric(service_name, family)
            self._samples.append(
                _get_collect_success_metric(service_name, True))

    def _add_uptime_metric(self, service_name, family):
        if (not family.metric
                or len(family.metric) == 0
                or not family.metric[0].gauge.value):
            logging.error("Could not parse start time metric: %s", family)
            return
        start_time = family.metric[0].gauge.value
        uptime = _get_uptime_metric(service_name, start_time)
        if uptime is not None:
            self._samples.append(uptime)


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


def example_metrics_postprocessor_fn(
    samples: typing.List[metrics_pb2.MetricFamily]
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
    _samples: typing.List[metrics_pb2.MetricFamily]
) -> None:
    """This metrics post processor does nothing for config examples"""
