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

from prometheus_client import Histogram

SUBSCRIBER_ICMP_LATENCY_MS = Histogram('subscriber_icmp_latency_ms',
                                  'Reported latency for subscriber '
                                  'in milliseconds',
                                  ['imsi'],
                                  buckets=[50, 100, 200, 500, 1000, 2000])
