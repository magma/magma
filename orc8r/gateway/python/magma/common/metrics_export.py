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

import logging
import time

import metrics_pb2
from orc8r.protos import metricsd_pb2
from prometheus_client import REGISTRY


def get_metrics(registry=REGISTRY, verbose=False):
    """
    Collects timeseries samples from prometheus metric collector registry
    adds a common timestamp, and encodes them to protobuf

    Arguments:
        regsitry: a prometheus CollectorRegistry instance
        verbose: whether to optimize for bandwidth and ignore metric name/help

    Returns:
        a prometheus MetricFamily protobuf stream
    """
    timestamp_ms = int(time.time() * 1000)
    for metric_family in registry.collect():
        if metric_family.type in ('counter', 'gauge'):
            family_proto = encode_counter_gauge(metric_family, timestamp_ms)
        elif metric_family.type == 'summary':
            family_proto = encode_summary(metric_family, timestamp_ms)
        elif metric_family.type == 'histogram':
            family_proto = encode_histogram(metric_family, timestamp_ms)

        if verbose:
            family_proto.help = metric_family.documentation
            family_proto.name = metric_family.name
        else:
            try:
                family_proto.name = \
                    str(metricsd_pb2.MetricName.Value(metric_family.name))
            except ValueError as e:
                logging.debug(e)  # If enum is not defined
                family_proto.name = metric_family.name
        yield family_proto


def encode_counter_gauge(family, timestamp_ms):
    """
    Takes a Counter/Gauge family which is a collection of timeseries
    samples that share a name (uniquely identified by labels) and yields
    equivalent protobufs.

    Each timeseries corresponds to a single sample tuple of the format:
        (NAME, LABELS, VALUE)

    Arguments:
        family: a prometheus gauge metric family
        timestamp_ms: the timestamp to attach to the samples
    Raises:
        ValueError if metric name is not defined in MetricNames protobuf
    Returns:
        A Counter or Gauge prometheus MetricFamily protobuf
    """
    family_proto = metrics_pb2.MetricFamily()
    family_proto.type = \
        metrics_pb2.MetricType.Value(family.type.upper())
    for sample in family.samples:
        metric_proto = metrics_pb2.Metric()
        if family_proto.type == metrics_pb2.COUNTER:
            metric_proto.counter.value = sample[2]
        elif family_proto.type == metrics_pb2.GAUGE:
            metric_proto.gauge.value = sample[2]
        # Add meta-data to the timeseries
        metric_proto.timestamp_ms = timestamp_ms
        metric_proto.label.extend(_convert_labels_to_enums(sample[1].items()))
        # Append metric sample to family
        family_proto.metric.extend([metric_proto])
    return family_proto


def encode_summary(family, timestamp_ms):
    """
    Takes a Summary Metric family which is a collection of timeseries
    samples that share a name (uniquely identified by labels) and yields
    equivalent protobufs.

    Each summary timeseries consists of sample tuples for the count, sum,
    and quantiles in the format (NAME,LABELS,VALUE). The NAME is suffixed
    with either _count, _sum to indicate count and sum respectively.
    Quantile samples will be of the same NAME with quantile label.

    Arguments:
        family: a prometheus summary metric family
        timestamp_ms: the timestamp to attach to the samples
    Raises:
        ValueError if metric name is not defined in MetricNames protobuf
    Returns:
        a Summary prometheus MetricFamily protobuf
    """
    family_proto = metrics_pb2.MetricFamily()
    family_proto.type = metrics_pb2.SUMMARY
    metric_protos = {}
    # Build a map of each of the summary timeseries from the samples
    for sample in family.samples:
        quantile = sample[1].pop('quantile', None)  # Remove from label set
        # Each time series identified by label set excluding the quantile
        metric_proto = \
            metric_protos.setdefault(
                frozenset(sample[1].items()),
                metrics_pb2.Metric(),
            )
        if sample[0].endswith('_count'):
            metric_proto.summary.sample_count = int(sample[2])
        elif sample[0].endswith('_sum'):
            metric_proto.summary.sample_sum = sample[2]
        elif quantile:
            quantile = metric_proto.summary.quantile.add()
            quantile.value = sample[2]
            quantile.quantile = _goStringToFloat(quantile)
    # Go back and add meta-data to the timeseries
    for labels, metric_proto in metric_protos.items():
        metric_proto.timestamp_ms = timestamp_ms
        metric_proto.label.extend(_convert_labels_to_enums(labels))
        # Add it to the family
        family_proto.metric.extend([metric_proto])
    return family_proto


def encode_histogram(family, timestamp_ms):
    """
    Takes a Histogram Metric family which is a collection of timeseries
    samples that share a name (uniquely identified by labels) and yields
    equivalent protobufs.

    Each summary timeseries consists of sample tuples for the count, sum,
    and quantiles in the format (NAME,LABELS,VALUE). The NAME is suffixed
    with either _count, _sum, _buckets to indicate count, sum and buckets
    respectively. Bucket samples will also contain a le to indicate its
    upper bound.

    Arguments:
        family: a prometheus histogram metric family
        timestamp_ms: the timestamp to attach to the samples
    Raises:
        ValueError if metric name is not defined in MetricNames protobuf
    Returns:
        a Histogram prometheus MetricFamily protobuf
    """
    family_proto = metrics_pb2.MetricFamily()
    family_proto.type = metrics_pb2.HISTOGRAM
    metric_protos = {}
    for sample in family.samples:
        upper_bound = sample[1].pop('le', None)  # Remove from label set
        metric_proto = \
            metric_protos.setdefault(
                frozenset(sample[1].items()),
                metrics_pb2.Metric(),
            )
        if sample[0].endswith('_count'):
            metric_proto.histogram.sample_count = int(sample[2])
        elif sample[0].endswith('_sum'):
            metric_proto.histogram.sample_sum = sample[2]
        elif sample[0].endswith('_bucket'):
            quantile = metric_proto.histogram.bucket.add()
            quantile.cumulative_count = int(sample[2])
            quantile.upper_bound = _goStringToFloat(upper_bound)
    # Go back and add meta-data to the timeseries
    for labels, metric_proto in metric_protos.items():
        metric_proto.timestamp_ms = timestamp_ms
        metric_proto.label.extend(_convert_labels_to_enums(labels))
        # Add it to the family
        family_proto.metric.extend([metric_proto])
    return family_proto


def _goStringToFloat(s):
    if s == '+Inf':
        return float("inf")
    elif s == '-Inf':
        return float("-inf")
    elif s == 'NaN':
        return float('nan')
    else:
        return float(s)


def _convert_labels_to_enums(labels):
    """
    Try to convert both the label names and label values to enum values.
    Defaults to the given name and value if it fails to convert.
    Arguments:
        labels: an array of label pairs that may contain enum names
    Returns:
        an array of label pairs with enum names converted to enum values
    """
    new_labels = []
    for name, value in labels:
        try:
            name = str(metricsd_pb2.MetricLabelName.Value(name))
        except ValueError as e:
            logging.debug(e)
        new_labels.append(metrics_pb2.LabelPair(name=name, value=value))
    return new_labels
