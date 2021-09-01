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

from prometheus_client import Counter

MAGMA_BYTES_SENT_TOTAL = Counter(
    'magma_bytes_sent_total',
    'Total bytes sent from Magma gateway services to Orc8r services',
    ['service_name', 'dest_service'],
)

LINUX_BYTES_SENT_TOTAL = Counter(
    'linux_bytes_sent_total',
    'Total bytes sent from non-Magma-service binaries to any destination',
    ['binary_name'],
)
