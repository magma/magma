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

from magma.common import metrics_export
from magma.monitord.cpe_monitoring import SUBSCRIBER_ICMP_LATENCY_MS


class MetricTests(unittest.TestCase):
    """
    Tests for the Service303 metrics interface
    """
    def test_metrics_defined(self):
        """ Test that all metrics are defined in proto enum """
        SUBSCRIBER_ICMP_LATENCY_MS.labels('IMSI00000001').observe(10.33)

        metrics_protos = list(metrics_export.get_metrics())
        for metrics_proto in metrics_protos:
            if metrics_proto.name == "subscriber_latency_ms":
                metric = metrics_proto.metric[0]
                self.assertEqual(metric.histogram.sample_sum, 10.33)
                self.assertEqual(metric.label[0].value, 'IMSI00000001')
