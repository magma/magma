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

# Counters for IP address management
IP_ALLOCATED_TOTAL = Counter(
    'ip_address_allocated',
    'Total IP addresses allocated',
)
IP_RELEASED_TOTAL = Counter(
    'ip_address_released',
    'Total IP addresses released',
)
